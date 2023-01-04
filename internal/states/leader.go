package states

import (
	"fmt"
	"leader-election/internal/nodes"
	"time"
)

func NewLeader() *Leader {
	return &Leader{
		base: base{
			status: StateLeader,
			time:   time.Now(),
		},
	}
}

type Leader struct {
	base
}

func (s *Leader) Handle(n *nodes.Node, b []byte) error {
	if b == nil {
		return nil
	}
	status, err := s.parseStatus(b)
	if err != nil {
		return fmt.Errorf("handle messsage failed: %w", err)
	}
	if status.statusType == StateLeader && status.time.Before(s.time) {
		n.SetState(NewPassive())
	}
	return nil
}
