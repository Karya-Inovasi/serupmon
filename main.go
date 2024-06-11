package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/karyainovasiab/serupmon/config"
	"github.com/karyainovasiab/serupmon/monitor"
	"github.com/karyainovasiab/serupmon/service"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
	commit  = "none"
	dateStr = "2024-06-09 18:16:00"

	configPath string
	prefixPath string
	monitors   []*monitor.Monitor
)

func main() {
	date, _ := time.Parse("2006-01-02 15:04:05", dateStr)
	app := &cli.App{
		Name:        "serupmon",
		Usage:       "A simple server up/down monitoring tool",
		Compiled:    date,
		Copyright:   "Copyright (c) 2024 Serupmon Authors",
		HideVersion: false,
		Version:     version + " (" + commit + ")" + " built at " + dateStr,
		Commands: []*cli.Command{
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "start supermon monitoring service",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "config",
						Usage:       "load configuration from `file`",
						Aliases:     []string{"c"},
						Required:    true,
						Destination: &configPath,
					},
					&cli.StringFlag{
						Name:        "prefix",
						Usage:       "prefix path for the running service",
						Aliases:     []string{"p"},
						Value:       "/tmp/serupmon",
						Required:    false,
						Destination: &prefixPath,
					},
				},
				Action: func(c *cli.Context) error {
					cfg, err := config.LoadConfig(configPath)
					if err != nil {
						return err
					}

					loadMonitor(cfg)

					fmt.Printf("=> prefix   : %s\n", prefixPath)
					fmt.Printf("=> config   : %s\n", configPath)
					fmt.Printf("=> monitors : %v\n", len(monitors))

					startService()

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func loadMonitor(cfg *config.Config) {
	for _, monitorCfg := range cfg.Monitor {
		for _, serviceCfg := range monitorCfg.Services {
			switch serviceCfg.Type {
			case "http":
				m := monitor.NewHTTPMonitor(
					monitorCfg.Name,
					serviceCfg.Upstream,
					serviceCfg.Interval,
					serviceCfg.Threshold,
					serviceCfg.Timeout,
				)

				if serviceCfg.Alert != nil {
					if serviceCfg.Alert.Telegram != nil && serviceCfg.Alert.Email != nil {
						if serviceCfg.Alert.Telegram.Enabled && serviceCfg.Alert.Email.Enabled {
							m.SetAlertChannel(monitor.ALL_NOTIFIER)

							tgToken := ""
							tgChatID := ""

							if serviceCfg.Alert.Telegram.Config != nil {
								tgToken = serviceCfg.Alert.Telegram.Config.Token
								tgChatID = serviceCfg.Alert.Telegram.Config.ChatID
							}

							if strings.TrimSpace(tgToken) == "" {
								// for now just use global telegram config
								tgToken = cfg.Global.Alert.Telegram.Config.Token
							}

							if strings.TrimSpace(tgChatID) == "" {
								// for now just use global telegram config
								tgChatID = cfg.Global.Alert.Telegram.Config.ChatID
							}

							m.SetTelegramConfig(tgToken, tgChatID)
						}
					} else if serviceCfg.Alert.Telegram != nil {
						if serviceCfg.Alert.Telegram.Enabled {
							m.SetAlertChannel(monitor.TELEGRAM_NOTIFIER)

							tgToken := ""
							tgChatID := ""

							if serviceCfg.Alert.Telegram.Config != nil {
								tgToken = serviceCfg.Alert.Telegram.Config.Token
								tgChatID = serviceCfg.Alert.Telegram.Config.ChatID
							}

							if strings.TrimSpace(tgToken) == "" {
								// for now just use global telegram config
								tgToken = cfg.Global.Alert.Telegram.Config.Token
							}

							if strings.TrimSpace(tgChatID) == "" {
								// for now just use global telegram config
								tgChatID = cfg.Global.Alert.Telegram.Config.ChatID
							}

							m.SetTelegramConfig(tgToken, tgChatID)
						}
					} else if serviceCfg.Alert.Email != nil {
						if serviceCfg.Alert.Email.Enabled {
							m.SetAlertChannel(monitor.EMAIL_NOTIFIER)
						}
					}
				}

				monitors = append(monitors, m)
			}
		}
	}
}

func startService() {
	service.EnsureSerupmonInitialized(prefixPath)

	monitor.StartMonitor(monitors)
}

func stopService() {
	service.EnsureSerupmonInitialized(prefixPath)
}

func restartService() {
	service.EnsureSerupmonInitialized(prefixPath)
}
