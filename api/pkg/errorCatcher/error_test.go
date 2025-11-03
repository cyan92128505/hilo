package errorCatcher

import (
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ErrorSuite struct {
	suite.Suite
}

func (suite *ErrorSuite) TestPanicIfErr() {
	suite.Panics(func() {
		PanicIfErr(
			errors.New("error"),
			errors.New("errorType"),
			errors.New("errorSubject"),
		)
	})
}

func (suite *ErrorSuite) TestPanicIfErrNoError() {
	suite.NotPanics(func() {
		PanicIfErr(nil, nil, nil)
	})
}

func (suite *ErrorSuite) TestReturnIfErr() {
	suite.Error(ReturnIfErr(errors.New("error"),
		errors.New("errorType"),
		errors.New("errorSubject"),
	),
	)
}

func (suite *ErrorSuite) TestReturnIfErrNoError() {
	suite.NoError(ReturnIfErr(nil, nil, nil))
}

func (suite *ErrorSuite) TestConcatError() {
	suite.Equal("errType: errSubject: err message", ConcatError(errors.New("errType"), errors.New("errSubject"), errors.New("err message")).Error())
}

func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(ErrorSuite))
}
