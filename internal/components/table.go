package components

type Table interface {
	Columns() []string
	Rows() []string
	OnTableResize(int)
}
