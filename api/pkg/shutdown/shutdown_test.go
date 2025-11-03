package shutdown

import (
	"github.com/stretchr/testify/suite"
	"os"
	"reflect"
	"syscall"
	"testing"
	"time"
)

type ShutdownTestSuite struct {
	suite.Suite
}

func (suite *ShutdownTestSuite) TestWithQuitOption() {
	suite.Equal("chan os.Signal", reflect.TypeOf(NewShutdown(WithQuit(make(chan os.Signal))).quit).String())
}

func (suite *ShutdownTestSuite) TestWithDoneOption() {
	suite.Equal("chan bool", reflect.TypeOf(NewShutdown(WithDone(make(chan bool))).done).String())
}

func (suite *ShutdownTestSuite) TestWithServerTimeoutOption() {
	suite.Equal("time.Duration", reflect.TypeOf(NewShutdown(WithServerTimeout(5*time.Second)).serverTimeout).String())
}

func (suite *ShutdownTestSuite) TestWithEndTaskOption() {
	suite.Equal("func()", reflect.TypeOf(NewShutdown(WithEndTask(func() {})).endTask).String())
}

func TestShutdownTestSuite(t *testing.T) {
	suite.Run(t, new(ShutdownTestSuite))
}

func Test_Shutdown(t *testing.T) {
	t.Run("test Shutdown", func(t *testing.T) {
		quit := make(chan os.Signal)
		defer close(quit)
		done := make(chan bool)
		defer close(done)

		go func() {
			(&Shutdown{
				quit: quit,
				done: done,
				endTask: func() {
					t.Log("endTask")
				},
			}).Shutdown()
		}()
		quit <- syscall.SIGTERM
		<-done
	})
}
