package states

import (
	"fmt"
	"leader-election/internal/nodes"
	"time"
)

func NewElection() *Election {
	return &Election{
		base: base{
			status: StateElection,
			time:   time.Now(),
		},
	}
}

type Election struct {
	base
}

func (s *Election) Handle(n *nodes.Node, b []byte) error {
	if b == nil {
		n.SetState(NewLeader())
		return nil
	}
	status, err := s.parseStatus(b)
	if err != nil {
		return fmt.Errorf("handle messsage failed: %w", err)
	}
	if status.statusType == StateLeader || (status.statusType == StateElection && status.time.Before(s.time)) {
		n.SetState(NewPassive())
	}
	return nil
}
