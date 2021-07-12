package routes

import (
	"github.com/edanko/moses/entities"
	"github.com/edanko/moses/service/nest"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gofiber/fiber/v2"
)

func getNestsAll(service nest.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		nests, err := service.GetAll()

		if err != nil {
			return c.JSON(&fiber.Map{
				"success": false,
				"error":   err,
			})
		}

		return c.JSON(&fiber.Map{
			"success": true,
			"nests":   nests,
		})
	}
}

func getNests(service nest.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		project := c.Params("project")
		dimension := c.Params("dimension")
		quality := c.Params("quality")

		nests, err := service.Get(project, dimension, quality)

		if err != nil {
			return c.JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		return c.JSON(&fiber.Map{
			"success": true,
			"nests":   nests,
		})
	}
}

func addOrUpdateNest(service nest.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestBody := new(entities.Nest)
		var result *entities.Nest
		err := c.BodyParser(requestBody)
		if err != nil {
			return c.JSON(&fiber.Map{
				"success": false,
				"error":   err,
			})
		}

		if requestBody.ID != primitive.NilObjectID {
			result, err = service.Update(requestBody)
		} else {
			result, err = service.Create(requestBody)
		}
		if err != nil {
			return c.JSON(&fiber.Map{
				"success": false,
				"error":   err,
			})
		}

		return c.JSON(&fiber.Map{
			"success": true,
			"data":    result,
			"error":   err,
		})
	}
}

func deleteNest(service nest.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dberr := service.Delete(c.Params("id"))
		if dberr != nil {
			return c.JSON(&fiber.Map{
				"success": false,
				"error":   dberr.Error(),
			})
		}
		return c.JSON(&fiber.Map{
			"success": true,
			"message": "deleted successfully",
		})
	}
}

func NestRouter(app fiber.Router, service nest.UseCase) {
	app.Get("/nests", getNestsAll(service))
	app.Get("/nest/:project/:dimension/:quality", getNests(service))
	app.Post("/nest", addOrUpdateNest(service))
	app.Delete("/nest/:id", deleteNest(service))
}
