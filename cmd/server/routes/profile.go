package routes

import (
	"github.com/edanko/moses/entities"
	"github.com/edanko/moses/service/profile"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gofiber/fiber/v2"
)

func getProfilesAll(service profile.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		profiles, err := service.GetAll()

		if err != nil {
			return c.JSON(&fiber.Map{
				"success": false,
				"error":   err,
			})
		}

		return c.JSON(&fiber.Map{
			"success":  true,
			"profiles": profiles,
		})
	}
}

func getProfiles(service profile.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		project := c.Params("project")
		dimension := c.Params("dimension")
		quality := c.Params("quality")

		profiles, err := service.Get(project, dimension, quality)

		if err != nil {
			return c.JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		return c.JSON(&fiber.Map{
			"success":  true,
			"profiles": profiles,
		})
	}
}

func addOrUpdateProfile(service profile.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestBody := new(entities.Profile)
		var result *entities.Profile
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

func deleteProfile(service profile.UseCase) fiber.Handler {
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

func ProfileRouter(app fiber.Router, service profile.UseCase) {
	app.Get("/profiles", getProfilesAll(service))
	app.Get("/profile/:project/:dimension/:quality", getProfiles(service))
	app.Post("/profile", addOrUpdateProfile(service))
	app.Delete("/profile/:id", deleteProfile(service))
}
