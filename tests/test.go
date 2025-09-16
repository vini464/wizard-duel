package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/vini464/wizard-duel/share"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Escolha quantos bots voce quer:")
  var wg sync.WaitGroup
	scanner.Scan()
	scanner.Scan()
	input := scanner.Text()
	bots, err := strconv.Atoi(input)
	end := make(chan bool)
	if err == nil {
		for range bots {
      wg.Add(1)
			go createClient(end, &wg)
		}
    end <-true
	}
  fmt.Println("[press enter to end test]")
  scanner.Scan()
  wg.Wait()
}

func createClient(end chan bool, wg *sync.WaitGroup) {
	conn, err := net.Dial(share.SERVERTYPE, net.JoinHostPort(share.SERVERNAME, share.SERVERPORT))
	for err == nil && !<-end {
		message := share.Message{Type: share.ECHO, Data: []byte("Hello")}
		err = share.SendMessage(conn, message)
		if err == nil {
			share.ReceiveMessage(conn, &message)
			println("Type;", message.Type, "data:", string(message.Data), "latency:", message.TimeStamp)
		}
	}
  wg.Done()
}
