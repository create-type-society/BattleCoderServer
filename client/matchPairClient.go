package client

import (
	"BattleCoderServer/client/connection"
)

func NewMatchPairClient(
	clientConnection1 *connection.ClientConnection,
	clientConnection2 *connection.ClientConnection) {

	finished := false

	go func() {
		for finished == false {
			clientConnection2.WriteChannel <- <-clientConnection1.ReadChannel
		}
	}()

	go func() {
		for finished == false {
			clientConnection1.WriteChannel <- <-clientConnection2.ReadChannel
		}
	}()

	go func() {
		<-clientConnection1.FinishedChannel
		finished = true
	}()

	go func() {
		<-clientConnection2.FinishedChannel
		finished = true
	}()

}
