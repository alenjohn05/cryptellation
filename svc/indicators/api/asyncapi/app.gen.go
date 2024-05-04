// Package "asyncapi" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version v0.39.0 DO NOT EDIT.
package asyncapi

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// AppSubscriber contains all handlers that are listening messages for App
type AppSubscriber interface {
	// SMAOperationReceived receive all SMARequest messages from SMARequest channel.
	SMAOperationReceived(ctx context.Context, msg SMARequestMessage) error

	// ServiceInfoOperationReceived receive all ServiceInfoRequest messages from ServiceInfoRequest channel.
	ServiceInfoOperationReceived(ctx context.Context, msg ServiceInfoRequestMessage) error
}

// AppController is the structure that provides sending capabilities to the
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
		errorHandler:  extensions.DefaultErrorHandler(),
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

func addAppContextValues(ctx context.Context, addr string) context.Context {
	ctx = context.WithValue(ctx, extensions.ContextKeyIsVersion, "1.0.0")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsProvider, "app")
	return context.WithValue(ctx, extensions.ContextKeyIsChannel, addr)
}

// Close will clean up any existing resources on the controller
func (c *AppController) Close(ctx context.Context) {
	// Unsubscribing remaining channels
	c.UnsubscribeFromAllChannels(ctx)

	c.logger.Info(ctx, "Closed app controller")
}

// SubscribeToAllChannels will receive messages from channels where channel has
// no parameter on which the app is expecting messages. For channels with parameters,
// they should be subscribed independently.
func (c *AppController) SubscribeToAllChannels(ctx context.Context, as AppSubscriber) error {
	if as == nil {
		return extensions.ErrNilAppSubscriber
	}

	if err := c.SubscribeToSMAOperation(ctx, as.SMAOperationReceived); err != nil {
		return err
	}
	if err := c.SubscribeToServiceInfoOperation(ctx, as.ServiceInfoOperationReceived); err != nil {
		return err
	}

	return nil
}

// UnsubscribeFromAllChannels will stop the subscription of all remaining subscribed channels
func (c *AppController) UnsubscribeFromAllChannels(ctx context.Context) {
	c.UnsubscribeFromSMAOperation(ctx)
	c.UnsubscribeFromServiceInfoOperation(ctx)
}

// SubscribeToSMAOperation will receive SMARequest messages from SMARequest channel.
//
// Callback function 'fn' will be called each time a new message is received.
//
// NOTE: for now, this only support the first message from AsyncAPI list.
//
// NOTE: for now, this only support the first message from AsyncAPI list.
// If you need support for other messages, please raise an issue.
func (c *AppController) SubscribeToSMAOperation(
	ctx context.Context,
	fn func(ctx context.Context, msg SMARequestMessage) error,
) error {
	// Get channel address
	addr := "cryptellation.indicators.sma"

	// Set context
	ctx = addAppContextValues(ctx, addr)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "reception")

	// Check if the controller is already subscribed
	_, exists := c.subscriptions[addr]
	if exists {
		err := fmt.Errorf("%w: controller is already subscribed on channel %q", extensions.ErrAlreadySubscribedChannel, addr)
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Subscribe to broker channel
	sub, err := c.broker.Subscribe(ctx, addr)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}
	c.logger.Info(ctx, "Subscribed to channel")

	// Asynchronously listen to new messages and pass them to app receiver
	go func() {
		for {
			// Wait for next message
			acknowledgeableBrokerMessage, open := <-sub.MessagesChannel()

			// If subscription is closed and there is no more message
			// (i.e. uninitialized message), then exit the function
			if !open && acknowledgeableBrokerMessage.IsUninitialized() {
				return
			}

			// Set broker message to context
			ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, acknowledgeableBrokerMessage.String())

			// Execute middlewares before handling the message
			if err := c.executeMiddlewares(ctx, &acknowledgeableBrokerMessage.BrokerMessage, func(ctx context.Context) error {
				// Process message
				msg, err := brokerMessageToSMARequestMessage(acknowledgeableBrokerMessage.BrokerMessage)
				if err != nil {
					return err
				}

				// Add correlation ID to context if it exists
				if id := msg.CorrelationID(); id != "" {
					ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, id)
				}

				// Execute the subscription function
				if err := fn(ctx, msg); err != nil {
					return err
				}

				acknowledgeableBrokerMessage.Ack()

				return nil
			}); err != nil {
				c.errorHandler(ctx, addr, &acknowledgeableBrokerMessage, err)
				// On error execute the acknowledgeableBrokerMessage nack() function and
				// let the BrokerAcknowledgment decide what is the right nack behavior for the broker
				acknowledgeableBrokerMessage.Nak()
			}
		}
	}()

	// Add the cancel channel to the inside map
	c.subscriptions[addr] = sub

	return nil
}

// ReplyToSMAOperation is a helper function to
// reply to a SMARequest message with a SMAResponse message on SMAResponse channel.
func (c *AppController) ReplyToSMAOperation(ctx context.Context, recvMsg SMARequestMessage, fn func(replyMsg *SMAResponseMessage)) error {
	// Create reply message
	replyMsg := NewSMAResponseMessage()
	replyMsg.SetAsResponseFrom(&recvMsg)

	// Execute callback function
	fn(&replyMsg)

	// Publish reply
	chanAddr := recvMsg.Headers.ReplyTo

	return c.SendAsReplyToSMAOperation(ctx, chanAddr, replyMsg)
}

// UnsubscribeFromSMAOperation will stop the reception of SMARequest messages from SMARequest channel.
// A timeout can be set in context to avoid blocking operation, if needed.
func (c *AppController) UnsubscribeFromSMAOperation(
	ctx context.Context,
) {
	// Get channel address
	addr := "cryptellation.indicators.sma"

	// Check if there receivers for this channel
	sub, exists := c.subscriptions[addr]
	if !exists {
		return
	}

	// Set context
	ctx = addAppContextValues(ctx, addr)

	// Stop the subscription
	sub.Cancel(ctx)

	// Remove if from the receivers
	delete(c.subscriptions, addr)

	c.logger.Info(ctx, "Unsubscribed from channel")
} // SubscribeToServiceInfoOperation will receive ServiceInfoRequest messages from ServiceInfoRequest channel.
// Callback function 'fn' will be called each time a new message is received.
//
// NOTE: for now, this only support the first message from AsyncAPI list.
//
// NOTE: for now, this only support the first message from AsyncAPI list.
// If you need support for other messages, please raise an issue.
func (c *AppController) SubscribeToServiceInfoOperation(
	ctx context.Context,
	fn func(ctx context.Context, msg ServiceInfoRequestMessage) error,
) error {
	// Get channel address
	addr := "cryptellation.indicators.info"

	// Set context
	ctx = addAppContextValues(ctx, addr)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "reception")

	// Check if the controller is already subscribed
	_, exists := c.subscriptions[addr]
	if exists {
		err := fmt.Errorf("%w: controller is already subscribed on channel %q", extensions.ErrAlreadySubscribedChannel, addr)
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Subscribe to broker channel
	sub, err := c.broker.Subscribe(ctx, addr)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}
	c.logger.Info(ctx, "Subscribed to channel")

	// Asynchronously listen to new messages and pass them to app receiver
	go func() {
		for {
			// Wait for next message
			acknowledgeableBrokerMessage, open := <-sub.MessagesChannel()

			// If subscription is closed and there is no more message
			// (i.e. uninitialized message), then exit the function
			if !open && acknowledgeableBrokerMessage.IsUninitialized() {
				return
			}

			// Set broker message to context
			ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, acknowledgeableBrokerMessage.String())

			// Execute middlewares before handling the message
			if err := c.executeMiddlewares(ctx, &acknowledgeableBrokerMessage.BrokerMessage, func(ctx context.Context) error {
				// Process message
				msg, err := brokerMessageToServiceInfoRequestMessage(acknowledgeableBrokerMessage.BrokerMessage)
				if err != nil {
					return err
				}

				// Add correlation ID to context if it exists
				if id := msg.CorrelationID(); id != "" {
					ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, id)
				}

				// Execute the subscription function
				if err := fn(ctx, msg); err != nil {
					return err
				}

				acknowledgeableBrokerMessage.Ack()

				return nil
			}); err != nil {
				c.errorHandler(ctx, addr, &acknowledgeableBrokerMessage, err)
				// On error execute the acknowledgeableBrokerMessage nack() function and
				// let the BrokerAcknowledgment decide what is the right nack behavior for the broker
				acknowledgeableBrokerMessage.Nak()
			}
		}
	}()

	// Add the cancel channel to the inside map
	c.subscriptions[addr] = sub

	return nil
}

// ReplyToServiceInfoOperation is a helper function to
// reply to a ServiceInfoRequest message with a ServiceInfoResponse message on ServiceInfoResponse channel.
func (c *AppController) ReplyToServiceInfoOperation(ctx context.Context, recvMsg ServiceInfoRequestMessage, fn func(replyMsg *ServiceInfoResponseMessage)) error {
	// Create reply message
	replyMsg := NewServiceInfoResponseMessage()
	replyMsg.SetAsResponseFrom(&recvMsg)

	// Execute callback function
	fn(&replyMsg)

	// Publish reply
	chanAddr := recvMsg.Headers.ReplyTo

	return c.SendAsReplyToServiceInfoOperation(ctx, chanAddr, replyMsg)
}

// UnsubscribeFromServiceInfoOperation will stop the reception of ServiceInfoRequest messages from ServiceInfoRequest channel.
// A timeout can be set in context to avoid blocking operation, if needed.
func (c *AppController) UnsubscribeFromServiceInfoOperation(
	ctx context.Context,
) {
	// Get channel address
	addr := "cryptellation.indicators.info"

	// Check if there receivers for this channel
	sub, exists := c.subscriptions[addr]
	if !exists {
		return
	}

	// Set context
	ctx = addAppContextValues(ctx, addr)

	// Stop the subscription
	sub.Cancel(ctx)

	// Remove if from the receivers
	delete(c.subscriptions, addr)

	c.logger.Info(ctx, "Unsubscribed from channel")
}

// SendAsReplyToSMAOperation will send a SMAResponse message on SMAResponse channel.
//
// NOTE: for now, this only support the first message from AsyncAPI list.
// If you need support for other messages, please raise an issue.
func (c *AppController) SendAsReplyToSMAOperation(
	ctx context.Context,
	chanAddr string,
	msg SMAResponseMessage,
) error {
	// Set channel address
	addr := chanAddr

	// Set correlation ID if it does not exist
	if id := msg.CorrelationID(); id == "" {
		c.logger.Error(ctx, extensions.ErrNoCorrelationIDSet.Error())
		return extensions.ErrNoCorrelationIDSet

	}

	// Set context
	ctx = addAppContextValues(ctx, addr)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "publication")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, msg.CorrelationID())

	// Convert to BrokerMessage
	brokerMsg, err := msg.toBrokerMessage()
	if err != nil {
		return err
	}

	// Set broker message to context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

	// Send the message on event-broker through middlewares
	return c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
		return c.broker.Publish(ctx, addr, brokerMsg)
	})
}

// SendAsReplyToServiceInfoOperation will send a ServiceInfoResponse message on ServiceInfoResponse channel.
//
// NOTE: for now, this only support the first message from AsyncAPI list.
// If you need support for other messages, please raise an issue.
func (c *AppController) SendAsReplyToServiceInfoOperation(
	ctx context.Context,
	chanAddr string,
	msg ServiceInfoResponseMessage,
) error {
	// Set channel address
	addr := chanAddr

	// Set correlation ID if it does not exist
	if id := msg.CorrelationID(); id == "" {
		c.logger.Error(ctx, extensions.ErrNoCorrelationIDSet.Error())
		return extensions.ErrNoCorrelationIDSet

	}

	// Set context
	ctx = addAppContextValues(ctx, addr)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "publication")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, msg.CorrelationID())

	// Convert to BrokerMessage
	brokerMsg, err := msg.toBrokerMessage()
	if err != nil {
		return err
	}

	// Set broker message to context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

	// Send the message on event-broker through middlewares
	return c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
		return c.broker.Publish(ctx, addr, brokerMsg)
	})
}
