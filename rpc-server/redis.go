package main

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	cli *redis.Client
}

type Message struct {
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func (c *RedisClient) InitRedis(ctx context.Context, address, password string) error {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0, // use default DB
	})

	_, err := client.Ping(ctx).Result()

	if err != nil {
		return err
	}

	c.cli = client

	return nil
}

func (c *RedisClient) SaveMessage(ctx context.Context, chatId string, message *Message) error {
	// https://golang.cafe/blog/golang-json-marshal-example.html
	text, err := json.Marshal(message) // returns JSON encoding of what's passed in

	if err != nil {
		return err
	}

	// Create a slice of redis.Z to represent sorted set members and scores
	pair := &redis.Z{
		Score:  float64(message.Timestamp),
		Member: text, // JSON encoded message
	}

	_, err = c.cli.ZAdd(ctx, chatId, *pair).Result()

	if err != nil {
		return err
	}

	return nil
}

func (c *RedisClient) PullMessage(ctx context.Context, chatId string, start, end int64, reverse bool) ([]*Message, error) {
	var (
		rawMessages     []string
		decodedMessages []*Message
		err             error
	)

	if reverse {
		rawMessages, err = c.cli.ZRevRange(ctx, chatId, start, end).Result()

		if err != nil {
			return nil, err
		}
	} else {
		rawMessages, err = c.cli.ZRange(ctx, chatId, start, end).Result()

		if err != nil {
			return nil, err
		}
	}

	for _, msg := range rawMessages {
		temp := &Message{}
		err := json.Unmarshal([]byte(msg), temp)

		if err != nil {
			return nil, err
		}

		decodedMessages = append(decodedMessages, temp)
	}

	return decodedMessages, nil
}
