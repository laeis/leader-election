package nodes

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"strconv"
	"time"
)

type Stater interface {
	Request(*Node) []byte
	Handle(*Node, []byte) error
	Status(*Node) []byte
	IsLeader() bool
	IsRequestMessage([]byte) bool
	GetState() string
}

type Node struct {
	State Stater
	Ip    net.IP
}

func (n *Node) InitConnection(port string) (*net.UDPConn, error) {
	// Get preferred outbound ip of this machine
	localAddress, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		log.Error().Err(err).Msg("attempt ResolvingUDPAddr failed")
		return nil, err
	}
	log.Info().Msgf("Listening to port: %d", localAddress.Port)
	return net.ListenUDP("udp", localAddress)
}

func (n *Node) Broadcast(conn *net.UDPConn) {
	//ask about main nodes

	remote, err := net.ResolveUDPAddr("udp4", "255.255.255.255:8829")
	if err != nil {
		panic(err)
	}
	_, err = conn.WriteTo(n.State.Request(n), remote)
	if err != nil {
		panic(err)
	}
}

func (n *Node) Handle(ctx context.Context, conn *net.UDPConn) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("listen udp connection closed by timeout")
			return
		default:
			buffer := make([]byte, 1024)
			if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
				log.Error().Err(err).Msg("set read deadline for connection failed")
				continue
			}

			length, addr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				opError, ok := err.(*net.OpError)
				if !ok || !opError.Timeout() {
					log.Error().Err(opError).Msg("read UDP response failed")
					continue
				}
				_ = n.State.Handle(n, nil)
				continue
			}

			localAddr := conn.LocalAddr().(*net.UDPAddr)
			portSting := strconv.Itoa(localAddr.Port)
			if addr.String() == n.Ip.String()+":"+portSting {
				continue
			}

			buffer = buffer[:length]
			if n.State.IsRequestMessage(buffer) {
				_, err = conn.WriteToUDP(n.State.Status(n), addr)
				if err != nil {
					log.Error().Err(err).Msg("send status failed")
				}
				continue
			}

			err = n.State.Handle(n, buffer)
			if err != nil {
				log.Error().Err(err).Msg("handle status message failed")
			}
		}
	}
}

func (n *Node) SetState(s Stater) {
	fmt.Println("Set new status: ", s.GetState())
	n.State = s
}

func (n *Node) GetIp() string {
	return n.Ip.String()
}

func (n *Node) IsLeader() bool {
	return n.State.IsLeader()
}

func (n *Node) Status() string {
	return n.State.GetState()
}
