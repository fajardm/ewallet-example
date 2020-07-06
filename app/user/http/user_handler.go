package http

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/fajardm/ewallet-example/app/base"
	"github.com/fajardm/ewallet-example/app/errorcode"
	"github.com/fajardm/ewallet-example/app/user"
	"github.com/fajardm/ewallet-example/app/user/model"
	"github.com/fajardm/ewallet-example/bootstrap"
	"github.com/fajardm/ewallet-example/middleware"
	"github.com/fajardm/ewallet-example/validator"
	"github.com/gofiber/fiber"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

type userHandler struct {
	userUsecase user.Usecase
}

func NewUserHandler(app *bootstrap.Bootstrap, userUsecase user.Usecase) {
	handler := userHandler{userUsecase: userUsecase}
	api := app.Group("/api")
	api.Post("/users/login", handler.Login)
	api.Post("/users", handler.Store)
	api.Use(middleware.Protected())
	api.Get("/users/:id", handler.GetByID)
	api.Put("/users/:id", handler.Update)
	api.Delete("/users/:id", handler.Delete)
}

func (u userHandler) Login(ctx *fiber.Ctx) {
	// Binds input
	type Input struct {
		UsernameOrEmail string `json:"username_or_email" validate:"required"`
		Password        string `json:"password" validate:"required"`
	}
	input := new(Input)
	if err := ctx.BodyParser(input); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errorcode.ErrBadParamInput.Error()})
		return
	}

	// Validate input
	if err := validator.Validate().Struct(input); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errorcode.ErrBadParamInput.Error(), "data": err.Error()})
		return
	}

	user, err := u.userUsecase.Login(ctx.Context(), input.UsernameOrEmail, input.UsernameOrEmail, input.Password)
	if err != nil {
		ctx.Status(errorcode.StatusCode(err)).JSON(fiber.Map{"status": "error", "message": err.Error()})
		return
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID.String()
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(viper.GetString("APP_SECRET")))
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(fiber.Map{"status": "success", "data": fiber.Map{"token": t, "expires": claims["exp"]}})
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

func (u userHandler) Update(ctx *fiber.Ctx) {
	// Preparing uuid
	id, err := uuid.FromString(ctx.Params("id"))
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errorcode.ErrBadParamInput.Error()})
		return
	}

	// Binds input
	type Input struct {
		Email    string `json:"email" validate:"required,email,max=128"`
		Password string `json:"password" validate:"required,max=10"`
	}
	input := new(Input)
	if err := ctx.BodyParser(input); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errorcode.ErrBadParamInput.Error()})
		return
	}

	// Validate input
	if err := validator.Validate().Struct(input); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errorcode.ErrBadParamInput.Error(), "data": err.Error()})
		return
	}

	// Preparing user model
	hashedPassword, err := model.GeneratePassword(input.Password)
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		return
	}
	now := time.Now()
	userModel := model.User{
		Model: base.Model{
			ID:        id,
			UpdatedBy: &uuid.UUID{},
			UpdatedAt: &now,
		},
		Email:          input.Email,
		HashedPassword: hashedPassword,
	}

	// Updating data
	if err := u.userUsecase.Update(ctx.Context(), userModel); err != nil {
		ctx.Status(errorcode.StatusCode(err)).JSON(fiber.Map{"status": "error", "message": err.Error()})
		return
	}

	// Resolve new data
	data, err := u.userUsecase.GetByID(ctx.Context(), id)
	if err != nil {
		ctx.Status(errorcode.StatusCode(err)).JSON(fiber.Map{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(fiber.Map{"status": "success", "data": data})
}

func (u userHandler) Delete(ctx *fiber.Ctx) {
	// Preparing uuid
	id, err := uuid.FromString(ctx.Params("id"))
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": errorcode.ErrBadParamInput.Error()})
		return
	}

	// Delete user
	if err := u.userUsecase.Delete(ctx.Context(), id); err != nil {
		ctx.Status(errorcode.StatusCode(err)).JSON(fiber.Map{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(fiber.Map{"status": "success", "data": true})
}
