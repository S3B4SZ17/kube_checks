# kube_checks

CLI app that uses a GUI to monitor kubernetes pods and performs API calls

### USAGE

1. Clone the repo: `git clone git@github.com:S3B4SZ17/kube_checks.git` \
2. Go to the new directory: `cd kube_checks`
3. Run `go install`
4. Source the current shell `source .` or open a new shell
5. The name of the binary is `kube_checks` \
   `kube_checks -t [vault.token-string]`

### Flags

`-t` / `--token` -> The token we use to authenticate to Vault \
`-k` / `--kubeconfig` -> The absolute path to the kubeconfig file for the kubernetes cluster and context
`-f` / `--file`
