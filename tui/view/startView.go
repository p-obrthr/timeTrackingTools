package view

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type StartView struct {
	list   list.Model
	choice string
}

func NewStartView() *StartView {
	list := *NewList()
	if list.Items() == nil {
		fmt.Println("err empty list")
	}
	return &StartView{
		list:   list,
		choice: "TimeTracker",
	}
}

func (v StartView) Render() string {
	return "\n" + v.list.View()
}

func (v StartView) HandleKey(msg interface{}) (View, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "enter":
			i, ok := v.list.SelectedItem().(item)
			if ok {
				v.choice = string(i)
			}
			if v.choice == "TimeTracker" {
				return NewTimeTrackerView(), nil
			} else {
				return NewPomodoroView(), NewPomodoroView().Init()
			}
		case "esc":
			return NewStartView(), nil
		}
	}
	var cmd tea.Cmd
	v.list, cmd = v.list.Update(msg)
	return v, cmd
}

type item string

type itemDelegate struct{}

func (i item) FilterValue() string {
	return ""
}

func (d itemDelegate) Height() int {
	return 1
}

func (d itemDelegate) Spacing() int {
	return 0
}

func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

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

func NewList() *list.Model {
	items := []list.Item{
		item("TimeTracker"),
		item("Pomodoro"),
	}

	const defaultWidth = 20

	l := list.New(
		items,
		itemDelegate{},
		defaultWidth,
		listHeight,
	)
	l.Title = "Choose tool."
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle

	return &l
}
