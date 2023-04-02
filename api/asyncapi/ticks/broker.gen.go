// Package "ticks" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package ticks

const (
	// CorrelationIDField is the name of the field that will contain the correlation ID
	CorrelationIDField = "correlation_id"
)

// UniversalMessage is a wrapper that will contain all information regarding a message
type UniversalMessage struct {
	CorrelationID *string
	Payload       []byte
}

// BrokerController represents the functions that should be implemented to connect
// the broker to the application or the client
type BrokerController interface {
	// Publish a message to the broker
	Publish(channel string, mw UniversalMessage) error

	// Subscribe to messages from the broker
	Subscribe(channel string) (msgs chan UniversalMessage, stop chan interface{}, err error)
}
