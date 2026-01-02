package model

import (
	"UserService/internal/config"
	"net/http"
	"time"
)

func NewHttpClient(cfg *config.Config) *http.Client {
	defaultClient := cfg.DefaultClient
	return &http.Client{
		Timeout: time.Second * time.Duration(defaultClient.Timeout),
		Transport: &http.Transport{
			MaxIdleConns:        defaultClient.MaxIdleConn,
			MaxIdleConnsPerHost: defaultClient.MaxIdleConnPerHost,
			IdleConnTimeout:     time.Duration(defaultClient.IdleConnTimeout) * time.Second,
		},
	}
}
