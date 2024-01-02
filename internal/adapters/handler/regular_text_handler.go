package handler

import (
	"github.com/gin-gonic/gin"
	"speed-typing-text-service/internal/core/ports"
	"speed-typing-text-service/pkg/logging"
)

type RegularTextHandler struct {
	svc ports.RegularTextService
	log logging.Logger
}

func NewRegularTextHandler(svc ports.RegularTextService, log logging.Logger) *RegularTextHandler {
	return &RegularTextHandler{
		svc: svc,
		log: log,
	}
}

func (h *RegularTextHandler) GetRegularTexts(c *gin.Context) {
	h.log.Debug("received get regular texts request")

	texts, err := h.svc.GetRegularTexts()
	if err != nil {
		err = c.Error(err)
		return
	}
	c.JSON(200, texts)
}
