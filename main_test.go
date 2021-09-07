package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	start()
}

// func TestToml(t *testing.T) {
// 	c := Config{}
// 	_, err := toml.DecodeFile("./configs/config.toml", &c)
// 	if err != nil {
// 		panic(err)
// 	}
// 	for k, v := range c.User[0].Postdata {
// 		fmt.Printf("k: %v\n", k)
// 		fmt.Printf("v: %v\n", v)
// 	}
// }
