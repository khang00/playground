package main

import "encoding/json"

type Message interface {
	Value(field string) string
	Add(flied, value string) Message
	Bytes() ([]byte, error)
}

type message struct {
	data map[string]string
}

func NewMsg() Message {
	return &message{data: make(map[string]string)}
}

func NewMsgFromByte(byte []byte) (Message, error) {
	data := make(map[string]string)
	err := json.Unmarshal(byte, &data)
	if err != nil {
		return nil, err
	}

	return &message{data: data}, err
}

func (m *message) Value(field string) string {
	return m.data[field]
}

func (m *message) Add(field, value string) Message {
	m.data[field] = value
	return m
}

func (m *message) Bytes() ([]byte, error) {
	bytes, err := json.Marshal(m.data)
	if err != nil {
		return nil, err
	}

	return bytes, err
}
