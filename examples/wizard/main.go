package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/atoms"
	"github.com/trippwill/polymer/gels/menu"
	"github.com/trippwill/polymer/poly"
	"github.com/trippwill/polymer/router"
	"github.com/trippwill/polymer/util"
)

// NamePromptScreen is an Atom that prompts the user to enter their name.
type NamePromptScreen struct {
	id    string // Unique identifier for the screen
	input string
}

func NewNamePromptScreen() NamePromptScreen {
	return NamePromptScreen{
		id: util.NewUniqeTypeId[NamePromptScreen](),
	}
}

var _ poly.Atomic[any] = (*NamePromptScreen)(nil)

func (n NamePromptScreen) Update(msg tea.Msg) (poly.Atomic[any], tea.Cmd) {
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
			return nil, util.Broadcast(NewGreetingScreen(n.input))
		}
	}
	return n, nil
}

func (n NamePromptScreen) View() string {
	return "Enter your name: " + n.input + "\n"
}

func (n NamePromptScreen) SetContext(ctx any) poly.Atomic[any] {
	return n // No context needed for NamePromptScreen
}

// GreetingScreen is an Model that greets the user by name.
type GreetingScreen struct {
	id    string // Unique identifier for the screen
	Value string
}

func NewGreetingScreen(name string) GreetingScreen {
	return GreetingScreen{
		id:    util.NewUniqeTypeId[any](),
		Value: name,
	}
}

var _ poly.Atomic[any] = (*GreetingScreen)(nil)

func (g GreetingScreen) Update(msg tea.Msg) (poly.Atomic[any], tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		return nil, nil
	}
	return g, nil
}

func (g GreetingScreen) View() string {
	return "\nHello, " + g.Value + "! (press any key to return)\n"
}

func (g GreetingScreen) SetContext(ctx any) poly.Atomic[any] {
	return g
}

type app struct {
	router router.Router[poly.Atomic[any], poly.Atomic[any]] // Router to manage screens
	id     string                                            // Unique identifier for the app
}

func NewApp() poly.Atomic[any] {
	// Create the menu with the name prompt and greeting screens
	menu := menu.NewMenu(
		"Name Wizard",
		menu.AdaptItem(NewNamePromptScreen(), "Name Wizard", "Enter your name"),
		menu.AdaptItem(atoms.NewQuitAtom(), "Quit", "Exit Application"),
	)

	return app{
		router: router.NewRouter[poly.Atomic[any], poly.Atomic[any]](menu, nil, router.SlotT),
		id:     util.NewUniqeTypeId[app](),
	}
}

func (a app) Update(msg tea.Msg) (poly.Atomic[any], tea.Cmd) {
	if greeting, ok := msg.(GreetingScreen); ok {
		a.router.ApplyU(func(slot *poly.Atomic[any]) {
			*slot = greeting
		})
		a.router.SetTarget(router.SlotU)
		return a, nil
	}

	a.router.SetTarget(router.SlotT)

	var cmd tea.Cmd
	a.router, cmd = a.router.Route(msg)
	if a.router.GetSlotAsT() == nil {
		return a, tea.Quit
	}

	return a, cmd
}

func (a app) View() string {
	return a.router.Render()
}

func (a app) SetContext(ctx any) poly.Atomic[any] {
	return a // No context needed for the app
}

func main() {
	// Set up a standard logger for tracing
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	root := NewApp()

	// Create the host and start the Bubble Tea program
	host := poly.NewHost(
		"Polymer Integration Example",
		root,
	)

	p := tea.NewProgram(host)
	if _, err := p.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
