package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	poly "github.com/trippwill/polymer"
	"github.com/trippwill/polymer/atom"
	"github.com/trippwill/polymer/gels/menu"
	"github.com/trippwill/polymer/trace"
)

// QuitAtom is a simple Atom that quits the application immediately.
type QuitAtom struct{}

func (q QuitAtom) Init() tea.Cmd                            { return tea.Quit }
func (q QuitAtom) Update(msg tea.Msg) (atom.Model, tea.Cmd) { return q, nil }
func (q QuitAtom) View() string                             { return "Goodbye!\n" }
func (q QuitAtom) Name() string                             { return "Quit" }

// NamePromptScreen is an Atom that prompts the user to enter their name.
type NamePromptScreen struct {
	input string
}

func (n NamePromptScreen) Init() tea.Cmd {
	return nil
}

func (n NamePromptScreen) Update(msg tea.Msg) (atom.Model, tea.Cmd) {
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
			// Push the GreetingScreen, passing along the entered name.
			return n, atom.Push(GreetingScreen{Value: n.input})
		}
	}
	return n, nil
}

func (n NamePromptScreen) View() string {
	return "Enter your name: " + n.input + "\n"
}

func (n NamePromptScreen) Name() string { return "Enter Name" }

// GreetingScreen is an Atom that greets the user by name.
type GreetingScreen struct {
	Value string
}

func (g GreetingScreen) Init() tea.Cmd {
	return nil
}

func (g GreetingScreen) Update(msg tea.Msg) (atom.Model, tea.Cmd) {
	// On any key, reset chain to initial state: exit greeting and name prompt back to main menu.
	if _, ok := msg.(tea.KeyMsg); ok {
		return g, atom.Reset(nil)
	}
	return g, nil
}

func (g GreetingScreen) View() string {
	return "\nHello, " + g.Value + "! (press any key to return)\n"
}

func (g GreetingScreen) Name() string { return "Greeting" }

func main() {
	// Set up a standard logger for tracing
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	logger := log.New(f, "polymer: ", log.LstdFlags)

	// Create the main menu using the menu gel
	// Wrap the menu in its own navigation chain
	root := atom.NewStack(
		menu.NewMenu(
			"Main Menu",
			menu.NewMenuItem(NamePromptScreen{}, "Run the Name Wizard"),
			menu.NewMenuItem(QuitAtom{}, "Exit Application"),
		),
	)

	// Create the host and start the Bubble Tea program
	host := poly.NewHost(
		"Polymer Integration Example",
		root,
		trace.WithLifecycleLogging(logger, trace.LevelTrace)...,
	)

	p := tea.NewProgram(host)
	if _, err := p.Run(); err != nil {
		logger.Fatalf("Application failed: %v", err)
	}
}
