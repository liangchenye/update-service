package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Con struct {
	Login         string
	MonitRepo     string
	Contributions int
}

func main() {
	data, _ := ioutil.ReadFile("githubdata/mapdata")
	contrib := make(map[string]int)

	json.Unmarshal(data, &contrib)

	var content string
	for k, v := range contrib {
		c := fmt.Sprintf("%s %d\n", k, v)
		content += c
	}
	ioutil.WriteFile("githubdata/mapdata", []byte(content), 0644)
}
