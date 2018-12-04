package shared

import (
	"fmt"
	"os"

	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
)

// Address of database
var Address interface{}

//Port of database
var Port interface{}

//GetConfig is used to get the configurations and use them accross packages
func GetConfig() int {
	workEnv := os.Getenv("WorkEnv")
	config.Load(file.NewSource(
		file.WithPath("./config/config.json"),
	))
	conf := config.Map()
	flag := 1
	if workEnv == "prod" {
		Address = conf["prod"].(map[string]interface{})["database"].(map[string]interface{})["address"]
		Port = conf["prod"].(map[string]interface{})["api"].(map[string]interface{})["port"]
		flag = 1
	} else if workEnv == "test" {
		Address = conf["test"].(map[string]interface{})["database"].(map[string]interface{})["address"]
		Port = conf["test"].(map[string]interface{})["api"].(map[string]interface{})["port"]
		flag = 1
	} else if workEnv == "dev" || workEnv == "" {
		Address = conf["dev"].(map[string]interface{})["database"].(map[string]interface{})["address"]
		Port = conf["dev"].(map[string]interface{})["api"].(map[string]interface{})["port"]

	} else {
		fmt.Println("Invalid WorkEnv. Valid options are dev, prod and  test")
		flag = 0
	}
	fmt.Println(Address)
	return flag
}
