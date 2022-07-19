package cmds

import (
	"fmt"
	"log"
	"os"

	"github.com/S3B4SZ17/kube_checks/healthcheck"
	"github.com/S3B4SZ17/kube_checks/terminal"

	"github.com/spf13/cobra"
)
var rootCmd = &cobra.Command{
  Use:   "kube_checks",
  Short: "A CLI app that automates and helps to debug issues in REST API calls env",
  Long: `A CLI app that automates and helps to debug issues in Kubernetes clusters env.
                With the help of a GUI we can easily detect issues in the cluster.
                Complete documentation is available at http://github.com/S3B4SZ17/kube_checks`,
  Run: func(cmd *cobra.Command, args []string) {
    // Do Stuff Here
    token, _ := cmd.Flags().GetString("token")
    kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
    file, _ := cmd.Flags().GetString("file")
    k8s := healthcheck.Connection(&kubeconfig)
    config, err := healthcheck.ReadYaml(file); if err != nil {
        log.Fatal(err.Error())
    }
    if token == "" {
      cmd.Usage()
      return
    }else{
      
      healthcheck.SetCredentials(token, config.Vault_url)
    }
    
    terminal.Run(k8s, config)
    // fmt.Println(k8s)

  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
    var path string
    var token string
    var values string
    rootCmd.Flags().StringVarP(&path, "kubeconfig", "k", "", "(Optiional) Absolute path to the kubeconfig file")
    rootCmd.Flags().StringVarP(&token, "token", "t", "", "The Vault temporary token")
    rootCmd.Flags().StringVarP(&values, "file", "f", "", "(Optional) A yaml file that defines custom values")
}
