package cli

const jsonStyle = "json"
const tableStyle = "table"

type outputStylesContainer struct {
	Table string
	JSON  string
}

var OutputStyles outputStylesContainer = outputStylesContainer{
	Table: tableStyle,
	JSON:  jsonStyle,
}
