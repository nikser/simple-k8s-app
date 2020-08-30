package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	envConfigPath  = "CONFIG_PATH"
	envDatabaseUrl = "DATABASE_URL"
)

type conf struct {
	DatabaseUrl string `yaml:"databaseUrl"`
}

func main() {
	var config conf
	err := config.readConfig()
	if err != nil {
		log.Fatal(err)
	}

	startServer(&config)
}

func (c *conf) readConfig() error {
	if databaseUrl, ok := os.LookupEnv(envDatabaseUrl); ok {
		c.DatabaseUrl = databaseUrl
		return nil
	}

	configPath, ok := os.LookupEnv(envConfigPath)
	if !ok {
		configPath = "config.yaml"
	}

	if yamlFile, err := ioutil.ReadFile(configPath); err == nil {
		if err = yaml.Unmarshal(yamlFile, c); err != nil {
			return errors.New(fmt.Sprintf("Unmarshal: %w", err))
		}
		return nil
	}

	return errors.New("ENV variable 'DATABASE_URL' has not found! Define the variable before start this application.")
}

func startServer(c *conf) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, fmt.Sprintf("Congratulation!\nAccess received by DATABASE_URL: %s", c.DatabaseUrl))
	})

	log.Println("Start server")
	http.ListenAndServe(":8000", nil)
}
