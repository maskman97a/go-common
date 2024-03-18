package main

import "go-common/config"

func main() {
	app := config.NewApp()
	app.Init()
	app.Run()
}
