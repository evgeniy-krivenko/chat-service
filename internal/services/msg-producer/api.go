package msgproducer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

type Message struct {
	ID         types.MessageID `json:"id"`
	ChatID     types.ChatID    `json:"chatId"`
	Body       string          `json:"body"`
	FromClient bool            `json:"fromClient"`
}

func (s *Service) ProduceMessage(ctx context.Context, msg Message) error {
	var sendMessage []byte

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal for produce: %v", err)
	}

	if s.cipher != nil {
		nonce, err := s.nonceFactory(s.cipher.NonceSize())
		if err != nil {
			return fmt.Errorf("create nonce for encrypt message")
		}

		sendMessage = s.cipher.Seal(nonce, nonce, data, nil)
	} else {
		sendMessage = data
	}

	key, err := msg.ChatID.MarshalText()
	if err != nil {
		return fmt.Errorf("marshal chat id: %v", err)
	}

	if err := s.wr.WriteMessages(ctx, kafka.Message{Key: key, Value: sendMessage}); err != nil {
		return fmt.Errorf("produce message: %v", err)
	}

	return nil
}

func (s *Service) Close() error {
	return s.wr.Close()
}