package tabcluster

import (
	"github.com/S3B4SZ17/kube_checks/healthcheck"

	"github.com/gizak/termui/v3"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type ControllerTabPods struct {
	//Tab pods
   GridTabPods *ui.Grid
   App1 *widgets.Paragraph
   Env *widgets.Paragraph
   K8s *healthcheck.K8s

   // Tabpane
   Tabpane *widgets.TabPane

   // Grid dimension
   TermWidth int
   TermHeight int
}

func NewControllerTabPods(env *widgets.Paragraph, tabpane *widgets.TabPane, k8s *healthcheck.K8s) *ControllerTabPods {
	ctl := &ControllerTabPods{
		GridTabPods: ui.NewGrid(),
		App1: widgets.NewParagraph(),
		Env: env,
		Tabpane: tabpane,
		K8s: k8s,
	}

	ctl.initUI()

	return ctl
}

func (p *ControllerTabPods) initUI() {

	p.App1.Text = "App1 healthcheck" //healthcheck.function()
	p.App1.Title = "App1 healthcheck"
	p.App1.SetRect(5, 5, 40, 15)
	p.App1.BorderStyle.Fg = ui.ColorYellow

	
	p.TermHeight, p.TermWidth = ui.TerminalDimensions()
	p.GridTabPods.SetRect(0, 0, p.TermWidth, p.TermHeight)
	p.GridTabPods.Set(
		ui.NewRow(2.0/2,
			ui.NewRow(1.0/10,
				ui.NewCol(.75/4, p.Tabpane),
				ui.NewCol(.75/4, p.Env),
			),
			ui.NewRow(9.0/20,
				ui.NewCol(1.0/2, p.App1),
				ui.NewCol(1.0/2, p.App1),
			),
			ui.NewRow(9.0/20,
				ui.NewCol(1.0/2, p.App1),
				ui.NewCol(1.0/2, p.App1),
			),
			
		),
	)
	p.resize()
	ui.Render(p.GridTabPods)
}

func (p *ControllerTabPods) Render(config *healthcheck.Config) {
	health := make(chan string)

	for _, v:= range config.Endpoints{
		go func() {
			healthcheck.Healthchecks(&v.Url, health)
    	}()
	}

	p.App1.Text = <-health
	p.resize()
	ui.Render(p.GridTabPods)
}

func (p *ControllerTabPods) Resize() {
	p.resize()
	ui.Render(p.GridTabPods)
}

func (p *ControllerTabPods) resize() {
	w, h := ui.TerminalDimensions()
	p.GridTabPods.SetRect(0, 0, w, h)
}

func CreateTabPane() (tabpane *widgets.TabPane) {
	tabpane = widgets.NewTabPane("Pods", "Offerings", "Cluster")
	tabpane.SetRect(0, 1, 50, 4)
	tabpane.Border = true
	tabpane.Title = "q = quit, h or l to switch tabs"
	tabpane.TitleStyle.Bg = termui.ColorBlue
	return
}