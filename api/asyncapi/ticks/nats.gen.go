// Package "ticks" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package ticks

import (
	"errors"

	"github.com/nats-io/nats.go"
)

// NATSController is the NATS implementation for asyncapi-codegen
type NATSController struct {
	connection *nats.Conn
	logger     Logger
}

// NewNATSController creates a new NATSController that fulfill the BrokerLinker interface
func NewNATSController(connection *nats.Conn) *NATSController {
	return &NATSController{
		connection: connection,
	}
}

// AttachLogger attaches a logger that will log operations on broker controller
func (c NATSController) AttachLogger(logger Logger) {
	c.logger = logger
}

// Publish a message to the broker
func (c *NATSController) Publish(channel string, um UniversalMessage) error {
	msg := nats.NewMsg(channel)

	// Set message content
	msg.Data = um.Payload
	if um.CorrelationID != nil {
		msg.Header.Add(CorrelationIDField, *um.CorrelationID)
	}

	// Publish message
	if err := c.connection.PublishMsg(msg); err != nil {
		return err
	}

	// Flush the queue
	return c.connection.Flush()
}

// Subscribe to messages from the broker
func (c *NATSController) Subscribe(channel string) (msgs chan UniversalMessage, stop chan interface{}, err error) {
	// Subscribe to channel
	natsMsgs := make(chan *nats.Msg, 64)
	sub, err := c.connection.ChanSubscribe(channel, natsMsgs)
	if err != nil {
		return nil, nil, err
	}

	// Handle events
	msgs = make(chan UniversalMessage, 64)
	stop = make(chan interface{}, 1)
	go func() {
		for {
			select {
			// Handle new message
			case msg := <-natsMsgs:
				var correlationID *string

				// Add correlation ID if not empty
				str := msg.Header.Get(CorrelationIDField)
				if str != "" {
					correlationID = &str
				}

				// Create message
				msgs <- UniversalMessage{
					Payload:       msg.Data,
					CorrelationID: correlationID,
				}
			// Handle closure request from function caller
			case _ = <-stop:
				if err := sub.Unsubscribe(); err != nil && !errors.Is(err, nats.ErrConnectionClosed) && c.logger != nil {
					c.logger.Error(err.Error(), "module", "asyncapi", "controller", "nats")
				}

				if err := sub.Drain(); err != nil && !errors.Is(err, nats.ErrConnectionClosed) && c.logger != nil {
					c.logger.Error(err.Error(), "module", "asyncapi", "controller", "nats")
				}

				close(msgs)
			}
		}
	}()

	return msgs, stop, nil
}
