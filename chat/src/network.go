package main

type Network interface {
	sendToClient(clientID int) error
}

type network struct {
}

func NewNetwork() Network {

}
