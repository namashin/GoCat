package app

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
	"github.com/go-ini/ini"
	"github.com/shirou/gopsutil/cpu"
)

const (
	ConfigPath           = "./config/config.ini"
	defaultDuration      = 100 * time.Millisecond
	cpuUsageUpdatePeriod = 1 * time.Second
)

// RunAnimal struct holds the state for the system tray application.
type RunAnimal struct {
	runningAnimal      string
	runningAnimalIcons map[string][][]byte          // Map of animal names to their corresponding icon data
	tray               *systray.MenuItem            // Tray item displaying CPU usage
	animalMenuItems    map[string]*systray.MenuItem // Menu items for each animal
	currentCPU         float64                      // Current CPU usage percentage
}

// NewRunAnimal creates a new RunAnimal instance.
// It initializes the running animal and its icons and prepares the menu items map.
func NewRunAnimal(runningAnimal string, runningAnimalIcons map[string][][]byte) *RunAnimal {
	return &RunAnimal{
		runningAnimal:      runningAnimal,
		runningAnimalIcons: runningAnimalIcons,
		animalMenuItems:    make(map[string]*systray.MenuItem),
	}
}

// Start initializes the system tray interface and starts background processes.
// It sets up the main tray item, animal switching menu items, and a quit option.
// It also launches goroutines for updating CPU usage and icon display based on the selected animal.
func (m *RunAnimal) Start() {
	m.tray = systray.AddMenuItem("CPU: "+fmt.Sprintf("%.2f%%", m.currentCPU), "GoCat CPU Usage")

	animalNamesMapping := map[string]string{
		"white_cat":    "White Cat",
		"black_cat":    "Black Cat",
		"white_horse":  "White Horse",
		"black_horse":  "Black Horse",
		"white_parrot": "White Parrot",
		"black_parrot": "Black Parrot",
	}

	for animal := range m.runningAnimalIcons {
		menu := systray.AddMenuItem(animalNamesMapping[animal], "Switch to "+animal)
		m.animalMenuItems[animal] = menu

		go func(animal string, menuItem *systray.MenuItem) {
			for range menuItem.ClickedCh {
				m.changeAnimal(animal)
			}
		}(animal, menu)
	}

	quit := systray.AddMenuItem("Quit", "Quit the application")
	go func() {
		<-quit.ClickedCh
		systray.Quit()
	}()

	go m.updateCPUUsage()
	go m.run()
}

// End saves the current state to a configuration file upon application exit.
// It handles errors during file loading and saving and updates the config with the current running animal.
func (m *RunAnimal) End() {
	cfg, err := ini.Load(ConfigPath)
	if err != nil {
		return
	}

	animalSection := cfg.Section("animal")
	animalSection.Key("run_animal").SetValue(m.runningAnimal)

	cfg.SaveTo(ConfigPath)
}

// changeAnimal updates the currently displayed animal icon.
// It is triggered by user interaction with the systray menu.
func (m *RunAnimal) changeAnimal(animal string) {
	m.runningAnimal = animal
}

// updateCPUUsage continuously updates the CPU usage information.
// It updates the systray item title to reflect the current CPU usage percentage.
func (m *RunAnimal) updateCPUUsage() {
	for {
		time.Sleep(cpuUsageUpdatePeriod)
		m.tray.SetTitle("CPU: " + fmt.Sprintf("%.2f%%", m.currentCPU))
	}
}

// run controls the timing and display of the animal icons based on CPU usage.
// It adjusts the refresh rate of icon updates based on CPU load to manage application responsiveness.
func (m *RunAnimal) run() {
	ticker := time.NewTicker(defaultDuration)
	defer ticker.Stop()

	updateTicker := func() {
		m.currentCPU = getCPUUsage()
		newDuration := time.Duration((1-m.currentCPU/100)*50+50) * time.Millisecond
		ticker.Reset(newDuration)
	}

	for {
		updateTicker()

		for _, icon := range m.runningAnimalIcons[m.runningAnimal] {
			<-ticker.C
			systray.SetIcon(icon)
		}
	}
}

// getCPUUsage retrieves the current CPU usage percentage.
// It handles potential errors in fetching CPU data and returns 0 if data is unavailable.
func getCPUUsage() float64 {
	percent, err := cpu.Percent(0, false)
	if err != nil || len(percent) <= 0 {
		return 0
	}

	return percent[0]
}
