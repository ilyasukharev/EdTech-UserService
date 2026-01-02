package client

import (
	"UserService/internal/model"
	"UserService/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	getSubscriptionResponseEndpoint = "/subscription/%s"
	reserveDraftServerEndpoint      = "/server:draft-reserve"
	clearDraftReserveServerEndpoint = "/server:draft-clear"
)

type SubscriptionServiceAPI interface {
	GetSubscription(string, int) (model.SubscriptionResponse, error)
	ReserveServerDraft(requestID string, subID uint) (model.ReserveDraftServerResponse, error)
	ClearServerDraftReserve(requestID string, serverID int) error
}

type SubscriptionServiceClient struct {
	Ctx     context.Context
	Client  *http.Client
	BaseURL string
}

func NewSubscriptionServiceAPI(ctx context.Context, httpClient *http.Client, baseURL string) *SubscriptionServiceClient {
	return &SubscriptionServiceClient{
		Ctx:     ctx,
		Client:  httpClient,
		BaseURL: baseURL,
	}
}

func (c *SubscriptionServiceClient) GetSubscription(requestID string, subID int) (model.SubscriptionResponse, error) {
	var subscriptionResponse model.SubscriptionResponse

	resp, err := c.makeRequest(requestID, http.MethodGet,
		c.BaseURL+fmt.Sprintf(getSubscriptionResponseEndpoint, strconv.Itoa(subID)), nil)
	if err != nil {
		return subscriptionResponse, err
	}

	if err = json.NewDecoder(resp.Body).Decode(&subscriptionResponse); err != nil {
		return subscriptionResponse, err
	}

	return subscriptionResponse, nil
}

func (c *SubscriptionServiceClient) ReserveServerDraft(requestID string, subID uint) (model.ReserveDraftServerResponse, error) {
	var reserveServerResponse model.ReserveDraftServerResponse

	resp, err := c.makeRequest(requestID, http.MethodGet,
		c.BaseURL+reserveDraftServerEndpoint+fmt.Sprintf("?subscriptionId=%d", subID), nil)
	if err != nil {
		return reserveServerResponse, err
	}

	if err = json.NewDecoder(resp.Body).Decode(&reserveServerResponse); err != nil {
		return reserveServerResponse, err
	}

	return reserveServerResponse, nil
}

func (c *SubscriptionServiceClient) ClearServerDraftReserve(requestID string, serverID int) error {
	_, err := c.makeRequest(requestID, http.MethodGet,
		c.BaseURL+clearDraftReserveServerEndpoint+fmt.Sprintf("?serverId=%d", serverID), nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *SubscriptionServiceClient) makeRequest(reqID string, method string, endpoint string, body []byte) (*http.Response, error) {
	request, err := http.NewRequestWithContext(c.Ctx, method,
		c.BaseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Add(model.XRequestIDHeader, reqID)

	resp, err := c.Client.Do(request)
	if err != nil {
		return nil, err
	}

	if err = utils.CheckHttpClientStatusCodeIsPositive(sourceErr, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
