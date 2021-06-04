package output_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/svartlfheim/ymir/internal/output"
)

func TestTableWriter_CreateAndPrint(t *testing.T) {
	tests := []struct {
		name        string
		h           []string
		r           [][]string
		w           *bytes.Buffer
		opts        []output.Opt
		expected    string
		expectedErr error
	}{
		{
			name: "Table with 1 row and 1 column",
			h: []string{
				"HEAD1",
			},
			r: [][]string{
				{
					"COL1",
				},
			},
			w: bytes.NewBuffer([]byte{}),
			expected: `+-------+
| HEAD1 |
+-------+
| COL1  |
+-------+
`,
		},
		{
			name: "Table with 1 rows and multiple columns",
			h: []string{
				"HEAD1",
				"HEAD2",
				"HEAD3",
				"HEAD4",
			},
			r: [][]string{
				{
					"COL1",
					"SOME VAL",
					"SOME LONGER VAL",
					"FINALLY A LONG SENTENCE CAUSE WHY NOT",
				},
			},
			w: bytes.NewBuffer([]byte{}),
			expected: `+-------+----------+-----------------+--------------------------------+
| HEAD1 |  HEAD2   |      HEAD3      |             HEAD4              |
+-------+----------+-----------------+--------------------------------+
| COL1  | SOME VAL | SOME LONGER VAL | FINALLY A LONG SENTENCE CAUSE  |
|       |          |                 | WHY NOT                        |
+-------+----------+-----------------+--------------------------------+
`,
		},
		{
			name: "Table with multiple rows and multiple columns",
			h: []string{
				"HEAD1",
				"HEAD2",
				"HEAD3",
				"HEAD4",
			},
			r: [][]string{
				{
					"ROW 1",
					"Lorem ipsum dolor sit amet, consectetur adipiscing elit",
					"Sed blandit quam eu mauris laoreet",
					"Ut at arcu vulputate.",
				},
				{
					"ROW 2",
					"Cras facilisis.",
					"Etiam aliquam mauris in feugiat venenatis.",
					"Donec aliquet neque accumsan nisi vulputate, nec porttitor justo auctor.",
				},
				{
					"ROW 3",
					"Etiam vel lacus placerat, scelerisque velit id, egestas nunc.",
					"Sed sed tellus in neque",
					"In id arcu",
				},
				{
					"ROW 4",
					"Fusce vulputate libero eget consectetur",
					"Donec lacinia ligula a mi porta elementum",
					"Suspendisse vulputate metus id elementum vulputate",
				},
				{
					"ROW 5",
					"In ac justo sit amet lectus consequat cursus eget in tortor.",
					"Quisque molestie odio et pulvinar interdum.",
					"Nulla",
				},
			},
			w: bytes.NewBuffer([]byte{}),
			expected: `+-------+--------------------------------+--------------------------------+--------------------------------+
| HEAD1 |             HEAD2              |             HEAD3              |             HEAD4              |
+-------+--------------------------------+--------------------------------+--------------------------------+
| ROW 1 | Lorem ipsum dolor sit amet,    | Sed blandit quam eu mauris     | Ut at arcu vulputate.          |
|       | consectetur adipiscing elit    | laoreet                        |                                |
+-------+--------------------------------+--------------------------------+--------------------------------+
| ROW 2 | Cras facilisis.                | Etiam aliquam mauris in        | Donec aliquet neque accumsan   |
|       |                                | feugiat venenatis.             | nisi vulputate, nec porttitor  |
|       |                                |                                | justo auctor.                  |
+-------+--------------------------------+--------------------------------+--------------------------------+
| ROW 3 | Etiam vel lacus placerat,      | Sed sed tellus in neque        | In id arcu                     |
|       | scelerisque velit id, egestas  |                                |                                |
|       | nunc.                          |                                |                                |
+-------+--------------------------------+--------------------------------+--------------------------------+
| ROW 4 | Fusce vulputate libero eget    | Donec lacinia ligula a mi      | Suspendisse vulputate metus id |
|       | consectetur                    | porta elementum                | elementum vulputate            |
+-------+--------------------------------+--------------------------------+--------------------------------+
| ROW 5 | In ac justo sit amet lectus    | Quisque molestie odio et       | Nulla                          |
|       | consequat cursus eget in       | pulvinar interdum.             |                                |
|       | tortor.                        |                                |                                |
+-------+--------------------------------+--------------------------------+--------------------------------+
`,
		},
		{
			name: "Row Indexes are added",
			h: []string{
				"HEAD1",
			},
			r: [][]string{
				{
					"ROW1",
				},
				{
					"ROW2",
				},
				{
					"ROW3",
				},
			},
			opts: []output.Opt{
				output.WithIndexColumn(),
			},
			w: bytes.NewBuffer([]byte{}),
			expected: `+---+-------+
| # | HEAD1 |
+---+-------+
| 1 | ROW1  |
+---+-------+
| 2 | ROW2  |
+---+-------+
| 3 | ROW3  |
+---+-------+
`,
		},
		{
			name: "Error is reported if row has too few columns",
			h: []string{
				"HEAD1",
				"HEAD2",
				"HEAD3",
			},
			r: [][]string{
				{
					"ROW1",
					"ROW1",
					"ROW1",
				},
				{
					// This requires 3 columns
					"ROW2",
					"ROW2",
				},
				{
					"ROW3",
					"ROW3",
					"ROW3",
				},
			},
			w:           bytes.NewBuffer([]byte{}),
			expectedErr: errors.New("row[2] has 2 columns, 3 expected"),
		},
		{
			name: "Error is reported if no rows are supplied",
			h: []string{
				"HEAD1",
				"HEAD2",
				"HEAD3",
			},
			r:           [][]string{},
			w:           bytes.NewBuffer([]byte{}),
			expectedErr: output.ErrEmptyRows,
		},
		{
			name: "Error is reported if no headers are supplied",
			h:    []string{},
			r: [][]string{
				{
					"ROW1",
					"ROW1",
					"ROW1",
				},
				{
					"ROW2",
					"ROW2",
					"ROW3",
				},
				{
					"ROW3",
					"ROW3",
					"ROW3",
				},
			},
			w:           bytes.NewBuffer([]byte{}),
			expectedErr: output.ErrEmptyHeaders,
		},
		{
			name: "Cells are merged by common values",
			h: []string{
				"HEAD1",
				"HEAD2",
				"HEAD3",
			},
			r: [][]string{
				{
					"ROW1COL1",
					"ROW1COL2",
					"ROW1COL3A",
				},
				{
					"ROW1COL1",
					"ROW1COL2",
					"ROW1COL3B",
				},
				{
					"ROW1COL1",
					"ROW1COL2",
					"ROW1COL3C",
				},
				{
					"ROW2COL1",
					"ROW2COL2",
					"ROW2COL3A",
				},
				{
					"ROW2COL1",
					"ROW2COL2",
					"ROW2COL3B",
				},
				{
					"ROW3COL1",
					"ROW3COL2",
					"ROW3COL3",
				},
			},
			opts: []output.Opt{
				output.WithAutoMergeByIndexes([]int{0, 1}),
			},
			w: bytes.NewBuffer([]byte{}),
			expected: `+----------+----------+-----------+
|  HEAD1   |  HEAD2   |   HEAD3   |
+----------+----------+-----------+
| ROW1COL1 | ROW1COL2 | ROW1COL3A |
+          +          +-----------+
|          |          | ROW1COL3B |
+          +          +-----------+
|          |          | ROW1COL3C |
+----------+----------+-----------+
| ROW2COL1 | ROW2COL2 | ROW2COL3A |
+          +          +-----------+
|          |          | ROW2COL3B |
+----------+----------+-----------+
| ROW3COL1 | ROW3COL2 | ROW3COL3  |
+----------+----------+-----------+
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tf := output.NewTableFactory(test.w)

			err := tf.CreateAndPrint(test.h, test.r, test.opts...)

			out := test.w.String()

			if test.expectedErr != nil {
				assert.Equal(tt, test.expectedErr, err)
			} else {
				assert.Equal(tt, test.expected, out)
			}

		})
	}
}
