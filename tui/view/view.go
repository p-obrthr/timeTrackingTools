package view

import (
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

type View interface {
	Render() string
	HandleKey(msg interface{}) (View, tea.Cmd)
}
