package main

import (
	"encoding/json"
	"os"
)

type config struct {
	// Realm is the Message displayed in the brower login box
	Realm string
	// Path of the directory to serve
	Path string
	// Users is a map of username=>[passwords]
	Users map[string][]string
}

func loadConfig(p string) (*config, error) {
	c := &config{}
	f, err := os.Open(p)
	if err != nil {
		return c, nil
	}
	defer f.Close()
	return c, json.NewDecoder(f).Decode(c)
}
