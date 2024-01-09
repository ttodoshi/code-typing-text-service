package servises

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"speed-typing-text-service/internal/adapters/dto"
	"speed-typing-text-service/internal/core/domain"
	"speed-typing-text-service/internal/core/errors"
	"speed-typing-text-service/internal/core/ports/mocks"
	"speed-typing-text-service/pkg/logging/discard"
	"testing"
)

func TestGetProgrammingLanguages(t *testing.T) {
	var log = discard.GetLogger()
	// repo mock
	repo := new(mocks.CodeExampleRepository)
	var expectedProgrammingLanguages []domain.ProgrammingLanguage
	for i := 0; i < gofakeit.IntRange(5, 10); i++ {
		expectedProgrammingLanguages = append(
			expectedProgrammingLanguages,
			domain.ProgrammingLanguage{
				UUID: gofakeit.UUID(),
				Name: gofakeit.ProgrammingLanguage(),
			},
		)
	}
	repo.
		On("GetProgrammingLanguages").
		Return(expectedProgrammingLanguages)

	// service
	service := NewCodeExampleService(repo, log)

	// expected
	var expectedResult []dto.GetProgrammingLanguageDto
	err := copier.Copy(&expectedResult, expectedProgrammingLanguages)
	require.NoError(t, err)
	// actual
	actualResult, err := service.GetProgrammingLanguages()

	// checks
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, actualResult)
	repo.AssertExpectations(t)
}

func TestGetCodeExampleByUUID(t *testing.T) {
	var log = discard.GetLogger()
	// repo mock
	repo := new(mocks.CodeExampleRepository)
	programmingLanguage := domain.ProgrammingLanguage{
		UUID: gofakeit.UUID(),
		Name: gofakeit.ProgrammingLanguage(),
	}
	expectedCodeExamples := map[string]domain.CodeExample{}
	for i := 0; i < gofakeit.IntRange(1, 3); i++ {
		uuid := gofakeit.UUID()
		expectedCodeExamples[uuid] = domain.CodeExample{
			UUID:                    uuid,
			Content:                 gofakeit.LoremIpsumSentence(gofakeit.IntRange(1, 10)),
			ProgrammingLanguageUUID: programmingLanguage.UUID,
			ProgrammingLanguage:     programmingLanguage,
		}
	}
	for uuid, codeExample := range expectedCodeExamples {
		repo.
			On("GetCodeExampleByUUID", uuid).
			Return(codeExample, nil)
	}
	repo.
		On(
			"GetCodeExampleByUUID",
			mock.AnythingOfType("string"),
		).Return(
		domain.CodeExample{},
		&errors.NotFoundError{},
	)

	// service
	service := NewCodeExampleService(repo, log)

	t.Run("successful retrieval code example by UUID", func(t *testing.T) {
		for uuid, codeExample := range expectedCodeExamples {
			// expected
			var expectedResult dto.GetCodeExampleDto
			err := copier.Copy(&expectedResult, &codeExample)
			require.NoError(t, err)
			// actual
			actualResult, err := service.GetCodeExampleByUUID(uuid)

			// checks
			assert.NoError(t, err)
			assert.Equal(t, expectedResult, actualResult)
		}
		repo.AssertExpectations(t)
	})

	t.Run("unsuccessful retrieval code example by UUID", func(t *testing.T) {
		_, err := service.GetCodeExampleByUUID(gofakeit.UUID())

		// checks
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestGetCodeExamples(t *testing.T) {
	var log = discard.GetLogger()
	// repo mock
	repo := new(mocks.CodeExampleRepository)
	var programmingLanguages []domain.ProgrammingLanguage
	for i := 0; i < gofakeit.IntRange(5, 10); i++ {
		programmingLanguages = append(
			programmingLanguages,
			domain.ProgrammingLanguage{
				UUID: gofakeit.UUID(),
				Name: gofakeit.ProgrammingLanguage(),
			},
		)
	}
	var expectedCodeExamples []domain.CodeExample
	for _, programmingLanguage := range programmingLanguages {
		for i := 0; i < gofakeit.IntRange(1, 3); i++ {
			expectedCodeExamples = append(
				expectedCodeExamples,
				domain.CodeExample{
					UUID:                    gofakeit.UUID(),
					Content:                 gofakeit.LoremIpsumSentence(gofakeit.IntRange(1, 10)),
					ProgrammingLanguageUUID: programmingLanguage.UUID,
					ProgrammingLanguage:     programmingLanguage,
				},
			)
		}
	}
	repo.
		On("GetCodeExamples").
		Return(expectedCodeExamples)

	// service
	service := NewCodeExampleService(repo, log)

	// expected
	var expectedResult []dto.GetCodeExampleDto
	err := copier.Copy(&expectedResult, &expectedCodeExamples)
	require.NoError(t, err)
	// actual
	actualResult, err := service.GetCodeExamples()

	// checks
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, actualResult)
	repo.AssertExpectations(t)
}

func TestGetCodeExamplesByProgrammingLanguageName(t *testing.T) {
	var log = discard.GetLogger()
	// repo mock
	repo := new(mocks.CodeExampleRepository)
	programmingLanguages := map[domain.ProgrammingLanguage][]domain.CodeExample{}
	for i := 0; i < gofakeit.IntRange(1, 3); i++ {
		programmingLanguages[domain.ProgrammingLanguage{
			UUID: gofakeit.UUID(),
			Name: gofakeit.ProgrammingLanguage(),
		}] = []domain.CodeExample{}
	}
	for programmingLanguage := range programmingLanguages {
		for i := 0; i < gofakeit.IntRange(1, 3); i++ {
			programmingLanguages[programmingLanguage] = append(
				programmingLanguages[programmingLanguage],
				domain.CodeExample{
					UUID:                    gofakeit.UUID(),
					Content:                 gofakeit.LoremIpsumSentence(gofakeit.IntRange(1, 10)),
					ProgrammingLanguageUUID: programmingLanguage.UUID,
					ProgrammingLanguage:     programmingLanguage,
				},
			)
		}
	}
	for programmingLanguage, codeExamples := range programmingLanguages {
		repo.
			On("GetCodeExamplesByProgrammingLanguageName", programmingLanguage.Name).
			Return(codeExamples, nil)
	}
	repo.
		On(
			"GetCodeExamplesByProgrammingLanguageName",
			mock.AnythingOfType("string"),
		).Return(
		nil,
		&errors.NotFoundError{},
	)

	// service
	service := NewCodeExampleService(repo, log)

	t.Run("successful retrieval by programming language name", func(t *testing.T) {
		for programmingLanguage, expectedCodeExamples := range programmingLanguages {
			// expected
			var expectedResult []dto.GetCodeExampleDto
			err := copier.Copy(&expectedResult, &expectedCodeExamples)
			require.NoError(t, err)
			// actual
			actualResult, err := service.GetCodeExamplesByProgrammingLanguageName(programmingLanguage.Name)

			// checks
			assert.NoError(t, err)
			assert.Equal(t, expectedResult, actualResult)
		}
		repo.AssertExpectations(t)
	})

	t.Run("unsuccessful retrieval with non-existent programming language name", func(t *testing.T) {
		_, err := service.GetCodeExamplesByProgrammingLanguageName(gofakeit.LoremIpsumWord())

		// checks
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
