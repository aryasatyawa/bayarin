package testutil

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
)

// SetupTestRouter creates test router
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// MakeRequest makes HTTP request for testing
func MakeRequest(method, url string, body interface{}, router *gin.Engine) *httptest.ResponseRecorder {
	var bodyReader io.Reader
	if body != nil {
		jsonData, _ := json.Marshal(body)
		bodyReader = strings.NewReader(string(jsonData))
	}

	req := httptest.NewRequest(method, url, bodyReader)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// MakeAuthRequest makes authenticated HTTP request
func MakeAuthRequest(method, url, token string, body interface{}, router *gin.Engine) *httptest.ResponseRecorder {
	var bodyReader io.Reader
	if body != nil {
		jsonData, _ := json.Marshal(body)
		bodyReader = strings.NewReader(string(jsonData))
	}

	req := httptest.NewRequest(method, url, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}
