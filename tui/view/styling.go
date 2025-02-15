package view

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

const (
	listHeight        = 14
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	padding           = 2
	maxWidth          = 80
)

var pad = strings.Repeat(" ", padding)

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

	// progressbar
	// greenProgress = progress.WithGradient("#00770c", "#0aff03")
	greenProgress = progress.WithGradient("#04B575", "#04B575")
	redProgress   = progress.WithGradient("#FF3131", "#FF3131")
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

// https://patorjk.com/software/taag font: Bulbhead

var intToAsciiArt = map[int]string{
	0: `
  ___  
 / _ \ 
( (_) )
 \___/ `,
	1: `
  __ 
 /  )
  )( 
 (__)`,
	2: `
 ___  
(__ \ 
 / _/ 
(____)`,
	3: `
 ___ 
(__ )
 (_ \
(___/`,
	4: `
  __  
 /. | 
(_  _)
  (_) `,
	5: `
 ___ 
| __)
|__ \
(___/`,
	6: `
  _  
 / ) 
/ _ \
\___/`,
	7: `
 ___ 
(__ )
 / / 
(_/  `,
	8: `
 ___ 
( _ )
/ _ \
\___/`,
	9: `
 ___ 
/ _ \
\_  /
 (_/ `,
}

var colon = `
  
()
  
()`

func GetAsciiTime(minutes int, seconds int) string {
	timeStr := fmt.Sprintf("%02d:%02d", minutes, seconds)

	var firstArt string
	if timeStr[0] == ':' {
		firstArt = strings.Trim(colon, "\n")
	} else {
		digit := int(timeStr[0] - '0')
		firstArt = strings.Trim(intToAsciiArt[digit], "\n")
	}
	artLines := strings.Split(firstArt, "\n")
	lineCount := len(artLines)

	result := make([]string, lineCount)

	for _, ch := range timeStr {
		var art string
		if ch == ':' {
			art = strings.Trim(colon, "\n")
		} else {
			digit := int(ch - '0')
			art = strings.Trim(intToAsciiArt[digit], "\n")
		}
		lines := strings.Split(art, "\n")
		for i, line := range lines {
			result[i] += line + "  "
		}
	}

	return pad + strings.Join(result, "\n"+pad) + "\n"
}
