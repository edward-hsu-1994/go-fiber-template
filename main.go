package main

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go-fiber-template/accesses"
	"go-fiber-template/base"
	docs "go-fiber-template/docs"
	"go-fiber-template/helpers"
	"go-fiber-template/routes"
	"go-fiber-template/services"
	"go.uber.org/fx"
	"log"
	"os"
	"strings"
)

// @title Go Fiber Template
// @version 1.0
// @description Go Fiber Template
// @license MIT
// @BasePath /
func main() {
	app := fx.New(
		fx.Provide(
			LoadConfig,
			accesses.NewMockPostRepository,
			services.NewPostService,
			routes.NewPostRouter,
			routes.NewNewsRouter,
			func(
				newsRouter *routes.NewsRouter,
				postRouter *routes.PostRouter,
			) []routes.FiberRouter {
				return []routes.FiberRouter{newsRouter, postRouter}
			},
			FiberConfig,
			NewFiberApp,
		),
		fx.Invoke(func(app *fiber.App) {
			log.Println("Starting Fiber app on port 3000...")
			if err := app.Listen(":3000"); err != nil {
				log.Fatalf("Failed to start Fiber app: %v", err)
			}
		}),
	)

	app.Start(context.Background())
}

func NewFiberApp(
	lc fx.Lifecycle,
	fiberConfig []fiber.Config,
	routers []routes.FiberRouter,
) *fiber.App {
	app := fiber.New(fiberConfig...)

	// Setting hostname for swagger
	alreadySettingSwaggerHostname := false
	app.Use(func(c *fiber.Ctx) error {
		if alreadySettingSwaggerHostname {
			return c.Next()
		}

		if strings.HasPrefix(c.Path(), "/swagger/") == false {
			return c.Next()
		}

		alreadySettingSwaggerHostname = true

		docs.SwaggerInfo.Host = c.GetRespHeader("X-Forwarded-For", c.Hostname())

		return c.Next()
	})
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			return helpers.ErrorResponse(c, err)
		}
		return err
	})

	for _, router := range routers {
		router.ConfigureRoutes(app)
	}

	return app
}

func FiberConfig() []fiber.Config {
	return []fiber.Config{}
}

func LoadConfig() *base.Config {
	file, err := os.Open("./config/config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	config := new(base.Config)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}
