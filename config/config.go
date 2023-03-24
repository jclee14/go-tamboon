package config

import "os"

type IConfig interface {
	OmiseBaseAPI() string
	OmiseCurrency() string
	OmisePublicKey() string
	OmiseSecretKey() string
}

type config struct {
	omiseCurrency  string
	omisePublicKey string
	omiseSecretKey string
	omiseBaseAPI   string
}

func NewConfig() IConfig {
	return &config{}
}

func (c *config) OmiseBaseAPI() string {
	if len(c.omiseBaseAPI) == 0 {
		c.omiseBaseAPI = os.Getenv("OMISE_BASE_API")
	}
	return c.omiseBaseAPI
}

func (c *config) OmiseCurrency() string {
	if len(c.omiseCurrency) == 0 {
		c.omiseCurrency = os.Getenv("OMISE_CURRENCY")
	}
	return c.omiseCurrency
}

func (c *config) OmisePublicKey() string {
	if len(c.omisePublicKey) == 0 {
		c.omisePublicKey = os.Getenv("OMISE_PUBLIC_KEY")
	}
	return c.omisePublicKey
}

func (c *config) OmiseSecretKey() string {
	if len(c.omiseSecretKey) == 0 {
		c.omiseSecretKey = os.Getenv("OMISE_SECRET_KEY")
	}
	return c.omiseSecretKey
}
