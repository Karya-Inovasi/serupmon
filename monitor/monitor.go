package monitor

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/karyainovasiab/serupmon/notifier"
)

type MonitorType string

type MonitorState string

type MonitorNotifierChannel string

type Monitor struct {
	Name           string
	Type           MonitorType
	Upstream       string
	Interval       int
	Threshold      int
	Timeout        int
	failCount      int
	lastState      MonitorState
	mu             sync.Mutex
	timeStart      time.Time
	timeEnd        time.Time
	alertChannel   MonitorNotifierChannel
	telegramConfig struct {
		ChatId string
		Token  string
	}
}

const (
	HTTP MonitorType = "http"
	TCP  MonitorType = "tcp"

	UP   MonitorState = "up"
	DOWN MonitorState = "down"

	TELEGRAM_NOTIFIER MonitorNotifierChannel = "telegram"
	EMAIL_NOTIFIER    MonitorNotifierChannel = "email"
	ALL_NOTIFIER      MonitorNotifierChannel = "all"
)

func New(name, upstream string, t MonitorType, interval, threshold, timeout *int) *Monitor {
	_timeout := 10
	_interval := 15
	_threshold := 3

	if interval != nil {
		_interval = *interval
		if _interval < 1 {
			_interval = 15
		}
	}

	if threshold != nil {
		_threshold = *threshold
		if _threshold < 1 {
			_threshold = 3
		}
	}

	if timeout != nil {
		_timeout = *timeout
		if _timeout < 1 {
			_timeout = 10
		}
	}

	return &Monitor{
		Name:      name,
		Type:      t,
		Upstream:  upstream,
		Interval:  _interval,
		Threshold: _threshold,
		Timeout:   _timeout,
		lastState: UP,
	}
}

func (m *Monitor) SetAlertChannel(notifier MonitorNotifierChannel) {
	m.alertChannel = notifier
}

func (m *Monitor) SetTelegramConfig(token, chatID string) {
	m.telegramConfig.Token = token
	m.telegramConfig.ChatId = chatID
}

func check(m *Monitor) error {
	switch m.Type {
	case HTTP:
		return checkHTTP(m)
	case TCP:
		return checkTCP(m)
	default:
		return nil
	}
}

func (m *Monitor) notifyAlertChannel(message string) error {
	switch m.alertChannel {
	case TELEGRAM_NOTIFIER:
		if m.telegramConfig.Token == "" || m.telegramConfig.ChatId == "" {
			return fmt.Errorf("telegram token or chat id is empty")
		}

		return notifier.TelegramNotifier.Notify(m.telegramConfig.ChatId, message, m.telegramConfig.Token)
	case EMAIL_NOTIFIER:
		// return emailNotifier.Notify(message)
	case ALL_NOTIFIER:
		// return allNotifier.Notify(message)
	default:
		return nil
	}

	return nil
}

func NewHTTPMonitor(name, upstream string, interval, threshold, timeout *int) *Monitor {
	return New(name, upstream, HTTP, interval, threshold, timeout)
}

func NewTCPMonitor(name, upstream string, interval, threshold, timeout *int) *Monitor {
	return New(name, upstream, TCP, interval, threshold, timeout)
}

func checkHTTP(m *Monitor) error {
	check, err := http.Get(m.Upstream)
	if err != nil {
		return err
	}

	defer check.Body.Close()

	if check.StatusCode < 200 || check.StatusCode >= 300 {
		return fmt.Errorf("error response code: %d", check.StatusCode)
	}

	return nil
}

func checkTCP(m *Monitor) error {
	return nil
}

func StartMonitor(monitors []*Monitor) {
	for _, m := range monitors {
		go func(m *Monitor) {
			ticker := time.NewTicker(time.Duration(m.Interval) * time.Second)
			for range ticker.C {
				msg := "[%s] monitor %s CHECK! | type=%s, upstream=%s, interval=%d, threshold=%d, timeout=%d\n"
				log := fmt.Sprintf(msg, time.Now().Format(time.RFC3339), m.Name, m.Type, m.Upstream, m.Interval, m.Threshold, m.Timeout)
				logCheck(log)
				if err := check(m); err != nil {
					m.mu.Lock()
					m.failCount++
					if m.failCount >= m.Threshold && m.lastState != DOWN {
						m.lastState = DOWN
						m.timeStart = time.Now()
						msg := "[%s] monitor %s DOWN! | type=%s, upstream=%s, error=%v, down_time=%s, fail_count=%d, threshold=%d\n"
						log := fmt.Sprintf(msg, time.Now().Format(time.RFC3339), m.Name, m.Type, m.Upstream, err, m.timeStart, m.failCount, m.Threshold)
						logDown(log, err, m)
					}
					m.mu.Unlock()
				} else {
					m.mu.Lock()
					if m.lastState == DOWN {
						m.timeEnd = time.Now()
						downDuration := m.timeEnd.Sub(m.timeStart)
						msg := "[%s] monitor %s UP! | type=%s, upstream=%s, down_time=%s, up_time=%s, down_duration=%s\n"
						log := fmt.Sprintf(msg, time.Now().Format(time.RFC3339), m.Name, m.Type, m.Upstream, m.timeStart, m.timeEnd, downDuration.String())
						logUp(log, m)
					}

					m.failCount = 0
					if m.lastState != UP {
						m.lastState = UP
					}

					m.mu.Unlock()
				}
			}
		}(m)
	}

	select {}
}

func logCheck(log string) {
	fmt.Fprint(os.Stdout, log)
}

func logDown(log string, err error, m *Monitor) {
	fmt.Fprint(os.Stderr, log)
	msg := "<b>🔴 Service DOWN !!</b>\n"
	msg += fmt.Sprintf("<i>We have detected that service <b>%s</b> is down and unreachable at <u>%s</u></i>\n\n", m.Name, normalizeDateTime(m.timeStart))
	msg += "<pre><code>"
	msg += fmt.Sprintf(`{
		"Name": "%s",
		"Type": "%s",
		"Upstream": "%s",
		"Error": "%v",
		"Down Time": "%s",
		"Fail Count": %d,
		"Threshold": %d
	}`, m.Name, m.Type, m.Upstream, err, normalizeDateTime(m.timeStart), m.failCount, m.Threshold)
	msg += "</code></pre>"
	// msg += "\n\n"
	// msg += "<pre><code>" + log + "</code></pre>"
	m.notifyAlertChannel(msg)
}

func logUp(log string, m *Monitor) {
	fmt.Fprint(os.Stdout, log)
	msg := "<b>🟢 Service UP !!</b>\n"
	msg += fmt.Sprintf("<i>We have detected that service <b>%s</b> is up and reachable at <u>%s</u></i>\n\n", m.Name, normalizeDateTime(m.timeEnd))
	msg += "<pre><code>"
	msg += fmt.Sprintf(`{
		"Name": "%s",
		"Type": "%s",
		"Upstream": "%s",
		"Down Time": "%s",
		"Up Time": "%s",
		"Down Duration": "%s"
	}`, m.Name, m.Type, m.Upstream, normalizeDateTime(m.timeStart), normalizeDateTime(m.timeEnd), m.timeEnd.Sub(m.timeStart).String())
	msg += "</code></pre>"
	// msg += "\n\n"
	// msg += "<pre><code>" + log + "</code></pre>"
	m.notifyAlertChannel(msg)
}

func normalizeDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
