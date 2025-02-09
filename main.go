package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

var db *TimeLogDb

type tickMsg time.Time

type model struct {
	list      list.Model
	choice    string
	quit      bool
	Ticks     int
	Frames    int
	progress  progress.Model
	Loaded    bool
	Tabs      []string
	tabLogs   [][]TimeLog
	activeTab int
	table     table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}
		cmd := m.progress.IncrPercent(0.05)
		return m, tea.Batch(tickCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		if m.choice == "TimeTracker" {
			return m, tea.Quit
		}
		m.quit = true
		return m, tea.Quit

	case "enter":
		i, ok := m.list.SelectedItem().(item)
		if ok {
			m.choice = string(i)
		}
		if m.choice == "Pomodoro" {
			return m, tickCmd()
		}

	case "right", "l", "n", "tab":
		if m.choice == "TimeTracker" {
			m.updateTabNavigation(1)
			return m, nil
		}

	case "left", "h", "p", "shift+tab":
		if m.choice == "TimeTracker" {
			m.updateTabNavigation(-1)
			return m, nil
		}
	default:
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *model) updateTabNavigation(direction int) {
	m.activeTab = min(m.activeTab+direction, len(m.Tabs)-1)
	m.activeTab = max(m.activeTab, 0)
	m.ActualizeTable()
}

func (m model) View() string {
	if m.choice != "" {
		return m.chosenView()
	}
	if m.quit {
		return quitTextStyle.Render("Goodbye.")
	}
	return "\n" + m.list.View()
}

func (m *model) chosenView() string {
	var msg string

	switch m.choice {
	case "TimeTracker":
		db.GetAllTimeLogs()
		// fmt.Sprint(NewTimeLogNow(0))
		msg = m.GetTabView()
	case "Pomodoro":
		pad := strings.Repeat(" ", padding)
		msg = "\n" +
			pad + m.progress.View() + "\n\n"
	default:
	}

	return msg
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func main() {

	var err error
	db, err = InitDb()
	if err != nil {
		fmt.Println("Error open db:", err)
	}

	var tabLogs [][]TimeLog

	currentWeek := GetCalendarWeek()

	for i := 0; i < 3; i++ {
		week := currentWeek - i
		timeLogs, err := db.GetTimeLogsByWeek(week)
		if err != nil {
			fmt.Printf("error fetching for week %d: %v\n", week, err)
			continue
		}
		tabLogs = append(tabLogs, timeLogs)
	}

	list := NewList()
	tabs := GetTabNames()

	m := model{
		list:     *list,
		progress: progress.New(progress.WithDefaultGradient()),
		Tabs:     tabs,
		tabLogs:  tabLogs,
	}

	m.ActualizeTable()
	//db.Insert(*NewTimeLogNow(0))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
