package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
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

var db *TimeLogDb

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
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

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

		if m.choice == "TimeTracker" {
			switch keypress := msg.String(); keypress {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "right", "l", "n", "tab":
				m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
				m.ActualizeTable()
				return m, nil
			case "left", "h", "p", "shift+tab":
				m.activeTab = max(m.activeTab-1, 0)
				m.ActualizeTable()
				return m, nil
			}
			return m, nil
		}

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
		msg = m.GetTabs()
	case "Pomodoro":
		pad := strings.Repeat(" ", padding)
		msg = "\n" +
			pad + m.progress.View() + "\n\n"
	default:
	}

	return msg
}

func (m model) GetTabs() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		// renderedTabs = append(renderedTabs, style.Render(t))
		renderedTabs = append(renderedTabs, style.Width(15).Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.GetTable()))
	return docStyle.Render(doc.String())
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) GetTable() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func (m *model) ActualizeTable() {

	rows := []table.Row{}

	for _, timeLog := range m.tabLogs[m.activeTab] {
		str := table.Row{
			strconv.Itoa(timeLog.day),
			strconv.Itoa(timeLog.month),
			strconv.Itoa(timeLog.year),
			strconv.Itoa(timeLog.week),
		}
		rows = append(rows, str)
	}

	columns := []table.Column{
		{Title: "Day", Width: 6},
		{Title: "Month", Width: 6},
		{Title: "Year", Width: 6},
		{Title: "Week", Width: 6},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		// table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m.table = t
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

	currentWeek := GetCalendarWeek()

	tabs := []string{
		fmt.Sprintf("this week (%d)", currentWeek),
		fmt.Sprintf("%d", currentWeek-1),
		fmt.Sprintf("%d", currentWeek-2),
	}

	var err error
	db, err = InitDb()
	if err != nil {
		fmt.Println("Error open db:", err)
		os.Exit(1)
	}

	all, _ := db.GetAllTimeLogs()
	if all == nil {
		db.InsertDummyData()
	}

	var tabLogs [][]TimeLog

	for i := 0; i < 3; i++ {
		week := currentWeek - i
		timeLogs, err := db.GetTimeLogsByWeek(week)
		if err != nil {
			fmt.Printf("error fetching for week %d: %v\n", week, err)
			continue
		}
		tabLogs = append(tabLogs, timeLogs)
	}

	m := model{list: l, progress: progress.New(progress.WithDefaultGradient()), Tabs: tabs, tabLogs: tabLogs}

	m.ActualizeTable()
	//db.Insert(*NewTimeLogNow(0))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
