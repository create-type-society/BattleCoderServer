package connection

import (
	"BattleCoderServer/client/connection/logger"
	"bufio"
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
	reader          *bufio.Reader
	writer          *bufio.Writer
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
		writer:          bufio.NewWriter(conn),
		reader:          bufio.NewReader(conn),
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

	go func() {
		defer t.close()
		for {
			t.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
			readBuf, err := t.reader.ReadBytes('\n')
			if err != nil {
				return
			}
			readStr := string(readBuf)
			if readStr != "empty\n" {
				t.logger.Print("受信:" + readStr)
				t.readChannel <- readBuf
			}
		}
	}()

	go func() {
		defer t.close()
		for {
			writeBuf := <-t.writeChannel
			t.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			_, err2 := t.writer.Write(writeBuf)
			t.writer.Flush()
			if err2 != nil {
				return
			}
			writeStr := string(writeBuf)
			t.logger.Print("送信:" + writeStr)
		}
	}()

	<-t.finishedChannel
	t.conn.Close()
	t.logger.PrintLn("接続終了")
}
