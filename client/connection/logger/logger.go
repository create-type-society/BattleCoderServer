package logger

import "fmt"

// Logger は接続してきたクライアントの状態を標準出力するためのもの
type Logger struct {
	RemoteAdderName string
}

// PrintLn は関連しているRemoteAdderNameのログを表示します
func (t Logger) PrintLn(s string) {
	fmt.Println("[" + t.RemoteAdderName + "] " + s)
}
