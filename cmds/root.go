package cmds

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/S3B4SZ17/kube_checks/terminal"

	"github.com/S3B4SZ17/kube_checks/healthcheck"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
    // file, _ := cmd.Flags().GetString("file")
    k8s := healthcheck.Connection(&kubeconfig)
    if token == "" {
      cmd.Usage()
      return
    }else{
      healthcheck.SetCredentials(token)
    }

    terminal.Run(k8s)
    // fmt.Println(k8s)
    // readYaml(file)

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

func readYaml(file string) {
  yfile, err := ioutil.ReadFile(file)
  if err != nil {
    log.Fatal(err)
  }
  data := make(map[string][]Endpoints)

  err2 := yaml.Unmarshal(yfile, &data)

  if err2 != nil {

      log.Fatal(err2)
  }

  for k, v := range data {

      fmt.Printf("%s: %s\n", k, v)
  }
}

type Endpoints struct {
  App string
  Url string
}