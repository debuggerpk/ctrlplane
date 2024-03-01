// Package slack provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen, a modified copy of github.com/deepmap/oapi-codegen/v2.
// It was modified to add support for the following features:
//  - Support for custom templates by filename.
//  - Supporting x-breu-entity in the schema to generate a struct for the entity.
//
// DO NOT EDIT!!

package slack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	"go.breu.io/quantm/internal/shared"
	externalRef0 "go.breu.io/quantm/internal/shared"
)

const (
	APIKeyAuthScopes = "APIKeyAuth.Scopes"
	BearerAuthScopes = "BearerAuth.Scopes"
)

var (
	ErrInvalidSlackStatus = errors.New("invalid SlackStatus value")
)

type (
	SlackStatusMapType map[string]SlackStatus // SlackStatusMapType is a quick lookup map for SlackStatus.
)

// Defines values for SlackStatus.
const (
	SlackStatusError SlackStatus = "Error"
	SlackStatusOk    SlackStatus = "Ok"
)

// SlackStatusMap returns all known values for SlackStatus.
var (
	SlackStatusMap = SlackStatusMapType{
		SlackStatusError.String(): SlackStatusError,
		SlackStatusOk.String():    SlackStatusOk,
	}
)

/*
 * Helper methods for SlackStatus for easy marshalling and unmarshalling.
 */
func (v SlackStatus) String() string               { return string(v) }
func (v SlackStatus) MarshalJSON() ([]byte, error) { return json.Marshal(v.String()) }
func (v *SlackStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, ok := SlackStatusMap[s]
	if !ok {
		return ErrInvalidSlackStatus
	}

	*v = val

	return nil
}

// SlackNotification defines model for SlackNotification.
type SlackNotification struct {
	// ChannelID The ID of the Slack channel where the message will be sent.
	ChannelID string `json:"channelID"`

	// Message The message content to send to the Slack channel.
	Message string `json:"message"`
}

// SlackResponse defines model for SlackResponse.
type SlackResponse struct {
	Errors *map[string]string `json:"errors,omitempty"`
	Status SlackStatus        `json:"status"`
}

// SlackStatus defines model for SlackStatus.
type SlackStatus string

// SlackOauthParams defines parameters for SlackOauth.
type SlackOauthParams struct {
	Code string `form:"code" json:"code"`
}

// SendNotificationToChannelParams defines parameters for SendNotificationToChannel.
type SendNotificationToChannelParams struct {
	// ChannelID The ID of the Slack channel to send the notification to.
	ChannelID string `form:"channelID" json:"channelID"`
}

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// Login request
	Login(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// SlackOauth request
	SlackOauth(ctx context.Context, params *SlackOauthParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// SendNotificationToChannel request
	SendNotificationToChannel(ctx context.Context, params *SendNotificationToChannelParams, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) Login(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewLoginRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SlackOauth(ctx context.Context, params *SlackOauthParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSlackOauthRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SendNotificationToChannel(ctx context.Context, params *SendNotificationToChannelParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSendNotificationToChannelRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewLoginRequest generates requests for Login
func NewLoginRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/auth/slack/login")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewSlackOauthRequest generates requests for SlackOauth
func NewSlackOauthRequest(server string, params *SlackOauthParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/auth/slack/login/callback")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "code", runtime.ParamLocationQuery, params.Code); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewSendNotificationToChannelRequest generates requests for SendNotificationToChannel
func NewSendNotificationToChannelRequest(server string, params *SendNotificationToChannelParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/slack/notification")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "channelID", runtime.ParamLocationQuery, params.ChannelID); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// LoginWithResponse request
	LoginWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*LoginResponse, error)

	// SlackOauthWithResponse request
	SlackOauthWithResponse(ctx context.Context, params *SlackOauthParams, reqEditors ...RequestEditorFn) (*SlackOauthResponse, error)

	// SendNotificationToChannelWithResponse request
	SendNotificationToChannelWithResponse(ctx context.Context, params *SendNotificationToChannelParams, reqEditors ...RequestEditorFn) (*SendNotificationToChannelResponse, error)
}

type LoginResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *externalRef0.BadRequest
	JSON500      *externalRef0.InternalServerError
}

// Status returns HTTPResponse.Status
func (r LoginResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r LoginResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type SlackOauthResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *SlackResponse
	JSON400      *externalRef0.BadRequest
	JSON500      *externalRef0.InternalServerError
}

// Status returns HTTPResponse.Status
func (r SlackOauthResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SlackOauthResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type SendNotificationToChannelResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *SlackNotification
	JSON400      *externalRef0.BadRequest
	JSON500      *externalRef0.InternalServerError
}

// Status returns HTTPResponse.Status
func (r SendNotificationToChannelResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SendNotificationToChannelResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// LoginWithResponse request returning *LoginResponse
func (c *ClientWithResponses) LoginWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*LoginResponse, error) {
	rsp, err := c.Login(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseLoginResponse(rsp)
}

// SlackOauthWithResponse request returning *SlackOauthResponse
func (c *ClientWithResponses) SlackOauthWithResponse(ctx context.Context, params *SlackOauthParams, reqEditors ...RequestEditorFn) (*SlackOauthResponse, error) {
	rsp, err := c.SlackOauth(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSlackOauthResponse(rsp)
}

// SendNotificationToChannelWithResponse request returning *SendNotificationToChannelResponse
func (c *ClientWithResponses) SendNotificationToChannelWithResponse(ctx context.Context, params *SendNotificationToChannelParams, reqEditors ...RequestEditorFn) (*SendNotificationToChannelResponse, error) {
	rsp, err := c.SendNotificationToChannel(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSendNotificationToChannelResponse(rsp)
}

// ParseLoginResponse parses an HTTP response from a LoginWithResponse call
func ParseLoginResponse(rsp *http.Response) (*LoginResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &LoginResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest externalRef0.BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest externalRef0.InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseSlackOauthResponse parses an HTTP response from a SlackOauthWithResponse call
func ParseSlackOauthResponse(rsp *http.Response) (*SlackOauthResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &SlackOauthResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest SlackResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest externalRef0.BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest externalRef0.InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseSendNotificationToChannelResponse parses an HTTP response from a SendNotificationToChannelWithResponse call
func ParseSendNotificationToChannelResponse(rsp *http.Response) (*SendNotificationToChannelResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &SendNotificationToChannelResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest SlackNotification
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest externalRef0.BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest externalRef0.InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Initiate slack login
	// (GET /v1/auth/slack/login)
	Login(ctx echo.Context) error

	// Callback after Slack login
	// (GET /v1/auth/slack/login/callback)
	SlackOauth(ctx echo.Context) error

	// sends a early warning notification to channel
	// (GET /v1/slack/notification)
	SendNotificationToChannel(ctx echo.Context) error

	// SecurityHandler returns the underlying Security Wrapper
	SecureHandler(handler echo.HandlerFunc, ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// Login converts echo context to params.

func (w *ServerInterfaceWrapper) Login(ctx echo.Context) error {
	var err error

	// Get the handler, get the secure handler if needed and then invoke with unmarshalled params.
	handler := w.Handler.Login
	err = handler(ctx)

	return err
}

// SlackOauth converts echo context to params.

func (w *ServerInterfaceWrapper) SlackOauth(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params SlackOauthParams
	// ------------- Required query parameter "code" -------------

	err = runtime.BindQueryParameter("form", true, true, "code", ctx.QueryParams(), &params.Code)
	if err != nil {
		return shared.NewAPIError(http.StatusBadRequest, fmt.Errorf("Invalid format for parameter code: %s", err))
	}

	// Get the handler, get the secure handler if needed and then invoke with unmarshalled params.
	handler := w.Handler.SlackOauth
	err = handler(ctx)

	return err
}

// SendNotificationToChannel converts echo context to params.

func (w *ServerInterfaceWrapper) SendNotificationToChannel(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	ctx.Set(APIKeyAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params SendNotificationToChannelParams
	// ------------- Required query parameter "channelID" -------------

	err = runtime.BindQueryParameter("form", true, true, "channelID", ctx.QueryParams(), &params.ChannelID)
	if err != nil {
		return shared.NewAPIError(http.StatusBadRequest, fmt.Errorf("Invalid format for parameter channelID: %s", err))
	}

	// Get the handler, get the secure handler if needed and then invoke with unmarshalled params.
	handler := w.Handler.SendNotificationToChannel
	secure := w.Handler.SecureHandler
	err = secure(handler, ctx)

	return err
}

// EchoRouter is an interface that wraps the methods of echo.Echo & echo.Group to provide a common interface
// for registering routes.
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/v1/auth/slack/login", wrapper.Login)
	router.GET(baseURL+"/v1/auth/slack/login/callback", wrapper.SlackOauth)
	router.GET(baseURL+"/v1/slack/notification", wrapper.SendNotificationToChannel)

}
