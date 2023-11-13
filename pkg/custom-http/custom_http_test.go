package customhttp_test

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mockCustomHttp "github.com/edwynrrangel/go-libraries/mocks/custom-http"
	customhttp "github.com/edwynrrangel/go-libraries/pkg/custom-http"
)

var (
	mockHttpClient *mockCustomHttp.MockHttpClient
	httpsTest      customhttp.Http
)

func casesSetup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockHttpClient = mockCustomHttp.NewMockHttpClient(ctrl)

	httpsTest = customhttp.NewCustomHttp(mockHttpClient, 0, true)

	return func() {
		mockHttpClient = nil
		httpsTest = nil
		ctrl.Finish()
	}
}

func TestDoRequest(t *testing.T) {
	paramsRequest := customhttp.ParamsRequest{
		Path:   "http://127.0.0.1:8080",
		Method: http.MethodGet,
		QueryParams: map[string]string{
			"param1": "value1",
		},
		Headers: map[string]string{
			"header1": "value1",
		},
	}
	type expectedBody struct {
		StatusCode int
		Message    string
	}
	type args struct {
		params             customhttp.ParamsRequest
		expectedStatusCode int
		body               interface{}
		bodyError          interface{}
	}
	testCases := []struct {
		desc         string
		args         args
		mockResponse func(args)
		expected     string
		expectedErr  error
	}{
		{
			desc: "return error when url.Parse fails",
			args: args{
				params: customhttp.ParamsRequest{
					Path:   "://",
					Method: http.MethodPost,
				},
			},
			mockResponse: func(args args) {},
			expected:     "",
			expectedErr: &url.Error{
				Op:  "parse",
				URL: "://",
				Err: errors.New("missing protocol scheme"),
			},
		},
		{
			desc: "return error when h.client.Do fails",
			args: args{
				params:             paramsRequest,
				expectedStatusCode: http.StatusOK,
			},
			mockResponse: func(args args) {
				mockHttpClient.
					EXPECT().
					Do(gomock.Any()).
					Return(nil, assert.AnError)
			},
			expected:    "",
			expectedErr: assert.AnError,
		},
		{
			desc: "return error when response status code is not the expected but bodyError is nil",
			args: args{
				params:             paramsRequest,
				expectedStatusCode: http.StatusOK,
			},
			mockResponse: func(args args) {
				mockHttpClient.
					EXPECT().
					Do(gomock.Any()).
					Return(&http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       io.NopCloser(strings.NewReader(`{"message": "hello world"}`)),
					}, nil)
			},
			expected:    "",
			expectedErr: errors.New(`{"status_code": 400, "body": "{"message": "hello world"}"}`),
		},
		{
			desc: "return error when response status code is not the expected and bodyError is not nil but error unmarshaling bodyError",
			args: args{
				params:             paramsRequest,
				expectedStatusCode: http.StatusOK,
				bodyError: &expectedBody{
					StatusCode: http.StatusBadRequest,
					Message:    "hello world",
				},
			},
			mockResponse: func(args args) {
				mockHttpClient.
					EXPECT().
					Do(gomock.Any()).
					Return(&http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       io.NopCloser(strings.NewReader(``)),
					}, nil)
			},
			expected:    "",
			expectedErr: errors.New(`{"status_code": 400, "body": "", "error": "unexpected end of JSON input"}`),
		},
		{
			desc: "return error when response status code is not the expected and bodyError is not nil",
			args: args{
				params:             paramsRequest,
				expectedStatusCode: http.StatusOK,
				bodyError: &expectedBody{
					StatusCode: http.StatusBadRequest,
					Message:    "hello world",
				},
			},
			mockResponse: func(args args) {
				mockHttpClient.
					EXPECT().
					Do(gomock.Any()).
					Return(&http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       io.NopCloser(strings.NewReader(`{"message": "hello world"}`)),
					}, nil)
			},
			expected:    "",
			expectedErr: errors.New(`{"status_code": 400, "body": "{"message": "hello world"}"}`),
		},
		{
			desc: "return error when response status code is the expected but error unmarshaling body",
			args: args{
				params:             paramsRequest,
				expectedStatusCode: http.StatusOK,
				body:               &expectedBody{},
			},
			mockResponse: func(args args) {
				mockHttpClient.
					EXPECT().
					Do(gomock.Any()).
					Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(``)),
					}, nil)
			},
			expected:    "",
			expectedErr: errors.New(`{"status_code": 200, "body": "", "error": "unexpected end of JSON input"}`),
		},
		{
			desc: "return success when response status code is the expected and body is not nil",
			args: args{
				params:             paramsRequest,
				expectedStatusCode: http.StatusOK,
				body: &expectedBody{
					StatusCode: http.StatusOK,
					Message:    "hello world",
				},
			},
			mockResponse: func(args args) {
				mockHttpClient.
					EXPECT().
					Do(gomock.Any()).
					Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"status_code": 200, "message": "hello world"}`)),
					}, nil)
			},
			expected:    `{"status_code": 200, "body": "{"status_code": 200, "message": "hello world"}"}`,
			expectedErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			defer casesSetup(t)()

			tc.mockResponse(tc.args)

			got, err := httpsTest.DoRequest(tc.args.params, tc.args.expectedStatusCode, tc.args.body, tc.args.bodyError)

			assert.Equal(t, tc.expected, got)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestNewCustomHttp(t *testing.T) {
	t.Run("custom client", func(t *testing.T) {
		customClient := &http.Client{}
		httpsClient := customhttp.NewCustomHttp(customClient, 0, false)
		assert.NotNil(t, httpsClient)
	})

	t.Run("default client", func(t *testing.T) {
		httpsClient := customhttp.NewCustomHttp(nil, 0, false)
		assert.NotNil(t, httpsClient)
	})

	t.Run("default client with insecure", func(t *testing.T) {
		httpsClient := customhttp.NewCustomHttp(nil, 0, true)
		assert.NotNil(t, httpsClient)
	})

	t.Run("default client with timeout", func(t *testing.T) {
		httpsClient := customhttp.NewCustomHttp(nil, 1, false)
		assert.NotNil(t, httpsClient)
	})
}
