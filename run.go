package main

import (
	"fmt"
	"log"
	"time"

	"github.com/getlantern/systray"
	"github.com/go-ini/ini"
)

const (
	defaultDuration      = 100 * time.Millisecond
	cpuUsageUpdatePeriod = 1 * time.Second
)

var appContext *RunAnimal

func init() {
	appContext = NewRunAnimal(Config.RunAnimal, setUpIcons())
}

func onReady() {
	appContext.start()
}

func onExit() {
	appContext.end()
}

// RunAnimal struct holds the state for the system tray application.
type RunAnimal struct {
	runningAnimal      string
	runningAnimalIcons map[string][][]byte
	taskTray           *systray.MenuItem
	animalMenuItems    map[string]*systray.MenuItem
	currentCPU         float64
}

// NewRunAnimal creates a new RunAnimal instance.
func NewRunAnimal(runningAnimal string, runningAnimalIcons map[string][][]byte) *RunAnimal {
	return &RunAnimal{
		runningAnimal:      runningAnimal,
		runningAnimalIcons: runningAnimalIcons,
		animalMenuItems:    make(map[string]*systray.MenuItem),
	}
}

// Start initializes the system tray interface and starts background processes.
func (m *RunAnimal) start() {
	m.taskTray = systray.AddMenuItem(fmt.Sprintf("CPU: %.2f%%", m.currentCPU), "GoCat CPU Usage")

	animalNamesMapping := map[string]string{
		"white_cat":    "White Cat",
		"black_cat":    "Black Cat",
		"white_horse":  "White Horse",
		"black_horse":  "Black Horse",
		"white_parrot": "White Parrot",
		"black_parrot": "Black Parrot",
	}

	for animal, name := range animalNamesMapping {
		if _, exists := m.runningAnimalIcons[animal]; exists {
			m.addMenuItem(animal, name)
		}
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
func (m *RunAnimal) end() {
	cfg, err := ini.Load(ConfigPath)
	if err != nil {
		log.Printf("action=End status=fail err=%s", err)
		return
	}

	cfg.Section("animal").Key("run_animal").SetValue(m.runningAnimal)
	if err = cfg.SaveTo(ConfigPath); err != nil {
		log.Printf("action=End status=fail err=%s", err)
	}
}

// addMenuItem adds a menu item for switching the running animal.
func (m *RunAnimal) addMenuItem(animal, name string) {
	menu := systray.AddMenuItem(name, "Switch to "+animal)
	m.animalMenuItems[animal] = menu

	go func() {
		for range menu.ClickedCh {
			m.changeAnimal(animal)
		}
	}()
}

// changeAnimal updates the currently displayed animal icon.
func (m *RunAnimal) changeAnimal(animal string) {
	m.runningAnimal = animal
	if icons, exists := m.runningAnimalIcons[animal]; exists && len(icons) > 0 {
		systray.SetIcon(icons[0])
	}
}

// updateCPUUsage continuously updates the CPU usage information.
// It updates the systray item title and tooltip to reflect the current CPU usage percentage.
func (m *RunAnimal) updateCPUUsage() {
	ticker := time.NewTicker(cpuUsageUpdatePeriod)
	defer ticker.Stop()

	for range ticker.C {
		m.currentCPU = getCPUUsage()
		cpuText := fmt.Sprintf("CPU: %.2f%%", m.currentCPU)
		m.taskTray.SetTitle(cpuText)
		systray.SetTooltip(cpuText)
	}
}

// run controls the display of animal icons based on CPU usage.
func (m *RunAnimal) run() {
	ticker := time.NewTicker(defaultDuration)
	defer ticker.Stop()

	for {
		m.currentCPU = getCPUUsage()
		newDuration := time.Duration((1-m.currentCPU/100)*50+50) * time.Millisecond
		ticker.Reset(newDuration)

		for _, icon := range m.runningAnimalIcons[m.runningAnimal] {
			<-ticker.C
			systray.SetIcon(icon)
		}
	}
}
