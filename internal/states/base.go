package states

import (
	"fmt"
	"leader-election/internal/nodes"
	"strconv"
	"strings"
	"time"
)

const (
	StateLeader   = "leader"
	StatePassive  = "passive"
	StateElection = "election"

	requestMessage = "leader_election"
	identifier     = "node_status"
)

type Status struct {
	statusType string
	time       time.Time
}

type base struct {
	status string
	time   time.Time
}

func (s *base) Status(n *nodes.Node) []byte {
	return []byte(fmt.Sprintf("%s %s %s %d", identifier, n.GetIp(), s.status, s.time.Unix()))
}

func (s *base) GetState() string {
	return s.status
}

func (s *base) IsRequestMessage(b []byte) bool {
	return string(b) == requestMessage
}

func (s *base) Request(_ *nodes.Node) []byte {
	return []byte(requestMessage)
}

func (s *base) IsLeader() bool {
	return s.status == StateLeader
}

func (s *base) parseStatus(b []byte) (*Status, error) {
	splitedB := strings.Split(string(b), " ")
	if len(splitedB) == 0 || len(splitedB) > 4 {
		return nil, fmt.Errorf("wrong message lenght")
	}
	if splitedB[0] != identifier {
		return nil, fmt.Errorf("wrong message identifier")
	}
	sec, err := strconv.Atoi(splitedB[3])
	if err != nil {
		return nil, fmt.Errorf("time parsing failed: %w", err)
	}
	t := time.Unix(int64(sec), 0)
	if !strings.Contains(fmt.Sprintf("%s %s %s", StateLeader, StatePassive, StateElection), splitedB[2]) {
		return nil, fmt.Errorf("unknown node status: %s", splitedB[2])
	}
	status := &Status{
		statusType: splitedB[2],
		time:       t,
	}
	return status, nil
}
