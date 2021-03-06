package client

import (
	"BattleCoderServer/client/connection"
)

// MatchWaitingClient はマッチ待機中のクライアントを表します
type MatchWaitingClient struct {
	clientConnection connection.ClientConnection
	matched          bool
}

// NewMatchWaitingClient を生成する
func NewMatchWaitingClient(clientConnection connection.ClientConnection) *MatchWaitingClient {
	matchWaitingClient := &MatchWaitingClient{
		clientConnection: clientConnection,
		matched:          false}
	return matchWaitingClient
}

// Match は対戦相手をマッチする
func (t *MatchWaitingClient) Match(other *MatchWaitingClient) {
	t.matched = true
	other.matched = true
	t.clientConnection.WriteChannel <- []byte("match_client\n")
	other.clientConnection.WriteChannel <- []byte("match_host\n")
	NewMatchPairClient(&t.clientConnection, &other.clientConnection)
}
