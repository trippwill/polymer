package router

import (
	tea "github.com/charmbracelet/bubbletea"
)

type mockModel struct {
	val string
}

func (m mockModel) Update(msg tea.Msg) (mockModel, tea.Cmd) {
	return mockModel{val: msg.(string)}, nil
}

func (m mockModel) View() string {
	return m.val
}
