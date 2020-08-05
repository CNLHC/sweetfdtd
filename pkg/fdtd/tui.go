package fdtd

import (
	"fmt"
	"strconv"

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
	Slurm    *widgets.Paragraph
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
	c.Grid.PaddingLeft = int(termWidth / 5)
	c.Grid.PaddingRight = int(termWidth / 5)
	c.Grid.PaddingTop = int(termHeight / 10)

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
	t1 := widgets.NewParagraph()
	t2 := widgets.NewParagraph()
	t1.Border = false
	t2.Border = false
	row := termui.NewRow(1.0/10,
		termui.NewCol(1.0/2, g),
		termui.NewCol(1.0/4, t1),
		termui.NewCol(1.0/4, t2),
	)
	return TaskRowView{
		Status:   t1,
		Slurm:    t2,
		Progress: g,
		Row:      row,
	}
}

func (c TaskRowView) Update(task *FDTDSlurmTask) {
	c.Progress.Title = task.FSPFile
	c.Slurm.Text = strconv.Itoa(int(task.JOBID))
	c.Progress.Percent = int(task.Statistic.Progress)

	switch task.Status.FDTD {
	case FDTDTimeout:
		c.Status.Text = "Timeout"
	case FDTDWaiting:
		c.Status.Text = "Waiting"
	case FDTDRunning:
		c.Status.Text = task.Status.Calculate
	case FDTDFailed:
		c.Status.Text = fmt.Sprintf("Failed(%d)", task.RetryInSec)
	}

}
