package main

import (
	"github.com/charmbracelet/bubbles/table"
	"strconv"
)

func (m model) GetTable() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func (m *model) ActualizeTable() {

	rows := []table.Row{}

	logs := m.tabLogs[m.activeTab]

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

	m.table = t
}
