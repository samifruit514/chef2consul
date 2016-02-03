package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"os"
	"os/exec"
	"reflect"
)

type jsonErrors struct {
	Errors []string `json:"errors"`
}

// Input ...
type Input struct {
	ChefNode      string
	ChefAttribute string
	KnifeRbFile   string
}

// ConsulConfig ...
type ConsulConfig struct {
	Host   string
	Token  string
	Prefix string
}

var input = &Input{}
var consulConfig = &ConsulConfig{}

func init() {
	input.ChefNode = os.Args[1]
	input.ChefAttribute = os.Args[2]
	input.KnifeRbFile = os.Getenv("KNIFERB_FILE")

	consulConfig.Prefix = os.Getenv("CONSUL_PREFIX")
	consulConfig.Host = os.Getenv("CONSUL_HOST")
	consulConfig.Token = os.Getenv("CONSUL_TOKEN")

}

func main() {

	if _, err := os.Stat(input.KnifeRbFile); os.IsNotExist(err) {
		log.Fatal("You must provide a valid path to the knife RB (is KNIFERB_FILE env var set?)")
	}

	cmd := exec.Command("knife", "node", "show", input.ChefNode, "-c", input.KnifeRbFile, "-F", "json", "-a", input.ChefAttribute)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("An error ocurred while running chef knife: %s\n", err)
		log.Fatal(string(out))
	} else {
		//body := string(out)
		var f map[string]interface{}
		err := json.Unmarshal(out, &f)
		if err != nil {
			fmt.Println("Cant read from knife")
		} else {
			chefConfig := f[input.ChefNode].(map[string]interface{})[input.ChefAttribute]
			fullStruct := make(map[string]string)
			processNode(consulConfig.Prefix, fullStruct, chefConfig)
			saveItem(fullStruct, consulConfig)
		}

	}
}

func saveItem(pairs map[string]string, consulConf *ConsulConfig) {

	os.Setenv("CONSUL_HTTP_SSL_VERIFY", "false")
	config := api.DefaultConfig()
	config.Address = consulConf.Host
	config.Token = consulConf.Token
	config.Scheme = "https"

	client, _ := api.NewClient(config)
	kv := client.KV()

	for k, v := range pairs {
		p := &api.KVPair{Key: k, Value: []byte(v)}
		_, err := kv.Put(p, nil)
		if err != nil {
			panic(err)
		}
	}

}

func processNode(fullPath string, fullStruct map[string]string, curNode interface{}) {

	node := curNode.(map[string]interface{})

	for k, v := range node {
		valType := reflect.TypeOf(v)
		if fmt.Sprintf("%v", valType) == "map[string]interface {}" {
			processNode(fullPath+"/"+k, fullStruct, v)
		} else if fmt.Sprintf("%v", valType) == "[]interface {}" {
			for ak, element := range v.([]interface{}) {
				fakeKey := make(map[string]interface{})
				fakeKey[fmt.Sprintf("%v", ak)] = element
				processNode(fmt.Sprintf("%v/%v", fullPath, k), fullStruct, fakeKey)
			}
		} else {
			finalPath := fullPath + "/" + k
			fullStruct[finalPath] = fmt.Sprintf("%v", v)
		}
	}

}
