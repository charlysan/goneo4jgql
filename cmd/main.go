package main

import (
	app "github.com/charlysan/goneo4jgql/internal/app"
)

func main() {
	myApp := app.Init()
	myApp.InitRoutes()
	myApp.Run()
}
