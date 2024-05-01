package app

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

type RunAnimal struct {
	runningAnimal      string
	runningAnimalIcons map[string][][]byte
	tray               *systray.MenuItem
}

func NewRunAnimal(runningAnimal string, runningAnimalIcons map[string][][]byte) *RunAnimal {
	return &RunAnimal{
		runningAnimal:      runningAnimal,
		runningAnimalIcons: runningAnimalIcons,
	}
}

func (m *RunAnimal) Start() {
	m.tray = systray.AddMenuItem("CPU: "+getCPUUsage(), "GoCat CPU Usage")

	go m.updateCPUUsage()
	go m.run()
}

func (m *RunAnimal) updateCPUUsage() {
	for {
		time.Sleep(100 * time.Millisecond)
		m.tray.SetTitle("CPU: " + getCPUUsage())
	}
}

func (m *RunAnimal) run() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		for _, icon := range m.runningAnimalIcons[m.runningAnimal] {
			<-ticker.C
			systray.SetIcon(icon)
		}
	}
}

func getCPUUsage() string {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return "N/A"
	}

	if len(percent) <= 0 {
		return "N/A"
	}

	return fmt.Sprintf("%.2f%%", percent[0])
}
