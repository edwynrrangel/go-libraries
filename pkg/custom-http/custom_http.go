package customhttp

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

//go:generate mockgen -destination=../../mocks/custom-http/http_client.go -package=customhttp github.com/edwynrrangel/go-libraries/pkg/custom-http HttpClient
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

//go:generate mockgen -destination=../../mocks/custom-http/http.go -package=customhttp github.com/edwynrrangel/go-libraries/pkg/custom-http Http
type Http interface {
	DoRequest(params ParamsRequest, expectedStatusCode int, body, bodyError interface{}) (string, error)
}

type customHttp struct {
	client HttpClient
}

type ParamsRequest struct {
	Method      string
	Path        string
	Data        io.Reader
	Headers     map[string]string
	QueryParams map[string]string
}

func NewCustomHttp(client HttpClient, timeout time.Duration, skipVerifyTLS bool) Http {

	if client != nil {
		return &customHttp{
			client: client,
		}
	}

	httpClient := new(http.Client)

	if timeout == 0 {
		timeout = 15 * time.Second
	}
	httpClient.Timeout = timeout

	if skipVerifyTLS {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}

		httpClient.Transport = tr
	}

	return &customHttp{
		client: httpClient,
	}
}

func (h *customHttp) DoRequest(params ParamsRequest, expectedStatusCode int, body, bodyError interface{}) (string, error) {
	// Crear un objeto URL a partir del campo Path
	path, err := url.Parse(params.Path)
	if err != nil {
		return "", err
	}

	// Añadir los parámetros de consulta al objeto URL
	query := path.Query()
	for key, value := range params.QueryParams {
		query.Set(key, value)
	}
	path.RawQuery = query.Encode()

	req, err := http.NewRequest(
		params.Method,
		path.String(),
		params.Data,
	)
	if err != nil {
		return "", err
	}

	// Añadir headers al request
	for key, value := range params.Headers {
		req.Header.Set(key, value)
	}

	response, err := h.client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	bodyByte, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// si el status code no es el esperado, se intenta unmarshal el body en el bodyError
	if response.StatusCode != expectedStatusCode {
		if bodyError != nil {
			err = json.Unmarshal(bodyByte, &bodyError)
			if err != nil {
				return "", fmt.Errorf(`{"status_code": %d, "body": "%s", "error": "%s"}`, response.StatusCode, string(bodyByte), err.Error())
			}
		}
		return "", fmt.Errorf(`{"status_code": %d, "body": "%s"}`, response.StatusCode, string(bodyByte))
	}

	// si el status code es el esperado, se intenta unmarshal el body en el body
	if body != nil {
		err = json.Unmarshal(bodyByte, &body)
		if err != nil {
			return "", fmt.Errorf(`{"status_code": %d, "body": "%s", "error": "%s"}`, response.StatusCode, string(bodyByte), err.Error())
		}
	}

	return fmt.Sprintf(`{"status_code": %d, "body": "%s"}`, response.StatusCode, string(bodyByte)), nil
}
