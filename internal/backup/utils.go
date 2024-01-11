package backup

import "fmt"

func safeCommandExecute(ctx Context, command Command) error {
	return safeCall(func() error {
		return command.Execute(ctx)
	})
}

func safeCommandCleanup(ctx Context, command Command) error {
	return safeCall(func() error {
		return command.Cleanup(ctx)
	})
}

func safeCall(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %s", r)
		}
	}()
	return f()
}
