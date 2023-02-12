package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
)

func main() {
	// 这个是在checkTx中起作用的
	// 需要改变ip地址。包括组内各个节点的ip地址
	// 存在问题，好像不需要这样
	// leader往本地ip端口写东西，然后其他replica监听该窗口，并处理该窗口发送的信息
	// 为什么可以，因为websocket本身就是双向的，理论上只需建立replica到leader的通道就好
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/"}
	// 这样可以建立联系，并准备发送消息
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("Dial error:", err)
		return
	}
	// 关闭conn，不再联系
	defer conn.Close()

	fmt.Println("Connected to WebSocket server")

	// 利用协程来读消息，并存起来
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Read error:", err)
				return
			}
			fmt.Printf("Received message: %s\n", msg)
		}
	}()

	// 这部分写到leader中，让leaderconn.WriteMessage
	//scanner := bufio.NewScanner(os.Stdin)
	for {
		err := conn.WriteMessage(websocket.BinaryMessage, []byte("1"))
		if err != nil {
			fmt.Println("Write error:", err)
			return
		}
	}
}
