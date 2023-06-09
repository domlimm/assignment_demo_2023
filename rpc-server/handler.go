package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	timestamp := time.Now().Unix()
	message := &Message{
		Message:   req.Message.GetText(),
		Sender:    req.Message.GetSender(),
		Timestamp: timestamp,
	}

	chatId, err := getChatId(req.Message.GetChat())

	if err != nil {
		return nil, err
	}

	err = rdb.SaveMessage(ctx, chatId, message)

	if err != nil {
		return nil, err
	}

	resp := rpc.NewSendResponse()
	resp.Code, resp.Msg = 0, "success"
	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	chatId, err := getChatId(req.GetChat())
	if err != nil {
		return nil, err
	}

	start := req.GetCursor()
	end := start + int64(req.GetLimit())
	messages, err := rdb.PullMessage(ctx, chatId, start, end, req.GetReverse())

	if err != nil {
		return nil, err
	}

	respMessages := make([]*rpc.Message, 0)
	var counter int32 = 0
	var nextCursor int64 = 0
	hasMore := false

	for _, msg := range messages {
		if counter+1 > req.GetLimit() {
			// having extra value here means it has more data
			hasMore = true
			nextCursor = end
			break // do not return the last message
		}

		temp := &rpc.Message{
			Chat:     req.GetChat(),
			Text:     msg.Message,
			Sender:   msg.Sender,
			SendTime: msg.Timestamp,
		}
		respMessages = append(respMessages, temp)
		counter += 1
	}

	resp := rpc.NewPullResponse()
	resp.Messages = respMessages
	resp.Code = 0
	resp.Msg = "success"
	resp.HasMore = &hasMore
	resp.NextCursor = &nextCursor

	return resp, nil
}

func getChatId(chat string) (string, error) {
	lowerCaseChatString := strings.ToLower(chat)
	people := strings.Split(lowerCaseChatString, ":")

	if len(people) != 2 {
		err := fmt.Errorf("Invalid Chat ID '%s', format should be as such: user1:user2", chat)
		return "", err
	}

	person1, person2 := people[0], people[1]

	var chatId string

	if comparator := strings.Compare(person1, person2); comparator == 1 {
		chatId = fmt.Sprintf("%s:%s", person2, person1)
	} else {
		chatId = fmt.Sprintf("%s:%s", person1, person2)
	}

	return chatId, nil
}
