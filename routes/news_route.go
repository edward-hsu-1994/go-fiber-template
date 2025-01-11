package routes

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-template/base"
	_ "go-fiber-template/domain"
)

type NewsRouter struct {
	Config *base.Config
}

func NewNewsRouter(config *base.Config) *NewsRouter {
	return &NewsRouter{
		Config: config,
	}
}

func (r *NewsRouter) ConfigureRoutes(app *fiber.App) {
	routes := app.Group("/api/v1/news")

	routes.Get("/", ListNews)
	routes.Get("/test", r.Test)
	routes.Get("/:newsId", GetNewsById)
}

func ListNews(c *fiber.Ctx) error {
	return c.SendString("List news")
}

func GetNewsById(c *fiber.Ctx) error {
	newsId := c.Params("newsId")

	return c.SendString("Get news by id: " + newsId)
}

// Test godoc
// @Summary Test
// @Description test
// @Tags news
// @Produce json
// @Success 200 {string} string "Test"
// @Router /api/v1/news/test [get]
func (r *NewsRouter) Test(c *fiber.Ctx) error {
	connStr := (*r.Config)["connectionString"].(string)

	return c.SendString("Get news by id: " + connStr)
}
