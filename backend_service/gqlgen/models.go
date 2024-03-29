// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gqlgen

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Message struct {
	ID       string    `json:"id"`
	Data     string    `json:"data"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

type MessageCountPayload struct {
	Count     int               `json:"count"`
	Providers []MessageProvider `json:"providers"`
}

type MessageCreatedPayload struct {
	Message  *Message        `json:"message"`
	Provider MessageProvider `json:"provider"`
}

type MessagePayload struct {
	Message  *Message        `json:"message"`
	Provider MessageProvider `json:"provider"`
}

type MessagesPayload struct {
	Messages  []*Message        `json:"messages"`
	Providers []MessageProvider `json:"providers"`
}

type NewMessageInput struct {
	Providers []MessageProvider `json:"providers"`
	Data      string            `json:"data"`
}

type NewMessagePayload struct {
	Status    string            `json:"status"`
	Providers []MessageProvider `json:"providers"`
}

type MessageProvider string

const (
	MessageProviderMongo    MessageProvider = "MONGO"
	MessageProviderPostgres MessageProvider = "POSTGRES"
)

var AllMessageProvider = []MessageProvider{
	MessageProviderMongo,
	MessageProviderPostgres,
}

func (e MessageProvider) IsValid() bool {
	switch e {
	case MessageProviderMongo, MessageProviderPostgres:
		return true
	}
	return false
}

func (e MessageProvider) String() string {
	return string(e)
}

func (e *MessageProvider) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MessageProvider(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MessageProvider", str)
	}
	return nil
}

func (e MessageProvider) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type MessageSortField string

const (
	MessageSortFieldID       MessageSortField = "ID"
	MessageSortFieldData     MessageSortField = "DATA"
	MessageSortFieldCreated  MessageSortField = "CREATED"
	MessageSortFieldModified MessageSortField = "MODIFIED"
)

var AllMessageSortField = []MessageSortField{
	MessageSortFieldID,
	MessageSortFieldData,
	MessageSortFieldCreated,
	MessageSortFieldModified,
}

func (e MessageSortField) IsValid() bool {
	switch e {
	case MessageSortFieldID, MessageSortFieldData, MessageSortFieldCreated, MessageSortFieldModified:
		return true
	}
	return false
}

func (e MessageSortField) String() string {
	return string(e)
}

func (e *MessageSortField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MessageSortField(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MessageSortField", str)
	}
	return nil
}

func (e MessageSortField) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
