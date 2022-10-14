package main

type Message interface {
	Value(field []byte)
	JSON() []byte
}
