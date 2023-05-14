// Package clientevents provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.5-0.20230403173426-fd06f5aed350 DO NOT EDIT.
package clientevents

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/getkin/kin-openapi/openapi3"
)

// Defines values for EventType.
const (
	EventTypeMessageBlockedEvent EventType = "MessageBlockedEvent"
	EventTypeMessageSentEvent    EventType = "MessageSentEvent"
	EventTypeNewMessageEvent     EventType = "NewMessageEvent"
)

// CommonMessage defines model for CommonMessage.
type CommonMessage struct {
	EventID   types.EventID   `json:"eventId"`
	EventType EventType       `json:"eventType"`
	MessageID types.MessageID `json:"messageId"`
	RequestID types.RequestID `json:"requestId"`
}

// EventType defines model for EventType.
type EventType string

// Message defines model for Message.
type Message struct {
	union json.RawMessage
}

// MessageBlockedEvent defines model for MessageBlockedEvent.
type MessageBlockedEvent = CommonMessage

// MessageSentEvent defines model for MessageSentEvent.
type MessageSentEvent = CommonMessage

// NewMessageEvent defines model for NewMessageEvent.
type NewMessageEvent struct {
	AuthorID  *types.UserID   `json:"authorId,omitempty"`
	Body      string          `json:"body"`
	CreatedAt time.Time       `json:"createdAt"`
	EventID   types.EventID   `json:"eventId"`
	EventType EventType       `json:"eventType"`
	IsService bool            `json:"isService"`
	MessageID types.MessageID `json:"messageId"`
	RequestID types.RequestID `json:"requestId"`
}

// AsNewMessageEvent returns the union data inside the Message as a NewMessageEvent
func (t Message) AsNewMessageEvent() (NewMessageEvent, error) {
	var body NewMessageEvent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromNewMessageEvent overwrites any union data inside the Message as the provided NewMessageEvent
func (t *Message) FromNewMessageEvent(v NewMessageEvent) error {
	v.EventType = "NewMessageEvent"
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeNewMessageEvent performs a merge with any union data inside the Message, using the provided NewMessageEvent
func (t *Message) MergeNewMessageEvent(v NewMessageEvent) error {
	v.EventType = "NewMessageEvent"
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t Message) Discriminator() (string, error) {
	var discriminator struct {
		Discriminator string `json:"eventType"`
	}
	err := json.Unmarshal(t.union, &discriminator)
	return discriminator.Discriminator, err
}

func (t Message) ValueByDiscriminator() (interface{}, error) {
	discriminator, err := t.Discriminator()
	if err != nil {
		return nil, err
	}
	switch discriminator {
	case "NewMessageEvent":
		return t.AsNewMessageEvent()
	default:
		return nil, errors.New("unknown discriminator value: " + discriminator)
	}
}

func (t Message) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *Message) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RUzW7bPBB8FWO/70hZCXoJeMvfoYemQNOeghxoaW2xJrkquZJrGHr3gpRiK6paNGl6",
	"srA/w9nZWR+gIFuTQ8cB5AFCUaFV6fOarCX3AUNQG4yB2lONnjWmNLbo+H0ZP9fkrWKQ0DS6BAG8rxEk",
	"BPbabUDA92xDmVM2Bm9T281TdCiNP2E5k8u0rclzel5xBRI2mqtmtSzI5thu0Ol9tvW6RbelvKgUZwF9",
	"qwvMtWP0Tpk8gUPXiZ7z5/TkAf73uAYJ/+UnBfJh/Pz2WNgJsL0EL5910O4X085m33Zej98aDK/Y0qeh",
	"cZ75bPYtmQ/UtccS5MNobeJou/FaxoM+dgJux2tG19gIcoe7QfGUBvG0nnt0PAldGSq2WPbRx6lU3bEu",
	"4pc6FF5b7RSTH53J/q6XEsdWIocf1yAffu++KdXusZunJg+gjPkDxOe3PMY7Tf83YFPGr8US038Z1XBF",
	"/uUGvuz75v37JaD/l2e3onIfoX6yTeFRMZaX/GycUjFmrC3CjNN0uO8fGgGuiAwqB9M7OcGP+wY+JxvT",
	"6isWyVURQLs1JWzNJmavlNsu7ps66rG4rhQvro1Gx4u02QACWvRBkwMJ7XkydY1O1RokvFueL89AJA0D",
	"SNcYIyAKhT4kK5QYj6Xmvv0GWzRU24jeV4GAxhuQsAsyzw0VylQUWF6cXZzluxA5/wgAAP//guX8n7IG",
	"AAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
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
	var res = make(map[string]func() ([]byte, error))
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
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
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
