package service

import (
	"fmt"
	"os"
)

const (
	pidFile = "$PREFIX/run/serupmon.pid"
	logFile = "$PREFIX/log/serupmon.log"
)

func InstallService() {
	// Install serupmon as a service
}

func UninstallService() {
	// Uninstall serupmon service
}

func EnsureSerupmonInitialized(prefix string) {
	if err := dirExistsOrMkdir(prefix); err != nil {
		fmt.Printf("failed to create directory: %v\n", err)
		os.Exit(1)
	}

	if err := dirExistsOrMkdir(prefix + "/run"); err != nil {
		fmt.Printf("failed to create directory: %v\n", err)
		os.Exit(1)
	}

	if err := dirExistsOrMkdir(prefix + "/log"); err != nil {
		fmt.Printf("failed to create directory: %v\n", err)
		os.Exit(1)
	}
}

func dirExistsOrMkdir(dir string) error {
	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}

	if !stat.IsDir() {
		return fmt.Errorf("unable to create directory: %s because a file with the same name already exists", dir)
	}

	return nil
}
