package overlay

import "context"

// Topic represents an overlay broadcast topic.
type Topic string

// Publisher can publish messages to an overlay topic.
type Publisher interface {
	Start(ctx context.Context) error
	Close(ctx context.Context) error
	Publish(ctx context.Context, topic Topic, data []byte) error
}

// Subscriber can subscribe to overlay topics and receive messages.
type Subscriber interface {
	Start(ctx context.Context) error
	Close(ctx context.Context) error
	Subscribe(ctx context.Context, topic Topic) (<-chan []byte, error)
	Unsubscribe(ctx context.Context, topic Topic) error
}

// Manager unifies publisher/subscriber aspects and membership controls.
type Manager interface {
	Publisher
	Subscriber

	Join(ctx context.Context, overlayID string) error
	Leave(ctx context.Context, overlayID string) error
}
