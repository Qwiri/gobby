package main

import (
	"github.com/Qwiri/gobby/pkg/gobby"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	g := gobby.New()
	if err := g.Hook(app); err != nil {
		panic(err)
	}
}
