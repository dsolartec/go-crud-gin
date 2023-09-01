package main

import "go-crud-gin/cmd/server/app"

func main() {
	application := app.NewAppBuilder().Build()

	err := application.Run()
	if err != nil {
		panic(err)
	}
}
