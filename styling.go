package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
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
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()

	// table
	tableHeaderBorderStyle      = lipgloss.NormalBorder()
	tableHeaderBorderForeground = lipgloss.Color("240")
	tableSelectedForeground     = lipgloss.Color("229")
	tableSelectedBackground     = lipgloss.Color("57")
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
