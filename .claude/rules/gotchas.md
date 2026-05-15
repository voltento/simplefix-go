# Gotchas

## All 8 message builders are mandatory for Session

Session construction requires ALL 8 builders: Header, Trailer, Logon, Logout, Reject, Heartbeat, TestRequest, ResendRequest. Missing any one fails with `ErrMissingMessageBuilder`. See [session/opts.go](../../session/opts.go).

## Conn framing reads until `10=` — messages without checksum hang

`Conn` reads TCP stream until the `10=` (checksum) tag delimiter appears. Any message not ending with a checksum tag causes an indefinite hang. See [conn.go](../../conn.go).

## Tags are string dependencies — mismatches silently break routing

Session requires tag numbers as integers in `messages.Tags` (MsgType=35, MsgSeqNum=34, etc.). A mismatch between actual message tags and the Tags struct silently breaks message routing and sequence tracking. See [session/messages/contstants.go](../../session/messages/contstants.go) (filename typo is in the source).

## `AllMsgTypes` constant is the string "ALL"

Handler subscription uses `"ALL"` as a catch-all. If any FIX message type is literally named "ALL", behavior is undefined. See [handler.go](../../handler.go).

## HandlerPool.Remove() is always a no-op for individual handlers

`Remove()` ignores the handler ID parameter, never removes anything from the internal slice, then calls `free()` which only deletes the msgType entry when the slice is empty — a condition `Remove()` itself never creates. Additionally, `Remove()` does not hold the pool mutex (`p.mu`) — concurrent calls alongside `Add()` or `Range()` are an unsynchronised data race. To stop routing a message type, replace the pool. Handler IDs are for tracking only. See [handler_func_pool.go](../../handler_func_pool.go).

## Generated code is overwritten on regeneration

`fixgen` writes to the output directory without checking for existing files. Custom edits to generated files in `tests/fix44/` will be lost. Never edit generated files.

## Session state is atomic but transitions are not

Session state uses `atomic.Int64`, but field mutations during state transitions are not atomic with the state change. A reader may see state X with stale field values.

## Custom message builders must reset state between messages

If a custom builder does not reset its state, old field values from the previous call may leak into the next message. fixgen-generated builders are not affected — they call `fix.NewMessage()` fresh on every `Build()`. Only matters when writing custom `MessageBuilder` implementations. See test patterns at [tests/acceptor_initiator_test.go](../../tests/acceptor_initiator_test.go).

## Context cancellation race between Initiator and DefaultHandler

`Initiator` and `DefaultHandler` each hold independent contexts (`Initiator.ctx` vs `DefaultHandler.ctx`). Simultaneous cancellation of both can cause goroutine leaks. See [initiator.go](../../initiator.go).
