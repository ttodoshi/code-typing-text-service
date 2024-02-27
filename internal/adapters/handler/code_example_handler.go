package handler

import (
	"github.com/gin-gonic/gin"
	"speed-typing-text-service/internal/adapters/dto"
	"speed-typing-text-service/internal/core/ports"
	"speed-typing-text-service/pkg/logging"
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

	UUID := c.Param("uuid")

	example, err := h.svc.GetCodeExampleByUUID(UUID)
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

	var examples []dto.GetCodeExampleDto
	var err error

	if programmingLanguageName == "" {
		examples, err = h.svc.GetCodeExamples()
	} else {
		examples, err = h.svc.GetCodeExamplesByProgrammingLanguageName(programmingLanguageName)
	}
	if err != nil {
		err = c.Error(err)
		return
	}
	c.JSON(200, examples)
}
