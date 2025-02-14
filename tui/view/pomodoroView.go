package view

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type PomodoroView struct {
	Ticks    int
	Frames   int
	progress progress.Model
	Loaded   bool
}

func NewPomodoroView() *PomodoroView {
	return &PomodoroView{
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

func (v PomodoroView) Render() string {
	pad := strings.Repeat(" ", padding)
	return "\n" + pad + v.progress.View() + "\n\n"
}

func (v *PomodoroView) HandleKey(msg interface{}) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if v.progress.Percent() == 1.0 {
			return v, tea.Quit
		}
		cmd := v.progress.IncrPercent(0.05)
		return v, tea.Batch(tickCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := v.progress.Update(msg)
		v.progress = progressModel.(progress.Model)
		return v, cmd

	}
	return v, nil
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
