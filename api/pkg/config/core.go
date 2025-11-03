package config

// Core type
type Core struct {
	LoggerMode    string `split_words:"true" default:"customized"`
	LogLevel      string `split_words:"true" default:"INFO"`
	IsReleaseMode bool   `split_words:"true" default:"false"`
	SystemName    string `split_words:"true" default:"system"`
}
