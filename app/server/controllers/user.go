package controllers

import (
	"net/http"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"github.com/horechek/teleport/app/database"
	"github.com/horechek/teleport/app/di"
)

type UserRequest struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Email   string `json:"email"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type UsersController struct {
	services *di.Services
}

func NewUsersController(services *di.Services) *UsersController {
	return &UsersController{
		services: services,
	}
}

// Authorization: Bearer AbCdEf12345sd6
func (c *UsersController) Login(ctx echo.Context) error {
	user := database.User{}
	if err := ctx.Bind(&user); err != nil {
		return err
	}

	if _, err := valid.ValidateStruct(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	err := c.services.Database.Model(user).
		Where("email = ?", user.Email).
		First(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if err == gorm.ErrRecordNotFound {
		return ctx.JSON(http.StatusBadRequest, UserResponse{Message: "Такого пользователя нет"})
	}

	if !user.CheckPassword(user.Password) {
		return ctx.JSON(http.StatusBadRequest, UserResponse{Message: "Неправильный пользователь или пароль"})
	}

	return ctx.JSON(http.StatusOK, UserResponse{
		Email: user.Email,
		Token: user.Token,
	})
}

func (c *UsersController) Register(ctx echo.Context) error {
	user := database.User{}
	err := ctx.Bind(&user)
	if err != nil {
		return err
	}

	if _, err := valid.ValidateStruct(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	var cnt int
	err = c.services.Database.Model(user).
		Where("email = ?", user.Email).
		Count(&cnt).Error

	if err != nil {
		return err
	}

	if cnt != 0 {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Такой пользователь уже есть",
		})
	}

	user.Token, err = user.GenerateToken()
	if err != nil {
		return err
	}

	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()

	if err := c.services.Database.Save(&user).Error; err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, UserResponse{
		Email: user.Email,
		Token: user.Token,
	})
}

type UserUpdateRequest struct {
	Password string `json:"password" valid:"required~'Пароль' обязательное поле"`
}

func (c *UsersController) Update(ctx echo.Context) error {
	req := UserUpdateRequest{}
	err := ctx.Bind(&req)
	if err != nil {
		return err
	}

	if _, err := valid.ValidateStruct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	id := ctx.Get("user").(int)
	user := database.User{}
	err = c.services.Database.Model(user).
		Where("id = ?", id).
		First(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	user.Password = req.Password
	user.UpdatedAt = time.Now().Unix()

	if err := c.services.Database.Save(&user).Error; err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "Пароль обновлен",
	})
}
