package ran

type Stack interface {
	Push(script Script)
	Pop() (Script, bool)
}

type stack struct {
	scripts []Script
}

func NewStack() Stack {
	return &stack{}
}

func (s *stack) Push(script Script) {
	s.scripts = append(s.scripts, script)
}

func (s *stack) Pop() (Script, bool) {
	if len(s.scripts) == 0 {
		return nil, false
	}
	script := s.scripts[len(s.scripts)-1]
	s.scripts = s.scripts[:len(s.scripts)-1]
	return script, true
}
