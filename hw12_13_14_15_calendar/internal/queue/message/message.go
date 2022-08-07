package message

import (
	"encoding/json"
	"time"
)

type Message struct {
	Title    string
	DateTime time.Time
}

func New(title string, dateTime time.Time) Message {
	return Message{
		Title:    title,
		DateTime: dateTime,
	}
}

func NewFromBytes(data []byte) (Message, error) {
	m := Message{}
	err := json.Unmarshal(data, &m)
	return m, err
}

func (m *Message) String() string {
	return m.DateTime.String() + ": " + m.Title
}

func (m *Message) Encode() ([]byte, error) {
	return json.Marshal(m)
}
