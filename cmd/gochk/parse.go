package gochk

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

// Config is converted data structure from config.json
type Config struct {
	DefaultTargetPath string
	DependencyOrders  []string
}

// Parse parses config.json and return its values
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
