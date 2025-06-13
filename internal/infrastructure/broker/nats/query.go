package bnats

import (
	"encoding/json"
	"fmt"
)

func push[T any](n *_nats, sub string, in T) error {
	data, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("ошибка маршала в json: %w", err)
	}

	err = n.conn.Publish(sub, data)
	if err != nil {
		return fmt.Errorf("ошибка отправки сообщения в топик %s: %w", sub, err)
	}

	return nil
}
