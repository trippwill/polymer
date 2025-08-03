package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	poly "github.com/trippwill/polymer"
	"github.com/trippwill/polymer/gels/menu"
	"github.com/trippwill/polymer/trace"
)

// QuitAtom is a simple Atom that quits the application immediately.
type QuitAtom struct {
	poly.Atom
}

func (q QuitAtom) Init() tea.Cmd                           { return tea.Quit }
func (q QuitAtom) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return q, nil }
func (q QuitAtom) View() string                            { return "Goodbye!\n" }

// NamePromptScreen is an Atom that prompts the user to enter their name.
type NamePromptScreen struct {
	poly.Atom
	input string
}

func (n NamePromptScreen) Init() tea.Cmd { return nil }

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
			return GreetingScreen{Value: n.input}, nil
		}
	}
	return n, nil
}

func (n NamePromptScreen) View() string {
	return "Enter your name: " + n.input + "\n"
}

// GreetingScreen is an Atom that greets the user by name.
type GreetingScreen struct {
	Value string
}

func (g GreetingScreen) Init() tea.Cmd {
	return nil
}

func (g GreetingScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// On any key, reset chain to initial state: exit greeting and name prompt back to main menu.
	if _, ok := msg.(tea.KeyMsg); ok {
		return nil, nil
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
	root := menu.NewMenu(
		"Main Menu",
		menu.NewItem(NamePromptScreen{Atom: poly.NewAtom("Enter Name")}, "Run the Name Wizard"),
		menu.NewItem(QuitAtom{Atom: poly.NewAtom("Quit")}, "Exit Application"),
	)

	// Create the host and start the Bubble Tea program
	host := poly.NewHost(
		"Polymer Integration Example",
		root,
		poly.WithLifecycleLogging(logger, trace.LevelTrace)...,
	)

	p := tea.NewProgram(host)
	if _, err := p.Run(); err != nil {
		logger.Fatalf("Application failed: %v", err)
	}
}
