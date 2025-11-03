package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type CommonHandlerSuite struct {
	suite.Suite
}

func (suite *CommonHandlerSuite) TestError404() {
	suite.Equal("gin.HandlerFunc", reflect.TypeOf(Error404()).String())
}

func (suite *CommonHandlerSuite) TestError404Run() {
	r := gin.Default()
	r.GET("/ping", Error404())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
}

func (suite *CommonHandlerSuite) TestQuickReply() {
	suite.Equal("gin.HandlerFunc", reflect.TypeOf(QuickReply()).String())
}

func (suite *CommonHandlerSuite) TestQuickReplyRun() {
	r := gin.Default()
	r.GET("/ping", QuickReply())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}

func TestCommonHandlerSuite(t *testing.T) {
	suite.Run(t, new(CommonHandlerSuite))
}
