package main

import (
	"BattleCoderServer/client"
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
	for {
		_, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			logger.PrintLn("接続終了")
			return
		}
		bufSliced := sliceBuf(buf)
		conn.Write(bufSliced)
		logger.PrintLn("受信:" + string(bufSliced))
	}
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
