package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
)

type JsonResponse struct {
	HelloMessage string `json:"helloMessage"`
}

func readiness(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}

func health(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}

func echo(ctx *fiber.Ctx) error {
	name := ctx.Query("name")
	if name == "" {
		name = "Anonymous"
	}

	hw := &JsonResponse{
		HelloMessage: "Hi " + name,
	}

	helloBytes, err := json.MarshalIndent(hw, "", "    ")
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.Status(fiber.StatusOK).SendString(string(helloBytes))
}

func main() {

	app := fiber.New()

	app.Get("/hi", func(ctx *fiber.Ctx) error {
		return echo(ctx)
	})

	app.Get("/readiness", func(ctx *fiber.Ctx) error {
		return readiness(ctx)
	})

	app.Get("/health", func(ctx *fiber.Ctx) error {
		return health(ctx)
	})

	log.Fatal(app.Listen(":3000"))
}
