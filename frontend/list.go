package frontend

import (
	"charm.land/bubbles/v2/list"
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
		MarginBottom(1).
		Bold(true)
)

const listHeight = 14

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
	s.selectedItem = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
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
	styles       styles
	quitting     bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) View() tea.View {
	s := titleStyle.Render("TODO MANAGER") + "\n"
	helpStr := "\ni: insert • enter: toggle • r: clear list • q: quit"
	if m.choice != "" {
		return tea.NewView(m.styles.quitText.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice)))
	}
	if m.quitting {
		return tea.NewView(m.styles.quitText.Render("Not hungry? That’s cool."))
	}
	return tea.NewView(s + "\n" + m.list.View() + helpStr + "\n")
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
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
			err := db.AddTask(m.dbconnection, "ceva")
			if err != nil {
				return m, nil
			}
			tasks, errs := db.GetList(m.dbconnection)
			if errs == nil {
				idx := m.list.Index()
				m.list.SetItems(toItems(tasks))
				m.list.Select(idx)
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
			id:          fmt.Sprintf("%d", t.Id),
			Description: t.Description,
			Done:        fmt.Sprintf("%t", t.Done),
		})
	}
	return items
}

func GenerateTable(dbconnection *sql.DB, tasks []db.Task) model {

	var tableRows []list.Item
	tableRows = toItems(tasks)
	const defaultWidth = 20

	l := list.New(tableRows, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What do you want for dinner?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	m := model{list: l, dbconnection: dbconnection}
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
