package states

import (
	"leader-election/internal/nodes"
	"time"
)

func NewPassive() *Passive {
	return &Passive{
		base: base{
			status: StatePassive,
			time:   time.Now(),
		},
	}
}

type Passive struct {
	base
}

func (s *Passive) Request(n *nodes.Node) []byte {
	n.SetState(NewElection())
	return []byte(requestMessage)
}

func (s *Passive) Handle(_ *nodes.Node, _ []byte) error {
	return nil
}
