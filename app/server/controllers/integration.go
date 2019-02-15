package controllers

import (
	"net/http"
	"strconv"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	db "github.com/horechek/teleport/app/database"
	"github.com/horechek/teleport/app/di"
	"github.com/horechek/teleport/pkg/params"
)

type IntegrationsController struct {
	services *di.Services
}

func NewIntegrationsController(services *di.Services) *IntegrationsController {
	return &IntegrationsController{
		services: services,
	}
}

func (c *IntegrationsController) List(ctx echo.Context) error {
	model := db.Integration{}
	condition := db.Condition{
		Sorting: db.Sorting{
			Sort:  params.GetString(ctx, "_sort", ""),
			Order: params.GetString(ctx, "_order", ""),
		},
	}

	count, err := model.Count(c.services.Database, condition)
	if err != nil {
		return err
	}

	condition.Pagination = db.Paginating{
		Start: params.GetInt(ctx, "_start", 0),
		End:   params.GetInt(ctx, "_end", 50),
	}

	list, err := model.FindAll(c.services.Database, condition)
	if err != nil {
		return err
	}

	ctx.Response().Header().Set("X-Total-Count", strconv.Itoa(count))
	return ctx.JSON(http.StatusOK, list)
}

func (c *IntegrationsController) Update(ctx echo.Context) error {
	model, err := findIntegration(c.services.Database, ctx)
	if err != nil {
		return err
	}

	if err := ctx.Bind(&model); err != nil {
		return err
	}

	model.UpdatedAt = time.Now().Unix()

	if _, err = valid.ValidateStruct(model); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	if err := c.services.Database.Save(&model).Error; err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, model)
}

func (c *IntegrationsController) Create(ctx echo.Context) error {
	model := db.Integration{}
	if err := ctx.Bind(&model); err != nil {
		c.services.Logger.Errorw("bind error", "err", err)
		return err
	}

	if _, err := valid.ValidateStruct(model); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	model.CreatedAt = c.services.Clock.Now().Unix()
	model.UpdatedAt = c.services.Clock.Now().Unix()

	if err := c.services.Database.Save(&model).Error; err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, model)
}

func (c *IntegrationsController) Remove(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return err
	}

	err = c.services.Database.Delete(&db.Integration{}, "id = ? ", id).Error
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, "removed")
}

func findIntegration(database *gorm.DB, c echo.Context) (*db.Integration, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return &db.Integration{}, err
	}

	condition := db.Condition{
		Params: []db.Param{
			{
				Field: "id",
				Value: id,
			},
		},
	}

	model := db.Integration{}
	return model.FindOne(database, condition)
}
