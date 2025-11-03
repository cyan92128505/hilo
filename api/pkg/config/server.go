package config

import (
	"time"
)

// Server type
type Server struct {
	ReleaseMode          bool          `split_words:"true" default:"true"`
	Port                 string        `split_words:"true" default:"3000"`
	ServerTimeout        time.Duration `split_words:"true" default:"5s"`
	PrefixMessage        string        `split_words:"true" default:"[Gin]"`
	CustomizedRender     bool          `split_words:"true" default:"false"`
	AllowAllOrigins      bool          `split_words:"true" default:"false"`
	AllowOrigins         []string      `split_words:"true" default:"http://localhost,https://localhost"`
	AllowedPaths         []string      `split_words:"true" default:"/favicon.ico,/ping,/api/v1/auth/register,/api/v1/auth/login"`
	JWTGuard             bool          `split_words:"true" default:"true"`
	MaxMultipartMemoryMB int64         `split_words:"true" default:"8"`
}
