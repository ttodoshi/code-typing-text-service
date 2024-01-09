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

func (h *CodeExampleHandler) GetProgrammingLanguages(c *gin.Context) {
	h.log.Debug("received get programming languages request")

	languages, err := h.svc.GetProgrammingLanguages()
	if err != nil {
		err = c.Error(err)
		return
	}
	c.JSON(200, languages)
}

func (h *CodeExampleHandler) GetCodeExampleByUUID(c *gin.Context) {
	UUID := c.Param("uuid")

	var example dto.GetCodeExampleDto
	var err error

	h.log.Debug("received get code example by UUID request")
	example, err = h.svc.GetCodeExampleByUUID(UUID)
	if err != nil {
		err = c.Error(err)
		return
	}
	c.JSON(200, example)
}

func (h *CodeExampleHandler) GetCodeExamples(c *gin.Context) {
	programmingLanguageName := c.Query("programming-language-name")

	var examples []dto.GetCodeExampleDto
	var err error

	h.log.Debug("received get code examples request")
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
