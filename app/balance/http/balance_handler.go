package http

import (
	"github.com/fajardm/ewallet-example/app/balance"
	"github.com/fajardm/ewallet-example/app/errorcode"
	"github.com/fajardm/ewallet-example/bootstrap"
	"github.com/fajardm/ewallet-example/middleware"
	"github.com/gofiber/fiber"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type balanceHandler struct {
	balanceUsecase balance.Usecase
}

func NewBalanceHandler(app *bootstrap.Bootstrap, balanceUsecase balance.Usecase) {
	handler := balanceHandler{balanceUsecase: balanceUsecase}
	api := app.Group("/api")
	api.Get("/balances/:user_id", middleware.Protected(), handler.GetBalanceByUserID)
	api.Get("/balances/histories/:user_id/", middleware.Protected(), handler.GetBalanceHistoriesByUserID)
}

func (b balanceHandler) GetBalanceByUserID(ctx *fiber.Ctx) {
	userID, err := uuid.FromString(ctx.Params("user_id"))
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errorcode.ErrBadParamInput.Error()})
		return
	}
	data, err := b.balanceUsecase.GetBalanceByUserID(ctx.Context(), userID)
	if err != nil {
		ctx.Status(errorcode.StatusCode(err)).JSON(fiber.Map{"status": "error", "message": err.Error()})
		return
	}
	ctx.JSON(fiber.Map{"status": "success", "data": data})
}

func (b balanceHandler) GetBalanceHistoriesByUserID(ctx *fiber.Ctx) {
	userID, err := uuid.FromString(ctx.Params("user_id"))
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errorcode.ErrBadParamInput.Error()})
		return
	}
	data, err := b.balanceUsecase.GetBalanceHistoriesByUserID(ctx.Context(), userID)
	if err != nil {
		ctx.Status(errorcode.StatusCode(err)).JSON(fiber.Map{"status": "error", "message": err.Error()})
		return
	}
	ctx.JSON(fiber.Map{"status": "success", "data": data})
}
