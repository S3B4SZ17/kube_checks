package tabcluster

import (
	"math"

	"github.com/S3B4SZ17/kube_checks/healthcheck"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Controller interface {
	Render()
	Resize()
}

type ControllerTabCluster struct {
	// Tab Cluster
	Sl             *widgets.Sparkline
	Slg            *widgets.SparklineGroup
	Lc             *widgets.Plot
	Gs             []*widgets.Gauge
	P              *widgets.Paragraph
	GridTabCluster *ui.Grid
	P_exec         *widgets.Paragraph

	// Tabpane
	Tabpane *widgets.TabPane
	Env     *widgets.Paragraph

	// Grid dimension
	TermWidth  int
	TermHeight int

	//Data
	SinFloat64 []float64
	K8s        *healthcheck.K8s
}

func NewControllerTabCluster(env *widgets.Paragraph, tabpane *widgets.TabPane, k8s *healthcheck.K8s) *ControllerTabCluster {
	ctl := &ControllerTabCluster{
		Sl:             widgets.NewSparkline(),
		Slg:            widgets.NewSparklineGroup(),
		Lc:             widgets.NewPlot(),
		Gs:             make([]*widgets.Gauge, 3),
		P:              widgets.NewParagraph(),
		GridTabCluster: ui.NewGrid(),
		Env:            env,
		P_exec:         widgets.NewParagraph(),

		Tabpane: tabpane,

		K8s: k8s,
	}

	ctl.initUI()

	return ctl
}

func (p *ControllerTabCluster) initUI() {

	p.SinFloat64 = (func() []float64 {
		n := 400
		data := make([]float64, n)
		for i := range data {
			data[i] = 1 + math.Sin(float64(i)/5)
		}
		return data
	})()

	p.Sl.Data = p.SinFloat64[:100]
	p.Sl.LineColor = ui.ColorCyan
	p.Sl.TitleStyle.Fg = ui.ColorWhite

	p.Slg = widgets.NewSparklineGroup(p.Sl)
	p.Slg.Title = "Sparkline"

	p.Lc.Title = "braille-mode Line Chart"
	p.Lc.Data = append(p.Lc.Data, p.SinFloat64)
	p.Lc.AxesColor = ui.ColorWhite
	p.Lc.LineColors[0] = ui.ColorYellow

	for i := range p.Gs {
		p.Gs[i] = widgets.NewGauge()
		p.Gs[i].Percent = i * 10
		p.Gs[i].BarColor = ui.ColorRed
	}

	p.P.Text = p.K8s.GetFailingPods()
	p.P.Title = p.Env.Text + " Pods with errors"
	p.P.TitleStyle.Fg = ui.ColorGreen
	p.P.PaddingLeft = 2
	p.P.PaddingTop = 1

	//Defining options to exec into a pod
	// commands := []string{"curl", "-ks", "https://localhost:443/summary"}
	// pod := p.K8s.GetPodName("name")
	// p.P_exec.Text = p.K8s.Exec("namespace", pod, commands)
	p.P_exec.Title = "Summary"
	p.P_exec.TitleStyle.Fg = ui.ColorGreen
	p.P_exec.PaddingLeft = 2
	p.P_exec.PaddingTop = 1

	p.TermHeight, p.TermWidth = ui.TerminalDimensions()
	p.GridTabCluster.SetRect(0, 0, p.TermWidth, p.TermHeight)
	p.GridTabCluster.Set(
		ui.NewRow(2.0/2,
			ui.NewRow(1.0/10,
				ui.NewCol(.75/4, p.Tabpane),
				ui.NewCol(.75/4, p.Env),
			),
			ui.NewRow(9.0/20,
				ui.NewCol(1.0/2, p.Slg),
				ui.NewCol(1.0/2, p.Lc),
			),
			ui.NewRow(9.0/20,
				// ui.NewCol(1.0/2,
				// 	ui.NewRow(.9/3, p.Gs[0]),
				// 	ui.NewRow(.9/3, p.Gs[1]),
				// 	ui.NewRow(1.2/3, p.Gs[2]),
				// ),
				ui.NewCol(1.0/2, p.P_exec),
				ui.NewCol(1.0/2, p.P),
			),
		),
	)
	p.resize()

}

func (p *ControllerTabCluster) Render() {
	failed_pods := make(chan string)
	go func() {
		result := p.K8s.GetFailingPods()
		failed_pods <- result
	}()
	p.P.Text = <-failed_pods
	p.resize()
	ui.Render(p.GridTabCluster)
}

func (p *ControllerTabCluster) Resize() {
	p.resize()
	ui.Render(p.GridTabCluster)
}

func (p *ControllerTabCluster) resize() {
	w, h := ui.TerminalDimensions()
	p.GridTabCluster.SetRect(0, 0, w, h)
}
