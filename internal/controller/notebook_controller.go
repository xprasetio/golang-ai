package controller

import (
	"golang-ai/internal/dto"
	"golang-ai/internal/pkg/serverutils"
	"golang-ai/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type INotebookController interface {
	RegisterRoutes(r fiber.Router)
	Create(ctx *fiber.Ctx) error
	Show(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
}

type notebookController struct {
	service service.INotebookService
}

func NewNotebookController(service service.INotebookService) INotebookController {
	return &notebookController{service: service}
}

func (c *notebookController) RegisterRoutes(r fiber.Router) {
	h := r.Group("/notebook/v1")
	h.Post("", c.Create)
	h.Get("/:id", c.Show)
	h.Put("/:id", c.Update)
}

func (c *notebookController) Create(ctx *fiber.Ctx) error {
	var req dto.CreateNotebookRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	err := serverutils.ValidateRequest(req)
	if err != nil {
		return err
	}

	res, err := c.service.Create(ctx.Context(), &req)
	if err != nil {
		return err
	}

	return ctx.JSON(serverutils.SuccessResponse("Success Create Notebook", res))
}
func (c *notebookController) Show(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return err
	}

	res, err := c.service.Show(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(serverutils.SuccessResponse("Success Show Notebook", res))
}
func (c *notebookController) Update(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")

	var req dto.UpdateNotebookRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		return err
	}

	req.Id = id
	res, err := c.service.Update(ctx.Context(), &req)
	if err != nil {
		return err
	}

	return ctx.JSON(serverutils.SuccessResponse("Success Update Notebook", res))
}
