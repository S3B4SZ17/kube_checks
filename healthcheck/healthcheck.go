package healthcheck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/dlclark/regexp2"
	vault "github.com/hashicorp/vault/api"
)

var token *string
var credentials map[string]interface{}
var url string
var basic_auth string
var env string

func GetEnv() string {
	cmd := "kubectl config current-context"
	rgx := regexp2.MustCompile("(?<=\\/).*(?=:\\d+)", 0)

	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatal("[Error] Please log in to the cluster: kubectl login ...")
	}
	if env, _ := rgx.FindStringMatch(string(out)); env != nil {
		return env.String()
	}else{
		log.Fatalf("Please check if you are logged in to cluster.")
		return ""
	}
	
}

func Healthchecks(api_endpoint *string, ch chan string) {
	var result string
	var mapObj map[string]interface{}
	var requestBody bytes.Buffer
	url_endpoint := url + *api_endpoint
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url_endpoint, &requestBody)
	req.Method = "GET"
	req.Header = map[string][]string{
		"Authorization": {"Basic " + basic_auth},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &mapObj); err != nil {
		result = "[Error] Didnt received a JSON response. Please check the "+ *api_endpoint + " endpoint. " + err.Error()
		return
	}
	for key, val := range mapObj{
		valStr := fmt.Sprint(val)
		newVal := strings.Replace(valStr, "map", "", -1)
		result += fmt.Sprintf("%v: %v\n", key, newVal)
	}

	ch <- result

}

func getVultCredentials(vault_url *string) {
	vault_path := checkIfProd()
	config := vault.DefaultConfig()
	config.Address = *vault_url
	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("[Vault Error]: Unable to initialize Vault client: %v", err)
	}
	client.SetToken(*token)

	secret, err := client.Logical().Read(*vault_path)
	if err != nil {
		log.Fatalf("[Vault Error]: Unable to read secret: %v \n\n Probably your token is expired.", err)
	}
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		log.Fatalf("Data type assertion failed: %T %#v", secret.Data["data"], secret.Data["data"])
	}
	credentials = data
}

func checkIfProd() *string {
	//[Invalid path for a versioned K/V secrets engine. See the API docs for the appropriate API endpoints to use. If using the Vault CLI, use 'vault kv get' for this operation.
	// Use the /secret/data/ path for V2 of the Vault engine
	var vault_path string
	if strings.Contains(env, "prod") {
		vault_path = "/secret/data/EC2/" + env + "/credentials"
	} else if strings.Contains(env, "dev") {
		vault_path = "/secret/data/EC2/dev/" + env + "/credentials"
	}else {
		vault_path = "/secret/data/EC2/credentials"
	}

	return &vault_path
}

func SetCredentials(token_ref string, vault_url string) {
	token = &token_ref
	env = GetEnv()
	getVultCredentials(&vault_url)
	url = credentials["url_app"].(string)
	basic_auth = credentials["API_KEY"].(string)
}