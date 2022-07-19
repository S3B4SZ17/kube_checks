package tabcluster

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type ControllerTabErrors struct {
	// Tab Errors
   GridTabErrors *ui.Grid
   Bc *widgets.BarChart

   // Tabpane
   Tabpane *widgets.TabPane
   Env *widgets.Paragraph

   // Grid dimension
   TermWidth int
   TermHeight int
}


func NewControllerTabErrors(env *widgets.Paragraph, tabpane *widgets.TabPane) *ControllerTabErrors {
	ctl := &ControllerTabErrors{
		GridTabErrors: ui.NewGrid(),
		Bc: widgets.NewBarChart(),
		Env: env,

		Tabpane: tabpane,
		
	}

	ctl.initUI()

	return ctl
}

func (p *ControllerTabErrors) initUI() {

	p.Bc.Title = "Bar Chart"
	p.Bc.Data = []float64{3, 2, 5, 3, 9, 5, 3, 2, 5, 8, 3, 2, 4, 5, 3, 2, 5, 7, 5, 3, 2, 6, 7, 4, 6, 3, 6, 7, 8, 3, 6, 4, 5, 3, 2, 4, 6, 4, 8, 5, 9, 4, 3, 6, 5, 3, 6}
	p.Bc.SetRect(5, 5, 40, 15)
	p.Bc.Labels = []string{"S0", "S1", "S2", "S3", "S4", "S5"}

	p.TermHeight, p.TermWidth = ui.TerminalDimensions()
	p.GridTabErrors.SetRect(0, 0, p.TermWidth, p.TermHeight)
	p.GridTabErrors.Set(
		ui.NewRow(2.0/2,
			ui.NewRow(1.0/10,
				ui.NewCol(.75/4, p.Tabpane),
				ui.NewCol(.75/4, p.Env),
			),
			ui.NewRow(9.0/20,
				ui.NewCol(1.0/2, p.Bc),
				ui.NewCol(1.0/2, p.Bc),
			),
			ui.NewRow(9.0/20,
				ui.NewCol(1.0/2, p.Bc),
				ui.NewCol(1.0/2, p.Bc),
			),
			
		),
	)
	p.resize()

}

func (p *ControllerTabErrors) Render() {
   p.Bc.Labels = []string{"Sebas", "S1", "S2", "S3", "S4", "S5"}
   p.resize()
   ui.Render(p.GridTabErrors)
}

func (p *ControllerTabErrors) Resize() {
	p.resize()
	ui.Render(p.GridTabErrors)
}

func (p *ControllerTabErrors) resize() {
	w, h := ui.TerminalDimensions()
	p.GridTabErrors.SetRect(0, 0, w, h)
}
