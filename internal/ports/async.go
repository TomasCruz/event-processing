package ports

import (
	"context"
	"io"
)

type AsyncMsgProducer interface {
	io.Closer
	SetTopic(topic string) error
	SendAsyncMsg(msg []byte) error
}

type AsyncMsgConsumer interface {
	io.Closer
	SubscribeTopic(topic string) error
	Consume(ctx context.Context, msgChan chan<- []byte)
}
