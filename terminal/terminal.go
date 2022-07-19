package terminal

import (
	"log"
	"time"

	"github.com/S3B4SZ17/kube_checks/tabcluster"

	"github.com/S3B4SZ17/kube_checks/healthcheck"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func Run(k8s *healthcheck.K8s) {
   if err := ui.Init(); err != nil {
      log.Fatalf("[Error] Failed to initialize termui: %v", err)
   }
   defer ui.Close()

 	env := *widgets.NewParagraph()  	
   	env.Text = healthcheck.GetEnv()
	env.TextStyle.Bg = ui.ColorCyan
	env.PaddingLeft = 1
	env.Title = "Cluster name"
   	
	tabpane := tabcluster.CreateTabPane()

	
   	cluster_controller := tabcluster.NewControllerTabCluster(&env, tabpane, k8s)
   	errors_controller := tabcluster.NewControllerTabErrors(&env, tabpane)
	pods_controller := tabcluster.NewControllerTabPods(&env, tabpane, k8s)

	renderTab := func() {
		switch tabpane.ActiveTabIndex {
		case 0:
			pods_controller.Render()
		case 1:
			errors_controller.Render()
		case 2:
			cluster_controller.Render()
		}
	}
	tickerCount := 1
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				renderTab()
			case "h":
				tabpane.FocusLeft()
				ui.Clear()
				ui.Render(tabpane)
				renderTab()
			case "l":
				tabpane.FocusRight()
				ui.Clear()
				ui.Render(tabpane)
				renderTab()
			}
		case <-ticker:
			if tickerCount % 3 == 0{
				for _, g := range cluster_controller.Gs {
					g.Percent = (g.Percent + 3) % 100
				}
				cluster_controller.Slg.Sparklines[0].Data = cluster_controller.SinFloat64[tickerCount : tickerCount+100]
				cluster_controller.Lc.Data[0] = cluster_controller.SinFloat64[2*tickerCount:]
				ui.Render(tabpane)
				renderTab()
				tickerCount++
			}
			tickerCount++
		}
	}
}
