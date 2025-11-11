package kafka

import (
    "context"
    "time"

    kfk "github.com/segmentio/kafka-go"
)

type Producer struct {
    writer *kfk.Writer
}

type Config struct {
    Brokers []string
    Topic   string
}

func NewProducer(cfg Config) *Producer {
    return &Producer{
        writer: &kfk.Writer{
            Addr:         kfk.TCP(cfg.Brokers...),
            Topic:        cfg.Topic,
            RequiredAcks: kfk.RequireAll,
            Balancer:     &kfk.Hash{},
        },
    }
}

func (p *Producer) Close() error { return p.writer.Close() }

func (p *Producer) Publish(ctx context.Context, key []byte, value []byte, headers map[string]string) error {
    var hs []kfk.Header
    for k, v := range headers {
        hs = append(hs, kfk.Header{Key: k, Value: []byte(v)})
    }
    msg := kfk.Message{Key: key, Value: value, Time: time.Now(), Headers: hs}
    return p.writer.WriteMessages(ctx, msg)
}

