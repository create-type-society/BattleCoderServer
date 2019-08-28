package main

import (
	"BattleCoderServer/client"
	"BattleCoderServer/client/connection"
	"fmt"
	"net"
)

const host = "localhost:3000"

func main() {
	listen()
}

//tcpサーバの起動と処理
func listen() {
	listen, _ := net.Listen("tcp", host)
	fmt.Println("サーバ起動@" + host)
	var firstMatchWaitingClient *client.MatchWaitingClient = nil
	for {
		conn, _ := listen.Accept()
		clientConnection := connection.NewClientConnection(conn)
		matchWaitingClient := client.NewMatchWaitingClient(clientConnection)
		if firstMatchWaitingClient == nil {
			firstMatchWaitingClient = matchWaitingClient
		} else {
			firstMatchWaitingClient.Match(matchWaitingClient)
			firstMatchWaitingClient = nil
		}
		go func() {
			<-clientConnection.FinishedChannel
			if firstMatchWaitingClient == matchWaitingClient {
				firstMatchWaitingClient = nil
			}
		}()
	}
}
