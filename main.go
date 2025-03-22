package main

import (
	"fmt"

	"github.com/markoc1120/gator/internal/config"
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
