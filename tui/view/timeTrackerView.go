package view

import (
	"fmt"
	"strconv"
	"strings"
	"timeTrackingTools/db"
	"timeTrackingTools/models"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TimeTrackerView struct {
	db        *db.TimeLogDb
	tabs      []string
	tabLogs   [][]models.TimeLog
	activeTab int
	table     table.Model
}

func NewTimeTrackerView() *TimeTrackerView {
	var err error
	db, err := db.InitDb()
	if err != nil {
		fmt.Println("Error open db:", err)
	}
	tabs := GetTabNames()

	var tabLogs [][]models.TimeLog

	currentWeek := 7

	for i := 0; i < 3; i++ {
		week := currentWeek - i
		timeLogs, err := db.GetTimeLogsByWeek(week)
		if err != nil {
			fmt.Printf("error fetching for week %d: %v\n", week, err)
			continue
		}
		tabLogs = append(tabLogs, timeLogs)
	}

	v := TimeTrackerView{
		db:        db,
		tabs:      tabs,
		tabLogs:   tabLogs,
		activeTab: 0,
	}

	v.ActualizeTable()

	return &v
}

func (v TimeTrackerView) Render() string {
	return v.GetTabView()
}

func (v *TimeTrackerView) HandleKey(msg interface{}) (View, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "right", "l", "n", "tab":
			v.updateTabNavigation(1)
		case "left", "h", "p", "shift+tab":
			v.updateTabNavigation(-1)
		case "esc":
			return NewStartView(), nil
		}
	}
	return v, nil
}

func (v *TimeTrackerView) updateTabNavigation(direction int) {
	v.activeTab = min(v.activeTab+direction, len(v.tabs)-1)
	v.activeTab = max(v.activeTab, 0)
	v.ActualizeTable()
}

func (v TimeTrackerView) GetTabView() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range v.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(v.tabs)-1, i == v.activeTab
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
		renderedTabs = append(renderedTabs, style.Width(20).Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(v.GetTable()))
	return docStyle.Render(doc.String())
}

func GetTabNames() []string {
	// currentWeek := GetCalendarWeekNow()

	tabs := []string{
		"7",
		"6",
		"5",
	}

	return tabs
}

func (v TimeTrackerView) GetTable() string {
	return baseStyle.Render(v.table.View()) + "\n"
}

func (v *TimeTrackerView) ActualizeTable() {

	rows := []table.Row{}

	logs := v.tabLogs[v.activeTab]

	for _, log := range logs {
		exist := false
		for _, r := range rows {
			if r[0] == strconv.Itoa(log.Timestamp.Day()) && r[1] == strconv.Itoa(int(log.Timestamp.Month())) {
				if log.Kind == 0 {
					r[4] = log.Timestamp.Format("15:04")
				}
				if log.Kind == 1 {
					r[5] = log.Timestamp.Format("15:04")
				}
				exist = true
			}
		}

		if exist {
			continue
		}
		str := table.Row{
			strconv.Itoa(log.Timestamp.Day()),
			strconv.Itoa(int(log.Timestamp.Month())),
			strconv.Itoa(log.Timestamp.Year()),
			strconv.Itoa(log.Week),
			"",
			"",
		}

		if log.Kind == 0 {
			str[4] = log.Timestamp.Format("15:04")
		}
		if log.Kind == 1 {
			str[5] = log.Timestamp.Format("15:04")
		}
		rows = append(rows, str)
	}

	columns := []table.Column{
		{Title: "Day", Width: 6},
		{Title: "Month", Width: 6},
		{Title: "Year", Width: 6},
		{Title: "Week", Width: 6},
		{Title: "In", Width: 6},
		{Title: "Out", Width: 6},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		// table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(tableHeaderBorderStyle).
		BorderForeground(tableHeaderBorderForeground).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(tableSelectedForeground).
		Background(tableSelectedBackground).
		Bold(false)
	t.SetStyles(s)

	v.table = t
}
