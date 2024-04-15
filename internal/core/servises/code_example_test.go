package servises

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ttodoshi/code-typing-text-service/internal/core/domain"
	"github.com/ttodoshi/code-typing-text-service/internal/core/ports/dto"
	"github.com/ttodoshi/code-typing-text-service/internal/core/ports/mocks"
	"github.com/ttodoshi/code-typing-text-service/pkg/logging/nop"
	"testing"
)

func TestGetProgrammingLanguages(t *testing.T) {
	var log = nop.GetLogger()
	// repo mock
	repo := new(mocks.CodeExampleRepository)
	var expectedProgrammingLanguages []domain.ProgrammingLanguage
	for i := 0; i < gofakeit.IntRange(5, 10); i++ {
		expectedProgrammingLanguages = append(
			expectedProgrammingLanguages,
			domain.ProgrammingLanguage{
				UUID: gofakeit.UUID(),
				Name: gofakeit.ProgrammingLanguage(),
				Logo: gofakeit.ImageURL(200, 200),
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

func TestGetDefaultCodeExampleByUUID(t *testing.T) {
	var log = nop.GetLogger()
	// repo mock
	repo := new(mocks.CodeExampleRepository)
	programmingLanguage := domain.ProgrammingLanguage{
		UUID: gofakeit.UUID(),
		Name: gofakeit.ProgrammingLanguage(),
		Logo: gofakeit.ImageURL(200, 200),
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
		fmt.Errorf(""),
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
			actualResult, err := service.GetCodeExampleByUUID("", uuid)

			// checks
			assert.NoError(t, err)
			assert.Equal(t, expectedResult, actualResult)
		}
	})

	t.Run("unsuccessful retrieval code example by UUID", func(t *testing.T) {
		_, err := service.GetCodeExampleByUUID("", gofakeit.UUID())

		// checks
		assert.Error(t, err)
	})
	repo.AssertExpectations(t)
}

func TestGetCodeExampleByUUIDAuthorized(t *testing.T) {
	var log = nop.GetLogger()
	// repo mock
	repo := new(mocks.CodeExampleRepository)
	programmingLanguage := domain.ProgrammingLanguage{
		UUID: gofakeit.UUID(),
		Name: gofakeit.ProgrammingLanguage(),
		Logo: gofakeit.ImageURL(200, 200),
	}
	var expectedCodeExamples []domain.CodeExample
	user1ID := gofakeit.UUID()
	expectedCodeExamples = append(expectedCodeExamples, domain.CodeExample{
		UUID:                    gofakeit.UUID(),
		Content:                 gofakeit.LoremIpsumSentence(gofakeit.IntRange(1, 10)),
		ProgrammingLanguageUUID: programmingLanguage.UUID,
		ProgrammingLanguage:     programmingLanguage,
		UserID:                  user1ID,
	})
	expectedCodeExamples = append(expectedCodeExamples, domain.CodeExample{
		UUID:                    gofakeit.UUID(),
		Content:                 gofakeit.LoremIpsumSentence(gofakeit.IntRange(1, 10)),
		ProgrammingLanguageUUID: programmingLanguage.UUID,
		ProgrammingLanguage:     programmingLanguage,
		UserID:                  gofakeit.UUID(),
	})

	for _, codeExample := range expectedCodeExamples {
		repo.
			On("GetCodeExampleByUUID", codeExample.UUID).
			Return(codeExample, nil)
	}

	// service
	service := NewCodeExampleService(repo, log)

	t.Run("successful retrieval code example by UUID", func(t *testing.T) {
		// expected
		var expectedResult dto.GetCodeExampleDto
		err := copier.Copy(&expectedResult, &expectedCodeExamples[0])
		require.NoError(t, err)
		// actual
		actualResult, err := service.GetCodeExampleByUUID(user1ID, expectedResult.UUID)

		// checks
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, actualResult)
	})

	t.Run("unsuccessful retrieval code example by UUID", func(t *testing.T) {
		_, err := service.GetCodeExampleByUUID(user1ID, expectedCodeExamples[1].UUID)

		// checks
		assert.Error(t, err)
	})
	repo.AssertExpectations(t)
}

func TestGetCodeExamplesUnauthorized(t *testing.T) {
	var log = nop.GetLogger()
	// repo mock
	repo := new(mocks.CodeExampleRepository)
	var programmingLanguages []domain.ProgrammingLanguage
	for i := 0; i < gofakeit.IntRange(5, 10); i++ {
		programmingLanguages = append(
			programmingLanguages,
			domain.ProgrammingLanguage{
				UUID: gofakeit.UUID(),
				Name: gofakeit.ProgrammingLanguage(),
				Logo: gofakeit.ImageURL(200, 200),
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
		On("GetCodeExamples", "").
		Return(expectedCodeExamples)

	// service
	service := NewCodeExampleService(repo, log)

	// expected
	var expectedResult []dto.GetCodeExampleDto
	err := copier.Copy(&expectedResult, &expectedCodeExamples)
	require.NoError(t, err)
	// actual
	actualResult, err := service.GetCodeExamples("")

	// checks
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, actualResult)
	repo.AssertExpectations(t)
}

func TestGetCodeExamplesAuthorized(t *testing.T) {
	var log = nop.GetLogger()
	// repo mock
	repo := new(mocks.CodeExampleRepository)

	programmingLanguage := domain.ProgrammingLanguage{
		UUID: gofakeit.UUID(),
		Name: gofakeit.ProgrammingLanguage(),
		Logo: gofakeit.ImageURL(200, 200),
	}

	standardCodeExamples := []domain.CodeExample{{
		UUID:                    gofakeit.UUID(),
		Content:                 "standard code example",
		ProgrammingLanguageUUID: programmingLanguage.UUID,
		ProgrammingLanguage:     programmingLanguage,
	}}
	user1ID := "user1"
	user1codeExample := domain.CodeExample{
		UUID:                    gofakeit.UUID(),
		Content:                 "user1 code example",
		ProgrammingLanguageUUID: programmingLanguage.UUID,
		ProgrammingLanguage:     programmingLanguage,
		UserID:                  user1ID,
	}

	repo.
		On("GetCodeExamples", user1ID).
		Return(append(standardCodeExamples, user1codeExample))

	// service
	service := NewCodeExampleService(repo, log)

	// expected
	var expectedResult []dto.GetCodeExampleDto
	err := copier.Copy(&expectedResult, append(standardCodeExamples, user1codeExample))

	// actual
	actualResult, err := service.GetCodeExamples(user1ID)

	// checks
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, actualResult)
	assert.Len(t, actualResult, 2)

	repo.AssertExpectations(t)
}

func TestGetCodeExamplesByProgrammingLanguageNameUnauthorized(t *testing.T) {
	var log = nop.GetLogger()
	// repo mock
	repo := new(mocks.CodeExampleRepository)
	programmingLanguages := map[domain.ProgrammingLanguage][]domain.CodeExample{}
	for i := 0; i < gofakeit.IntRange(1, 3); i++ {
		programmingLanguages[domain.ProgrammingLanguage{
			UUID: gofakeit.UUID(),
			Name: gofakeit.ProgrammingLanguage(),
			Logo: gofakeit.ImageURL(200, 200),
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
			On("GetCodeExamplesByProgrammingLanguageName", "", programmingLanguage.Name).
			Return(codeExamples, nil)
	}
	repo.
		On(
			"GetCodeExamplesByProgrammingLanguageName",
			"",
			mock.AnythingOfType("string"),
		).Return(
		nil,
		fmt.Errorf(""),
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
			actualResult, err := service.GetCodeExamplesByProgrammingLanguageName("", programmingLanguage.Name)

			// checks
			assert.NoError(t, err)
			assert.Equal(t, expectedResult, actualResult)
		}
	})

	t.Run("unsuccessful retrieval with non-existent programming language name", func(t *testing.T) {
		_, err := service.GetCodeExamplesByProgrammingLanguageName("", gofakeit.LoremIpsumWord())

		// checks
		assert.Error(t, err)
	})
	repo.AssertExpectations(t)
}

func TestGetCodeExamplesByProgrammingLanguageNameAuthorized(t *testing.T) {
	var log = nop.GetLogger()
	// repo mock
	repo := new(mocks.CodeExampleRepository)

	programmingLanguage := domain.ProgrammingLanguage{
		UUID: gofakeit.UUID(),
		Name: gofakeit.ProgrammingLanguage(),
		Logo: gofakeit.ImageURL(200, 200),
	}
	defaultCodeExamples := []domain.CodeExample{{
		UUID:                    gofakeit.UUID(),
		Content:                 gofakeit.LoremIpsumSentence(gofakeit.IntRange(1, 10)),
		ProgrammingLanguageUUID: programmingLanguage.UUID,
		ProgrammingLanguage:     programmingLanguage,
	}}
	user1ID := gofakeit.UUID()
	user1CodeExample := domain.CodeExample{
		UUID:                    gofakeit.UUID(),
		Content:                 gofakeit.LoremIpsumSentence(gofakeit.IntRange(1, 10)),
		ProgrammingLanguageUUID: programmingLanguage.UUID,
		ProgrammingLanguage:     programmingLanguage,
		UserID:                  user1ID,
	}
	user2ID := gofakeit.UUID()
	user2CodeExample := domain.CodeExample{
		UUID:                    gofakeit.UUID(),
		Content:                 gofakeit.LoremIpsumSentence(gofakeit.IntRange(1, 10)),
		ProgrammingLanguageUUID: programmingLanguage.UUID,
		ProgrammingLanguage:     programmingLanguage,
		UserID:                  user2ID,
	}

	repo.
		On("GetCodeExamplesByProgrammingLanguageName", user1ID, programmingLanguage.Name).
		Return(append(defaultCodeExamples, user1CodeExample), nil)
	repo.
		On("GetCodeExamplesByProgrammingLanguageName", user2ID, programmingLanguage.Name).
		Maybe().
		Return(append(defaultCodeExamples, user2CodeExample), nil)

	// service
	service := NewCodeExampleService(repo, log)

	// expected
	var expectedResult []dto.GetCodeExampleDto
	err := copier.Copy(&expectedResult, append(defaultCodeExamples, user1CodeExample))
	require.NoError(t, err)
	// actual
	actualResult, err := service.GetCodeExamplesByProgrammingLanguageName(user1ID, programmingLanguage.Name)

	// checks
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, actualResult)
	assert.Len(t, expectedResult, 2)
	repo.AssertExpectations(t)
}

func TestCreateCodeExample(t *testing.T) {
	var log = nop.GetLogger()
	// repo mock
	repo := new(mocks.CodeExampleRepository)

	repo.
		On("CreateCodeExample", mock.Anything).
		Return(gofakeit.UUID(), nil)

	// service
	service := NewCodeExampleService(repo, log)

	// actual
	_, err := service.CreateCodeExample(gofakeit.UUID(), dto.CreateCodeExampleDto{
		Content:                 gofakeit.LoremIpsumSentence(gofakeit.IntRange(1, 10)),
		ProgrammingLanguageUUID: gofakeit.UUID(),
	})

	// checks
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteCodeExample(t *testing.T) {
	var log = nop.GetLogger()
	// repo mock
	repo := new(mocks.CodeExampleRepository)

	programmingLanguage := domain.ProgrammingLanguage{
		UUID: gofakeit.UUID(),
		Name: gofakeit.ProgrammingLanguage(),
		Logo: gofakeit.ImageURL(200, 200),
	}
	codeExamples := []domain.CodeExample{{
		UUID:                    gofakeit.UUID(),
		Content:                 gofakeit.LoremIpsumSentence(gofakeit.IntRange(1, 10)),
		ProgrammingLanguageUUID: programmingLanguage.UUID,
		ProgrammingLanguage:     programmingLanguage,
	}}
	user1ID := gofakeit.UUID()
	codeExamples = append(codeExamples, domain.CodeExample{
		UUID:                    gofakeit.UUID(),
		Content:                 gofakeit.LoremIpsumSentence(gofakeit.IntRange(1, 10)),
		ProgrammingLanguageUUID: programmingLanguage.UUID,
		ProgrammingLanguage:     programmingLanguage,
		UserID:                  user1ID,
	})
	user2ID := gofakeit.UUID()
	codeExamples = append(codeExamples, domain.CodeExample{
		UUID:                    gofakeit.UUID(),
		Content:                 gofakeit.LoremIpsumSentence(gofakeit.IntRange(1, 10)),
		ProgrammingLanguageUUID: programmingLanguage.UUID,
		ProgrammingLanguage:     programmingLanguage,
		UserID:                  user2ID,
	})

	for _, example := range codeExamples {
		repo.
			On("GetCodeExampleByUUID", example.UUID).
			Return(example, nil)
		repo.
			On("DeleteCodeExample", example.UUID).
			Maybe().
			Return(nil)
	}
	repo.
		On(
			"GetCodeExampleByUUID",
			mock.AnythingOfType("string"),
		).Return(
		domain.CodeExample{},
		fmt.Errorf(""),
	)

	// service
	service := NewCodeExampleService(repo, log)

	t.Run("successful deletion", func(t *testing.T) {
		err := service.DeleteCodeExample(user1ID, codeExamples[1].UUID)

		// checks
		assert.NoError(t, err)
	})

	t.Run("unsuccessful deletion due to code example not found", func(t *testing.T) {
		err := service.DeleteCodeExample(user1ID, gofakeit.UUID())

		// checks
		assert.Error(t, err)
	})

	t.Run("unsuccessful deletion standard example", func(t *testing.T) {
		err := service.DeleteCodeExample(user1ID, codeExamples[0].UUID)

		// checks
		assert.Error(t, err)
	})

	t.Run("unsuccessful deletion other user's example", func(t *testing.T) {
		err := service.DeleteCodeExample(user1ID, codeExamples[2].UUID)

		// checks
		assert.Error(t, err)
	})

	repo.AssertExpectations(t)
}
