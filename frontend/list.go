package frontend

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"database/sql"
	"example/speedrun/db"
	"fmt"
	"io"
	"log"
	"strings"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))
var (
	titleStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1).
		MarginBottom(0).
		MarginTop(0).
		Bold(true)
)

const listHeight = 14

type (
	errMsg error
)
type item struct {
	id, Description, Done string
}
type styles struct {
	title        lipgloss.Style
	item         lipgloss.Style
	selectedItem lipgloss.Style
	pagination   lipgloss.Style
	help         lipgloss.Style
	quitText     lipgloss.Style
}

func newStyles(darkBG bool) styles {
	var s styles
	s.title = lipgloss.NewStyle().MarginLeft(2)
	s.item = lipgloss.NewStyle().PaddingLeft(4)
	s.selectedItem = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("140"))
	s.pagination = list.DefaultStyles(darkBG).PaginationStyle.PaddingLeft(4)
	s.help = list.DefaultStyles(darkBG).HelpStyle.PaddingLeft(4).PaddingBottom(1)
	s.quitText = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	return s
}

func (i item) FilterValue() string       { return "" }
func (i item) Id() string                { return i.id }
func (i item) DescriptionMethod() string { return i.Description }
func (i item) DoneMethod() string        { return i.Done }

type itemDelegate struct {
	styles *styles
}

func (m *model) updateStyles(isDark bool) {
	m.styles = newStyles(isDark)
	m.list.Styles.Title = m.styles.title
	m.list.Styles.PaginationStyle = m.styles.pagination
	m.list.Styles.HelpStyle = m.styles.help
	m.list.SetDelegate(itemDelegate{styles: &m.styles})
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := d.styles.item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return d.styles.selectedItem.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	items        []db.Task
	list         list.Model
	dbconnection *sql.DB
	choice       string
	textInput    textinput.Model
	err          error
	styles       styles
	quitting     bool
	showInput    bool
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) View() tea.View {
	var c *tea.Cursor
	if !m.textInput.VirtualCursor() {
		c = m.textInput.Cursor()
		c.Y += lipgloss.Height(m.headerView())
	}

	str := lipgloss.JoinVertical(lipgloss.Top, m.headerView(), m.textInput.View(), m.footerView())
	if m.quitting {
		str += "\n"
	}

	s := titleStyle.Render("TODO MANAGER")
	helpStr := "\ni: insert • enter: toggle • r: clear list • q: quit"
	if m.choice != "" {
		return tea.NewView(m.styles.quitText.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice)))
	}
	if m.showInput {
		pos := lipgloss.JoinVertical(lipgloss.Center, m.textInput.View()+helpStr)
		return tea.NewView(s + "\n" + pos)
	} else {
		pos := lipgloss.JoinVertical(lipgloss.Center, m.list.View()+helpStr)
		return tea.NewView(s + "\n" + pos)
	}
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	m.textInput, cmd = m.textInput.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if !m.showInput {
				m.quitting = true
				return m, tea.Quit
			}
		case "r":
			sel, err := m.list.SelectedItem().(item)
			if !err {
				return m, nil
			}
			if err := db.MarkDone(sel.id, m.dbconnection); err != nil {
				log.Println(err)
				return m, nil
			}
			tasks, errs := db.GetList(m.dbconnection)
			if errs == nil {
				idx := m.list.Index()
				m.list.SetItems(toItems(tasks))
				m.list.Select(idx)
			}
		case "i":
			if !m.showInput {
				m.textInput.Focus()
				m.showInput = true
				err := db.AddTask(m.dbconnection, "ceva")
				if err != nil {
					return m, nil
				}
				tasks, errs := db.GetList(m.dbconnection)
				for _, t := range tasks {
					fmt.Println(t.Description)
				}
				if errs == nil {
					idx := m.list.Index()
					m.list.SetItems(toItems(tasks))
					m.list.Select(idx)
				}
				fmt.Printf("m.textInput: %v\n", m.textInput)
			}
		case "esc":
			if m.showInput {
				m.textInput.SetValue("")
				m.showInput = false
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}
	return m, cmd
}

func toItems(tasks []db.Task) []list.Item {
	items := make([]list.Item, 0, len(tasks))
	for _, t := range tasks {
		items = append(items, item{
			Description: t.Description,
			Done:        fmt.Sprintf("%t", t.Done),
		})
	}
	return items
}

func GenerateTable(dbconnection *sql.DB, tasks []db.Task) model {
	ti := textinput.New()
	ti.Placeholder = "Task name..."
	ti.CharLimit = 156

	var tableRows []list.Item
	tableRows = toItems(tasks)
	const defaultWidth = 20

	l := list.New(tableRows, itemDelegate{}, defaultWidth, listHeight)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	m := model{list: l, dbconnection: dbconnection, textInput: ti}
	m.updateStyles(true) // default to dark styles.
	return m
}

func LoadTasks(dbconnection *sql.DB) ([]db.Task, error) {
	tasks, err := db.GetList(dbconnection)
	if err != nil {
		log.Fatal("Error inserting inorders table", err)
	}
	return tasks, nil
}
func (m model) headerView() string { return "What is on your mind ?\n" }
func (m model) footerView() string { return "\n(esc to quit)" }
