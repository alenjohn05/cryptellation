// Package "candlesticks" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package candlesticks

import (
	"fmt"
)

// AppSubscriber represents all handlers that are expecting messages for App
type AppSubscriber interface {
	// CryptellationCandlesticksListRequest
	CryptellationCandlesticksListRequest(msg CandlesticksListRequestMessage, done bool)
}

// AppController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the App
type AppController struct {
	brokerController BrokerController
	stopSubscribers  map[string]chan interface{}
	logger           Logger
}

// NewAppController links the App to the broker
func NewAppController(bs BrokerController) (*AppController, error) {
	if bs == nil {
		return nil, ErrNilBrokerController
	}

	return &AppController{
		brokerController: bs,
		stopSubscribers:  make(map[string]chan interface{}),
	}, nil
}

// AttachLogger attaches a logger that will log operations on controller
func (c *AppController) AttachLogger(logger Logger) {
	c.logger = logger
	c.brokerController.AttachLogger(logger)
}

// logError logs error if the logger has been set
func (c AppController) logError(msg string, keyvals ...interface{}) {
	if c.logger != nil {
		keyvals = append(keyvals, "module", "asyncapi", "controller", "App")
		c.logger.Error(msg, keyvals...)
	}
}

// logInfo logs information if the logger has been set
func (c AppController) logInfo(msg string, keyvals ...interface{}) {
	if c.logger != nil {
		keyvals = append(keyvals, "module", "asyncapi", "controller", "App")
		c.logger.Info(msg, keyvals...)
	}
}

// Close will clean up any existing resources on the controller
func (c *AppController) Close() {
	// Unsubscribing remaining channels
	c.logInfo("Closing App controller")
	c.UnsubscribeAll()
}

// SubscribeAll will subscribe to channels without parameters on which the app is expecting messages.
// For channels with parameters, they should be subscribed independently.
func (c *AppController) SubscribeAll(as AppSubscriber) error {
	if as == nil {
		return ErrNilAppSubscriber
	}

	if err := c.SubscribeCryptellationCandlesticksListRequest(as.CryptellationCandlesticksListRequest); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *AppController) UnsubscribeAll() {
	// Unsubscribe channels with no parameters (if any)
	c.UnsubscribeCryptellationCandlesticksListRequest()

	// Unsubscribe remaining channels
	for n, stopChan := range c.stopSubscribers {
		stopChan <- true
		delete(c.stopSubscribers, n)
	}
}

// SubscribeCryptellationCandlesticksListRequest will subscribe to new messages from 'cryptellation.candlesticks.list.request' channel.
//
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *AppController) SubscribeCryptellationCandlesticksListRequest(fn func(msg CandlesticksListRequestMessage, done bool)) error {
	// Get channel path
	path := "cryptellation.candlesticks.list.request"

	// Check if there is already a subscription
	_, exists := c.stopSubscribers[path]
	if exists {
		err := fmt.Errorf("%w: %q channel is already subscribed", ErrAlreadySubscribedChannel, path)
		c.logError(err.Error(), "channel", path)
		return err
	}

	// Subscribe to broker channel
	c.logInfo("Subscribing to channel", "channel", path, "operation", "subscribe")
	msgs, stop, err := c.brokerController.Subscribe(path)
	if err != nil {
		c.logError(err.Error(), "channel", path, "operation", "subscribe")
		return err
	}

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			um, open := <-msgs

			// Process message
			msg, err := newCandlesticksListRequestMessageFromUniversalMessage(um)
			if err != nil {
				c.logError(err.Error(), "channel", path, "operation", "subscribe", "message", msg)
			}

			// Send info if message is correct or susbcription is closed
			if err == nil || !open {
				c.logInfo("Received new message", "channel", path, "operation", "subscribe", "message", msg)
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

// UnsubscribeCryptellationCandlesticksListRequest will unsubscribe messages from 'cryptellation.candlesticks.list.request' channel
func (c *AppController) UnsubscribeCryptellationCandlesticksListRequest() {
	// Get channel path
	path := "cryptellation.candlesticks.list.request"

	// Get stop channel
	stopChan, exists := c.stopSubscribers[path]
	if !exists {
		return
	}

	// Stop the channel and remove the entry
	c.logInfo("Unsubscribing from channel", "channel", path, "operation", "unsubscribe")
	stopChan <- true
	delete(c.stopSubscribers, path)
}

// PublishCryptellationCandlesticksListResponse will publish messages to 'cryptellation.candlesticks.list.response' channel
func (c *AppController) PublishCryptellationCandlesticksListResponse(msg CandlesticksListResponseMessage) error {
	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Get channel path
	path := "cryptellation.candlesticks.list.response"

	// Publish on event broker
	c.logInfo("Publishing to channel", "channel", path, "operation", "publish", "message", msg)
	return c.brokerController.Publish(path, um)
}
