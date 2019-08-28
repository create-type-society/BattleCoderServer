package connection

import (
	"BattleCoderServer/client/connection/logger"
	"net"
	"time"
)

// ClientConnection は接続してきたクライアントの状態を標準出力するためのもの
type ClientConnection struct {
	conn            net.Conn
	logger          logger.Logger
	writeChannel    chan []byte
	readChannel     chan []byte
	WriteChannel    chan<- []byte
	ReadChannel     <-chan []byte
	finishedChannel chan bool
	FinishedChannel <-chan bool
	isFinished      bool
}

// NewClientConnection はClientConnectionを生成します
func NewClientConnection(conn net.Conn) ClientConnection {
	writeChannel := make(chan []byte, 10)
	readChannel := make(chan []byte, 10)
	finishedChannel := make(chan bool)
	logger := logger.Logger{RemoteAdderName: conn.RemoteAddr().String()}
	logger.PrintLn("接続")
	clientConnection := ClientConnection{
		conn:            conn,
		logger:          logger,
		writeChannel:    writeChannel,
		readChannel:     readChannel,
		WriteChannel:    writeChannel,
		ReadChannel:     readChannel,
		finishedChannel: finishedChannel,
		FinishedChannel: finishedChannel,
		isFinished:      false,
	}
	go clientConnection.clientProcess()
	return clientConnection
}

func (t *ClientConnection) close() {
	if t.isFinished == false {
		close(t.finishedChannel)
		t.isFinished = true
	}
}

//接続してきたクライアントに対する処理
func (t *ClientConnection) clientProcess() {

	buf := make([]byte, 4*1024)

	go func() {
		defer t.close()
		for {
			readBuf, err := read(t.conn, buf)
			if err != nil {
				return
			}

			readStr := string(readBuf)
			if readStr[:len(readStr)-1] != "empty" {
				t.logger.PrintLn("受信:" + readStr)
				t.readChannel <- readBuf
			}
		}
	}()

	go func() {
		defer t.close()
		for {
			writeBuf := <-t.writeChannel
			t.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			_, err2 := t.conn.Write(writeBuf)
			if err2 != nil {
				return
			}
			t.logger.PrintLn("送信:" + string(writeBuf))
		}
	}()

	<-t.finishedChannel
	t.conn.Close()
	t.logger.PrintLn("接続終了")
}

func read(conn net.Conn, buf []byte) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, err := conn.Read(buf)
	bufSliced := sliceBuf(buf)
	return bufSliced, err
}

func sliceBuf(buf []byte) []byte {
	return buf[:getBufEndIndex(buf)]
}

func getBufEndIndex(buf []byte) int {
	for index, b := range buf {
		if b == 0 {
			return index + 1
		}
	}
	return len(buf)
}
