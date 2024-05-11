package main

import (
	"github.com/fvbock/endless"
	"github.com/nazerkhh/shop/routers"
)

func main() {
	router := routers.Setup()

	err := endless.ListenAndServe(":8080", router) // instead of router.Run()
	if err != nil {
		return
	}
}
