package frontend

import (
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"database/sql"
	"example/speedrun/db"
	"fmt"
	"log"
	"os"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table        table.Model
	dbconnection *sql.DB
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			err := db.MarkDone(m.table.SelectedRow()[0], m.dbconnection)
			if err != nil {
				fmt.Println(err)
			}
			m.updateTask()

		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) updateTask() {
	highLightedRow := m.table.SelectedRow()
	highLightedRow[2] = "false"
	fmt.Println(highLightedRow)
}

func (m model) View() tea.View {
	return tea.NewView(baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n")
}

func GenerateTable(dbconnection *sql.DB) {
	columns := []table.Column{
		{Title: "ID", Width: 7},
		{Title: "Description", Width: 15},
		{Title: "Done", Width: 10},
	}

	tasks, err := db.GetList(dbconnection)
	if err != nil {
		log.Fatal("Error inserting inorders table", err)
	}
	var tableRows []table.Row
	for _, task := range tasks {
		row := table.Row{
			fmt.Sprintf("%d", task.Id),
			task.Description,
			fmt.Sprintf("%t", task.Done),
		}
		tableRows = append(tableRows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(tableRows),
		table.WithFocused(true),
		table.WithHeight(20),
		table.WithWidth(42),
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

	m := model{t, dbconnection}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
