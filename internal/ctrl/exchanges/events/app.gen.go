// Package "events" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version v0.15.0 DO NOT EDIT.
package events

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/log"
)

// AppSubscriber represents all handlers that are expecting messages for App
type AppSubscriber interface {
	// CryptellationExchangesListRequest
	CryptellationExchangesListRequest(msg ExchangesRequestMessage, done bool)
}

// AppController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the App
type AppController struct {
	brokerController BrokerController
	stopSubscribers  map[string]chan interface{}
	logger           log.Logger
}

// NewAppController links the App to the broker
func NewAppController(bs BrokerController) (*AppController, error) {
	if bs == nil {
		return nil, ErrNilBrokerController
	}

	return &AppController{
		brokerController: bs,
		stopSubscribers:  make(map[string]chan interface{}),
		logger:           log.Silent{},
	}, nil
}

// SetLogger attaches a logger that will log operations on controller
func (c *AppController) SetLogger(logger log.Logger) {
	c.logger = logger
	c.brokerController.SetLogger(logger)
}

// LogError logs error if the logger has been set
func (c AppController) LogError(ctx log.Context, msg string) {
	// Add more context
	ctx.Module = "asyncapi"
	ctx.Provider = "app"

	// Log error
	c.logger.Error(ctx, msg)
}

// LogInfo logs information if the logger has been set
func (c AppController) LogInfo(ctx log.Context, msg string) {
	// Add more context
	ctx.Module = "asyncapi"
	ctx.Provider = "app"

	// Log info
	c.logger.Info(ctx, msg)
}

// Close will clean up any existing resources on the controller
func (c *AppController) Close() {
	// Unsubscribing remaining channels
	c.LogInfo(log.Context{}, "Closing App controller")
	c.UnsubscribeAll()
}

// SubscribeAll will subscribe to channels without parameters on which the app is expecting messages.
// For channels with parameters, they should be subscribed independently.
func (c *AppController) SubscribeAll(as AppSubscriber) error {
	if as == nil {
		return ErrNilAppSubscriber
	}

	if err := c.SubscribeCryptellationExchangesListRequest(as.CryptellationExchangesListRequest); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *AppController) UnsubscribeAll() {
	// Unsubscribe channels with no parameters (if any)
	c.UnsubscribeCryptellationExchangesListRequest()

	// Unsubscribe remaining channels
	for n, stopChan := range c.stopSubscribers {
		stopChan <- true
		delete(c.stopSubscribers, n)
	}
}

// SubscribeCryptellationExchangesListRequest will subscribe to new messages from 'cryptellation.exchanges.list.request' channel.
//
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *AppController) SubscribeCryptellationExchangesListRequest(fn func(msg ExchangesRequestMessage, done bool)) error {
	// Get channel path
	path := "cryptellation.exchanges.list.request"

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
			msg, err := newExchangesRequestMessageFromUniversalMessage(um)
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

// UnsubscribeCryptellationExchangesListRequest will unsubscribe messages from 'cryptellation.exchanges.list.request' channel
func (c *AppController) UnsubscribeCryptellationExchangesListRequest() {
	// Get channel path
	path := "cryptellation.exchanges.list.request"

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

// PublishCryptellationExchangesListResponse will publish messages to 'cryptellation.exchanges.list.response' channel
func (c *AppController) PublishCryptellationExchangesListResponse(msg ExchangesResponseMessage) error {
	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Get channel path
	path := "cryptellation.exchanges.list.response"

	// Publish on event broker
	c.LogInfo(log.Context{Action: path, Operation: "publish", Message: msg}, "Publishing to channel")
	return c.brokerController.Publish(path, um)
}
