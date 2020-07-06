package http

import (
	"github.com/fajardm/ewallet-example/app/errorcode"
	"github.com/fajardm/ewallet-example/app/user"
	"github.com/fajardm/ewallet-example/app/user/model"
	"github.com/fajardm/ewallet-example/bootstrap"
	"github.com/gofiber/fiber"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type userHandler struct {
	userUsecase user.Usecase
}

func NewUserHandler(app *bootstrap.Bootstrap, userUsecase user.Usecase) {
	handler := userHandler{userUsecase: userUsecase}
	api := app.Group("/api")
	api.Post("/users", handler.Store)
	api.Get("/users/:id", handler.GetByID)
}

func (u userHandler) Store(ctx *fiber.Ctx) {
	input := new(model.Input)
	if err := ctx.BodyParser(input); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errorcode.ErrBadParamInput.Error()})
		return
	}
	if err := input.Validate(); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errorcode.ErrBadParamInput.Error(), "data": err.Error()})
		return
	}
	user, err := input.NewUser()
	if err != nil {
		ctx.Status(errorcode.StatusCode(err)).JSON(fiber.Map{"status": "error", "message": err.Error()})
		return
	}
	if err := u.userUsecase.Store(ctx.Context(), *user); err != nil {
		ctx.Status(errorcode.StatusCode(err)).JSON(fiber.Map{"status": "error", "message": err.Error()})
		return
	}
	ctx.Status(http.StatusCreated).JSON(fiber.Map{"status": "success", "data": user})
}

func (u userHandler) GetByID(ctx *fiber.Ctx) {
	id, err := uuid.FromString(ctx.Params("id"))
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errorcode.ErrBadParamInput.Error()})
		return
	}
	user, err := u.userUsecase.GetByID(ctx.Context(), id)
	if err != nil {
		ctx.Status(errorcode.StatusCode(err)).JSON(fiber.Map{"status": "error", "message": err.Error()})
		return
	}
	ctx.JSON(fiber.Map{"status": "success", "data": user})
}
