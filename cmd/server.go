package main

import (
	"github.com/alash3al/katch/pkg/katch"
	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"
)

func server() cli.ActionFunc {
	return func(ctx *cli.Context) error {
		app := fiber.New()

		app.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("Hello, World!")
		})

		app.Get("/export", func(c *fiber.Ctx) error {
			var input katch.Input

			if err := c.QueryParser(&input); err != nil {
				return c.Status(400).SendString(err.Error())
			}

			if input.URL == "" {
				return c.Status(400).SendString("empty url specified")
			}

			output, err := katch.Katch(c.Context(), input)
			if err != nil {
				return c.Status(500).SendString(err.Error())
			}

			switch input.OutputFormat {
			case katch.OutputFormatPNG:
				c.Set("Content-Type", "image/png")
			case katch.OutputFormatJPEG:
				c.Set("Content-Type", "image/jpeg")
			case katch.OutputFormatPDF:
				c.Set("Content-Type", "application/pdf")
			case katch.OutputFormatHTML:
				c.Set("Content-Type", "text/html")
			}

			return c.Status(200).Send(output)
		})

		app.Post("/export", func(c *fiber.Ctx) error {
			var input katch.Input

			if err := c.BodyParser(&input); err != nil {
				return c.Status(400).SendString(err.Error())
			}

			if input.URL == "" {
				return c.Status(400).SendString("empty url specified")
			}

			output, err := katch.Katch(c.Context(), input)
			if err != nil {
				return c.Status(500).SendString(err.Error())
			}

			switch input.OutputFormat {
			case katch.OutputFormatPNG:
				c.Set("Content-Type", "image/png")
			case katch.OutputFormatJPEG:
				c.Set("Content-Type", "image/jpeg")
			case katch.OutputFormatPDF:
				c.Set("Content-Type", "application/pdf")
			case katch.OutputFormatHTML:
				c.Set("Content-Type", "text/html")
			}

			return c.Status(200).Send(output)
		})

		return app.Listen(ctx.String("listen"))
	}
}
