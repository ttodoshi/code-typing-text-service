package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ttodoshi/code-typing-text-service/internal/core/domain"
	"github.com/ttodoshi/code-typing-text-service/internal/core/ports/dto"
	"github.com/ttodoshi/code-typing-text-service/internal/core/ports/mocks"
	"github.com/ttodoshi/code-typing-text-service/internal/core/servises"
	"github.com/ttodoshi/code-typing-text-service/pkg/logging/nop"
	"net/http"
	"net/http/httptest"
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

func TestCodeExampleHandler_GetDefaultCodeExampleByUUID(t *testing.T) {
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
		Return(domain.CodeExample{}, fmt.Errorf(""))

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

func TestCodeExampleHandler_GetCodeExampleByUUIDAuthorized(t *testing.T) {
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
			Maybe().
			Return(codeExample, nil)
	}

	// service
	service := servises.NewCodeExampleService(repo, log)

	handler := NewCodeExampleHandler(service, log)

	// gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Run("successful retrieval code example by UUID", func(t *testing.T) {
		request := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/texts/code-examples/%s", expectedCodeExamples[0].UUID),
			nil,
		)
		c.AddParam("uuid", expectedCodeExamples[0].UUID)
		c.Set("userID", user1ID)
		c.Request = request

		handler.GetCodeExampleByUUID(c)

		// check
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Nil(t, c.Errors)
	})

	t.Run("unsuccessful retrieval code example by UUID", func(t *testing.T) {
		request := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/texts/code-examples/%s", expectedCodeExamples[1].UUID),
			nil,
		)
		c.AddParam("uuid", expectedCodeExamples[1].UUID)
		c.Set("userID", user1ID)
		c.Request = request

		handler.GetCodeExampleByUUID(c)

		// check
		assert.Error(t, c.Errors.Last())
	})
	repo.AssertExpectations(t)
}

func TestCodeExampleHandler_GetCodeExamplesUnauthorized(t *testing.T) {
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

func TestCodeExampleHandler_GetCodeExamplesAuthorized(t *testing.T) {
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
	c.Set("userID", user1ID)
	c.Request = request

	handler.GetCodeExamples(c)

	// checks
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, c.Errors)

	var actualCodeExamples []interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &actualCodeExamples)
	assert.Len(t, actualCodeExamples, 2)

	repo.AssertExpectations(t)
}

func TestCodeExampleHandler_GetCodeExamplesByProgrammingLanguageNameUnauthorized(t *testing.T) {
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
		On("GetCodeExamplesByProgrammingLanguageName", "", programmingLanguage.Name).
		Return(expectedCodeExamples, nil)
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

func TestCodeExampleHandler_GetCodeExamplesByProgrammingLanguageNameAuthorized(t *testing.T) {
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

	repo.
		On("GetCodeExamplesByProgrammingLanguageName", user1ID, programmingLanguage.Name).
		Return(append(defaultCodeExamples, user1CodeExample), nil)

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
	c.Set("userID", user1ID)
	c.Request = request
	c.Request.URL.RawQuery = fmt.Sprintf("programming-language-name=%s", programmingLanguage.Name)

	handler.GetCodeExamples(c)

	// check
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, c.Errors)

	var actualCodeExamples []interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &actualCodeExamples)
	assert.Len(t, actualCodeExamples, 2)

	repo.AssertExpectations(t)
}

func TestCodeExampleHandler_CreateCodeExample(t *testing.T) {
	var log = nop.GetLogger()
	// repo mock
	repo := new(mocks.CodeExampleRepository)

	repo.
		On("CreateCodeExample", mock.Anything).
		Return(gofakeit.UUID(), nil)

	// service
	service := servises.NewCodeExampleService(repo, log)

	handler := NewCodeExampleHandler(service, log)

	// gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	createDto := dto.CreateCodeExampleDto{
		Content:                 gofakeit.LoremIpsumSentence(gofakeit.IntRange(1, 10)),
		ProgrammingLanguageUUID: gofakeit.UUID(),
	}

	marshal, _ := json.Marshal(createDto)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/texts/code-examples",
		bytes.NewBuffer(marshal),
	)
	c.Set("userID", gofakeit.UUID())
	c.Request = request

	handler.CreateCodeExample(c)

	// check
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Nil(t, c.Errors)

	repo.AssertExpectations(t)
}

func TestCodeExampleHandler_DeleteCodeExample(t *testing.T) {
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
			Maybe().
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
	service := servises.NewCodeExampleService(repo, log)

	handler := NewCodeExampleHandler(service, log)

	// gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Run("successful deletion", func(t *testing.T) {
		request := httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/v1/texts/code-examples/%s", codeExamples[1].UUID),
			nil,
		)
		c.AddParam("uuid", codeExamples[1].UUID)
		c.Set("userID", user1ID)
		c.Request = request

		handler.DeleteCodeExample(c)
		c.Writer.WriteHeaderNow()

		// check
		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Nil(t, c.Errors)
	})

	t.Run("unsuccessful deletion due to code example not found", func(t *testing.T) {
		uuid := gofakeit.UUID()
		request := httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/v1/texts/code-examples/%s", uuid),
			nil,
		)
		c.AddParam("uuid", uuid)
		c.Set("userID", user1ID)
		c.Request = request

		handler.DeleteCodeExample(c)

		// check
		assert.Error(t, c.Errors.Last())
	})

	t.Run("unsuccessful deletion standard example", func(t *testing.T) {
		request := httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/v1/texts/code-examples/%s", codeExamples[0].UUID),
			nil,
		)
		c.AddParam("uuid", codeExamples[0].UUID)
		c.Set("userID", user1ID)
		c.Request = request

		handler.DeleteCodeExample(c)

		// check
		assert.Error(t, c.Errors.Last())
	})

	t.Run("unsuccessful deletion other user's example", func(t *testing.T) {
		request := httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/v1/texts/code-examples/%s", codeExamples[2].UUID),
			nil,
		)
		c.AddParam("uuid", codeExamples[2].UUID)
		c.Set("userID", user1ID)
		c.Request = request

		handler.DeleteCodeExample(c)

		// check
		assert.Error(t, c.Errors.Last())
	})

	repo.AssertExpectations(t)
}
