// Package "events" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version v0.24.3 DO NOT EDIT.
package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"

	"github.com/google/uuid"
)

// controller is the controller that will be used to communicate with the broker
// It will be used internally by AppController and UserController
type controller struct {
	// broker is the broker controller that will be used to communicate
	broker extensions.BrokerController
	// subscriptions is a map of all subscriptions
	subscriptions map[string]extensions.BrokerChannelSubscription
	// logger is the logger that will be used² to log operations on controller
	logger extensions.Logger
	// middlewares are the middlewares that will be executed when sending or
	// receiving messages
	middlewares []extensions.Middleware
}

// ControllerOption is the type of the options that can be passed
// when creating a new Controller
type ControllerOption func(controller *controller)

// WithLogger attaches a logger to the controller
func WithLogger(logger extensions.Logger) ControllerOption {
	return func(controller *controller) {
		controller.logger = logger
	}
}

// WithMiddlewares attaches middlewares that will be executed when sending or receiving messages
func WithMiddlewares(middlewares ...extensions.Middleware) ControllerOption {
	return func(controller *controller) {
		controller.middlewares = middlewares
	}
}

type MessageWithCorrelationID interface {
	CorrelationID() string
	SetCorrelationID(id string)
}

type Error struct {
	Channel string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("channel %q: err %v", e.Channel, e.Err)
}

// SmaRequestMessage is the message expected for 'SmaRequest' channel
type SmaRequestMessage struct {
	// Headers will be used to fill the message headers
	Headers struct {
		// Description: Correlation ID set by client
		CorrelationID *string `json:"correlation_id"`
	}

	// Payload will be inserted in the message payload
	Payload struct {
		// Description: Date-time for the newest candlestick (RFC3339)
		End *DateSchema `json:"end"`

		// Description: Requested candlesticks exchange name
		ExchangeName ExchangeNameSchema `json:"exchange_name"`

		// Description: The maximum candlesticks to retrieve (0 = unlimited)
		Limit LimitSchema `json:"limit"`

		// Description: Requested candlesticks pair symbol
		PairSymbol PairSymbolSchema `json:"pair_symbol"`

		// Description: Number of periods used
		PeriodNumber NumberOfPeriodsSchema `json:"period_number"`

		// Description: Requested candlesticks period symbol
		PeriodSymbol PeriodSymbolSchema `json:"period_symbol"`

		// Description: Type of price from a candlestick
		PriceType *PriceTypeSchema `json:"price_type"`

		// Description: Date-time for the oldest candlestick (RFC3339)
		Start *DateSchema `json:"start"`
	}
}

func NewSmaRequestMessage() SmaRequestMessage {
	var msg SmaRequestMessage

	// Set correlation ID
	u := uuid.New().String()
	msg.Headers.CorrelationID = &u

	return msg
}

// newSmaRequestMessageFromBrokerMessage will fill a new SmaRequestMessage with data from generic broker message
func newSmaRequestMessageFromBrokerMessage(bMsg extensions.BrokerMessage) (SmaRequestMessage, error) {
	var msg SmaRequestMessage

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(bMsg.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// Get each headers from broker message
	for k, v := range bMsg.Headers {
		switch {
		case k == "correlationId": // Retrieving CorrelationID header
			h := string(v)
			msg.Headers.CorrelationID = &h
		default:
			// TODO: log unknown error
		}
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from SmaRequestMessage data
func (msg SmaRequestMessage) toBrokerMessage() (extensions.BrokerMessage, error) {
	// TODO: implement checks on message

	// Marshal payload to JSON
	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return extensions.BrokerMessage{}, err
	}

	// Add each headers to broker message
	headers := make(map[string][]byte, 1)

	// Adding CorrelationID header
	if msg.Headers.CorrelationID != nil {
		headers["correlationId"] = []byte(*msg.Headers.CorrelationID)
	}

	return extensions.BrokerMessage{
		Headers: headers,
		Payload: payload,
	}, nil
}

// CorrelationID will give the correlation ID of the message, based on AsyncAPI spec
func (msg SmaRequestMessage) CorrelationID() string {
	if msg.Headers.CorrelationID != nil {
		return *msg.Headers.CorrelationID
	}

	return ""
}

// SetCorrelationID will set the correlation ID of the message, based on AsyncAPI spec
func (msg *SmaRequestMessage) SetCorrelationID(id string) {
	msg.Headers.CorrelationID = &id
}

// SetAsResponseFrom will correlate the message with the one passed in parameter.
// It will assign the 'req' message correlation ID to the message correlation ID,
// both specified in AsyncAPI spec.
func (msg *SmaRequestMessage) SetAsResponseFrom(req MessageWithCorrelationID) {
	id := req.CorrelationID()
	msg.Headers.CorrelationID = &id
}

// SmaResponseMessage is the message expected for 'SmaResponse' channel
type SmaResponseMessage struct {
	// Headers will be used to fill the message headers
	Headers struct {
		// Description: Correlation ID set by client on corresponding request
		CorrelationID *string `json:"correlation_id"`
	}

	// Payload will be inserted in the message payload
	Payload struct {
		// Description: A list of timed numbers
		Data *NumericTimeSerieSchema `json:"data"`

		// Description: Response to a failed call
		Error *ErrorSchema `json:"error"`
	}
}

func NewSmaResponseMessage() SmaResponseMessage {
	var msg SmaResponseMessage

	// Set correlation ID
	u := uuid.New().String()
	msg.Headers.CorrelationID = &u

	return msg
}

// newSmaResponseMessageFromBrokerMessage will fill a new SmaResponseMessage with data from generic broker message
func newSmaResponseMessageFromBrokerMessage(bMsg extensions.BrokerMessage) (SmaResponseMessage, error) {
	var msg SmaResponseMessage

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(bMsg.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// Get each headers from broker message
	for k, v := range bMsg.Headers {
		switch {
		case k == "correlationId": // Retrieving CorrelationID header
			h := string(v)
			msg.Headers.CorrelationID = &h
		default:
			// TODO: log unknown error
		}
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from SmaResponseMessage data
func (msg SmaResponseMessage) toBrokerMessage() (extensions.BrokerMessage, error) {
	// TODO: implement checks on message

	// Marshal payload to JSON
	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return extensions.BrokerMessage{}, err
	}

	// Add each headers to broker message
	headers := make(map[string][]byte, 1)

	// Adding CorrelationID header
	if msg.Headers.CorrelationID != nil {
		headers["correlationId"] = []byte(*msg.Headers.CorrelationID)
	}

	return extensions.BrokerMessage{
		Headers: headers,
		Payload: payload,
	}, nil
}

// CorrelationID will give the correlation ID of the message, based on AsyncAPI spec
func (msg SmaResponseMessage) CorrelationID() string {
	if msg.Headers.CorrelationID != nil {
		return *msg.Headers.CorrelationID
	}

	return ""
}

// SetCorrelationID will set the correlation ID of the message, based on AsyncAPI spec
func (msg *SmaResponseMessage) SetCorrelationID(id string) {
	msg.Headers.CorrelationID = &id
}

// SetAsResponseFrom will correlate the message with the one passed in parameter.
// It will assign the 'req' message correlation ID to the message correlation ID,
// both specified in AsyncAPI spec.
func (msg *SmaResponseMessage) SetAsResponseFrom(req MessageWithCorrelationID) {
	id := req.CorrelationID()
	msg.Headers.CorrelationID = &id
}

// DateSchema is a schema from the AsyncAPI specification required in messages
// Description: Date-Time format according to RFC3339
type DateSchema time.Time

// MarshalJSON will override the marshal as this is not a normal 'time.Time' type
func (t DateSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t))
}

// UnmarshalJSON will override the unmarshal as this is not a normal 'time.Time' type
func (t *DateSchema) UnmarshalJSON(data []byte) error {
	var timeFormat time.Time
	if err := json.Unmarshal(data, &timeFormat); err != nil {
		return err
	}

	*t = DateSchema(timeFormat)
	return nil
}

// ErrorSchema is a schema from the AsyncAPI specification required in messages
// Description: Response to a failed call
type ErrorSchema struct {
	// Description: Code to identify the error type, based on HTTP errors
	Code int64 `json:"code"`

	// Description: Main error reason
	Message string `json:"message"`
}

// ExchangeNameSchema is a schema from the AsyncAPI specification required in messages
// Description: Exchange name
type ExchangeNameSchema string

// LimitSchema is a schema from the AsyncAPI specification required in messages
// Description: The maximum quantity to retrieve (0 = unlimited)
type LimitSchema int32

// NumberOfPeriodsSchema is a schema from the AsyncAPI specification required in messages
// Description: Number of periods used
type NumberOfPeriodsSchema int32

// NumericTimeSerieSchema is a schema from the AsyncAPI specification required in messages
// Description: A list of timed numbers
type NumericTimeSerieSchema []struct {
	// Description: Date-Time format according to RFC3339
	Time DateSchema `json:"time"`

	// Description: Numerical value
	Value float64 `json:"value"`
}

// PairSymbolSchema is a schema from the AsyncAPI specification required in messages
// Description: Pair symbol
type PairSymbolSchema string

// PeriodSymbolSchema is a schema from the AsyncAPI specification required in messages
// Description: Period symbol
type PeriodSymbolSchema string

// PriceTypeSchema is a schema from the AsyncAPI specification required in messages
// Description: Type of price from a candlestick
type PriceTypeSchema string
