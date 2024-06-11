package monitor

import (
	"testing"
	"time"
)

func TestMonitor_CheckHTTP(t *testing.T) {
	m := &Monitor{
		Name:      "TestMonitor",
		Type:      HTTP,
		Upstream:  "http://example.com",
		Interval:  15,
		Threshold: 3,
		Timeout:   10,
		lastState: UP,
	}

	err := checkHTTP(m)

	if err != nil {
		t.Errorf("checkHTTP() returned an error: %v", err)
	}
}

func TestMonitor_CheckTCP(t *testing.T) {
	m := &Monitor{
		Name:      "TestMonitor",
		Type:      TCP,
		Upstream:  "example.com:80",
		Interval:  15,
		Threshold: 3,
		Timeout:   10,
		lastState: UP,
	}

	err := checkTCP(m)

	if err != nil {
		t.Errorf("checkTCP() returned an error: %v", err)
	}
}

func TestMonitor_StartMonitor(t *testing.T) {
	monitors := []*Monitor{
		{
			Name:      "TestMonitor1",
			Type:      HTTP,
			Upstream:  "http://example.com",
			Interval:  15,
			Threshold: 3,
			Timeout:   10,
			lastState: UP,
		},
		{
			Name:      "TestMonitor2",
			Type:      TCP,
			Upstream:  "example.com:80",
			Interval:  15,
			Threshold: 3,
			Timeout:   10,
			lastState: UP,
		},
	}

	go StartMonitor(monitors)

	// Wait for some time to allow monitoring to happen
	time.Sleep(30 * time.Second)
}
