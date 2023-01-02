package adapters

import (
	"github.com/gofiber/fiber/v2"
	"github.com/luminosita/bee/internal/infra/http"
	"github.com/luminosita/bee/internal/infra/http/handlers"
)

type errorResponse struct {
	error string
}

func Convert(handler *handlers.BaseHandler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		response := handler.Handle(&http.HttpRequest{
			Body:    ctx.Body(),
			Params:  ctx.AllParams(),
			Headers: ctx.GetReqHeaders(),
		})

		ctx.SendStatus(response.StatusCode)
		if response.StatusCode >= 200 && response.StatusCode <= 299 {
			return ctx.JSON(response.Body)
		} else {
			return ctx.JSON(&errorResponse{
				//				error: response.Errors,
			})
		}
	}
}
