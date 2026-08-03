package utils

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// SetUpGinTest creates a gin.Context and ResponseRecorder for testing a handler directly,
// without running through a full HTTP router.
func SetUpGinTest(method, url, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

// SetUpGinTestWithParams is like SetUpGinTest but also sets URL path parameters on
// the context, useful for handlers that read from c.Param (e.g. /:id routes).
func SetUpGinTestWithParams(method, url, body string, params gin.Params) (*gin.Context, *httptest.ResponseRecorder) {
	c, w := SetUpGinTest(method, url, body)
	c.Params = params
	return c, w
}

// CheckData asserts that the JSON response body contains all expected key-value pairs.
// Values are compared after JSON round-tripping to avoid type mismatch (e.g. int vs float64).
func CheckData(t *testing.T, w *httptest.ResponseRecorder, expected map[string]interface{}) {
	t.Helper()
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("response body is not valid JSON: %v (raw: %s)", err, w.Body.String())
	}
	for key, expectedVal := range expected {
		gotVal, ok := response[key]
		if !ok {
			t.Errorf("expected key %q in response body, got: %v", key, response)
			continue
		}
		expectedJSON, _ := json.Marshal(expectedVal)
		gotJSON, _ := json.Marshal(gotVal)
		if string(expectedJSON) != string(gotJSON) {
			t.Errorf("key %q: expected %s, got %s", key, expectedJSON, gotJSON)
		}
	}
}

// CheckError asserts that the JSON response body contains an "error" key matching expectedError.
func CheckError(t *testing.T, w *httptest.ResponseRecorder, expectedError string) {
	t.Helper()
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("response body is not valid JSON: %v (raw: %s)", err, w.Body.String())
	}
	got, ok := response["error"]
	if !ok {
		t.Errorf("expected \"error\" key in response body, got: %v", response)
		return
	}
	if got != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, got)
	}
}
