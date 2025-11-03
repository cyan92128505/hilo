package restful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hilo-api/pkg/errorCatcher"
	"hilo-api/pkg/restful"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewMockGinServer(logger *zap.Logger, guard *restful.JWTGuarder, whitelist ...string) (engine *gin.Engine, err error) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(
		gin.Logger(),
		errorCatcher.GinPanicErrorHandler(logger, "Mock Test Gin Server"),
		guard.JWTGuarder(whitelist...))
	return router, nil
}

// Get method
func Get(uri string, headers map[string]string, router *gin.Engine) ([]byte, error) {
	req := httptest.NewRequest(http.MethodGet, uri, nil)
	return getBody(req, headers, router)
}

// PostJSON method
func PostJSON(uri string, param map[string]interface{}, headers map[string]string, router *gin.Engine) ([]byte, error) {
	jsonByte, _ := json.Marshal(param)
	req := httptest.NewRequest(http.MethodPost, uri, bytes.NewReader(jsonByte))
	return getBody(req, headers, router)
}

// PutJSON method
func PutJSON(uri string, param map[string]interface{}, headers map[string]string, router *gin.Engine) ([]byte, error) {
	jsonByte, _ := json.Marshal(param)
	req := httptest.NewRequest(http.MethodPut, uri, bytes.NewReader(jsonByte))
	return getBody(req, headers, router)
}

// DeleteJSON method
func DeleteJSON(uri string, param map[string]interface{}, headers map[string]string, router *gin.Engine) ([]byte, error) {
	jsonByte, _ := json.Marshal(param)
	req := httptest.NewRequest(http.MethodDelete, uri, bytes.NewReader(jsonByte))
	return getBody(req, headers, router)
}

func getBody(req *http.Request, headers map[string]string, router *gin.Engine) ([]byte, error) {
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	log.Println(result.StatusCode)
	if result.StatusCode >= 400 {
		return nil, fmt.Errorf("request error by code: %d", result.StatusCode)
	}

	body, _ := io.ReadAll(result.Body)
	return body, nil
}
