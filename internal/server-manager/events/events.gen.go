// Package managerevents provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.5-0.20230513000919-14548c7e7bbe DO NOT EDIT.
package managerevents

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
	EventTypeChatClosedEvent EventType = "ChatClosedEvent"
	EventTypeNewChatEvent    EventType = "NewChatEvent"
	EventTypeNewMessageEvent EventType = "NewMessageEvent"
)

// ChatClosedEvent defines model for ChatClosedEvent.
type ChatClosedEvent struct {
	CanTakeMoreProblems bool         `json:"canTakeMoreProblems"`
	ChatId              types.ChatID `json:"chatId"`
}

// ChatId defines model for ChatId.
type ChatId struct {
	ChatId types.ChatID `json:"chatId"`
}

// Event defines model for Event.
type Event struct {
	EventId   types.EventID   `json:"eventId"`
	EventType EventType       `json:"eventType"`
	RequestId types.RequestID `json:"requestId"`
	union     json.RawMessage
}

// EventType defines model for EventType.
type EventType string

// MessageId defines model for MessageId.
type MessageId struct {
	MessageId types.MessageID `json:"messageId"`
}

// NewChatEvent defines model for NewChatEvent.
type NewChatEvent struct {
	CanTakeMoreProblems bool         `json:"canTakeMoreProblems"`
	ChatId              types.ChatID `json:"chatId"`
	ClientId            types.UserID `json:"clientId"`
}

// NewMessageEvent defines model for NewMessageEvent.
type NewMessageEvent struct {
	AuthorId  types.UserID    `json:"authorId"`
	Body      string          `json:"body"`
	ChatId    types.ChatID    `json:"chatId"`
	CreatedAt time.Time       `json:"createdAt"`
	MessageId types.MessageID `json:"messageId"`
}

// AsNewChatEvent returns the union data inside the Event as a NewChatEvent
func (t Event) AsNewChatEvent() (NewChatEvent, error) {
	var body NewChatEvent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromNewChatEvent overwrites any union data inside the Event as the provided NewChatEvent
func (t *Event) FromNewChatEvent(v NewChatEvent) error {
	t.EventType = "NewChatEvent"

	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeNewChatEvent performs a merge with any union data inside the Event, using the provided NewChatEvent
func (t *Event) MergeNewChatEvent(v NewChatEvent) error {
	t.EventType = "NewChatEvent"

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsNewMessageEvent returns the union data inside the Event as a NewMessageEvent
func (t Event) AsNewMessageEvent() (NewMessageEvent, error) {
	var body NewMessageEvent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromNewMessageEvent overwrites any union data inside the Event as the provided NewMessageEvent
func (t *Event) FromNewMessageEvent(v NewMessageEvent) error {
	t.EventType = "NewMessageEvent"

	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeNewMessageEvent performs a merge with any union data inside the Event, using the provided NewMessageEvent
func (t *Event) MergeNewMessageEvent(v NewMessageEvent) error {
	t.EventType = "NewMessageEvent"

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsChatClosedEvent returns the union data inside the Event as a ChatClosedEvent
func (t Event) AsChatClosedEvent() (ChatClosedEvent, error) {
	var body ChatClosedEvent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromChatClosedEvent overwrites any union data inside the Event as the provided ChatClosedEvent
func (t *Event) FromChatClosedEvent(v ChatClosedEvent) error {
	t.EventType = "ChatClosedEvent"

	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeChatClosedEvent performs a merge with any union data inside the Event, using the provided ChatClosedEvent
func (t *Event) MergeChatClosedEvent(v ChatClosedEvent) error {
	t.EventType = "ChatClosedEvent"

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t Event) Discriminator() (string, error) {
	var discriminator struct {
		Discriminator string `json:"eventType"`
	}
	err := json.Unmarshal(t.union, &discriminator)
	return discriminator.Discriminator, err
}

func (t Event) ValueByDiscriminator() (interface{}, error) {
	discriminator, err := t.Discriminator()
	if err != nil {
		return nil, err
	}
	switch discriminator {
	case "ChatClosedEvent":
		return t.AsChatClosedEvent()
	case "NewChatEvent":
		return t.AsNewChatEvent()
	case "NewMessageEvent":
		return t.AsNewMessageEvent()
	default:
		return nil, errors.New("unknown discriminator value: " + discriminator)
	}
}

func (t Event) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	if err != nil {
		return nil, err
	}
	object := make(map[string]json.RawMessage)
	if t.union != nil {
		err = json.Unmarshal(b, &object)
		if err != nil {
			return nil, err
		}
	}

	object["eventId"], err = json.Marshal(t.EventId)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'eventId': %w", err)
	}

	object["eventType"], err = json.Marshal(t.EventType)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'eventType': %w", err)
	}

	object["requestId"], err = json.Marshal(t.RequestId)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'requestId': %w", err)
	}

	b, err = json.Marshal(object)
	return b, err
}

func (t *Event) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	if err != nil {
		return err
	}
	object := make(map[string]json.RawMessage)
	err = json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["eventId"]; found {
		err = json.Unmarshal(raw, &t.EventId)
		if err != nil {
			return fmt.Errorf("error reading 'eventId': %w", err)
		}
	}

	if raw, found := object["eventType"]; found {
		err = json.Unmarshal(raw, &t.EventType)
		if err != nil {
			return fmt.Errorf("error reading 'eventType': %w", err)
		}
	}

	if raw, found := object["requestId"]; found {
		err = json.Unmarshal(raw, &t.RequestId)
		if err != nil {
			return fmt.Errorf("error reading 'requestId': %w", err)
		}
	}

	return err
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xWPW/bMBD9K8a1I2Ul6BJwa5MMHZIWbToFHmjpLLHmV0lKrmHovxekZUm23DhJUSOT",
	"jSPf8d67x6M2kGlptELlHdANuKxEyeLf65L5a6Ed5rc1Kh9CTIgvC6CPG3hvcQEU3qU9PG2xaQB+zqEh",
	"GzBWG7SeY8yYMfXAlninLX61ei5QxrBfGwQKc60FMgVNQ8Dir4pbzIE+HkXNyA6l5z8x89DMGgLtwXR0",
	"bhdfaCuZBwpVxXPokjhvuSqAwO+k0EkbDD9uGnPeDJcSLo22UQ/DfAkUCu7Laj7NtEyxLlDxdbK0vEa1",
	"1Gk4O3Foa55hypVHq5hIY25oRlS3hQYuneY5d5nlkivmtQ0ByYwJ1R5r0d8bMtxG4B5XIfgkam9PhNyh",
	"c6zAU6i9bQ3ZNWN9z2TQFUP8IWjcENAKn+GovVKCr05sPqjgtF2H6jQzcuCfWPFrDRST/j8HkYGe9Gmi",
	"t0Phg+3QvZrVtxZ+rpvRkyRdO4YkuiuzEwJVJQPwhIvJ6A7NDuk3BFrIsdkih0sv13GX+Vw69uUGxfbH",
	"wLnmO4FM8H+4UT8c2vON5F2p5AUP0WhWPlfa3mhjdVnlS23fpmgE5jpfD5rd3523+/gSyCwyj/lHv1de",
	"zjwmnksc1TgyR4dvBSB9lzrmxywSEnG10FEx7kVY/cTUcvK9MoHfJDCf3DHFCrSTaCIHBGq0jmsFFOrL",
	"+HoaVMxwoPBhejm9ABJFcUBVJQSBwByti67LMXxGGL+F32CNQhuJyk+2u4BAZQVQWDmapkJnTJTaeXp1",
	"cXWZrlwo+k8AAAD//4mEBdssCgAA",
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
