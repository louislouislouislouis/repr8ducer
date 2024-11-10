package utils

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/rs/zerolog"
)

type Command struct {
	*exec.Cmd
}

// NewCmdWrapper creates and returns a new CmdWrapper with the given command and arguments.
func NewCommand(name string, args ...string) *Command {
	cmd := exec.Command(name, args...)
	return &Command{cmd}
}

func (cmd Command) ExecuteWithLiveOutput() error {
	// Log the Command
	Log.WithLevel(zerolog.InfoLevel).
		Msg(
			fmt.Sprintf("%sCommand started --> %s%s",
				COLOR_CYAN,
				COLOR_RESET,
				cmd.String(),
			),
		)

	// Create pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		Log.WithLevel(zerolog.FatalLevel).Err(err).Msg("Cannot create stdout pipe")
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		Log.WithLevel(zerolog.FatalLevel).Err(err).Msg("Cannot create stderr pipe")
		return err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		Log.WithLevel(zerolog.FatalLevel).Err(err).Msg("Cannot start the command")
		return err
	}

	// Create channels to handle stdout and stderr
	outputCh := make(chan string)
	errCh := make(chan string)

	// Function to read from a reader and send the output to a channel
	readOutput := func(r io.Reader, ch chan<- string) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			ch <- scanner.Text() // Send each line to the channel
		}
		close(ch) // Close the channel when done
	}

	// Go routines to read output
	go readOutput(stdoutPipe, outputCh)
	go readOutput(stderrPipe, errCh)

	// Use select to handle output from both channels
	for outputCh != nil || errCh != nil {
		select {
		case line, ok := <-outputCh:
			if !ok {
				outputCh = nil
			} else {
				Log.
					WithLevel(zerolog.InfoLevel).
					Msg(fmt.Sprintf("%sStdout command --> %s%s", COLOR_YELLOW, COLOR_RESET, line))
			}
		case line, ok := <-errCh:
			if !ok {
				errCh = nil
			} else {
				Log.WithLevel(zerolog.InfoLevel).Msg(fmt.Sprintf("%sStderr command --> %s%s", COLOR_PURPLE, COLOR_RESET, line))
			}
		}
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		Log.WithLevel(zerolog.FatalLevel).
			Err(err).
			Msg("Unexpected error while waiting for the command to finish")
		return err
	}

	// Log the command
	Log.WithLevel(zerolog.InfoLevel).
		Msg(
			fmt.Sprintf("%sCommand finished: --> %s%s",
				COLOR_CYAN,
				COLOR_RESET,
				cmd.String(),
			),
		)

	return nil
}
