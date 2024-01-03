package servises

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"speed-typing-text-service/internal/adapters/dto"
	"speed-typing-text-service/internal/core/domain"
	"speed-typing-text-service/internal/core/ports/mocks"
	"speed-typing-text-service/pkg/logging/discard"
	"testing"
)

func TestGetRegularTexts(t *testing.T) {
	var log = discard.GetLogger()
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
	service := NewRegularTextService(repo, log)

	// expected
	var expectedResult []dto.GetRegularTextDto
	err := copier.Copy(&expectedResult, expectedTexts)
	require.NoError(t, err)
	// actual
	actualResult, err := service.GetRegularTexts()

	// checks
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, actualResult)
	repo.AssertExpectations(t)
}
