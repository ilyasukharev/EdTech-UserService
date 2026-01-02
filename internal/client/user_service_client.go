package client

import (
	"UserService/internal/model"
	"UserService/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const getUserResponseEndpoint = "/user/%s"

type UserServiceAPI interface {
	GetUser(string, string) (model.UserResponse, error)
}

type UserServiceClient struct {
	Ctx     context.Context
	Cl      *http.Client
	BaseURL string
}

func NewUserServiceClient(ctx context.Context, httpClient *http.Client, baseURL string) *UserServiceClient {
	return &UserServiceClient{
		Ctx:     ctx,
		Cl:      httpClient,
		BaseURL: baseURL,
	}
}

func (c *UserServiceClient) GetUser(requestID string, userID string) (model.UserResponse, error) {
	var user model.UserResponse

	req, err := http.NewRequestWithContext(c.Ctx, http.MethodGet,
		c.BaseURL+fmt.Sprintf(getUserResponseEndpoint, userID), nil)
	if err != nil {
		return user, err
	}

	req.Header.Add(model.XRequestIDHeader, requestID)

	resp, err := c.Cl.Do(req)
	if err != nil {
		return user, err
	}

	if err = utils.CheckHttpClientStatusCodeIsPositive(model.SourceApiErr, resp); err != nil {
		return user, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}
