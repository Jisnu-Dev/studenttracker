package utils_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers/utils"
	"github.com/gin-gonic/gin"
)

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	password := "my_secure_password"
	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if hash == "" {
		t.Fatal("expected hash to not be empty")
	}

	if !utils.CheckPasswordHash(password, hash) {
		t.Error("expected CheckPasswordHash to return true for correct password")
	}

	if utils.CheckPasswordHash("wrong_password", hash) {
		t.Error("expected CheckPasswordHash to return false for incorrect password")
	}
}

func TestRespondWithError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	utils.RespondWithError(c, http.StatusBadRequest, "an error occurred")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
	if response["error"] != "an error occurred" {
		t.Errorf("expected error message 'an error occurred', got '%s'", response["error"])
	}
}

func TestRespondWithJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	payload := gin.H{"message": "success", "id": 123}
	utils.RespondWithJSON(c, http.StatusOK, payload)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
	if response["message"] != "success" {
		t.Errorf("expected message 'success', got '%v'", response["message"])
	}
	// json unmarshals numbers to float64
	if response["id"] != float64(123) {
		t.Errorf("expected id 123, got '%v'", response["id"])
	}
}

func TestBindJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"test name"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		var target struct {
			Name string `json:"name"`
		}

		if !utils.BindJSON(c, &target) {
			t.Error("expected BindJSON to return true")
		}
		if target.Name != "test name" {
			t.Errorf("expected name 'test name', got '%s'", target.Name)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{invalid`))
		c.Request.Header.Set("Content-Type", "application/json")

		var target struct {
			Name string `json:"name"`
		}

		if utils.BindJSON(c, &target) {
			t.Error("expected BindJSON to return false")
		}
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestParseParamID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "42"}}

		id, ok := utils.ParseParamID(c, "id")
		if !ok {
			t.Error("expected ok to be true")
		}
		if id != 42 {
			t.Errorf("expected id 42, got %d", id)
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		id, ok := utils.ParseParamID(c, "id")
		if ok {
			t.Error("expected ok to be false")
		}
		if id != 0 {
			t.Errorf("expected id 0, got %d", id)
		}
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}
