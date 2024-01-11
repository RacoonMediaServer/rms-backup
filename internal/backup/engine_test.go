package backup

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var testChan chan struct{}

func init() {
	testChan = make(chan struct{})
}

type testCommand struct {
	t            *testing.T
	executed     bool
	cleaned      bool
	executeErr   error
	cleanErr     error
	wait         bool
	executePanic bool
	cleanUpPanic bool
}

func (c *testCommand) Title() string {
	return fmt.Sprintf("Command %p", c)
}

func (c *testCommand) Execute(ctx Context) error {
	c.executed = true
	if c.wait {
		select {
		case <-testChan:
		case <-ctx.Done():
		}
	}
	if c.executePanic {
		panic("test")
	}
	return c.executeErr
}

func (c *testCommand) Cleanup(ctx Context) error {
	c.cleaned = true
	if c.cleanUpPanic {
		panic("test")
	}
	return c.cleanErr
}

func makeInstructionSet(t *testing.T) Instruction {
	return Instruction{
		Title: "Test set",
		Stages: []Stage{
			{
				Title: "Stage 1",
				Commands: []Command{
					&testCommand{t: t},
					&testCommand{t: t},
					&testCommand{t: t},
				},
			},
			{
				Title: "Stage 2",
				Commands: []Command{
					&testCommand{t: t},
					&testCommand{t: t},
					&testCommand{t: t, wait: true},
				},
			},
		},
	}
}

func wait(e *Engine) {
	for {
		report := e.GetReport()
		if report.Status != InProgress {
			return
		}
		<-time.After(10 * time.Millisecond)
	}
}

func TestEngine_Launch(t *testing.T) {
	e := NewEngine()
	set1 := makeInstructionSet(t)
	set2 := makeInstructionSet(t)
	assert.True(t, e.Launch(Context{}, set1))
	assert.False(t, e.Launch(Context{}, set2))
	testChan <- struct{}{}
	wait(e)

	report := e.GetReport()
	assert.Equal(t, Ready, report.Status)

	for _, stage := range set1.Stages {
		for _, cmd := range stage.Commands {
			tcmd := cmd.(*testCommand)
			assert.True(t, tcmd.executed)
			assert.True(t, tcmd.cleaned)
		}
	}
}

func TestEngine_Launch2(t *testing.T) {
	e := NewEngine()
	set := makeInstructionSet(t)
	failedCommand := set.Stages[1].Commands[1].(*testCommand)
	failedCommand.executeErr = errors.New("error")

	assert.True(t, e.Launch(Context{}, set))
	wait(e)

	report := e.GetReport()
	assert.Equal(t, Failed, report.Status)

	for _, cmd := range set.Stages[0].Commands {
		tcmd := cmd.(*testCommand)
		assert.True(t, tcmd.executed)
		assert.True(t, tcmd.cleaned)
	}

	for i, cmd := range set.Stages[1].Commands {
		tcmd := cmd.(*testCommand)
		if i == 0 {
			assert.True(t, tcmd.executed)
			assert.True(t, tcmd.cleaned)
		} else if i == 1 {
			assert.True(t, tcmd.executed)
			assert.False(t, tcmd.cleaned)
		} else {
			assert.False(t, tcmd.executed)
			assert.False(t, tcmd.cleaned)
		}
	}
}

func TestEngine_Launch3(t *testing.T) {
	e := NewEngine()
	set := makeInstructionSet(t)
	failedCommand := set.Stages[0].Commands[2].(*testCommand)
	failedCommand.executeErr = errors.New("error")

	assert.True(t, e.Launch(Context{}, set))
	testChan <- struct{}{}
	wait(e)

	report := e.GetReport()
	assert.Equal(t, ReadyWithErrors, report.Status)

	for i, cmd := range set.Stages[0].Commands {
		tcmd := cmd.(*testCommand)
		if i != 2 {
			assert.True(t, tcmd.executed)
			assert.True(t, tcmd.cleaned)
		} else {
			assert.True(t, tcmd.executed)
			assert.False(t, tcmd.cleaned)
		}
	}

	for _, cmd := range set.Stages[1].Commands {
		tcmd := cmd.(*testCommand)

		assert.True(t, tcmd.executed)
		assert.True(t, tcmd.cleaned)
	}
}

func TestEngine_Launch4(t *testing.T) {
	e := NewEngine()
	set := makeInstructionSet(t)
	failedCommand1 := set.Stages[0].Commands[2].(*testCommand)
	failedCommand2 := set.Stages[1].Commands[0].(*testCommand)
	failedCommand1.executePanic = true
	failedCommand2.cleanUpPanic = true

	assert.True(t, e.Launch(Context{}, set))
	testChan <- struct{}{}
	wait(e)

	report := e.GetReport()
	assert.Equal(t, ReadyWithErrors, report.Status)

	for i, cmd := range set.Stages[0].Commands {
		tcmd := cmd.(*testCommand)
		if i != 2 {
			assert.True(t, tcmd.executed)
			assert.True(t, tcmd.cleaned)
		} else {
			assert.True(t, tcmd.executed)
			assert.False(t, tcmd.cleaned)
		}
	}

	for _, cmd := range set.Stages[1].Commands {
		tcmd := cmd.(*testCommand)

		assert.True(t, tcmd.executed)
		assert.True(t, tcmd.cleaned)
	}
}

func TestEngine_GetReport(t *testing.T) {
	e := NewEngine()
	report := e.GetReport()
	assert.Equal(t, NeverRun, report.Status)

	set := makeInstructionSet(t)
	assert.True(t, e.Launch(Context{}, set))

	<-time.After(100 * time.Millisecond)
	report = e.GetReport()
	assert.Equal(t, InProgress, report.Status)

	testChan <- struct{}{}
	wait(e)
	report = e.GetReport()
	assert.Equal(t, Ready, report.Status)
}

func TestEngine_Shutdown(t *testing.T) {
	e := NewEngine()
	set := makeInstructionSet(t)
	waitCommand := set.Stages[0].Commands[2].(*testCommand)
	waitCommand.wait = true

	assert.True(t, e.Launch(Context{}, set))
	<-time.After(100 * time.Millisecond)
	e.Shutdown()
	wait(e)
}
