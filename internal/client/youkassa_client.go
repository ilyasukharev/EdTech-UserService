package client

import (
	"UserService/internal/model"
	"UserService/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

const (
	createPaymentEndpoint = "/v3/payments"
	getPaymentEndpoint    = "/v3/payments/%s"
)

const idempotenceHeader = "Idempotence-Key"

const sourceErr = "Youkassa"

type YoukassaAPI interface {
	CreatePayment(model.CreateYoukassaPaymentRequest) (model.CreateYoukassaPaymentResponse, error)
	GetPayment(paymentID string) (model.GetYoukassaPaymentResponse, error)
}

type YoukassaClient struct {
	C        *http.Client
	Ctx      context.Context
	BaseAuth YoukassaAuth
	BaseURL  string
}

type YoukassaAuth struct {
	ShopID string
	Token  string
}

func NewYoukassaClient(ctx context.Context, client *http.Client, url string, auth YoukassaAuth) *YoukassaClient {
	return &YoukassaClient{
		C:        client,
		Ctx:      ctx,
		BaseURL:  url,
		BaseAuth: auth,
	}
}

func (cl *YoukassaClient) CreatePayment(payment model.CreateYoukassaPaymentRequest) (model.CreateYoukassaPaymentResponse, error) {
	var response model.CreateYoukassaPaymentResponse

	bbytes, err := json.Marshal(payment)
	if err != nil {
		return response, err
	}

	rawResponse, err := cl.makeYoukassaRequest(http.MethodPost, createPaymentEndpoint, bbytes)
	if err != nil {
		return response, err
	}
	defer rawResponse.Body.Close()

	if err = json.NewDecoder(rawResponse.Body).Decode(&response); err != nil {
		return response, err
	}

	return response, nil
}

func (cl *YoukassaClient) GetPayment(paymentID string) (model.GetYoukassaPaymentResponse, error) {
	var response model.GetYoukassaPaymentResponse

	rawResponse, err := cl.makeYoukassaRequest(http.MethodGet,
		fmt.Sprintf(getPaymentEndpoint, paymentID), nil)
	if err != nil {
		return response, err
	}
	defer rawResponse.Body.Close()

	if err = json.NewDecoder(rawResponse.Body).Decode(&response); err != nil {
		return response, err
	}

	return response, nil
}

func (cl *YoukassaClient) makeYoukassaRequest(method string, endpoint string, body []byte) (*http.Response, error) {
	request, err := http.NewRequestWithContext(cl.Ctx, method,
		cl.BaseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(cl.BaseAuth.ShopID, cl.BaseAuth.Token)
	request.Header.Set(idempotenceHeader, uuid.NewString())

	resp, err := cl.C.Do(request)
	if err != nil {
		return nil, err
	}

	if err = utils.CheckHttpClientStatusCodeIsPositive(sourceErr, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
