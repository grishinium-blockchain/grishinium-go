package adnl

import "context"

// Address represents an ADNL-style address.
type Address struct {
	ID   string // peer ID or similar identity
	Host string
	Port int
}

// Message is a raw ADNL message payload.
type Message []byte

// Transport abstracts low-level connectivity.
type Transport interface {
	Start(ctx context.Context) error
	Close(ctx context.Context) error
	Dial(ctx context.Context, addr Address) (Conn, error)
	Listen(addr Address) error
}

// Conn is a bidirectional stream abstraction.
type Conn interface {
	Send(ctx context.Context, msg Message) error
	Recv(ctx context.Context) (Message, error)
	Close(ctx context.Context) error
}

// Node provides higher-level node operations over Transport.
type Node interface {
	Transport() Transport
	LocalAddr() Address
	SendTo(ctx context.Context, to Address, msg Message) error
}
