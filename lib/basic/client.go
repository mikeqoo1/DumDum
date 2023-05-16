package basic

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// 定義 TCPClient 結構
type TCPClient struct {
	conn      net.Conn
	SendCh    chan string
	ReceiveCh chan string
}

// 連接到服務器
func (c *TCPClient) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *TCPClient) Close() {
	c.conn.Close()
}

// 接收消息
func (c *TCPClient) ReceiveMessages() {
	fmt.Println(2)
	defer close(c.ReceiveCh)
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		c.ReceiveCh <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("ReceiveMessages解析scanner錯誤:", err)
	}
}

// 發送消息
func (c *TCPClient) SendMessages() {
	fmt.Println(1)
	msgCh := <-c.SendCh
	fmt.Println("要送的Msg:" + msgCh)
	writer := bufio.NewWriter(c.conn)
	scanner := bufio.NewScanner(strings.NewReader(msgCh))
	for scanner.Scan() {
		msg := scanner.Text()
		_, err := fmt.Fprintln(writer, msg)
		if err != nil {
			fmt.Println("Error writing:", err)
			break
		}
		err = writer.Flush()
		if err != nil {
			fmt.Println("Error flushing writer:", err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("SendMessages解析scanner錯誤:", err)
	}
}
