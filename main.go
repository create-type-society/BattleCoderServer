package main

import (
	"BattleCoderServer/client"
	"fmt"
	"net"
	"time"
)

const host = "localhost:3000"

func main() {
	listen()
}

//tcpサーバの起動と処理
func listen() {
	listen, _ := net.Listen("tcp", host)
	fmt.Println("サーバ起動@" + host)
	for {
		conn, _ := listen.Accept()
		go clientProcess(conn)
	}
}

//接続してきたクライアントに対する処理
func clientProcess(conn net.Conn) {

	logger := client.Logger{RemoteAdderName: conn.RemoteAddr().String()}
	logger.PrintLn("接続")
	buf := make([]byte, 4*1024)
	defer conn.Close()
	defer logger.PrintLn("接続終了")
	for {
		readBuf, err := read(conn, buf)
		if err != nil {
			return
		}
		conn.Write(readBuf)
		logger.PrintLn("受信:" + string(readBuf))
	}
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
			return index
		}
	}
	return len(buf)
}
