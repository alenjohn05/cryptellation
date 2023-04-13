// Package "backtests" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version v0.13.1 DO NOT EDIT.
package backtests

import (
	"fmt"
)

// AppSubscriber represents all handlers that are expecting messages for App
type AppSubscriber interface {
	// CryptellationBacktestsAccountsListRequest
	CryptellationBacktestsAccountsListRequest(msg BacktestsAccountsListRequestMessage, done bool)

	// CryptellationBacktestsAdvanceRequest
	CryptellationBacktestsAdvanceRequest(msg BacktestsAdvanceRequestMessage, done bool)

	// CryptellationBacktestsCreateRequest
	CryptellationBacktestsCreateRequest(msg BacktestsCreateRequestMessage, done bool)

	// CryptellationBacktestsOrdersCreateRequest
	CryptellationBacktestsOrdersCreateRequest(msg BacktestsOrdersCreateRequestMessage, done bool)

	// CryptellationBacktestsOrdersListRequest
	CryptellationBacktestsOrdersListRequest(msg BacktestsOrdersListRequestMessage, done bool)

	// CryptellationBacktestsSubscribeRequest
	CryptellationBacktestsSubscribeRequest(msg BacktestsSubscribeRequestMessage, done bool)
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

	if err := c.SubscribeCryptellationBacktestsAccountsListRequest(as.CryptellationBacktestsAccountsListRequest); err != nil {
		return err
	}
	if err := c.SubscribeCryptellationBacktestsAdvanceRequest(as.CryptellationBacktestsAdvanceRequest); err != nil {
		return err
	}
	if err := c.SubscribeCryptellationBacktestsCreateRequest(as.CryptellationBacktestsCreateRequest); err != nil {
		return err
	}
	if err := c.SubscribeCryptellationBacktestsOrdersCreateRequest(as.CryptellationBacktestsOrdersCreateRequest); err != nil {
		return err
	}
	if err := c.SubscribeCryptellationBacktestsOrdersListRequest(as.CryptellationBacktestsOrdersListRequest); err != nil {
		return err
	}
	if err := c.SubscribeCryptellationBacktestsSubscribeRequest(as.CryptellationBacktestsSubscribeRequest); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *AppController) UnsubscribeAll() {
	// Unsubscribe channels with no parameters (if any)
	c.UnsubscribeCryptellationBacktestsAccountsListRequest()
	c.UnsubscribeCryptellationBacktestsAdvanceRequest()
	c.UnsubscribeCryptellationBacktestsCreateRequest()
	c.UnsubscribeCryptellationBacktestsOrdersCreateRequest()
	c.UnsubscribeCryptellationBacktestsOrdersListRequest()
	c.UnsubscribeCryptellationBacktestsSubscribeRequest()

	// Unsubscribe remaining channels
	for n, stopChan := range c.stopSubscribers {
		stopChan <- true
		delete(c.stopSubscribers, n)
	}
}

// SubscribeCryptellationBacktestsAccountsListRequest will subscribe to new messages from 'cryptellation.backtests.accounts.list.request' channel.
//
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *AppController) SubscribeCryptellationBacktestsAccountsListRequest(fn func(msg BacktestsAccountsListRequestMessage, done bool)) error {
	// Get channel path
	path := "cryptellation.backtests.accounts.list.request"

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
			msg, err := newBacktestsAccountsListRequestMessageFromUniversalMessage(um)
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

// UnsubscribeCryptellationBacktestsAccountsListRequest will unsubscribe messages from 'cryptellation.backtests.accounts.list.request' channel
func (c *AppController) UnsubscribeCryptellationBacktestsAccountsListRequest() {
	// Get channel path
	path := "cryptellation.backtests.accounts.list.request"

	// Get stop channel
	stopChan, exists := c.stopSubscribers[path]
	if !exists {
		return
	}

	// Stop the channel and remove the entry
	c.logInfo("Unsubscribing from channel", "channel", path, "operation", "unsubscribe")
	stopChan <- true
	delete(c.stopSubscribers, path)
} // SubscribeCryptellationBacktestsAdvanceRequest will subscribe to new messages from 'cryptellation.backtests.advance.request' channel.
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *AppController) SubscribeCryptellationBacktestsAdvanceRequest(fn func(msg BacktestsAdvanceRequestMessage, done bool)) error {
	// Get channel path
	path := "cryptellation.backtests.advance.request"

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
			msg, err := newBacktestsAdvanceRequestMessageFromUniversalMessage(um)
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

// UnsubscribeCryptellationBacktestsAdvanceRequest will unsubscribe messages from 'cryptellation.backtests.advance.request' channel
func (c *AppController) UnsubscribeCryptellationBacktestsAdvanceRequest() {
	// Get channel path
	path := "cryptellation.backtests.advance.request"

	// Get stop channel
	stopChan, exists := c.stopSubscribers[path]
	if !exists {
		return
	}

	// Stop the channel and remove the entry
	c.logInfo("Unsubscribing from channel", "channel", path, "operation", "unsubscribe")
	stopChan <- true
	delete(c.stopSubscribers, path)
} // SubscribeCryptellationBacktestsCreateRequest will subscribe to new messages from 'cryptellation.backtests.create.request' channel.
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *AppController) SubscribeCryptellationBacktestsCreateRequest(fn func(msg BacktestsCreateRequestMessage, done bool)) error {
	// Get channel path
	path := "cryptellation.backtests.create.request"

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
			msg, err := newBacktestsCreateRequestMessageFromUniversalMessage(um)
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

// UnsubscribeCryptellationBacktestsCreateRequest will unsubscribe messages from 'cryptellation.backtests.create.request' channel
func (c *AppController) UnsubscribeCryptellationBacktestsCreateRequest() {
	// Get channel path
	path := "cryptellation.backtests.create.request"

	// Get stop channel
	stopChan, exists := c.stopSubscribers[path]
	if !exists {
		return
	}

	// Stop the channel and remove the entry
	c.logInfo("Unsubscribing from channel", "channel", path, "operation", "unsubscribe")
	stopChan <- true
	delete(c.stopSubscribers, path)
} // SubscribeCryptellationBacktestsOrdersCreateRequest will subscribe to new messages from 'cryptellation.backtests.orders.create.request' channel.
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *AppController) SubscribeCryptellationBacktestsOrdersCreateRequest(fn func(msg BacktestsOrdersCreateRequestMessage, done bool)) error {
	// Get channel path
	path := "cryptellation.backtests.orders.create.request"

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
			msg, err := newBacktestsOrdersCreateRequestMessageFromUniversalMessage(um)
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

// UnsubscribeCryptellationBacktestsOrdersCreateRequest will unsubscribe messages from 'cryptellation.backtests.orders.create.request' channel
func (c *AppController) UnsubscribeCryptellationBacktestsOrdersCreateRequest() {
	// Get channel path
	path := "cryptellation.backtests.orders.create.request"

	// Get stop channel
	stopChan, exists := c.stopSubscribers[path]
	if !exists {
		return
	}

	// Stop the channel and remove the entry
	c.logInfo("Unsubscribing from channel", "channel", path, "operation", "unsubscribe")
	stopChan <- true
	delete(c.stopSubscribers, path)
} // SubscribeCryptellationBacktestsOrdersListRequest will subscribe to new messages from 'cryptellation.backtests.orders.list.request' channel.
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *AppController) SubscribeCryptellationBacktestsOrdersListRequest(fn func(msg BacktestsOrdersListRequestMessage, done bool)) error {
	// Get channel path
	path := "cryptellation.backtests.orders.list.request"

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
			msg, err := newBacktestsOrdersListRequestMessageFromUniversalMessage(um)
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

// UnsubscribeCryptellationBacktestsOrdersListRequest will unsubscribe messages from 'cryptellation.backtests.orders.list.request' channel
func (c *AppController) UnsubscribeCryptellationBacktestsOrdersListRequest() {
	// Get channel path
	path := "cryptellation.backtests.orders.list.request"

	// Get stop channel
	stopChan, exists := c.stopSubscribers[path]
	if !exists {
		return
	}

	// Stop the channel and remove the entry
	c.logInfo("Unsubscribing from channel", "channel", path, "operation", "unsubscribe")
	stopChan <- true
	delete(c.stopSubscribers, path)
} // SubscribeCryptellationBacktestsSubscribeRequest will subscribe to new messages from 'cryptellation.backtests.subscribe.request' channel.
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *AppController) SubscribeCryptellationBacktestsSubscribeRequest(fn func(msg BacktestsSubscribeRequestMessage, done bool)) error {
	// Get channel path
	path := "cryptellation.backtests.subscribe.request"

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
			msg, err := newBacktestsSubscribeRequestMessageFromUniversalMessage(um)
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

// UnsubscribeCryptellationBacktestsSubscribeRequest will unsubscribe messages from 'cryptellation.backtests.subscribe.request' channel
func (c *AppController) UnsubscribeCryptellationBacktestsSubscribeRequest() {
	// Get channel path
	path := "cryptellation.backtests.subscribe.request"

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

// PublishCryptellationBacktestsAccountsListResponse will publish messages to 'cryptellation.backtests.accounts.list.response' channel
func (c *AppController) PublishCryptellationBacktestsAccountsListResponse(msg BacktestsAccountsListResponseMessage) error {
	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Get channel path
	path := "cryptellation.backtests.accounts.list.response"

	// Publish on event broker
	c.logInfo("Publishing to channel", "channel", path, "operation", "publish", "message", msg)
	return c.brokerController.Publish(path, um)
}

// PublishCryptellationBacktestsAdvanceResponse will publish messages to 'cryptellation.backtests.advance.response' channel
func (c *AppController) PublishCryptellationBacktestsAdvanceResponse(msg BacktestsAdvanceResponseMessage) error {
	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Get channel path
	path := "cryptellation.backtests.advance.response"

	// Publish on event broker
	c.logInfo("Publishing to channel", "channel", path, "operation", "publish", "message", msg)
	return c.brokerController.Publish(path, um)
}

// PublishCryptellationBacktestsCreateResponse will publish messages to 'cryptellation.backtests.create.response' channel
func (c *AppController) PublishCryptellationBacktestsCreateResponse(msg BacktestsCreateResponseMessage) error {
	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Get channel path
	path := "cryptellation.backtests.create.response"

	// Publish on event broker
	c.logInfo("Publishing to channel", "channel", path, "operation", "publish", "message", msg)
	return c.brokerController.Publish(path, um)
}

// PublishCryptellationBacktestsEventsID will publish messages to 'cryptellation.backtests.events.{id}' channel
func (c *AppController) PublishCryptellationBacktestsEventsID(params CryptellationBacktestsEventsIDParameters, msg BacktestsEventMessage) error {
	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Get channel path
	path := fmt.Sprintf("cryptellation.backtests.events.%v", params.ID)

	// Publish on event broker
	c.logInfo("Publishing to channel", "channel", path, "operation", "publish", "message", msg)
	return c.brokerController.Publish(path, um)
}

// PublishCryptellationBacktestsOrdersCreateResponse will publish messages to 'cryptellation.backtests.orders.create.response' channel
func (c *AppController) PublishCryptellationBacktestsOrdersCreateResponse(msg BacktestsOrdersCreateResponseMessage) error {
	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Get channel path
	path := "cryptellation.backtests.orders.create.response"

	// Publish on event broker
	c.logInfo("Publishing to channel", "channel", path, "operation", "publish", "message", msg)
	return c.brokerController.Publish(path, um)
}

// PublishCryptellationBacktestsOrdersListResponse will publish messages to 'cryptellation.backtests.orders.list.response' channel
func (c *AppController) PublishCryptellationBacktestsOrdersListResponse(msg BacktestsOrdersListResponseMessage) error {
	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Get channel path
	path := "cryptellation.backtests.orders.list.response"

	// Publish on event broker
	c.logInfo("Publishing to channel", "channel", path, "operation", "publish", "message", msg)
	return c.brokerController.Publish(path, um)
}

// PublishCryptellationBacktestsSubscribeResponse will publish messages to 'cryptellation.backtests.subscribe.response' channel
func (c *AppController) PublishCryptellationBacktestsSubscribeResponse(msg BacktestsSubscribeResponseMessage) error {
	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Get channel path
	path := "cryptellation.backtests.subscribe.response"

	// Publish on event broker
	c.logInfo("Publishing to channel", "channel", path, "operation", "publish", "message", msg)
	return c.brokerController.Publish(path, um)
}
