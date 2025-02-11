package main

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func (m model) GetTabView() string {
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
		renderedTabs = append(renderedTabs, style.Width(20).Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.GetTable()))
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
