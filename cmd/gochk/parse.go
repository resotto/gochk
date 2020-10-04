package gochk

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

// Config is data converted from config.json
type Config struct {
	TargetPath       string
	DependencyOrders []string
	IgnoreDirs       []string
}

// Parse parses config.json
func Parse() Config {
	absPath, _ := filepath.Abs("configs/config.json") // NOTICE: from root directory
	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err)
	}
	var config Config
	json.Unmarshal(bytes, &config)
	return config
}
