package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/mattn/go-sqlite3"
)

const (
	listHeight        = 14
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	padding           = 2
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#04B575"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	keywordStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	checkboxStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	progressEmpty     = subtleStyle.Render(progressEmptyChar)
	helpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
)

type item string

type tickMsg time.Time

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	choice   string
	quit     bool
	Ticks    int
	Frames   int
	progress progress.Model
	Loaded   bool
}

type TimeLog struct {
	id    int
	day   int
	month int
	year  int
	kind  int
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
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
		}

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
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return chosenView(m)
	}
	if m.quit {
		return quitTextStyle.Render("Goodbye.")
	}
	return "\n" + m.list.View()
}

func chosenView(m model) string {
	var msg string

	switch m.choice {
	case "TimeTracker":
		msg = "timetracker time"
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
	items := []list.Item{
		item("TimeTracker"),
		item("Pomodoro"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Choose tool."
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle

	m := model{list: l, progress: progress.New(progress.WithDefaultGradient())}

	c, err := InitDb()
	if err != nil {
		fmt.Println("Error open db:", err)
		os.Exit(1)
	}
	timeLog := TimeLog{
		day:   1,
		month: 1,
		year:  1,
		kind:  1,
	}
	fmt.Println(c.Insert(timeLog))
	fmt.Println(c.GetTimeLog(1))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
