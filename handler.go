package simplefixgo

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"gitlab.b2broker.tech/b2connect/b2connect/libs/go/simplefix-go/fix"
	"gitlab.b2broker.tech/b2connect/b2connect/libs/go/simplefix-go/fix/buffer"
	"gitlab.b2broker.tech/b2connect/b2connect/libs/go/simplefix-go/session/messages"
	"gitlab.b2broker.tech/b2connect/b2connect/libs/go/simplefix-go/utils"
)

// AllMsgTypes is a all msg types.
const AllMsgTypes = "ALL"

// SendingMessage provides a basic method for sending messages.
type SendingMessage interface {
	HeaderBuilder() messages.HeaderBuilder
	MsgType() string
	ToBytes() ([]byte, error)
	ToBytesBuffered(buffers *buffer.MessageByteBuffers) ([]byte, error)
}

// DefaultHandler is a standard handler for the Acceptor and Initiator objects.
type DefaultHandler struct {
	mu sync.Mutex

	out      chan []byte
	incoming chan []byte

	incomingHandlers IncomingHandlerPool
	outgoingHandlers OutgoingHandlerPool

	eventHandlers    *utils.EventHandlerPool
	messageConverter *fix.MessageByteConverter

	msgTypeTag string

	// sendDeadline bounds how long sendRaw waits for space in the outbound
	// channel. Zero keeps the historical behavior: wait until the handler
	// context is canceled. Set per acceptor so trading and quoting can differ.
	sendDeadline time.Duration

	ctx    context.Context
	cancel context.CancelFunc
	errors chan error
}

// NewAcceptorHandler creates a handler for an Acceptor object.
func NewAcceptorHandler(ctx context.Context, msgTypeTag string, bufferSize int) *DefaultHandler {
	sh := &DefaultHandler{
		msgTypeTag:    msgTypeTag,
		eventHandlers: utils.NewEventHandlerPool(),

		out:      make(chan []byte, bufferSize),
		incoming: make(chan []byte, bufferSize),
		errors:   make(chan error),

		incomingHandlers: NewIncomingHandlerPool(),
		outgoingHandlers: NewOutgoingHandlerPool(),
		messageConverter: fix.NewMessageByteConverter(500),
	}

	sh.ctx, sh.cancel = context.WithCancel(ctx)

	return sh
}

// NewInitiatorHandler creates a handler for the Initiator object.
func NewInitiatorHandler(ctx context.Context, msgTypeTag string, bufferSize int) *DefaultHandler {
	sh := &DefaultHandler{
		msgTypeTag:    msgTypeTag,
		eventHandlers: utils.NewEventHandlerPool(),

		out:      make(chan []byte, bufferSize),
		incoming: make(chan []byte, bufferSize),
		errors:   make(chan error),

		incomingHandlers: NewIncomingHandlerPool(),
		outgoingHandlers: NewOutgoingHandlerPool(),
		messageConverter: fix.NewMessageByteConverter(500),
	}

	sh.ctx, sh.cancel = context.WithCancel(ctx)

	return sh
}

func (h *DefaultHandler) sendRaw(data []byte) error {
	var timeout <-chan time.Time
	if h.sendDeadline > 0 {
		timer := time.NewTimer(h.sendDeadline)
		defer timer.Stop()
		timeout = timer.C
	}

	select {
	case h.out <- data:
	case <-h.ctx.Done():
		return fmt.Errorf("the handler is stopped")
	case <-timeout:
		// Cancelling the handler context is what tears the connection down:
		// Run returns, and the acceptor's errgroup defer closes the socket.
		// Run fires EventDisconnect on this path so the session state machine
		// moves to Disconnect instead of reporting its last live state.
		h.cancel()
		return fmt.Errorf("%w after %s", ErrSendTimeout, h.sendDeadline)
	}
	return nil
}

// SetSendDeadline bounds how long a send waits for space in the outbound
// channel before the handler context is canceled. Zero disables the bound.
// Call it before the handler starts serving; it is not safe to change later.
func (h *DefaultHandler) SetSendDeadline(d time.Duration) {
	h.sendDeadline = d
}

func (h *DefaultHandler) send(msg SendingMessage) error {
	ok := h.outgoingHandlers.Range(AllMsgTypes, func(handle OutgoingHandlerFunc) bool {
		return handle(msg)
	})
	if !ok {
		return errors.New("the handler for all message types has refused the message and returned false")
	}

	ok = h.outgoingHandlers.Range(msg.MsgType(), func(handle OutgoingHandlerFunc) bool {
		return handle(msg)
	})
	if !ok {
		return errors.New("the handler for the current type has refused the message and returned false")
	}

	data, err := msg.ToBytes()
	if err != nil {
		return err
	}

	return h.sendRaw(data)
}

func (h *DefaultHandler) sendBuffered(msg SendingMessage) error {
	ok := h.outgoingHandlers.Range(AllMsgTypes, func(handle OutgoingHandlerFunc) bool {
		return handle(msg)
	})
	if !ok {
		return errors.New("the handler for all message types has refused the message and returned false")
	}

	ok = h.outgoingHandlers.Range(msg.MsgType(), func(handle OutgoingHandlerFunc) bool {
		return handle(msg)
	})
	if !ok {
		return errors.New("the handler for the current type has refused the message and returned false")
	}

	data, err := h.messageConverter.ConvertToBytes(msg)
	if err != nil {
		return errors.New("failed to convert message to bytes: " + err.Error())
	}

	return h.sendRaw(data)
}

// SendRaw sends a message in the byte array format
// without involving any additional handlers.
func (h *DefaultHandler) SendRaw(data []byte) error {
	return h.sendRaw(data)
}

// Send is a function that sends a previously prepared message.
func (h *DefaultHandler) Send(message SendingMessage) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.send(message)
}

// SendBuffered sends msg using a buffered encoder.
func (h *DefaultHandler) SendBuffered(message SendingMessage) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.sendBuffered(message)
}

// SendBatch is a function that sends previously prepared messages.
func (h *DefaultHandler) SendBatch(messages []SendingMessage) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, message := range messages {
		err := h.send(message)
		if err != nil {
			return err
		}
	}

	return nil
}

// RemoveIncomingHandler removes an existing handler for incoming messages.
func (h *DefaultHandler) RemoveIncomingHandler(msgType string, id int64) (err error) {
	return h.incomingHandlers.Remove(msgType, id)
}

// RemoveOutgoingHandler removes an existing handler for outgoing messages.
func (h *DefaultHandler) RemoveOutgoingHandler(msgType string, id int64) (err error) {
	return h.outgoingHandlers.Remove(msgType, id)
}

// HandleIncoming subscribes a handler function to incoming messages with a specific msgType.
// To subscribe to all messages, specify the AllMsgTypes constant for the msgType field
// (such messages will have a higher priority than the ones assigned to specific handlers).
func (h *DefaultHandler) HandleIncoming(msgType string, handle IncomingHandlerFunc) (id int64) {
	return h.incomingHandlers.Add(msgType, handle)
}

// HandleOutgoing subscribes a handler function to outgoing messages with a specific msgType
// (this may be required for modifying messages before sending).
// To subscribe to all messages, specify the AllMsgTypes constant for the msgType field
// (such messages will have a higher priority than the ones assigned to specific handlers).
func (h *DefaultHandler) HandleOutgoing(msgType string, handle OutgoingHandlerFunc) (id int64) {
	return h.outgoingHandlers.Add(msgType, handle)
}

// ServeIncoming is an internal method for handling incoming messages.
func (h *DefaultHandler) ServeIncoming(msg []byte) {
	h.incoming <- msg
}

func (h *DefaultHandler) serve(msg []byte) (err error) {
	msgTypeB, err := fix.ValueByTag(msg, h.msgTypeTag)
	if err != nil {
		return fmt.Errorf("msg type: %w", err)
	}
	msgType := string(msgTypeB)

	ok := h.incomingHandlers.Range(AllMsgTypes, func(handle IncomingHandlerFunc) bool {
		return handle(msg)
	})
	if !ok {
		return fmt.Errorf("failed to handle the incoming message, msg: %v", string(msg))
	}

	ok = h.incomingHandlers.Range(msgType, func(handle IncomingHandlerFunc) bool {
		return handle(msg)
	})
	if !ok {
		return fmt.Errorf("failed to handle the incoming message by tag, msg: %v", string(msg))
	}

	return nil
}

// Run is a function that is used for listening and processing messages.
func (h *DefaultHandler) Run() (err error) {
	h.eventHandlers.Trigger(utils.EventConnect)
	defer h.processRemainingErrors()

	for {
		select {
		case msg, ok := <-h.incoming:
			if !ok {
				return ErrConnClosed
			}

			err = h.serve(msg)
			if err != nil {
				h.eventHandlers.Trigger(utils.EventDisconnect)
				return err
			}

		case <-h.ctx.Done():
			h.processRemainingIncoming()

			// The context is cancelled both by an explicit Stop and by a dead
			// connection (send deadline, read error via the acceptor). Either
			// way the counterparty is gone, so this is a disconnect — not just
			// a stop — for the session state machine.
			h.eventHandlers.Trigger(utils.EventDisconnect)
			h.eventHandlers.Trigger(utils.EventStopped)

			return

		case err := <-h.errors:
			h.processRemainingIncoming()

			if errors.Is(err, ErrConnClosed) {
				h.eventHandlers.Trigger(utils.EventDisconnect)
			}

			return err
		}
	}
}

func (h *DefaultHandler) processRemainingIncoming() {
	for {
		select {
		case msg, ok := <-h.incoming:
			if ok {
				_ = h.serve(msg)
			}
		default:
			return
		}
	}
}

func (h *DefaultHandler) processRemainingErrors() {
	go func() {
		for {
			_, ok := <-h.errors
			if !ok {
				return
			}
		}
	}()
}

// Context returns default handler's context.
func (h *DefaultHandler) Context() context.Context {
	return h.ctx
}

// Outgoing is a service method that provides an outgoing channel
// to the server or client connection manager.
func (h *DefaultHandler) Outgoing() <-chan []byte {
	return h.out
}

// Stop is a function that enables graceful termination of a session.
func (h *DefaultHandler) Stop() {
	h.cancel()
}

// StopWithError is a function that enables graceful termination of a session with throwing an error.
func (h *DefaultHandler) StopWithError(err error) {
	h.errors <- err
}

// CloseErrorChan is a function that closes the handler's error chan.
func (h *DefaultHandler) CloseErrorChan() {
	close(h.errors)
}

// OnDisconnect handles disconnection events.
func (h *DefaultHandler) OnDisconnect(handlerFunc utils.EventHandlerFunc) {
	h.eventHandlers.Handle(utils.EventDisconnect, handlerFunc)
}

// OnConnect handles connection events.
func (h *DefaultHandler) OnConnect(handlerFunc utils.EventHandlerFunc) {
	h.eventHandlers.Handle(utils.EventConnect, handlerFunc)
}

// OnStopped handles session termination events.
func (h *DefaultHandler) OnStopped(handlerFunc utils.EventHandlerFunc) {
	h.eventHandlers.Handle(utils.EventStopped, handlerFunc)
}
