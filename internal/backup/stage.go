package backup

// Command is an interface one 1 backup command
type Command interface {
	Title() string
	Execute(ctx Context) error
	Cleanup(ctx Context) error
}

// Stage is one atomic backup operation
type Stage struct {
	Title    string
	Commands []Command
}

// Instruction is set of instructions
type Instruction struct {
	Title  string
	Stages []Stage
}

func (s *Stage) Add(cmd Command) {
	s.Commands = append(s.Commands, cmd)
}

func (i *Instruction) Add(stage Stage) {
	i.Stages = append(i.Stages, stage)
}

func (i *Instruction) Operations() int {
	cnt := 0
	for j := range i.Stages {
		cnt += len(i.Stages[j].Commands)
	}
	return cnt
}
