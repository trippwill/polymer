package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/atoms"
	"github.com/trippwill/polymer/gels/menu"
	"github.com/trippwill/polymer/poly"
	"github.com/trippwill/polymer/router/auto"
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

var _ poly.Atomic[string] = NamePromptScreen{}

func (n NamePromptScreen) Update(msg tea.Msg) (poly.Atomic[string], tea.Cmd) {
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
			return nil, util.Broadcast(poly.ContextMsg[string]{Context: n.input})
		}
	}
	return n, nil
}

func (n NamePromptScreen) View() string {
	return "Enter your name: " + n.input + "\n"
}

// GreetingScreen is an Model that greets the user by name.
type GreetingScreen struct {
	id    string // Unique identifier for the screen
	value string
}

func NewGreetingScreen() GreetingScreen {
	return GreetingScreen{
		id: util.NewUniqeTypeId[GreetingScreen](),
	}
}

var _ poly.Atomic[string] = GreetingScreen{}

func (g GreetingScreen) Update(msg tea.Msg) (poly.Atomic[string], tea.Cmd) {
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
	router auto.Auto[poly.Atomic[string]] // Router to manage screens
	id     string                         // Unique identifier for the app
}

func NewApp() poly.Atomic[string] {
	// Create the menu with the name prompt and greeting screens
	menu := menu.NewMenu(
		"Name Wizard",
		menu.AdaptItem(NewNamePromptScreen(), "Name Wizard", "Enter your name"),
		menu.AdaptItem(atoms.NewQuitAtom[string](), "Quit", "Exit Application"),
	)

	return app{
		router: auto.NewAuto(menu, nil),
		id:     util.NewUniqeTypeId[app](),
	}
}

func (a app) Update(msg tea.Msg) (poly.Atomic[string], tea.Cmd) {
	if name, ok := msg.(poly.ContextMsg[string]); ok {
		greetingScreen := NewGreetingScreen()
		greetingScreen.SetContext(name.Context)
		a.router = a.router.Set(auto.SlotOverride, greetingScreen)
		return a, nil
	}

	a.router = a.router.Set(auto.SlotOverride, nil)

	var cmd tea.Cmd
	a.router, cmd = a.router.Route(msg)

	// The menu has quit
	if !a.router.IsSet(auto.SlotPrimary) {
		return a, tea.Quit
	}

	return a, cmd
}

func (a app) View() string {
	return a.router.Render()
}

func (a *app) SetContext(ctx string) {
	auto.SetContext(&a.router, ctx)
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
