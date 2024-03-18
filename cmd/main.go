package main

import "github.com/maskman97a/go-common/config"

func main() {
	app := config.NewApp()
	app.Init()
	app.Run()
}
