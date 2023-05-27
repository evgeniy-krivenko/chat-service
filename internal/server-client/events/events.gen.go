// Package clientevents provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.5-0.20230513000919-14548c7e7bbe DO NOT EDIT.
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

// Event defines model for Event.
type Event struct {
	EventId   types.EventID   `json:"eventId"`
	EventType EventType       `json:"eventType"`
	RequestId types.RequestID `json:"requestId"`
	union     json.RawMessage
}

// EventType defines model for EventType.
type EventType string

// MessageBlockedEvent defines model for MessageBlockedEvent.
type MessageBlockedEvent = MessageId

// MessageId defines model for MessageId.
type MessageId struct {
	MessageId types.MessageID `json:"messageId"`
}

// MessageSentEvent defines model for MessageSentEvent.
type MessageSentEvent = MessageId

// NewMessageEvent defines model for NewMessageEvent.
type NewMessageEvent struct {
	AuthorId  *types.UserID   `json:"authorId,omitempty"`
	Body      string          `json:"body"`
	CreatedAt time.Time       `json:"createdAt"`
	IsService bool            `json:"isService"`
	MessageId types.MessageID `json:"messageId"`
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

// AsMessageSentEvent returns the union data inside the Event as a MessageSentEvent
func (t Event) AsMessageSentEvent() (MessageSentEvent, error) {
	var body MessageSentEvent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromMessageSentEvent overwrites any union data inside the Event as the provided MessageSentEvent
func (t *Event) FromMessageSentEvent(v MessageSentEvent) error {
	t.EventType = "MessageSentEvent"

	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeMessageSentEvent performs a merge with any union data inside the Event, using the provided MessageSentEvent
func (t *Event) MergeMessageSentEvent(v MessageSentEvent) error {
	t.EventType = "MessageSentEvent"

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsMessageBlockedEvent returns the union data inside the Event as a MessageBlockedEvent
func (t Event) AsMessageBlockedEvent() (MessageBlockedEvent, error) {
	var body MessageBlockedEvent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromMessageBlockedEvent overwrites any union data inside the Event as the provided MessageBlockedEvent
func (t *Event) FromMessageBlockedEvent(v MessageBlockedEvent) error {
	t.EventType = "MessageBlockedEvent"

	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeMessageBlockedEvent performs a merge with any union data inside the Event, using the provided MessageBlockedEvent
func (t *Event) MergeMessageBlockedEvent(v MessageBlockedEvent) error {
	t.EventType = "MessageBlockedEvent"

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
	case "MessageBlockedEvent":
		return t.AsMessageBlockedEvent()
	case "MessageSentEvent":
		return t.AsMessageSentEvent()
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

	"H4sIAAAAAAAC/9xWTY/aPBD+K2je9+gQVr2sfOvucthDt1JpTysOJhmIi79qO6EI5b9XY9JN+BBCSNvD",
	"nrDsmSfPRzxhB4XVzho0MQDfQSgq1CItpw2aSItShsJLLY2I1tOGFs5Js6LlFwxBrPBB2WKNZdcC/+U9",
	"at5B5udK2V+AGZp4TXdfx+AFN93uxc7jspaB89ahj9sXoRE4IO1/3zqkM2vw6xL46w7+97i8HvRy/Qn9",
	"KxsO3Grnb9QlpowS8+eSlkvrtSAT6lqWwCCSHg4heoqKwe9sZbNuk37COIE+Pw3PMqmd9Sl1J2IFHFYy",
	"VvViXFidY7NCI7fZ2ssGzdrmRSViFtA3ssBcmojeCJUncGhbNvCVXxY7HQbg8VeN4WZV37r2d9PVUZQe",
	"S+CvA5HsLY6hiHnLYDo0Ak2tqfH4FTpzF9jZ+zU/tqE9X0diD94WvS+61dnuGf/M2Z7uvFe4J//xdA0G",
	"4IeSdzKkdyCUumLC9nnTqDx0RNSxsv5WQ34E9O859ha23BLUySUtPIqI5ed4wLsUEbMoNcKZey3DbP+g",
	"AeDCWoXCwLHxPfywr+PTDw27+IkFfUxSctIsbcKWUdHpgzDr0ax25MfosRJx9KgkmjhK8QVg0KAP0hrg",
	"0NylD6ZDI5wEDp/Gd+MJsORhAG5qpRiQUehDyrtE+ifh4r79CRtU1mlC31cBg9or4LAJPM+VLYSqbIj8",
	"fnI/yTeBOP8JAAD//+j+D3GxCAAA",
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
