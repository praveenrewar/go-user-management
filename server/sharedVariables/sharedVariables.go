package shared

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//Environments is an array of Environment
type Environments struct {
	Environments []Environment `json:"environments"`
}

// Environment describes the config values for various work environments
type Environment struct {
	Env      string   `json:"env"`
	Database Database `json:"database"`
	API      API      `json:"api"`
}

// Database defines database config values
type Database struct {
	Address    string `json:"address"`
	DBName     string `json:"db_name"`
	DBUsername string `json:"db_username"`
	DBPassword string `json:"db_password"`
}

//API defines the api config values
type API struct {
	Port int `json:"port"`
}

// Address of database
var Address string

//DbName is the name of the database
var DbName string

//DbUsername is the username for the db
var DbUsername string

//DbPassword is the password for the db
var DbPassword string

//Port where the app runs
var Port int

//WorkEnv is used to store Wroking Environment
var WorkEnv string

//GetConfig is used to get the configurations and use them accross packages
func init() {
	WorkEnv = os.Getenv("WorkEnv")
	jsonFile, err := os.Open("./config/config.json")
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var environments Environments
	json.Unmarshal(byteValue, &environments)
	var checkWorkEnv = 0
	for i := 0; i < len(environments.Environments); i++ {
		if WorkEnv == environments.Environments[i].Env || (WorkEnv == "" && environments.Environments[i].Env == "dev") {
			Address = environments.Environments[i].Database.Address
			DbName = environments.Environments[i].Database.DBName
			DbUsername = environments.Environments[i].Database.DBUsername
			DbPassword = environments.Environments[i].Database.DBPassword
			Port = environments.Environments[i].API.Port
			checkWorkEnv = 1
		}
	}
	if checkWorkEnv == 0 {
	}
}
