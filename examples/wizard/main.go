package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/atoms"
	"github.com/trippwill/polymer/gels/menu"
	"github.com/trippwill/polymer/host"
	"github.com/trippwill/polymer/router"
	"github.com/trippwill/polymer/util"
)

// NamePromptScreen prompts the user to enter their name.
type NamePromptScreen struct {
	atoms.NilInit
	id    string
	input string
}

func NewNamePromptScreen() NamePromptScreen {
	return NamePromptScreen{
		id: util.NewUniqeTypeId[NamePromptScreen](),
	}
}

func (n NamePromptScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes:
			n.input += string(msg.Runes)
			return n, nil
		case tea.KeyBackspace:
			if len(n.input) > 0 {
				n.input = n.input[:len(n.input)-1]
			}
			return n, nil
		case tea.KeyEnter:
			return nil, host.Context(n.input)
		}
	}
	return n, nil
}

func (n NamePromptScreen) View() string {
	return "Enter your name: " + n.input + "\n"
}

// GreetingScreen greets the user by name.
type GreetingScreen struct {
	atoms.NilInit
	id    string // Unique identifier for the screen
	value string
}

func NewGreetingScreen(value string) GreetingScreen {
	return GreetingScreen{
		id:    util.NewUniqeTypeId[GreetingScreen](),
		value: value,
	}
}

func (g GreetingScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		return nil, nil
	}
	return g, nil
}

func (g GreetingScreen) View() string {
	return "\nHello, " + g.value + "! (press any key to return)\n"
}

func (g *GreetingScreen) SetContext(ctx string) {
	g.value = ctx // Set the context to the name entered
}

type app struct {
	atoms.NilInit
	router.Auto[tea.Model]        // Router to manage screens
	id                     string // Unique identifier for the app
}

func NewApp() app {
	menu := menu.NewMenu(
		"Name Wizard",
		menu.AdaptItem(NewNamePromptScreen(), "Name Wizard", "Enter your name"),
		menu.AdaptItem(menu.NewQuitAtom(), "Quit", "Exit Application"),
	)

	return app{
		Auto: router.NewAuto(menu, nil),
		id:   util.NewUniqeTypeId[app](),
	}
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if name, ok := msg.(host.ContextMsg[string]); ok {
		greetingScreen := NewGreetingScreen(*name.Context)
		a.Auto = a.Set(router.AutoSlotOverride, greetingScreen)
		return a, nil
	}

	a.Auto = a.Set(router.AutoSlotOverride, nil)

	var cmd tea.Cmd
	a.Auto, cmd = a.Route(msg)

	// The menu has quit
	if !a.IsSet(router.AutoSlotPrimary) {
		return a, tea.Quit
	}

	return a, cmd
}

func (a app) View() string { return a.Render() }

func main() {
	// Set up a standard logger for tracing
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	// Create the host and start the Bubble Tea program
	host := host.NewContextHost(
		"Polymer Integration Example",
		NewApp(),
		"",
	)

	p := tea.NewProgram(host)
	if _, err := p.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
