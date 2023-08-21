// Package "events" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version v0.15.0 DO NOT EDIT.
package events

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/log"
)

// ClientSubscriber represents all handlers that are expecting messages for Client
type ClientSubscriber interface {
	// CryptellationTicksListenExchangePair
	CryptellationTicksListenExchangePair(msg TickMessage, done bool)

	// CryptellationTicksRegisterResponse
	CryptellationTicksRegisterResponse(msg RegisteringResponseMessage, done bool)

	// CryptellationTicksUnregisterResponse
	CryptellationTicksUnregisterResponse(msg RegisteringResponseMessage, done bool)
}

// ClientController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the Client
type ClientController struct {
	brokerController BrokerController
	stopSubscribers  map[string]chan interface{}
	logger           log.Logger
}

// NewClientController links the Client to the broker
func NewClientController(bs BrokerController) (*ClientController, error) {
	if bs == nil {
		return nil, ErrNilBrokerController
	}

	return &ClientController{
		brokerController: bs,
		stopSubscribers:  make(map[string]chan interface{}),
		logger:           log.Silent{},
	}, nil
}

// SetLogger attaches a logger that will log operations on controller
func (c *ClientController) SetLogger(logger log.Logger) {
	c.logger = logger
	c.brokerController.SetLogger(logger)
}

// LogError logs error if the logger has been set
func (c ClientController) LogError(ctx log.Context, msg string) {
	// Add more context
	ctx.Module = "asyncapi"
	ctx.Provider = "client"

	// Log error
	c.logger.Error(ctx, msg)
}

// LogInfo logs information if the logger has been set
func (c ClientController) LogInfo(ctx log.Context, msg string) {
	// Add more context
	ctx.Module = "asyncapi"
	ctx.Provider = "client"

	// Log info
	c.logger.Info(ctx, msg)
}

// Close will clean up any existing resources on the controller
func (c *ClientController) Close() {
	// Unsubscribing remaining channels
	c.LogInfo(log.Context{}, "Closing Client controller")
	c.UnsubscribeAll()
}

// SubscribeAll will subscribe to channels without parameters on which the app is expecting messages.
// For channels with parameters, they should be subscribed independently.
func (c *ClientController) SubscribeAll(as ClientSubscriber) error {
	if as == nil {
		return ErrNilClientSubscriber
	}

	if err := c.SubscribeCryptellationTicksRegisterResponse(as.CryptellationTicksRegisterResponse); err != nil {
		return err
	}
	if err := c.SubscribeCryptellationTicksUnregisterResponse(as.CryptellationTicksUnregisterResponse); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *ClientController) UnsubscribeAll() {
	// Unsubscribe channels with no parameters (if any)
	c.UnsubscribeCryptellationTicksRegisterResponse()
	c.UnsubscribeCryptellationTicksUnregisterResponse()

	// Unsubscribe remaining channels
	for n, stopChan := range c.stopSubscribers {
		stopChan <- true
		delete(c.stopSubscribers, n)
	}
}

// SubscribeCryptellationTicksListenExchangePair will subscribe to new messages from 'cryptellation.ticks.listen.{exchange}.{pair}' channel.
//
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *ClientController) SubscribeCryptellationTicksListenExchangePair(params CryptellationTicksListenExchangePairParameters, fn func(msg TickMessage, done bool)) error {
	// Get channel path
	path := fmt.Sprintf("cryptellation.ticks.listen.%v.%v", params.Exchange, params.Pair)

	// Check if there is already a subscription
	_, exists := c.stopSubscribers[path]
	if exists {
		err := fmt.Errorf("%w: %q channel is already subscribed", ErrAlreadySubscribedChannel, path)
		c.LogError(log.Context{Action: path}, err.Error())
		return err
	}

	// Subscribe to broker channel
	c.LogInfo(log.Context{Action: path, Operation: "subscribe"}, "Subscribing to channel")
	msgs, stop, err := c.brokerController.Subscribe(path)
	if err != nil {
		c.LogError(log.Context{Action: path, Operation: "subscribe"}, err.Error())
		return err
	}

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			um, open := <-msgs

			// Process message
			msg, err := newTickMessageFromUniversalMessage(um)
			if err != nil {
				c.LogError(log.Context{Action: path, Operation: "subscribe", Message: msg}, err.Error())
			}

			// Send info if message is correct or susbcription is closed
			if err == nil || !open {
				c.LogInfo(log.Context{Action: path, Operation: "subscribe", Message: msg}, "Received new message")
				fn(msg, !open)
			}

			// If subscription is closed, then exit the function
			if !open {
				return
			}
		}
	}()

	// Add the stop channel to the inside map
	c.stopSubscribers[path] = stop

	return nil
}

// UnsubscribeCryptellationTicksListenExchangePair will unsubscribe messages from 'cryptellation.ticks.listen.{exchange}.{pair}' channel
func (c *ClientController) UnsubscribeCryptellationTicksListenExchangePair(params CryptellationTicksListenExchangePairParameters) {
	// Get channel path
	path := fmt.Sprintf("cryptellation.ticks.listen.%v.%v", params.Exchange, params.Pair)

	// Get stop channel
	stopChan, exists := c.stopSubscribers[path]
	if !exists {
		return
	}

	// Stop the channel and remove the entry
	c.LogInfo(log.Context{Action: path, Operation: "unsubscribe"}, "Unsubscribing from channel")
	stopChan <- true
	delete(c.stopSubscribers, path)
} // SubscribeCryptellationTicksRegisterResponse will subscribe to new messages from 'cryptellation.ticks.register.response' channel.
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *ClientController) SubscribeCryptellationTicksRegisterResponse(fn func(msg RegisteringResponseMessage, done bool)) error {
	// Get channel path
	path := "cryptellation.ticks.register.response"

	// Check if there is already a subscription
	_, exists := c.stopSubscribers[path]
	if exists {
		err := fmt.Errorf("%w: %q channel is already subscribed", ErrAlreadySubscribedChannel, path)
		c.LogError(log.Context{Action: path}, err.Error())
		return err
	}

	// Subscribe to broker channel
	c.LogInfo(log.Context{Action: path, Operation: "subscribe"}, "Subscribing to channel")
	msgs, stop, err := c.brokerController.Subscribe(path)
	if err != nil {
		c.LogError(log.Context{Action: path, Operation: "subscribe"}, err.Error())
		return err
	}

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			um, open := <-msgs

			// Process message
			msg, err := newRegisteringResponseMessageFromUniversalMessage(um)
			if err != nil {
				c.LogError(log.Context{Action: path, Operation: "subscribe", Message: msg}, err.Error())
			}

			// Send info if message is correct or susbcription is closed
			if err == nil || !open {
				c.LogInfo(log.Context{Action: path, Operation: "subscribe", Message: msg}, "Received new message")
				fn(msg, !open)
			}

			// If subscription is closed, then exit the function
			if !open {
				return
			}
		}
	}()

	// Add the stop channel to the inside map
	c.stopSubscribers[path] = stop

	return nil
}

// UnsubscribeCryptellationTicksRegisterResponse will unsubscribe messages from 'cryptellation.ticks.register.response' channel
func (c *ClientController) UnsubscribeCryptellationTicksRegisterResponse() {
	// Get channel path
	path := "cryptellation.ticks.register.response"

	// Get stop channel
	stopChan, exists := c.stopSubscribers[path]
	if !exists {
		return
	}

	// Stop the channel and remove the entry
	c.LogInfo(log.Context{Action: path, Operation: "unsubscribe"}, "Unsubscribing from channel")
	stopChan <- true
	delete(c.stopSubscribers, path)
} // SubscribeCryptellationTicksUnregisterResponse will subscribe to new messages from 'cryptellation.ticks.unregister.response' channel.
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *ClientController) SubscribeCryptellationTicksUnregisterResponse(fn func(msg RegisteringResponseMessage, done bool)) error {
	// Get channel path
	path := "cryptellation.ticks.unregister.response"

	// Check if there is already a subscription
	_, exists := c.stopSubscribers[path]
	if exists {
		err := fmt.Errorf("%w: %q channel is already subscribed", ErrAlreadySubscribedChannel, path)
		c.LogError(log.Context{Action: path}, err.Error())
		return err
	}

	// Subscribe to broker channel
	c.LogInfo(log.Context{Action: path, Operation: "subscribe"}, "Subscribing to channel")
	msgs, stop, err := c.brokerController.Subscribe(path)
	if err != nil {
		c.LogError(log.Context{Action: path, Operation: "subscribe"}, err.Error())
		return err
	}

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			um, open := <-msgs

			// Process message
			msg, err := newRegisteringResponseMessageFromUniversalMessage(um)
			if err != nil {
				c.LogError(log.Context{Action: path, Operation: "subscribe", Message: msg}, err.Error())
			}

			// Send info if message is correct or susbcription is closed
			if err == nil || !open {
				c.LogInfo(log.Context{Action: path, Operation: "subscribe", Message: msg}, "Received new message")
				fn(msg, !open)
			}

			// If subscription is closed, then exit the function
			if !open {
				return
			}
		}
	}()

	// Add the stop channel to the inside map
	c.stopSubscribers[path] = stop

	return nil
}

// UnsubscribeCryptellationTicksUnregisterResponse will unsubscribe messages from 'cryptellation.ticks.unregister.response' channel
func (c *ClientController) UnsubscribeCryptellationTicksUnregisterResponse() {
	// Get channel path
	path := "cryptellation.ticks.unregister.response"

	// Get stop channel
	stopChan, exists := c.stopSubscribers[path]
	if !exists {
		return
	}

	// Stop the channel and remove the entry
	c.LogInfo(log.Context{Action: path, Operation: "unsubscribe"}, "Unsubscribing from channel")
	stopChan <- true
	delete(c.stopSubscribers, path)
}

// PublishCryptellationTicksRegisterRequest will publish messages to 'cryptellation.ticks.register.request' channel
func (c *ClientController) PublishCryptellationTicksRegisterRequest(msg RegisteringRequestMessage) error {
	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Get channel path
	path := "cryptellation.ticks.register.request"

	// Publish on event broker
	c.LogInfo(log.Context{Action: path, Operation: "publish", Message: msg}, "Publishing to channel")
	return c.brokerController.Publish(path, um)
}

// PublishCryptellationTicksUnregisterRequest will publish messages to 'cryptellation.ticks.unregister.request' channel
func (c *ClientController) PublishCryptellationTicksUnregisterRequest(msg RegisteringRequestMessage) error {
	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Get channel path
	path := "cryptellation.ticks.unregister.request"

	// Publish on event broker
	c.LogInfo(log.Context{Action: path, Operation: "publish", Message: msg}, "Publishing to channel")
	return c.brokerController.Publish(path, um)
}

// WaitForCryptellationTicksRegisterResponse will wait for a specific message by its correlation ID
//
// The pub function is the publication function that should be used to send the message
// It will be called after subscribing to the channel to avoid race condition, and potentially loose the message
func (cc *ClientController) WaitForCryptellationTicksRegisterResponse(ctx context.Context, publishMsg MessageWithCorrelationID, pub func() error) (RegisteringResponseMessage, error) {
	// Get channel path
	path := "cryptellation.ticks.register.response"

	// Subscribe to broker channel
	cc.LogInfo(log.Context{Action: path, Operation: "wait-for", CorrelationID: publishMsg.CorrelationID()}, "Wait for response")
	msgs, stop, err := cc.brokerController.Subscribe(path)
	if err != nil {
		cc.LogError(log.Context{Action: path, Operation: "wait-for"}, err.Error())
		return RegisteringResponseMessage{}, err
	}

	// Close subscriber on leave
	defer func() { stop <- true }()

	// Execute publication
	cc.LogInfo(log.Context{Action: path, Operation: "wait-for", Message: publishMsg, CorrelationID: publishMsg.CorrelationID()},
		"Sending request",
	)
	if err := pub(); err != nil {
		return RegisteringResponseMessage{}, err
	}

	// Wait for corresponding response
	for {
		select {
		case um, open := <-msgs:
			// Get new message
			msg, err := newRegisteringResponseMessageFromUniversalMessage(um)
			if err != nil {
				cc.LogError(log.Context{Action: path, Operation: "wait-for"}, err.Error())
			}

			// If valid message with corresponding correlation ID, return message
			if err == nil && publishMsg.CorrelationID() == msg.CorrelationID() {
				cc.LogInfo(log.Context{Action: path, Operation: "wait-for", Message: msg, CorrelationID: msg.CorrelationID()},
					"Received expected message",
				)
				return msg, nil
			} else if !open { // If message is invalid or not corresponding and the subscription is closed, then return error
				cc.LogError(log.Context{Action: path, Operation: "wait-for", CorrelationID: publishMsg.CorrelationID()},
					"Channel closed before getting message",
				)
				return RegisteringResponseMessage{}, ErrSubscriptionCanceled
			}
		case <-ctx.Done(): // Return error if context is done
			cc.LogError(log.Context{Action: path, Operation: "wait-for", CorrelationID: publishMsg.CorrelationID()},
				"Context done before getting message",
			)
			return RegisteringResponseMessage{}, ErrContextCanceled
		}
	}
}

// WaitForCryptellationTicksUnregisterResponse will wait for a specific message by its correlation ID
//
// The pub function is the publication function that should be used to send the message
// It will be called after subscribing to the channel to avoid race condition, and potentially loose the message
func (cc *ClientController) WaitForCryptellationTicksUnregisterResponse(ctx context.Context, publishMsg MessageWithCorrelationID, pub func() error) (RegisteringResponseMessage, error) {
	// Get channel path
	path := "cryptellation.ticks.unregister.response"

	// Subscribe to broker channel
	cc.LogInfo(log.Context{Action: path, Operation: "wait-for", CorrelationID: publishMsg.CorrelationID()}, "Wait for response")
	msgs, stop, err := cc.brokerController.Subscribe(path)
	if err != nil {
		cc.LogError(log.Context{Action: path, Operation: "wait-for"}, err.Error())
		return RegisteringResponseMessage{}, err
	}

	// Close subscriber on leave
	defer func() { stop <- true }()

	// Execute publication
	cc.LogInfo(log.Context{Action: path, Operation: "wait-for", Message: publishMsg, CorrelationID: publishMsg.CorrelationID()},
		"Sending request",
	)
	if err := pub(); err != nil {
		return RegisteringResponseMessage{}, err
	}

	// Wait for corresponding response
	for {
		select {
		case um, open := <-msgs:
			// Get new message
			msg, err := newRegisteringResponseMessageFromUniversalMessage(um)
			if err != nil {
				cc.LogError(log.Context{Action: path, Operation: "wait-for"}, err.Error())
			}

			// If valid message with corresponding correlation ID, return message
			if err == nil && publishMsg.CorrelationID() == msg.CorrelationID() {
				cc.LogInfo(log.Context{Action: path, Operation: "wait-for", Message: msg, CorrelationID: msg.CorrelationID()},
					"Received expected message",
				)
				return msg, nil
			} else if !open { // If message is invalid or not corresponding and the subscription is closed, then return error
				cc.LogError(log.Context{Action: path, Operation: "wait-for", CorrelationID: publishMsg.CorrelationID()},
					"Channel closed before getting message",
				)
				return RegisteringResponseMessage{}, ErrSubscriptionCanceled
			}
		case <-ctx.Done(): // Return error if context is done
			cc.LogError(log.Context{Action: path, Operation: "wait-for", CorrelationID: publishMsg.CorrelationID()},
				"Context done before getting message",
			)
			return RegisteringResponseMessage{}, ErrContextCanceled
		}
	}
}
