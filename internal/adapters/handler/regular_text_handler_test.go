package handler

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"speed-typing-text-service/internal/core/domain"
	"speed-typing-text-service/internal/core/ports/mocks"
	"speed-typing-text-service/internal/core/servises"
	"speed-typing-text-service/pkg/logging/nop"
	"testing"
)

func TestRegularTextHandler_GetRegularTexts(t *testing.T) {
	var log = nop.GetLogger()
	// repo mock
	repo := new(mocks.RegularTextRepository)
	var expectedTexts []domain.RegularText
	for i := 0; i < gofakeit.IntRange(5, 10); i++ {
		expectedTexts = append(
			expectedTexts,
			domain.RegularText{
				UUID:    gofakeit.UUID(),
				Content: gofakeit.LoremIpsumSentence(gofakeit.IntRange(1, 10)),
			},
		)
	}
	repo.
		On("GetRegularTexts").
		Return(expectedTexts)

	// service
	service := servises.NewRegularTextService(repo, log)

	handler := NewRegularTextHandler(service, log)

	// gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/texts",
		nil,
	)
	c.Request = request

	handler.GetRegularTexts(c)

	// checks
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, c.Errors)
	repo.AssertExpectations(t)
}
