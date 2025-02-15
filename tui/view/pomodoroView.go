package view

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type PomodoroView struct {
	progress progress.Model
	stat     stat
	workSec  int
	pauseSec int
	stopped  bool
	round    int
	keymap   keymap
	help     help.Model
}

type stat struct {
	mode          Mode
	total         int
	remaining     int
	progressColor progress.Option
}

type Mode int

const (
	Work Mode = iota
	Pause
)

func (m Mode) GetString() string {
	switch m {
	case Work:
		return "work"
	case Pause:
		return "pause"
	}
	return ""
}

func (v *PomodoroView) SetStat(mode Mode) {
	durations := map[Mode]int{
		Work:  v.workSec,
		Pause: v.pauseSec,
	}

	progressColors := map[Mode]progress.Option{
		Work:  redProgress,
		Pause: greenProgress,
	}

	if mode == Work {
		v.round++
	}

	v.stat = stat{
		mode:          mode,
		total:         durations[mode],
		remaining:     durations[mode],
		progressColor: progressColors[mode],
	}
}

func GetStatTime(totalSec int) string {
	minutes := totalSec / 60
	seconds := totalSec % 60

	return GetAsciiTime(minutes, seconds)
}

func NewPomodoroView() *PomodoroView {
	workMin := 25
	pauseMin := 5

	v := &PomodoroView{
		workSec:  workMin * 60,
		pauseSec: pauseMin * 60,
		progress: progress.New(redProgress),
		stopped:  false,
		round:    0,
		keymap: keymap{
			escape: key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "escape"),
			),
			pauseTimer: key.NewBinding(
				key.WithKeys("p"),
				key.WithHelp("p", "stop/continue timer"),
			),
		},
		help: help.New(),
	}

	v.SetStat(Work)
	return v
}

func (v *PomodoroView) Init() tea.Cmd {
	return tickCmd()
}

func (v PomodoroView) Render() string {
	timeStr := GetStatTime(v.stat.remaining)
	mode := v.stat.mode.GetString()

	return "\n" + pad + "--" + mode + "-- " + fmt.Sprintf("(round %d)", v.round) + "\n" + timeStr + "\n" + pad + v.progress.View() + "\n\n" + pad + v.helpView() + "\n\n"
}

func (v *PomodoroView) HandleKey(msg interface{}) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return NewStartView(), nil
		case "p":
			v.stopped = !v.stopped
			if !v.stopped {
				return v, tickCmd()
			}
			return v, nil
		}
	case tea.WindowSizeMsg:
		v.progress.Width = msg.Width - padding*2 - 4
		if v.progress.Width > maxWidth {
			v.progress.Width = maxWidth
		}
		return v, nil

	case tickMsg:
		if v.stopped {
			return v, nil
		}

		if v.stat.remaining <= 0 || v.progress.Percent() >= 1.0 {
			nextMode := 1 - v.stat.mode
			v.SetStat(nextMode)
			v.progress = progress.New(v.stat.progressColor)
			return v, tickCmd()
		}

		v.stat.remaining--
		increment := float64(1) / float64(v.stat.total)
		cmd := v.progress.IncrPercent(increment)
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

type keymap struct {
	escape     key.Binding
	pauseTimer key.Binding
}

func (v PomodoroView) helpView() string {
	return "\n" + v.help.ShortHelpView([]key.Binding{
		v.keymap.escape,
		v.keymap.pauseTimer,
	})
}
