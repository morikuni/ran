package ran

import (
	"os/exec"
)

type Stack struct {
	cmds []*exec.Cmd
}

func NewStack() *Stack {
	return &Stack{}
}

func (s *Stack) Push(cmd *exec.Cmd) {
	s.cmds = append(s.cmds, cmd)
}

func (s *Stack) Pop() (*exec.Cmd, bool) {
	if len(s.cmds) == 0 {
		return nil, false
	}
	cmd := s.cmds[len(s.cmds)-1]
	s.cmds = s.cmds[:len(s.cmds)-1]
	return cmd, true
}
