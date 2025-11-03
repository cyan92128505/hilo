package config

// JWT type
type JWT struct {
	PrivateKeyPath string `split_words:"true" default:"./es256_private.pem"`
	PrivateKey     string `split_words:"true" default:""`
}
