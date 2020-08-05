package fdtd

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type TuiView struct {
	Rows []TaskRowView
	Grid *termui.Grid
	ctl  *FDTDTaskSet
}

type TaskRowView struct {
	Progress *widgets.Gauge
	Status   *widgets.Paragraph
	Row      termui.GridItem
}

func NewTuiView(ctl *FDTDTaskSet) *TuiView {
	var c TuiView
	c.ctl = ctl
	c.Rows = make([]TaskRowView, len(ctl.Tasks))
	for i := range ctl.Tasks {
		row_view := BuildTaskRow()
		c.Rows[i] = row_view
	}
	c.Grid = termui.NewGrid()
	termWidth, termHeight := termui.TerminalDimensions()
	c.Grid.SetRect(0, 0, termWidth, termHeight)

	rows := make([]interface{}, len(ctl.Tasks))
	for i, row := range c.Rows {
		rows[i] = row.Row
	}
	c.Grid.Set(rows...)
	return &c
}

func (c *TuiView) Update() {
	for i, task := range c.ctl.Tasks {
		c.Rows[i].Update(task)
	}
}

func BuildTaskRow() TaskRowView {
	g := widgets.NewGauge()
	g.Percent = 50
	t := widgets.NewParagraph()
	t.Text = "Unknown"
	t.Border = false
	row := termui.NewRow(1.0/2,
		termui.NewCol(1.0/2, g),
		termui.NewCol(1.0/2, t),
	)
	return TaskRowView{
		Status:   t,
		Progress: g,
		Row:      row,
	}
}

func (c TaskRowView) Update(task *FDTDSlurmTask) {
	c.Status.Text = task.Status.Calculate
	c.Progress.Percent = int(task.Statistic.Progress)
}
