// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// ChatMessage defines model for ChatMessage.
type ChatMessage struct {
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	Id        string     `json:"id"`
	Message   string     `json:"message"`
	RoomId    string     `json:"roomId"`
	UserId    string     `json:"userId"`
}

// ChatRoom defines model for ChatRoom.
type ChatRoom struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Error defines model for Error.
type Error struct {
	// Code Error code
	Code int32 `json:"code"`

	// Message Error message
	Message string `json:"message"`
}

// Events defines model for Events.
type Events struct {
	Description string `json:"description"`
	ImageUrl    string `json:"imageUrl"`
	Name        string `json:"name"`
}

// EventsWithID defines model for EventsWithID.
type EventsWithID struct {
	Description string `json:"description"`
	Id          string `json:"id"`
	ImageUrl    string `json:"imageUrl"`
	Name        string `json:"name"`
}

// GetTokenFormdataBody defines parameters for GetToken.
type GetTokenFormdataBody struct {
	Password string `form:"password" json:"password"`
	Username string `form:"username" json:"username"`
}

// GetTokenFormdataRequestBody defines body for GetToken for application/x-www-form-urlencoded ContentType.
type GetTokenFormdataRequestBody GetTokenFormdataBody

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
	// GetTokenWithBody request with any body
	GetTokenWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	GetTokenWithFormdataBody(ctx context.Context, body GetTokenFormdataRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ListChatRooms request
	ListChatRooms(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ListChatMessages request
	ListChatMessages(ctx context.Context, roomId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ListEvents request
	ListEvents(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// HealthCheck request
	HealthCheck(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ChatWebSocket request
	ChatWebSocket(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetTokenWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetTokenRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetTokenWithFormdataBody(ctx context.Context, body GetTokenFormdataRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetTokenRequestWithFormdataBody(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ListChatRooms(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListChatRoomsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ListChatMessages(ctx context.Context, roomId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListChatMessagesRequest(c.Server, roomId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ListEvents(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListEventsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) HealthCheck(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewHealthCheckRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ChatWebSocket(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewChatWebSocketRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetTokenRequestWithFormdataBody calls the generic GetToken builder with application/x-www-form-urlencoded body
func NewGetTokenRequestWithFormdataBody(server string, body GetTokenFormdataRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	bodyStr, err := runtime.MarshalForm(body, nil)
	if err != nil {
		return nil, err
	}
	bodyReader = strings.NewReader(bodyStr.Encode())
	return NewGetTokenRequestWithBody(server, "application/x-www-form-urlencoded", bodyReader)
}

// NewGetTokenRequestWithBody generates requests for GetToken with any type of body
func NewGetTokenRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/auth/token")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewListChatRoomsRequest generates requests for ListChatRooms
func NewListChatRoomsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/chats")
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

// NewListChatMessagesRequest generates requests for ListChatMessages
func NewListChatMessagesRequest(server string, roomId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "roomId", runtime.ParamLocationPath, roomId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/chats/%s/messages", pathParam0)
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

// NewListEventsRequest generates requests for ListEvents
func NewListEventsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/events")
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

// NewHealthCheckRequest generates requests for HealthCheck
func NewHealthCheckRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/health")
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

// NewChatWebSocketRequest generates requests for ChatWebSocket
func NewChatWebSocketRequest(server string, id string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "id", runtime.ParamLocationPath, id)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/ws/%s", pathParam0)
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
	// GetTokenWithBodyWithResponse request with any body
	GetTokenWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*GetTokenResponse, error)

	GetTokenWithFormdataBodyWithResponse(ctx context.Context, body GetTokenFormdataRequestBody, reqEditors ...RequestEditorFn) (*GetTokenResponse, error)

	// ListChatRoomsWithResponse request
	ListChatRoomsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListChatRoomsResponse, error)

	// ListChatMessagesWithResponse request
	ListChatMessagesWithResponse(ctx context.Context, roomId string, reqEditors ...RequestEditorFn) (*ListChatMessagesResponse, error)

	// ListEventsWithResponse request
	ListEventsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListEventsResponse, error)

	// HealthCheckWithResponse request
	HealthCheckWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*HealthCheckResponse, error)

	// ChatWebSocketWithResponse request
	ChatWebSocketWithResponse(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*ChatWebSocketResponse, error)
}

type GetTokenResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Token *string `json:"token,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r GetTokenResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetTokenResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ListChatRoomsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]ChatRoom
}

// Status returns HTTPResponse.Status
func (r ListChatRoomsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListChatRoomsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ListChatMessagesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]ChatMessage
}

// Status returns HTTPResponse.Status
func (r ListChatMessagesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListChatMessagesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ListEventsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]EventsWithID
}

// Status returns HTTPResponse.Status
func (r ListEventsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListEventsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type HealthCheckResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r HealthCheckResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r HealthCheckResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ChatWebSocketResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r ChatWebSocketResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ChatWebSocketResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetTokenWithBodyWithResponse request with arbitrary body returning *GetTokenResponse
func (c *ClientWithResponses) GetTokenWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*GetTokenResponse, error) {
	rsp, err := c.GetTokenWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetTokenResponse(rsp)
}

func (c *ClientWithResponses) GetTokenWithFormdataBodyWithResponse(ctx context.Context, body GetTokenFormdataRequestBody, reqEditors ...RequestEditorFn) (*GetTokenResponse, error) {
	rsp, err := c.GetTokenWithFormdataBody(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetTokenResponse(rsp)
}

// ListChatRoomsWithResponse request returning *ListChatRoomsResponse
func (c *ClientWithResponses) ListChatRoomsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListChatRoomsResponse, error) {
	rsp, err := c.ListChatRooms(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListChatRoomsResponse(rsp)
}

// ListChatMessagesWithResponse request returning *ListChatMessagesResponse
func (c *ClientWithResponses) ListChatMessagesWithResponse(ctx context.Context, roomId string, reqEditors ...RequestEditorFn) (*ListChatMessagesResponse, error) {
	rsp, err := c.ListChatMessages(ctx, roomId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListChatMessagesResponse(rsp)
}

// ListEventsWithResponse request returning *ListEventsResponse
func (c *ClientWithResponses) ListEventsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListEventsResponse, error) {
	rsp, err := c.ListEvents(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListEventsResponse(rsp)
}

// HealthCheckWithResponse request returning *HealthCheckResponse
func (c *ClientWithResponses) HealthCheckWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*HealthCheckResponse, error) {
	rsp, err := c.HealthCheck(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseHealthCheckResponse(rsp)
}

// ChatWebSocketWithResponse request returning *ChatWebSocketResponse
func (c *ClientWithResponses) ChatWebSocketWithResponse(ctx context.Context, id string, reqEditors ...RequestEditorFn) (*ChatWebSocketResponse, error) {
	rsp, err := c.ChatWebSocket(ctx, id, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseChatWebSocketResponse(rsp)
}

// ParseGetTokenResponse parses an HTTP response from a GetTokenWithResponse call
func ParseGetTokenResponse(rsp *http.Response) (*GetTokenResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetTokenResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Token *string `json:"token,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseListChatRoomsResponse parses an HTTP response from a ListChatRoomsWithResponse call
func ParseListChatRoomsResponse(rsp *http.Response) (*ListChatRoomsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListChatRoomsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []ChatRoom
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseListChatMessagesResponse parses an HTTP response from a ListChatMessagesWithResponse call
func ParseListChatMessagesResponse(rsp *http.Response) (*ListChatMessagesResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListChatMessagesResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []ChatMessage
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseListEventsResponse parses an HTTP response from a ListEventsWithResponse call
func ParseListEventsResponse(rsp *http.Response) (*ListEventsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListEventsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []EventsWithID
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseHealthCheckResponse parses an HTTP response from a HealthCheckWithResponse call
func ParseHealthCheckResponse(rsp *http.Response) (*HealthCheckResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &HealthCheckResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseChatWebSocketResponse parses an HTTP response from a ChatWebSocketWithResponse call
func ParseChatWebSocketResponse(rsp *http.Response) (*ChatWebSocketResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ChatWebSocketResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /auth/token)
	GetToken(ctx echo.Context) error
	// List all chat rooms
	// (GET /chats)
	ListChatRooms(ctx echo.Context) error
	// List all messages in a chat room
	// (GET /chats/{roomId}/messages)
	ListChatMessages(ctx echo.Context, roomId string) error

	// (GET /events)
	ListEvents(ctx echo.Context) error

	// (GET /health)
	HealthCheck(ctx echo.Context) error
	// WebSocket endpoint for chat
	// (GET /ws/{id})
	ChatWebSocket(ctx echo.Context, id string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetToken converts echo context to params.
func (w *ServerInterfaceWrapper) GetToken(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetToken(ctx)
	return err
}

// ListChatRooms converts echo context to params.
func (w *ServerInterfaceWrapper) ListChatRooms(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListChatRooms(ctx)
	return err
}

// ListChatMessages converts echo context to params.
func (w *ServerInterfaceWrapper) ListChatMessages(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "roomId" -------------
	var roomId string

	err = runtime.BindStyledParameterWithOptions("simple", "roomId", ctx.Param("roomId"), &roomId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter roomId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListChatMessages(ctx, roomId)
	return err
}

// ListEvents converts echo context to params.
func (w *ServerInterfaceWrapper) ListEvents(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListEvents(ctx)
	return err
}

// HealthCheck converts echo context to params.
func (w *ServerInterfaceWrapper) HealthCheck(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.HealthCheck(ctx)
	return err
}

// ChatWebSocket converts echo context to params.
func (w *ServerInterfaceWrapper) ChatWebSocket(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", ctx.Param("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ChatWebSocket(ctx, id)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
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

	router.POST(baseURL+"/auth/token", wrapper.GetToken)
	router.GET(baseURL+"/chats", wrapper.ListChatRooms)
	router.GET(baseURL+"/chats/:roomId/messages", wrapper.ListChatMessages)
	router.GET(baseURL+"/events", wrapper.ListEvents)
	router.GET(baseURL+"/health", wrapper.HealthCheck)
	router.GET(baseURL+"/ws/:id", wrapper.ChatWebSocket)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7xWTW/bOBD9KwR3j4rltJdCN7fN7qbbokXjRQ6JgTLU2GIjkVpyZNcw/N8XQ+rDsuSs",
	"ixS5WeZw3sx788Edl6YojQaNjic77mQGhfA/32UCP4FzYgX0WVpTgkUF/lBaEAjpDOljaWwhkCc8FQgX",
	"qArgEcdtCTzhDq3SK76PuErJdvB30UEMzqwxxfX4tcqBHT2ia/BvpSykPLkj1NZPe6tDXewjn+hXY4ph",
	"lidC1qKAM5G9KYFcWWvsCI8m9a5ScNKqEpXRPAnGzJ9FHbtK4+tXHbNKI6zAHnE45qg5jv4n4hrwkJur",
	"dVMZ/bh7MCMUqUKs4B+bP4M/bxX1kA78dsHdKsyu35M/keeflzy52/HfLSx5wn+Lu+qO69KO65T20Vlq",
	"D0Vd7Bf0rwNZWYXbG3IbHLwFYcHOKszo68F//dGo9+F2zqPQX+Q/nHaSZIgl35NjpZdmqORMM/ghijIH",
	"NvtyzTaZkhmrHDgWPDE0j6CZk6YEx4RO2YfbORMUS8RRYU4gFBpoVJJ61/u5Cj55xNdgXYC6nEwnU9LJ",
	"lKBFqXjCX/u/Il4KzHyqMXmOPaYvDuNwGPNXwMpqx4SPJQS4NJYJtlJr0BS+JZV9uKVwbmNsOmHzTDkG",
	"Oi2N0nivUwOOaYOs1qHnDg0TUoJzk3sqD5JTEDqNBv4n4NxHGDQEh29Nug1tpxG0D1mUZU6EKKPjHxeb",
	"zeaCGu6isjloaoi0G4vDNmiCPjmiziv11jLqPC7a0jAP30FiKI6johin0LNyIDU/RENbgYd3pdEupPFq",
	"On2Cle8uNPkpEtoqGCZ5RgatlnS8j3gsMxEmzgqerKlcOWRmyUSeM7rDaMq7yaAKPiqHzYh3/JmZK4TC",
	"X3xqwrQLpWNAWCu24wQ0eXQ5eOpcVRTCbusEjrLkHVPxLmy3fVzP7Z/lrrnGlGaiwzhN5KcGh6rVigIQ",
	"rPNjtw83z4BdvycYzKBzTFOcjmmWNBsy6VZ0v1KjA+6Py2vxUlo2r6Cfl7OV5ISiJ7gP8kK7e89UM1w4",
	"GqDN2OyP4Y3CjGS5198siDQJV7+F/fH0VKXg6x36EgL0lvwZChzTwQ+Xta/SwzV9V2/zOAORh709SvaM",
	"OeWXb7BjMgP52HI84Ogvb/WOjE6Q1Hf/+e9m+m1cvFPp/iCQvmeqx1t4uDHyEfCXdqB6XvddTi+Hed1s",
	"FMpM6RX7Yg0aafLjVmhz6QqWnggUa1DXgV2PZ/fRSJGzFNaQm7IAjSzY0jufXp7+TZXEcU52mXGYvJm+",
	"oVfMWlglHvJ6hRtb670UVU4vtdpqSCWZUm/YSntCAxwzfnUt9v8FAAD//8b/OoNMDQAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
