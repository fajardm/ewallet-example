package http

import (
	"github.com/fajardm/ewallet-example/app/balance"
	"github.com/fajardm/ewallet-example/bootstrap"
	"github.com/fajardm/ewallet-example/errorcode"
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
	api.Get("/balances", middleware.Protected(), handler.GetBalance)
	api.Get("/balances/histories", middleware.Protected(), handler.GetBalanceHistories)
	api.Post("/balances/transfer", middleware.Protected(), handler.TransferBalance)
	api.Post("/balances/topup", middleware.Protected(), handler.TopUp)
}

func (b balanceHandler) GetBalance(ctx *fiber.Ctx) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errorcode.ErrBadParamInput.Error()})
		return
	}
	data, err := b.balanceUsecase.GetBalanceByUserID(ctx.Context(), *userID)
	if err != nil {
		ctx.Status(errorcode.StatusCode(err)).JSON(fiber.Map{"status": "error", "message": err.Error()})
		return
	}
	ctx.JSON(fiber.Map{"status": "success", "data": data})
}

func (b balanceHandler) GetBalanceHistories(ctx *fiber.Ctx) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errorcode.ErrBadParamInput.Error()})
		return
	}
	data, err := b.balanceUsecase.GetBalanceHistoriesByUserID(ctx.Context(), *userID)
	if err != nil {
		ctx.Status(errorcode.StatusCode(err)).JSON(fiber.Map{"status": "error", "message": err.Error()})
		return
	}
	ctx.JSON(fiber.Map{"status": "success", "data": data})
}

func (b balanceHandler) TransferBalance(ctx *fiber.Ctx) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errorcode.ErrBadParamInput.Error()})
		return
	}

	// Binds input
	type Input struct {
		ToUserID uuid.UUID `json:"to_user_id" validate:"required,max=36"`
		Amount   float64   `json:"amount" validate:"required"`
	}
	input := new(Input)
	if err := ctx.BodyParser(input); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errorcode.ErrBadParamInput.Error()})
		return
	}

	if err := b.balanceUsecase.TransferBalance(ctx.Context(), *userID, input.ToUserID, input.Amount); err != nil {
		ctx.Status(errorcode.StatusCode(err)).JSON(fiber.Map{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(fiber.Map{"status": "success", "data": true})
}

func (b balanceHandler) TopUp(ctx *fiber.Ctx) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errorcode.ErrBadParamInput.Error()})
		return
	}

	// Binds input
	type Input struct {
		Amount float64 `json:"amount" validate:"required"`
	}
	input := new(Input)
	if err := ctx.BodyParser(input); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errorcode.ErrBadParamInput.Error()})
		return
	}

	if err := b.balanceUsecase.TopUp(ctx.Context(), *userID, input.Amount); err != nil {
		ctx.Status(errorcode.StatusCode(err)).JSON(fiber.Map{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(fiber.Map{"status": "success", "data": true})
}
