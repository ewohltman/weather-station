package format

import (
	"bytes"
	"fmt"
	"html/template"
)

const (
	tableFmt = `
<table class="table">
  <tbody class="align-middle">{{range .Rows}}
    <tr>{{range .Columns}}
      <td id="{{.ID}}"></td>{{end}}
    </tr>{{end}}
  </tbody>
</table>
`
)

// Table is used in an HTML template.
type Table struct {
	Rows []Row
}

// Row is a dimension of a Table.
type Row struct {
	Columns []Column
}

// Column is a dimension of a Table.
type Column struct {
	ID string
}

func ExecuteTemplate(tablePrefix string, rows, columns int) (*bytes.Buffer, error) {
	table := Table{Rows: make([]Row, rows)}

	for i := 0; i < rows; i++ {
		table.Rows[i].Columns = make([]Column, columns)

		for j := 0; j < columns; j++ {
			table.Rows[i].Columns[j].ID = CellID(tablePrefix, i, j)
		}
	}

	tmpl, err := template.New("").Parse(tableFmt)
	if err != nil {
		return nil, fmt.Errorf("error creating new template: %w", err)
	}

	index := &bytes.Buffer{}

	err = tmpl.Execute(index, table)
	if err != nil {
		return nil, fmt.Errorf("error executing template: %w", err)
	}

	return index, nil
}

func CellID(tablePrefix string, row, column int) string {
	return fmt.Sprintf("%s_%d_%d", tablePrefix, row, column)
}
