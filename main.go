package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"leader-election/internal/nodes"
	"leader-election/internal/states"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal().Err(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {
	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)

	n := &nodes.Node{
		State: states.NewPassive(),
		Ip:    GetOutboundIP(),
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	connection, err := n.InitConnection("8829")
	if err != nil {
		log.Error().Err(err).Msg("attempt ListenUDP failed")
		return
	}
	defer connection.Close()

	go n.Handle(ctx, connection)
	iteration := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-iteration.C:
			fmt.Println("Current status:", n.Status())
			n.Broadcast(connection)
		case signalReceived := <-s:
			fmt.Printf("Exit with signal: %s", signalReceived.String())
			return
		}
	}
	os.Exit(0)
}
