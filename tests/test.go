package main

import (
  "net"
  "github.com/vini464/wizard-duel/share"
)


func main(){
  conn, err := net.Dial(share.SERVERTYPE,net.JoinHostPort(share.SERVERNAME, share.SERVERPORT))
  for err == nil {
    message := share.Message{Type: share.ECHO, Data: []byte("Hello")}
    err = share.SendMessage(conn, message)
    if err == nil {
      share.ReceiveMessage(conn, &message)
      println("Type;", message.Type, "data:", string(message.Data), "latency:", message.TimeStamp)
    }
  }
}
