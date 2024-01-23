// Package "asyncapi" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version v0.30.2 DO NOT EDIT.
package asyncapi

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"

	"github.com/google/uuid"
)

// AppSubscriber represents all handlers that are expecting messages for App
type AppSubscriber interface {
	// ListBacktestAccountsRequest subscribes to messages placed on the 'cryptellation.backtests.accounts.list.request' channel
	ListBacktestAccountsRequest(ctx context.Context, msg ListBacktestAccountsRequestMessage)

	// AdvanceBacktestRequest subscribes to messages placed on the 'cryptellation.backtests.advance.request' channel
	AdvanceBacktestRequest(ctx context.Context, msg AdvanceBacktestRequestMessage)

	// CreateBacktestRequest subscribes to messages placed on the 'cryptellation.backtests.create.request' channel
	CreateBacktestRequest(ctx context.Context, msg CreateBacktestRequestMessage)

	// CreateBacktestOrderRequest subscribes to messages placed on the 'cryptellation.backtests.orders.create.request' channel
	CreateBacktestOrderRequest(ctx context.Context, msg CreateBacktestOrderRequestMessage)

	// ListBacktestOrdersRequest subscribes to messages placed on the 'cryptellation.backtests.orders.list.request' channel
	ListBacktestOrdersRequest(ctx context.Context, msg ListBacktestOrdersRequestMessage)

	// ServiceInfoRequest subscribes to messages placed on the 'cryptellation.backtests.service.info.request' channel
	ServiceInfoRequest(ctx context.Context, msg ServiceInfoRequestMessage)

	// SubscribeBacktestRequest subscribes to messages placed on the 'cryptellation.backtests.subscribe.request' channel
	SubscribeBacktestRequest(ctx context.Context, msg SubscribeBacktestRequestMessage)
}

// AppController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the App
type AppController struct {
	controller
}

// NewAppController links the App to the broker
func NewAppController(bc extensions.BrokerController, options ...ControllerOption) (*AppController, error) {
	// Check if broker controller has been provided
	if bc == nil {
		return nil, extensions.ErrNilBrokerController
	}

	// Create default controller
	controller := controller{
		broker:        bc,
		subscriptions: make(map[string]extensions.BrokerChannelSubscription),
		logger:        extensions.DummyLogger{},
		middlewares:   make([]extensions.Middleware, 0),
	}

	// Apply options
	for _, option := range options {
		option(&controller)
	}

	return &AppController{controller: controller}, nil
}

func (c AppController) wrapMiddlewares(
	middlewares []extensions.Middleware,
	callback extensions.NextMiddleware,
) func(ctx context.Context, msg *extensions.BrokerMessage) error {
	var called bool

	// If there is no more middleware
	if len(middlewares) == 0 {
		return func(ctx context.Context, msg *extensions.BrokerMessage) error {
			// Call the callback if it exists and it has not been called already
			if callback != nil && !called {
				called = true
				return callback(ctx)
			}

			// Nil can be returned, as the callback has already been called
			return nil
		}
	}

	// Get the next function to call from next middlewares or callback
	next := c.wrapMiddlewares(middlewares[1:], callback)

	// Wrap middleware into a check function that will call execute the middleware
	// and call the next wrapped middleware if the returned function has not been
	// called already
	return func(ctx context.Context, msg *extensions.BrokerMessage) error {
		// Call the middleware and the following if it has not been done already
		if !called {
			// Create the next call with the context and the message
			nextWithArgs := func(ctx context.Context) error {
				return next(ctx, msg)
			}

			// Call the middleware and register it as already called
			called = true
			if err := middlewares[0](ctx, msg, nextWithArgs); err != nil {
				return err
			}

			// If next has already been called in middleware, it should not be executed again
			return nextWithArgs(ctx)
		}

		// Nil can be returned, as the next middleware has already been called
		return nil
	}
}

func (c AppController) executeMiddlewares(ctx context.Context, msg *extensions.BrokerMessage, callback extensions.NextMiddleware) error {
	// Wrap middleware to have 'next' function when calling them
	wrapped := c.wrapMiddlewares(c.middlewares, callback)

	// Execute wrapped middlewares
	return wrapped(ctx, msg)
}

func addAppContextValues(ctx context.Context, path string) context.Context {
	ctx = context.WithValue(ctx, extensions.ContextKeyIsVersion, "1.0.0")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsProvider, "app")
	return context.WithValue(ctx, extensions.ContextKeyIsChannel, path)
}

// Close will clean up any existing resources on the controller
func (c *AppController) Close(ctx context.Context) {
	// Unsubscribing remaining channels
	c.UnsubscribeAll(ctx)

	c.logger.Info(ctx, "Closed app controller")
}

// SubscribeAll will subscribe to channels without parameters on which the app is expecting messages.
// For channels with parameters, they should be subscribed independently.
func (c *AppController) SubscribeAll(ctx context.Context, as AppSubscriber) error {
	if as == nil {
		return extensions.ErrNilAppSubscriber
	}

	if err := c.SubscribeListBacktestAccountsRequest(ctx, as.ListBacktestAccountsRequest); err != nil {
		return err
	}
	if err := c.SubscribeAdvanceBacktestRequest(ctx, as.AdvanceBacktestRequest); err != nil {
		return err
	}
	if err := c.SubscribeCreateBacktestRequest(ctx, as.CreateBacktestRequest); err != nil {
		return err
	}
	if err := c.SubscribeCreateBacktestOrderRequest(ctx, as.CreateBacktestOrderRequest); err != nil {
		return err
	}
	if err := c.SubscribeListBacktestOrdersRequest(ctx, as.ListBacktestOrdersRequest); err != nil {
		return err
	}
	if err := c.SubscribeServiceInfoRequest(ctx, as.ServiceInfoRequest); err != nil {
		return err
	}
	if err := c.SubscribeSubscribeBacktestRequest(ctx, as.SubscribeBacktestRequest); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *AppController) UnsubscribeAll(ctx context.Context) {
	c.UnsubscribeListBacktestAccountsRequest(ctx)
	c.UnsubscribeAdvanceBacktestRequest(ctx)
	c.UnsubscribeCreateBacktestRequest(ctx)
	c.UnsubscribeCreateBacktestOrderRequest(ctx)
	c.UnsubscribeListBacktestOrdersRequest(ctx)
	c.UnsubscribeServiceInfoRequest(ctx)
	c.UnsubscribeSubscribeBacktestRequest(ctx)
}

// SubscribeListBacktestAccountsRequest will subscribe to new messages from 'cryptellation.backtests.accounts.list.request' channel.
//
// Callback function 'fn' will be called each time a new message is received.
func (c *AppController) SubscribeListBacktestAccountsRequest(ctx context.Context, fn func(ctx context.Context, msg ListBacktestAccountsRequestMessage)) error {
	// Get channel path
	path := "cryptellation.backtests.accounts.list.request"

	// Set context
	ctx = addAppContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "reception")

	// Check if there is already a subscription
	_, exists := c.subscriptions[path]
	if exists {
		err := fmt.Errorf("%w: %q channel is already subscribed", extensions.ErrAlreadySubscribedChannel, path)
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Subscribe to broker channel
	sub, err := c.broker.Subscribe(ctx, path)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}
	c.logger.Info(ctx, "Subscribed to channel")

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			brokerMsg, open := <-sub.MessagesChannel()

			// If subscription is closed and there is no more message
			// (i.e. uninitialized message), then exit the function
			if !open && brokerMsg.IsUninitialized() {
				return
			}

			// Set broker message to context
			ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

			// Execute middlewares before handling the message
			if err := c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
				// Process message
				msg, err := newListBacktestAccountsRequestMessageFromBrokerMessage(brokerMsg)
				if err != nil {
					return err
				}

				// Add correlation ID to context if it exists
				if id := msg.CorrelationID(); id != "" {
					ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, id)
				}

				// Execute the subscription function
				fn(ctx, msg)

				return nil
			}); err != nil {
				c.logger.Error(ctx, err.Error())
			}
		}
	}()

	// Add the cancel channel to the inside map
	c.subscriptions[path] = sub

	return nil
}

// UnsubscribeListBacktestAccountsRequest will unsubscribe messages from 'cryptellation.backtests.accounts.list.request' channel.
// A timeout can be set in context to avoid blocking operation, if needed.
func (c *AppController) UnsubscribeListBacktestAccountsRequest(ctx context.Context) {
	// Get channel path
	path := "cryptellation.backtests.accounts.list.request"

	// Check if there subscribers for this channel
	sub, exists := c.subscriptions[path]
	if !exists {
		return
	}

	// Set context
	ctx = addAppContextValues(ctx, path)

	// Stop the subscription
	sub.Cancel(ctx)

	// Remove if from the subscribers
	delete(c.subscriptions, path)

	c.logger.Info(ctx, "Unsubscribed from channel")
} // SubscribeAdvanceBacktestRequest will subscribe to new messages from 'cryptellation.backtests.advance.request' channel.
// Callback function 'fn' will be called each time a new message is received.
func (c *AppController) SubscribeAdvanceBacktestRequest(ctx context.Context, fn func(ctx context.Context, msg AdvanceBacktestRequestMessage)) error {
	// Get channel path
	path := "cryptellation.backtests.advance.request"

	// Set context
	ctx = addAppContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "reception")

	// Check if there is already a subscription
	_, exists := c.subscriptions[path]
	if exists {
		err := fmt.Errorf("%w: %q channel is already subscribed", extensions.ErrAlreadySubscribedChannel, path)
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Subscribe to broker channel
	sub, err := c.broker.Subscribe(ctx, path)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}
	c.logger.Info(ctx, "Subscribed to channel")

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			brokerMsg, open := <-sub.MessagesChannel()

			// If subscription is closed and there is no more message
			// (i.e. uninitialized message), then exit the function
			if !open && brokerMsg.IsUninitialized() {
				return
			}

			// Set broker message to context
			ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

			// Execute middlewares before handling the message
			if err := c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
				// Process message
				msg, err := newAdvanceBacktestRequestMessageFromBrokerMessage(brokerMsg)
				if err != nil {
					return err
				}

				// Add correlation ID to context if it exists
				if id := msg.CorrelationID(); id != "" {
					ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, id)
				}

				// Execute the subscription function
				fn(ctx, msg)

				return nil
			}); err != nil {
				c.logger.Error(ctx, err.Error())
			}
		}
	}()

	// Add the cancel channel to the inside map
	c.subscriptions[path] = sub

	return nil
}

// UnsubscribeAdvanceBacktestRequest will unsubscribe messages from 'cryptellation.backtests.advance.request' channel.
// A timeout can be set in context to avoid blocking operation, if needed.
func (c *AppController) UnsubscribeAdvanceBacktestRequest(ctx context.Context) {
	// Get channel path
	path := "cryptellation.backtests.advance.request"

	// Check if there subscribers for this channel
	sub, exists := c.subscriptions[path]
	if !exists {
		return
	}

	// Set context
	ctx = addAppContextValues(ctx, path)

	// Stop the subscription
	sub.Cancel(ctx)

	// Remove if from the subscribers
	delete(c.subscriptions, path)

	c.logger.Info(ctx, "Unsubscribed from channel")
} // SubscribeCreateBacktestRequest will subscribe to new messages from 'cryptellation.backtests.create.request' channel.
// Callback function 'fn' will be called each time a new message is received.
func (c *AppController) SubscribeCreateBacktestRequest(ctx context.Context, fn func(ctx context.Context, msg CreateBacktestRequestMessage)) error {
	// Get channel path
	path := "cryptellation.backtests.create.request"

	// Set context
	ctx = addAppContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "reception")

	// Check if there is already a subscription
	_, exists := c.subscriptions[path]
	if exists {
		err := fmt.Errorf("%w: %q channel is already subscribed", extensions.ErrAlreadySubscribedChannel, path)
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Subscribe to broker channel
	sub, err := c.broker.Subscribe(ctx, path)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}
	c.logger.Info(ctx, "Subscribed to channel")

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			brokerMsg, open := <-sub.MessagesChannel()

			// If subscription is closed and there is no more message
			// (i.e. uninitialized message), then exit the function
			if !open && brokerMsg.IsUninitialized() {
				return
			}

			// Set broker message to context
			ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

			// Execute middlewares before handling the message
			if err := c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
				// Process message
				msg, err := newCreateBacktestRequestMessageFromBrokerMessage(brokerMsg)
				if err != nil {
					return err
				}

				// Add correlation ID to context if it exists
				if id := msg.CorrelationID(); id != "" {
					ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, id)
				}

				// Execute the subscription function
				fn(ctx, msg)

				return nil
			}); err != nil {
				c.logger.Error(ctx, err.Error())
			}
		}
	}()

	// Add the cancel channel to the inside map
	c.subscriptions[path] = sub

	return nil
}

// UnsubscribeCreateBacktestRequest will unsubscribe messages from 'cryptellation.backtests.create.request' channel.
// A timeout can be set in context to avoid blocking operation, if needed.
func (c *AppController) UnsubscribeCreateBacktestRequest(ctx context.Context) {
	// Get channel path
	path := "cryptellation.backtests.create.request"

	// Check if there subscribers for this channel
	sub, exists := c.subscriptions[path]
	if !exists {
		return
	}

	// Set context
	ctx = addAppContextValues(ctx, path)

	// Stop the subscription
	sub.Cancel(ctx)

	// Remove if from the subscribers
	delete(c.subscriptions, path)

	c.logger.Info(ctx, "Unsubscribed from channel")
} // SubscribeCreateBacktestOrderRequest will subscribe to new messages from 'cryptellation.backtests.orders.create.request' channel.
// Callback function 'fn' will be called each time a new message is received.
func (c *AppController) SubscribeCreateBacktestOrderRequest(ctx context.Context, fn func(ctx context.Context, msg CreateBacktestOrderRequestMessage)) error {
	// Get channel path
	path := "cryptellation.backtests.orders.create.request"

	// Set context
	ctx = addAppContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "reception")

	// Check if there is already a subscription
	_, exists := c.subscriptions[path]
	if exists {
		err := fmt.Errorf("%w: %q channel is already subscribed", extensions.ErrAlreadySubscribedChannel, path)
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Subscribe to broker channel
	sub, err := c.broker.Subscribe(ctx, path)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}
	c.logger.Info(ctx, "Subscribed to channel")

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			brokerMsg, open := <-sub.MessagesChannel()

			// If subscription is closed and there is no more message
			// (i.e. uninitialized message), then exit the function
			if !open && brokerMsg.IsUninitialized() {
				return
			}

			// Set broker message to context
			ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

			// Execute middlewares before handling the message
			if err := c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
				// Process message
				msg, err := newCreateBacktestOrderRequestMessageFromBrokerMessage(brokerMsg)
				if err != nil {
					return err
				}

				// Add correlation ID to context if it exists
				if id := msg.CorrelationID(); id != "" {
					ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, id)
				}

				// Execute the subscription function
				fn(ctx, msg)

				return nil
			}); err != nil {
				c.logger.Error(ctx, err.Error())
			}
		}
	}()

	// Add the cancel channel to the inside map
	c.subscriptions[path] = sub

	return nil
}

// UnsubscribeCreateBacktestOrderRequest will unsubscribe messages from 'cryptellation.backtests.orders.create.request' channel.
// A timeout can be set in context to avoid blocking operation, if needed.
func (c *AppController) UnsubscribeCreateBacktestOrderRequest(ctx context.Context) {
	// Get channel path
	path := "cryptellation.backtests.orders.create.request"

	// Check if there subscribers for this channel
	sub, exists := c.subscriptions[path]
	if !exists {
		return
	}

	// Set context
	ctx = addAppContextValues(ctx, path)

	// Stop the subscription
	sub.Cancel(ctx)

	// Remove if from the subscribers
	delete(c.subscriptions, path)

	c.logger.Info(ctx, "Unsubscribed from channel")
} // SubscribeListBacktestOrdersRequest will subscribe to new messages from 'cryptellation.backtests.orders.list.request' channel.
// Callback function 'fn' will be called each time a new message is received.
func (c *AppController) SubscribeListBacktestOrdersRequest(ctx context.Context, fn func(ctx context.Context, msg ListBacktestOrdersRequestMessage)) error {
	// Get channel path
	path := "cryptellation.backtests.orders.list.request"

	// Set context
	ctx = addAppContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "reception")

	// Check if there is already a subscription
	_, exists := c.subscriptions[path]
	if exists {
		err := fmt.Errorf("%w: %q channel is already subscribed", extensions.ErrAlreadySubscribedChannel, path)
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Subscribe to broker channel
	sub, err := c.broker.Subscribe(ctx, path)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}
	c.logger.Info(ctx, "Subscribed to channel")

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			brokerMsg, open := <-sub.MessagesChannel()

			// If subscription is closed and there is no more message
			// (i.e. uninitialized message), then exit the function
			if !open && brokerMsg.IsUninitialized() {
				return
			}

			// Set broker message to context
			ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

			// Execute middlewares before handling the message
			if err := c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
				// Process message
				msg, err := newListBacktestOrdersRequestMessageFromBrokerMessage(brokerMsg)
				if err != nil {
					return err
				}

				// Add correlation ID to context if it exists
				if id := msg.CorrelationID(); id != "" {
					ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, id)
				}

				// Execute the subscription function
				fn(ctx, msg)

				return nil
			}); err != nil {
				c.logger.Error(ctx, err.Error())
			}
		}
	}()

	// Add the cancel channel to the inside map
	c.subscriptions[path] = sub

	return nil
}

// UnsubscribeListBacktestOrdersRequest will unsubscribe messages from 'cryptellation.backtests.orders.list.request' channel.
// A timeout can be set in context to avoid blocking operation, if needed.
func (c *AppController) UnsubscribeListBacktestOrdersRequest(ctx context.Context) {
	// Get channel path
	path := "cryptellation.backtests.orders.list.request"

	// Check if there subscribers for this channel
	sub, exists := c.subscriptions[path]
	if !exists {
		return
	}

	// Set context
	ctx = addAppContextValues(ctx, path)

	// Stop the subscription
	sub.Cancel(ctx)

	// Remove if from the subscribers
	delete(c.subscriptions, path)

	c.logger.Info(ctx, "Unsubscribed from channel")
} // SubscribeServiceInfoRequest will subscribe to new messages from 'cryptellation.backtests.service.info.request' channel.
// Callback function 'fn' will be called each time a new message is received.
func (c *AppController) SubscribeServiceInfoRequest(ctx context.Context, fn func(ctx context.Context, msg ServiceInfoRequestMessage)) error {
	// Get channel path
	path := "cryptellation.backtests.service.info.request"

	// Set context
	ctx = addAppContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "reception")

	// Check if there is already a subscription
	_, exists := c.subscriptions[path]
	if exists {
		err := fmt.Errorf("%w: %q channel is already subscribed", extensions.ErrAlreadySubscribedChannel, path)
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Subscribe to broker channel
	sub, err := c.broker.Subscribe(ctx, path)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}
	c.logger.Info(ctx, "Subscribed to channel")

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			brokerMsg, open := <-sub.MessagesChannel()

			// If subscription is closed and there is no more message
			// (i.e. uninitialized message), then exit the function
			if !open && brokerMsg.IsUninitialized() {
				return
			}

			// Set broker message to context
			ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

			// Execute middlewares before handling the message
			if err := c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
				// Process message
				msg, err := newServiceInfoRequestMessageFromBrokerMessage(brokerMsg)
				if err != nil {
					return err
				}

				// Add correlation ID to context if it exists
				if id := msg.CorrelationID(); id != "" {
					ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, id)
				}

				// Execute the subscription function
				fn(ctx, msg)

				return nil
			}); err != nil {
				c.logger.Error(ctx, err.Error())
			}
		}
	}()

	// Add the cancel channel to the inside map
	c.subscriptions[path] = sub

	return nil
}

// UnsubscribeServiceInfoRequest will unsubscribe messages from 'cryptellation.backtests.service.info.request' channel.
// A timeout can be set in context to avoid blocking operation, if needed.
func (c *AppController) UnsubscribeServiceInfoRequest(ctx context.Context) {
	// Get channel path
	path := "cryptellation.backtests.service.info.request"

	// Check if there subscribers for this channel
	sub, exists := c.subscriptions[path]
	if !exists {
		return
	}

	// Set context
	ctx = addAppContextValues(ctx, path)

	// Stop the subscription
	sub.Cancel(ctx)

	// Remove if from the subscribers
	delete(c.subscriptions, path)

	c.logger.Info(ctx, "Unsubscribed from channel")
} // SubscribeSubscribeBacktestRequest will subscribe to new messages from 'cryptellation.backtests.subscribe.request' channel.
// Callback function 'fn' will be called each time a new message is received.
func (c *AppController) SubscribeSubscribeBacktestRequest(ctx context.Context, fn func(ctx context.Context, msg SubscribeBacktestRequestMessage)) error {
	// Get channel path
	path := "cryptellation.backtests.subscribe.request"

	// Set context
	ctx = addAppContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "reception")

	// Check if there is already a subscription
	_, exists := c.subscriptions[path]
	if exists {
		err := fmt.Errorf("%w: %q channel is already subscribed", extensions.ErrAlreadySubscribedChannel, path)
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Subscribe to broker channel
	sub, err := c.broker.Subscribe(ctx, path)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}
	c.logger.Info(ctx, "Subscribed to channel")

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			brokerMsg, open := <-sub.MessagesChannel()

			// If subscription is closed and there is no more message
			// (i.e. uninitialized message), then exit the function
			if !open && brokerMsg.IsUninitialized() {
				return
			}

			// Set broker message to context
			ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

			// Execute middlewares before handling the message
			if err := c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
				// Process message
				msg, err := newSubscribeBacktestRequestMessageFromBrokerMessage(brokerMsg)
				if err != nil {
					return err
				}

				// Add correlation ID to context if it exists
				if id := msg.CorrelationID(); id != "" {
					ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, id)
				}

				// Execute the subscription function
				fn(ctx, msg)

				return nil
			}); err != nil {
				c.logger.Error(ctx, err.Error())
			}
		}
	}()

	// Add the cancel channel to the inside map
	c.subscriptions[path] = sub

	return nil
}

// UnsubscribeSubscribeBacktestRequest will unsubscribe messages from 'cryptellation.backtests.subscribe.request' channel.
// A timeout can be set in context to avoid blocking operation, if needed.
func (c *AppController) UnsubscribeSubscribeBacktestRequest(ctx context.Context) {
	// Get channel path
	path := "cryptellation.backtests.subscribe.request"

	// Check if there subscribers for this channel
	sub, exists := c.subscriptions[path]
	if !exists {
		return
	}

	// Set context
	ctx = addAppContextValues(ctx, path)

	// Stop the subscription
	sub.Cancel(ctx)

	// Remove if from the subscribers
	delete(c.subscriptions, path)

	c.logger.Info(ctx, "Unsubscribed from channel")
}

// PublishListBacktestAccountsResponse will publish messages to 'cryptellation.backtests.accounts.list.response' channel
func (c *AppController) PublishListBacktestAccountsResponse(ctx context.Context, msg ListBacktestAccountsResponseMessage) error {
	// Get channel path
	path := "cryptellation.backtests.accounts.list.response"

	// Set correlation ID if it does not exist
	if id := msg.CorrelationID(); id == "" {
		msg.SetCorrelationID(uuid.New().String())
	}

	// Set context
	ctx = addAppContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "publication")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, msg.CorrelationID())

	// Convert to BrokerMessage
	brokerMsg, err := msg.toBrokerMessage()
	if err != nil {
		return err
	}

	// Set broker message to context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

	// Publish the message on event-broker through middlewares
	return c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
		return c.broker.Publish(ctx, path, brokerMsg)
	})
}

// PublishAdvanceBacktestResponse will publish messages to 'cryptellation.backtests.advance.response' channel
func (c *AppController) PublishAdvanceBacktestResponse(ctx context.Context, msg AdvanceBacktestResponseMessage) error {
	// Get channel path
	path := "cryptellation.backtests.advance.response"

	// Set correlation ID if it does not exist
	if id := msg.CorrelationID(); id == "" {
		msg.SetCorrelationID(uuid.New().String())
	}

	// Set context
	ctx = addAppContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "publication")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, msg.CorrelationID())

	// Convert to BrokerMessage
	brokerMsg, err := msg.toBrokerMessage()
	if err != nil {
		return err
	}

	// Set broker message to context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

	// Publish the message on event-broker through middlewares
	return c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
		return c.broker.Publish(ctx, path, brokerMsg)
	})
}

// PublishCreateBacktestResponse will publish messages to 'cryptellation.backtests.create.response' channel
func (c *AppController) PublishCreateBacktestResponse(ctx context.Context, msg CreateBacktestResponseMessage) error {
	// Get channel path
	path := "cryptellation.backtests.create.response"

	// Set correlation ID if it does not exist
	if id := msg.CorrelationID(); id == "" {
		msg.SetCorrelationID(uuid.New().String())
	}

	// Set context
	ctx = addAppContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "publication")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, msg.CorrelationID())

	// Convert to BrokerMessage
	brokerMsg, err := msg.toBrokerMessage()
	if err != nil {
		return err
	}

	// Set broker message to context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

	// Publish the message on event-broker through middlewares
	return c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
		return c.broker.Publish(ctx, path, brokerMsg)
	})
}

// PublishBacktestEvent will publish messages to 'cryptellation.backtests.events.{id}' channel
func (c *AppController) PublishBacktestEvent(ctx context.Context, params CryptellationBacktestsEventsParameters, msg BacktestsEventMessage) error {
	// Get channel path
	path := fmt.Sprintf("cryptellation.backtests.events.%v", params.Id)

	// Set context
	ctx = addAppContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "publication")

	// Convert to BrokerMessage
	brokerMsg, err := msg.toBrokerMessage()
	if err != nil {
		return err
	}

	// Set broker message to context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

	// Publish the message on event-broker through middlewares
	return c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
		return c.broker.Publish(ctx, path, brokerMsg)
	})
}

// PublishCreateBacktestOrderResponse will publish messages to 'cryptellation.backtests.orders.create.response' channel
func (c *AppController) PublishCreateBacktestOrderResponse(ctx context.Context, msg CreateBacktestOrderResponseMessage) error {
	// Get channel path
	path := "cryptellation.backtests.orders.create.response"

	// Set correlation ID if it does not exist
	if id := msg.CorrelationID(); id == "" {
		msg.SetCorrelationID(uuid.New().String())
	}

	// Set context
	ctx = addAppContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "publication")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, msg.CorrelationID())

	// Convert to BrokerMessage
	brokerMsg, err := msg.toBrokerMessage()
	if err != nil {
		return err
	}

	// Set broker message to context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

	// Publish the message on event-broker through middlewares
	return c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
		return c.broker.Publish(ctx, path, brokerMsg)
	})
}

// PublishListBacktestOrdersResponse will publish messages to 'cryptellation.backtests.orders.list.response' channel
func (c *AppController) PublishListBacktestOrdersResponse(ctx context.Context, msg ListBacktestOrdersResponseMessage) error {
	// Get channel path
	path := "cryptellation.backtests.orders.list.response"

	// Set correlation ID if it does not exist
	if id := msg.CorrelationID(); id == "" {
		msg.SetCorrelationID(uuid.New().String())
	}

	// Set context
	ctx = addAppContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "publication")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, msg.CorrelationID())

	// Convert to BrokerMessage
	brokerMsg, err := msg.toBrokerMessage()
	if err != nil {
		return err
	}

	// Set broker message to context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

	// Publish the message on event-broker through middlewares
	return c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
		return c.broker.Publish(ctx, path, brokerMsg)
	})
}

// PublishServiceInfoResponse will publish messages to 'cryptellation.backtests.service.info.response' channel
func (c *AppController) PublishServiceInfoResponse(ctx context.Context, msg ServiceInfoResponseMessage) error {
	// Get channel path
	path := "cryptellation.backtests.service.info.response"

	// Set correlation ID if it does not exist
	if id := msg.CorrelationID(); id == "" {
		msg.SetCorrelationID(uuid.New().String())
	}

	// Set context
	ctx = addAppContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "publication")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, msg.CorrelationID())

	// Convert to BrokerMessage
	brokerMsg, err := msg.toBrokerMessage()
	if err != nil {
		return err
	}

	// Set broker message to context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

	// Publish the message on event-broker through middlewares
	return c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
		return c.broker.Publish(ctx, path, brokerMsg)
	})
}

// PublishSubscribeBacktestResponse will publish messages to 'cryptellation.backtests.subscribe.response' channel
func (c *AppController) PublishSubscribeBacktestResponse(ctx context.Context, msg SubscribeBacktestResponseMessage) error {
	// Get channel path
	path := "cryptellation.backtests.subscribe.response"

	// Set correlation ID if it does not exist
	if id := msg.CorrelationID(); id == "" {
		msg.SetCorrelationID(uuid.New().String())
	}

	// Set context
	ctx = addAppContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "publication")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, msg.CorrelationID())

	// Convert to BrokerMessage
	brokerMsg, err := msg.toBrokerMessage()
	if err != nil {
		return err
	}

	// Set broker message to context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

	// Publish the message on event-broker through middlewares
	return c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
		return c.broker.Publish(ctx, path, brokerMsg)
	})
}
