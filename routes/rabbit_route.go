package routes

import (
	"rabbitmq/controllers"

	"github.com/gofiber/fiber/v2"
)

func CreateProofs(a *fiber.App) {
	a.Get("/proofs", controllers.Proofs)
}
