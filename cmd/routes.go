package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/greybluesea/dockerised-fullstack-webapp-gofiber-gorm-postgres/database"
	"github.com/greybluesea/dockerised-fullstack-webapp-gofiber-gorm-postgres/models"
)

func setupRoutes(app *fiber.App) {
	app.Get("/hello", helloHandler)
	app.Get("/", homeHandler)
	app.Get("/newfact", newfactHandler)
	app.Post("/create", createHandler)
	app.Post("/delete", deleteHandler)
	app.Get("/edit/:id", editHandler)
	app.Post("/update/:id", updateHandler)
}

func helloHandler(c *fiber.Ctx) error {
	return c.SendString("Hello, World! Tony here 👋!!")
}
func homeHandler(c *fiber.Ctx) error {
	facts := []models.Fact{}
	result := database.DB.Find(&facts)
	//	fmt.Println(reflect.TypeOf(facts[0].ID))
	// 	return c.Status(200).JSON(facts)

	fmt.Println(result)
	fmt.Println(facts)

	return c.Render("index", fiber.Map{"Title": "Fun Facts", "Subtitle": "Dockerised Fullstack WebApp(GoFiber + GORM + Postgres) - learned from Div Rhino", "Facts": facts})
}

func newfactHandler(c *fiber.Ctx) error {
	return c.Render("newfact", fiber.Map{"Title": "New Fact", "Subtitle": "Add an interesting fact"})
}

func editHandler(c *fiber.Ctx) error {
	fact := &models.Fact{}
	id := c.Params("id")

	result := database.DB.Where("id = ?", id).First(&fact)
	if result.Error != nil {
		return result.Error
		//return NotFound(c)
	}

	return c.Render("edit", fiber.Map{
		"Title": "Edit Fact",
		"Fact":  fact,
	})
}

func updateHandler(c *fiber.Ctx) error {
	fact := &models.Fact{}
	id := c.Params("id")

	if err := c.BodyParser(fact); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
	}

	result := database.DB.Model(fact).Where("id = ?", id).Updates(fact)
	if result.Error != nil {
		return editHandler(c)
	}

	return c.Redirect("/")
}

func createHandler(c *fiber.Ctx) error {
	fact := &models.Fact{}
	if err := c.BodyParser(fact); err != nil {
		/* return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		}) */
		return newfactHandler(c)
	}
	result := database.DB.Create(fact)
	if result.Error != nil {
		// return result.Error
		return newfactHandler(c)
	}
	//return c.Status(200).JSON(fact)
	// return  successHandler(c)
	return c.Redirect("/")
}

/* func successHandler(c *fiber.Ctx) error {
	return c.Render("success", fiber.Map{
		"Title": "Fact added successfully",
	})
} */

func deleteHandler(c *fiber.Ctx) error {
	factID := c.FormValue("ID")
	// fmt.Println(reflect.TypeOf(factID))
	result := database.DB.Delete(&models.Fact{}, factID)
	if result.Error != nil {
		return result.Error
	}
	c.Redirect("/")
	return nil
}
