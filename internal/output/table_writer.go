package output

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

var ErrEmptyHeaders error = errors.New("headers must supplied to build a table")
var ErrEmptyRows error = errors.New("at least one must supplied to build a table")

type Opt func(*internalTable)

func WithIndexColumn() Opt {
	return func(t *internalTable) {
		t.addIndexColumn = true
	}
}

func WithAutoMergeByIndexes(indexes []int) Opt {
	return func(t *internalTable) {
		t.autoMergeIndexes = indexes
	}
}

type internalTable struct {
	headers          []string
	rows             [][]string
	addIndexColumn   bool
	autoMergeIndexes []int
}

func (t *internalTable) Write(w io.Writer) {
	tbl := tablewriter.NewWriter(w)

	headers := t.headers
	if t.addIndexColumn {
		headers = append([]string{"#"}, headers...)
	}

	for i, r := range t.rows {
		if t.addIndexColumn {
			r = append([]string{strconv.Itoa(i + 1)}, r...)
		}
		tbl.Append(r)
	}

	tbl.SetHeader(headers)
	tbl.SetRowLine(true)
	tbl.SetRowSeparator("-")

	if len(t.autoMergeIndexes) > 0 {
		tbl.SetAutoMergeCellsByColumnIndex(t.autoMergeIndexes)
	}

	tbl.Render()
}

func (t *internalTable) Validate() error {
	if len(t.headers) == 0 {
		return ErrEmptyHeaders
	}

	if len(t.rows) == 0 {
		return ErrEmptyRows
	}

	headerCount := len(t.headers)

	for i, r := range t.rows {
		if columnCount := len(r); columnCount != headerCount {
			return fmt.Errorf("row[%d] has %d columns, %d expected", i+1, columnCount, headerCount)
		}
	}

	return nil
}

type TableFactory struct {
	writer io.Writer
}

func (tf *TableFactory) CreateAndPrint(headers []string, rows [][]string, opts ...Opt) error {
	t := &internalTable{
		headers: headers,
		rows:    rows,
	}

	for _, opt := range opts {
		opt(t)
	}

	if err := t.Validate(); err != nil {
		return err
	}

	t.Write(tf.writer)

	return nil
}

func NewTableFactory(w io.Writer) *TableFactory {
	return &TableFactory{
		writer: w,
	}
}
