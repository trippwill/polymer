package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/gels/menu"
	"github.com/trippwill/polymer/poly"
	"github.com/trippwill/polymer/util"
)

// QuitAtom is a simple Atom that quits the application immediately.
type QuitAtom struct {
	id string
}

func NewQuitAtom() poly.Atomic[any] {
	return QuitAtom{
		id: util.NewUniqeTypeId[QuitAtom](),
	}
}

var _ poly.Atomic[any] = (*QuitAtom)(nil)

func (q QuitAtom) Init() tea.Cmd                                  { return tea.Quit }
func (q QuitAtom) Update(msg tea.Msg) (poly.Atomic[any], tea.Cmd) { return q, nil }
func (q QuitAtom) View() string                                   { return "Goodbye!\n" }
func (q QuitAtom) SetContext(ctx any) poly.Atomic[any] {
	// No context needed for QuitAtom
	return q
}

// NamePromptScreen is an Atom that prompts the user to enter their name.
type NamePromptScreen struct {
	id    string // Unique identifier for the screen
	input string
}

func NewNamePromptScreen() poly.Atomic[any] {
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

func NewGreetingScreen(name string) poly.Atomic[any] {
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
	id       string // Unique identifier for the app
	menu     poly.Atomic[any]
	greeting poly.Atomic[any]
}

func NewApp() poly.Atomic[any] {
	// Create the menu with the name prompt and greeting screens
	menu := menu.NewMenu(
		"Name Wizard",
		menu.AdaptItem(NewNamePromptScreen(), "Name Wizard", "Enter your name"),
		menu.AdaptItem(NewQuitAtom(), "Quit", "Exit Application"),
	)

	return app{menu: menu, id: util.NewUniqeTypeId[app](), greeting: nil}
}

// func (a app) GetCurrent() poly.Identifier {
// 	if a.greeting != nil {
// 		return a.greeting
// 	}
//
// 	return a.menu
// }

func (a app) Update(msg tea.Msg) (poly.Atomic[any], tea.Cmd) {
	if greeting, ok := msg.(GreetingScreen); ok {
		a.greeting = greeting
		return a, nil
	}

	if a.greeting != nil {
		next, cmd := a.greeting.Update(msg)
		a.greeting = next
		return a, cmd
	}

	next, cmd := a.menu.Update(msg)
	if next == nil {
		return a, tea.Sequence(cmd, tea.Quit)
	}

	a.menu = next
	return a, cmd
}

func (a app) View() string {
	if a.greeting != nil {
		return a.greeting.View()
	}

	return a.menu.View()
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
	logger := log.New(f, "polymer", log.Lmsgprefix)

	root := NewApp()

	// Create the host and start the Bubble Tea program
	host := poly.NewHost(
		"Polymer Integration Example",
		root,
	)

	p := tea.NewProgram(host)
	if _, err := p.Run(); err != nil {
		logger.Fatalf("Application failed: %v", err)
	}
}
