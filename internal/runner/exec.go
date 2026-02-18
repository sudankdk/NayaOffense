package runner

import (
	"bufio"
	"context"
	"os/exec"
	"time"
)

// func Execute(ctx context.Context, bin string, args ...string) (string, error) {
// 	defaultTimeout := 60 * time.Second
// 	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
// 	defer cancel()
// 	cmd := exec.CommandContext(ctx, bin, args...)
// 	output, err := cmd.CombinedOutput()
// 	return string(output), err
// }

// event based banaunu parxa
// mathi ko fail hunxa if the reponse from the tool is taking long time or dherai lamo response xa vane

func Execute(ctx context.Context, bin string, args ...string) chan Event {
	events := make(chan Event) // unbuffred chan for events
	go func() {
		cmd := exec.CommandContext(ctx, bin, args...)

		//pipes for attaching
		// stdout or stderr
		stdOut, err := cmd.StdoutPipe()
		if err != nil {
			events <- Event{Type: EventError, Data: err.Error(), Timestamp: time.Now()}
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			events <- Event{Type: EventError, Data: err.Error(), Timestamp: time.Now()}
			return
		}

		// aba we take the process hai
		if err := cmd.Start(); err != nil {
			events <- Event{Type: EventError, Data: err.Error(), Timestamp: time.Now()}
			return
		}

		events <- Event{
			Type:      EventStart,
			Data:      bin,
			Timestamp: time.Now(),
		}

		go func() {
			scanner := bufio.NewScanner(stdOut)
			for scanner.Scan() {
				events <- Event{
					Type:      EventStdout,
					Data:      scanner.Text(),
					Timestamp: time.Now(),
				}
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				events <- Event{
					Type:      EventError,
					Data:      scanner.Text(),
					Timestamp: time.Now(),
				}
			}
		}()

		if err := cmd.Wait(); err != nil {
			events <- Event{
				Type:      EventError,
				Data:      err.Error(),
				Timestamp: time.Now(),
			}
		}

		events <- Event{
			Type:      EventComplete,
			Data:      "process finished",
			Timestamp: time.Now(),
		}

	}()
	return events
}
