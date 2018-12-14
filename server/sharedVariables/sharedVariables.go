package shared

import (
	"fmt"
	"os"

	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
)

// Address of database
var Address interface{}

//DbName is the name of the database
var DbName interface{}

//Port where the app runs
var Port interface{}

var workEnv string

//GetConfig is used to get the configurations and use them accross packages
func init() {
	workEnv = os.Getenv("WorkEnv")
	config.Load(file.NewSource(
		file.WithPath("./config/config.json"),
	))
	conf := config.Map()
	if workEnv == "prod" {
		Address = conf["prod"].(map[string]interface{})["database"].(map[string]interface{})["address"]
		DbName = conf["prod"].(map[string]interface{})["database"].(map[string]interface{})["db_name"]
		Port = conf["prod"].(map[string]interface{})["api"].(map[string]interface{})["port"]
	} else if workEnv == "test" {
		Address = conf["test"].(map[string]interface{})["database"].(map[string]interface{})["address"]
		DbName = conf["test"].(map[string]interface{})["database"].(map[string]interface{})["db_name"]
		Port = conf["test"].(map[string]interface{})["api"].(map[string]interface{})["port"]
	} else if workEnv == "dev" || workEnv == "" {
		Address = conf["dev"].(map[string]interface{})["database"].(map[string]interface{})["address"]
		DbName = conf["dev"].(map[string]interface{})["database"].(map[string]interface{})["db_name"]
		Port = conf["dev"].(map[string]interface{})["api"].(map[string]interface{})["port"]
	} else {
		fmt.Println("Invalid WorkEnv. Valid options are dev, prod and  test")
	}
}
