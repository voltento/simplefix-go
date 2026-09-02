package tests

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"gitlab.b2broker.tech/b2connect/b2connect/libs/go/simplefix-go/storages/memory"
	"net"
	"testing"
	"time"

	simplefixgo "gitlab.b2broker.tech/b2connect/b2connect/libs/go/simplefix-go"
	"gitlab.b2broker.tech/b2connect/b2connect/libs/go/simplefix-go/fix"
	"gitlab.b2broker.tech/b2connect/b2connect/libs/go/simplefix-go/session"
	fixgen "gitlab.b2broker.tech/b2connect/b2connect/libs/go/simplefix-go/tests/fix44"
)

// RunNewInitiator ...
func RunNewInitiator(addr string, t *testing.T, settings *session.LogonSettings, logon bool) (s *session.Session, handler *simplefixgo.DefaultHandler) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("could not dial: %s", err)
	}

	handler = simplefixgo.NewInitiatorHandler(context.Background(), fixgen.FieldMsgType, 10)
	client := simplefixgo.NewInitiator(conn, handler, 10, time.Minute*50)

	testStorage := memory.NewStorage()

	s, err = session.NewInitiatorSession(
		handler,
		&pseudoGeneratedOpts,
		settings,
		testStorage,
		testStorage,
	)
	if err != nil {
		panic(err)
	}

	// logging messages:
	handler.HandleIncoming(simplefixgo.AllMsgTypes, func(msg []byte) bool {
		fmt.Println("incoming:", string(bytes.ReplaceAll(msg, fix.Delimiter, []byte("|"))))
		return true
	})
	handler.HandleOutgoing(simplefixgo.AllMsgTypes, func(msg simplefixgo.SendingMessage) bool {
		data, mErr := msg.ToBytes()
		if mErr != nil {
			panic(mErr)
		}
		fmt.Println("outgoing:", string(bytes.ReplaceAll(data, fix.Delimiter, []byte("|"))))
		return true
	})
	if !logon {
		handler.HandleOutgoing(fixgen.MsgTypeLogon, func(msg simplefixgo.SendingMessage) bool {
			return false
		})
	}

	// todo move
	go func() {
		time.Sleep(time.Second * 10)
		fmt.Println("resending the request after 10 seconds")
		sErr := s.Send(fixgen.ResendRequest{}.New().SetFieldBeginSeqNo(2).SetFieldEndSeqNo(3))
		if sErr != nil {
			panic(sErr)
		}
	}()

	err = s.Run()
	if err != nil {
		t.Fatalf("could not run the session: %s", err)
	}

	go func() {
		err := client.Serve()
		if err != nil && !errors.Is(err, simplefixgo.ErrConnClosed) {
			panic(fmt.Errorf("could not serve the client: %s", err))
		}
	}()

	return s, handler
}
