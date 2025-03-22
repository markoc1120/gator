package main

import (
	"fmt"
	"gator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		panic(err)
	}

	config.SetUser("mark", conf)

	conf, err = config.Read()
	if err != nil {
		panic(err)
	}

	fmt.Println(conf)
}
