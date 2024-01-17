package backup

import (
	"context"
	"errors"
	"fmt"
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var testChan chan struct{}

const backupSize uint64 = 133

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

type testCompressor struct {
	fail          bool
	artifacts     []string
	targetArchive string
}

func (t *testCompressor) Compress(ctx context.Context, targetArchive string, files []string) (uint64, error) {
	t.artifacts = files
	t.targetArchive = targetArchive
	if t.fail {
		return 0, errors.New("error")
	}
	return backupSize, nil
}

func (t *testCompressor) Extension() string {
	return "tc"
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
				Artifacts: []string{
					"stage1_1.bak", "stage1_2.bak",
				},
			},
			{
				Title: "Stage 2",
				Commands: []Command{
					&testCommand{t: t},
					&testCommand{t: t},
					&testCommand{t: t, wait: true},
				},
				Artifacts: []string{
					"stage2_1.bak",
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

func TestEngine_Launch_AllOK(t *testing.T) {
	tc := testCompressor{}
	e := NewEngine(&tc)
	set1 := makeInstructionSet(t)
	set2 := makeInstructionSet(t)
	assert.True(t, e.Launch(Context{}, rms_backup.BackupType_Full, set1))
	assert.False(t, e.Launch(Context{}, rms_backup.BackupType_Full, set2))
	testChan <- struct{}{}
	wait(e)

	report := e.GetReport()
	assert.Equal(t, Ready, report.Status)
	assert.Equal(t, tc.targetArchive, report.FileName)
	assert.Equal(t, rms_backup.BackupType_Full, report.Type)
	assert.Equal(t, backupSize, report.Size)
	assert.Nil(t, report.Errors)

	for _, stage := range set1.Stages {
		for _, cmd := range stage.Commands {
			tcmd := cmd.(*testCommand)
			assert.True(t, tcmd.executed)
			assert.True(t, tcmd.cleaned)
		}
	}
	assert.Equal(t, []string{"stage1_1.bak", "stage1_2.bak", "stage2_1.bak"}, tc.artifacts)
}

func TestEngine_Launch2_SecondFailed(t *testing.T) {
	tc := testCompressor{}
	e := NewEngine(&tc)
	set := makeInstructionSet(t)
	failedCommand := set.Stages[1].Commands[1].(*testCommand)
	failedCommand.executeErr = errors.New("error")

	assert.True(t, e.Launch(Context{}, rms_backup.BackupType_Partial, set))
	wait(e)

	report := e.GetReport()
	assert.Equal(t, ReadyWithErrors, report.Status)
	assert.Equal(t, tc.targetArchive, report.FileName)
	assert.Equal(t, rms_backup.BackupType_Partial, report.Type)
	assert.Equal(t, backupSize, report.Size)
	assert.Equal(t, 1, len(report.Errors))

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
	assert.Equal(t, []string{"stage1_1.bak", "stage1_2.bak"}, tc.artifacts)
}

func TestEngine_Launch3_FirstFailed(t *testing.T) {
	tc := testCompressor{}
	e := NewEngine(&tc)
	set := makeInstructionSet(t)
	failedCommand := set.Stages[0].Commands[2].(*testCommand)
	failedCommand.executeErr = errors.New("error")

	assert.True(t, e.Launch(Context{}, rms_backup.BackupType_Full, set))
	testChan <- struct{}{}
	wait(e)

	report := e.GetReport()
	assert.Equal(t, ReadyWithErrors, report.Status)
	assert.Equal(t, tc.targetArchive, report.FileName)
	assert.Equal(t, rms_backup.BackupType_Full, report.Type)
	assert.Equal(t, backupSize, report.Size)
	assert.Equal(t, 1, len(report.Errors))

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

	assert.Equal(t, []string{"stage2_1.bak"}, tc.artifacts)
}

func TestEngine_Launch4_Panics(t *testing.T) {
	tc := testCompressor{}
	e := NewEngine(&tc)
	set := makeInstructionSet(t)
	failedCommand1 := set.Stages[0].Commands[2].(*testCommand)
	failedCommand2 := set.Stages[1].Commands[0].(*testCommand)
	failedCommand1.executePanic = true
	failedCommand2.cleanUpPanic = true

	assert.True(t, e.Launch(Context{}, rms_backup.BackupType_Full, set))
	testChan <- struct{}{}
	wait(e)

	report := e.GetReport()
	assert.Equal(t, ReadyWithErrors, report.Status)
	assert.Equal(t, tc.targetArchive, report.FileName)
	assert.Equal(t, rms_backup.BackupType_Full, report.Type)
	assert.Equal(t, backupSize, report.Size)
	assert.Equal(t, 1, len(report.Errors))

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

	assert.Equal(t, []string{"stage2_1.bak"}, tc.artifacts)
}

func TestEngine_Launch5_AllFailed(t *testing.T) {
	tc := testCompressor{}
	e := NewEngine(&tc)
	set := makeInstructionSet(t)
	failedCommand1 := set.Stages[0].Commands[2].(*testCommand)
	failedCommand2 := set.Stages[1].Commands[0].(*testCommand)
	failedCommand1.executePanic = true
	failedCommand2.executeErr = errors.New("error")

	assert.True(t, e.Launch(Context{}, rms_backup.BackupType_Full, set))
	wait(e)

	report := e.GetReport()
	assert.Equal(t, Failed, report.Status)
	assert.Equal(t, rms_backup.BackupType_Full, report.Type)
	assert.Equal(t, uint64(0), report.Size)
	assert.Equal(t, 2, len(report.Errors))

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

	assert.Nil(t, tc.artifacts)
}

func TestEngine_Launch_CompressFailed(t *testing.T) {
	tc := testCompressor{fail: true}
	e := NewEngine(&tc)
	set1 := makeInstructionSet(t)
	set2 := makeInstructionSet(t)
	assert.True(t, e.Launch(Context{}, rms_backup.BackupType_Full, set1))
	assert.False(t, e.Launch(Context{}, rms_backup.BackupType_Full, set2))
	testChan <- struct{}{}
	wait(e)

	report := e.GetReport()
	assert.Equal(t, Failed, report.Status)
	assert.Equal(t, rms_backup.BackupType_Full, report.Type)
	assert.Equal(t, uint64(0), report.Size)
	assert.Equal(t, 1, len(report.Errors))

	for _, stage := range set1.Stages {
		for _, cmd := range stage.Commands {
			tcmd := cmd.(*testCommand)
			assert.True(t, tcmd.executed)
			assert.True(t, tcmd.cleaned)
		}
	}
	assert.Equal(t, []string{"stage1_1.bak", "stage1_2.bak", "stage2_1.bak"}, tc.artifacts)
}

func TestEngine_GetReport(t *testing.T) {
	e := NewEngine(&testCompressor{})
	report := e.GetReport()
	assert.Equal(t, NeverRun, report.Status)

	set := makeInstructionSet(t)
	assert.True(t, e.Launch(Context{}, rms_backup.BackupType_Full, set))

	<-time.After(100 * time.Millisecond)
	report = e.GetReport()
	assert.Equal(t, InProgress, report.Status)

	testChan <- struct{}{}
	wait(e)
	report = e.GetReport()
	assert.Equal(t, Ready, report.Status)
}

func TestEngine_Shutdown(t *testing.T) {
	e := NewEngine(&testCompressor{})
	set := makeInstructionSet(t)
	waitCommand := set.Stages[0].Commands[2].(*testCommand)
	waitCommand.wait = true

	assert.True(t, e.Launch(Context{}, rms_backup.BackupType_Full, set))
	<-time.After(100 * time.Millisecond)
	e.Shutdown()
	wait(e)
}
