package ran

type Cmd interface {
	Run() error
}

type Stack interface {
	Push(cmd Cmd)
	Pop() (Cmd, bool)
}

type stack struct {
	cmds []Cmd
}

func NewStack() Stack {
	return &stack{}
}

func (s *stack) Push(cmd Cmd) {
	s.cmds = append(s.cmds, cmd)
}

func (s *stack) Pop() (Cmd, bool) {
	if len(s.cmds) == 0 {
		return nil, false
	}
	cmd := s.cmds[len(s.cmds)-1]
	s.cmds = s.cmds[:len(s.cmds)-1]
	return cmd, true
}
