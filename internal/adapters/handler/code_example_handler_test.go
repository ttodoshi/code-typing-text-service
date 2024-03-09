package handler

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"speed-typing-text-service/internal/core/domain"
	"speed-typing-text-service/internal/core/errors"
	"speed-typing-text-service/internal/core/ports/mocks"
	"speed-typing-text-service/internal/core/servises"
	"speed-typing-text-service/pkg/logging/nop"
	"testing"
)

func TestCodeExampleHandler_GetProgrammingLanguages(t *testing.T) {
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
	service := servises.NewCodeExampleService(repo, log)

	handler := NewCodeExampleHandler(service, log)

	// gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/texts/programming-languages",
		nil,
	)
	c.Request = request

	handler.GetProgrammingLanguages(c)

	// checks
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, c.Errors)
	repo.AssertExpectations(t)
}

func TestCodeExampleHandler_GetCodeExampleByUUID(t *testing.T) {
	var log = nop.GetLogger()
	// repo mock
	repo := new(mocks.CodeExampleRepository)
	programmingLanguage := domain.ProgrammingLanguage{
		UUID: gofakeit.UUID(),
		Name: gofakeit.ProgrammingLanguage(),
		Logo: gofakeit.ImageURL(200, 200),
	}
	codeExample := domain.CodeExample{
		UUID:                    gofakeit.UUID(),
		Content:                 gofakeit.LoremIpsumSentence(gofakeit.IntRange(1, 10)),
		ProgrammingLanguageUUID: programmingLanguage.UUID,
		ProgrammingLanguage:     programmingLanguage,
	}
	repo.
		On("GetCodeExampleByUUID", codeExample.UUID).
		Return(codeExample, nil)
	repo.
		On(
			"GetCodeExampleByUUID",
			mock.AnythingOfType("string"),
		).
		Return(domain.CodeExample{}, &errors.NotFoundError{})

	// service
	service := servises.NewCodeExampleService(repo, log)

	handler := NewCodeExampleHandler(service, log)

	// gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Run("successful retrieval code example by UUID", func(t *testing.T) {
		request := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/texts/code-examples/%s", codeExample.UUID),
			nil,
		)
		c.AddParam("uuid", codeExample.UUID)
		c.Request = request

		handler.GetCodeExampleByUUID(c)

		// check
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Nil(t, c.Errors)
	})

	t.Run("unsuccessful retrieval code example by UUID", func(t *testing.T) {
		request := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/texts/code-examples/%s", gofakeit.UUID()),
			nil,
		)
		c.Request = request

		handler.GetCodeExampleByUUID(c)

		// check
		assert.Error(t, c.Errors.Last())
	})
	repo.AssertExpectations(t)
}

func TestCodeExampleHandler_GetCodeExamples(t *testing.T) {
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
		On("GetCodeExamples").
		Return(expectedCodeExamples)

	// service
	service := servises.NewCodeExampleService(repo, log)

	handler := NewCodeExampleHandler(service, log)

	// gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/texts/code-examples",
		nil,
	)
	c.Request = request

	handler.GetCodeExamples(c)

	// checks
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, c.Errors)
	repo.AssertExpectations(t)
}

func TestCodeExampleHandler_GetCodeExamplesByProgrammingLanguageName(t *testing.T) {
	var log = nop.GetLogger()
	// repo mock
	repo := new(mocks.CodeExampleRepository)
	programmingLanguage := domain.ProgrammingLanguage{
		UUID: gofakeit.UUID(),
		Name: gofakeit.ProgrammingLanguage(),
		Logo: gofakeit.ImageURL(200, 200),
	}
	var expectedCodeExamples []domain.CodeExample
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
	repo.
		On("GetCodeExamplesByProgrammingLanguageName", programmingLanguage.Name).
		Return(expectedCodeExamples, nil)
	repo.
		On(
			"GetCodeExamplesByProgrammingLanguageName",
			mock.AnythingOfType("string"),
		).Return(
		nil,
		&errors.NotFoundError{},
	)

	// service
	service := servises.NewCodeExampleService(repo, log)

	handler := NewCodeExampleHandler(service, log)

	// gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Run("successful retrieval by programming language name", func(t *testing.T) {
		request := httptest.NewRequest(
			http.MethodGet,
			"/api/v1/texts/code-examples",
			nil,
		)
		c.Request = request
		c.Request.URL.RawQuery = fmt.Sprintf("programming-language-name=%s", programmingLanguage.Name)

		handler.GetCodeExamples(c)

		// check
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Nil(t, c.Errors)
	})

	t.Run("unsuccessful retrieval with non-existent programming language name", func(t *testing.T) {
		request := httptest.NewRequest(
			http.MethodGet,
			"/api/v1/texts/code-examples",
			nil,
		)
		c.Request = request
		c.Request.URL.RawQuery = fmt.Sprintf("programming-language-name=%s", gofakeit.LoremIpsumWord())

		handler.GetCodeExamples(c)

		// check
		assert.Error(t, c.Errors.Last())
	})
	repo.AssertExpectations(t)
}
