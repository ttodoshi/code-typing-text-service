package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ttodoshi/code-typing-text-service/internal/core/ports"
	"github.com/ttodoshi/code-typing-text-service/internal/core/ports/dto"
	"github.com/ttodoshi/code-typing-text-service/pkg/logging"
)

type CodeExampleHandler struct {
	svc ports.CodeExampleService
	log logging.Logger
}

func NewCodeExampleHandler(svc ports.CodeExampleService, log logging.Logger) *CodeExampleHandler {
	return &CodeExampleHandler{
		svc: svc,
		log: log,
	}
}

// GetProgrammingLanguages godoc
//
//	@Summary		Get programming languages
//	@Description	Get all programming languages
//	@Tags			code examples
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	dto.GetProgrammingLanguageDto
//	@Router			/texts/programming-languages [get]
func (h *CodeExampleHandler) GetProgrammingLanguages(c *gin.Context) {
	h.log.Debug("received get programming languages request")

	languages, err := h.svc.GetProgrammingLanguages()
	if err != nil {
		err = c.Error(err)
		return
	}
	c.JSON(200, languages)
}

// GetCodeExampleByUUID godoc
//
//	@Summary		Get code example by UUID
//	@Description	Get code example by UUID
//	@Tags			code examples
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string	true	"Code example UUID"
//	@Success		200		{object}	dto.GetCodeExampleDto
//	@Failure		404		{object}	dto.ErrorResponseDto
//	@Router			/texts/code-examples/{uuid} [get]
func (h *CodeExampleHandler) GetCodeExampleByUUID(c *gin.Context) {
	h.log.Debug("received get code example by UUID request")

	example, err := h.svc.GetCodeExampleByUUID(
		c.GetString("userID"),
		c.Param("uuid"),
	)
	if err != nil {
		err = c.Error(err)
		return
	}
	c.JSON(200, example)
}

// GetCodeExamples godoc
//
//	@Summary		Get code examples
//	@Description	Get code examples by programming language name
//	@Tags			code examples
//	@Accept			json
//	@Produce		json
//	@Param			programming-language-name	query		string	false	"Programming language name"
//	@Success		200							{array}		dto.GetCodeExampleDto
//	@Failure		404							{object}	dto.ErrorResponseDto
//	@Router			/texts/code-examples [get]
func (h *CodeExampleHandler) GetCodeExamples(c *gin.Context) {
	h.log.Debug("received get code examples request")

	programmingLanguageName := c.Query("programming-language-name")
	userID := c.GetString("userID")

	var examples []dto.GetCodeExampleDto
	var err error

	if programmingLanguageName == "" {
		examples, err = h.svc.GetCodeExamples(userID)
	} else {
		examples, err = h.svc.GetCodeExamplesByProgrammingLanguageName(
			userID, programmingLanguageName,
		)
	}
	if err != nil {
		err = c.Error(err)
		return
	}
	c.JSON(200, examples)
}

// CreateCodeExample godoc
//
//	@Summary		Create code example
//	@Description	Create code example
//	@Tags			code examples
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateCodeExampleDto	true	"create code example request"
//	@Success		201		{object}	string
//	@Failure		404		{object}	dto.ErrorResponseDto
//	@Router			/texts/code-examples [post]
func (h *CodeExampleHandler) CreateCodeExample(c *gin.Context) {
	h.log.Debug("received create code example request")

	userID := c.GetString("userID")
	var err error
	if userID == "" {
		err = c.Error(
			fmt.Errorf("user not authenticated: %w", ports.UnauthorizedError),
		)
		return
	}
	var createCodeExampleDto dto.CreateCodeExampleDto
	if err = c.ShouldBindJSON(&createCodeExampleDto); err != nil {
		h.log.Warn("error in request body")
		err = c.Error(
			fmt.Errorf("error in request body: %w", ports.BadRequestError),
		)
		return
	}

	UUID, err := h.svc.CreateCodeExample(userID, createCodeExampleDto)
	if err != nil {
		err = c.Error(err)
		return
	}

	c.Data(201, "text/html; charset=utf-8", []byte(UUID))
}

// DeleteCodeExample godoc
//
//	@Summary		Delete code example
//	@Description	Delete code example
//	@Tags			code examples
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path	string	true	"Code example UUID"
//	@Success		204
//	@Failure		404	{object}	dto.ErrorResponseDto
//	@Router			/texts/code-examples/{uuid} [delete]
func (h *CodeExampleHandler) DeleteCodeExample(c *gin.Context) {
	h.log.Debug("received delete code example request")

	userID := c.GetString("userID")
	var err error
	if userID == "" {
		err = c.Error(
			fmt.Errorf("user not authenticated: %w", ports.UnauthorizedError),
		)
		return
	}

	err = h.svc.DeleteCodeExample(userID, c.Param("uuid"))
	if err != nil {
		err = c.Error(err)
		return
	}

	c.Status(204)
}
