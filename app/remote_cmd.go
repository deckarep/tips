package app

import (
	"bufio"
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"io"
	"os/exec"
	"sync/atomic"
	"time"
)

const (
	maxLinesToProcess    = 10
	maxCompletionTimeout = time.Millisecond * 1

	// TODO: consider prioritizing the built-in ssh into Tailscale which accounts for authentication via the Tailscale api.
	sshBin = "/usr/bin/ssh"
)

type hostLine struct {
	hostname string
	idx      int
	line     string
}

type chanCompletions struct {
	host      string
	idx       int
	completed bool
	ch        chan hostLine
}

func ExecuteClusterRemoteCmd(ctx context.Context, w io.Writer, hosts []string, remoteCmd string) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)
	startTime := time.Now()

	const (
		// TODO: Make this configurable.
		chanBuffer = 10
	)

	var (
		allCompletions []*chanCompletions
		sem            = make(chan struct{}, cfg.Concurrency)

		totalErrors  atomic.Uint32
		totalSuccess atomic.Uint32
	)

	// For each host, kick-off a goroutine to execute the remote command.
	for idx, host := range hosts {
		resultsChan := make(chan hostLine, chanBuffer)

		allCompletions = append(allCompletions, &chanCompletions{
			ch:        resultsChan,
			host:      host,
			idx:       idx,
			completed: false,
		})

		go func(i int, hn string, rch chan hostLine) {
			sem <- struct{}{}
			if err := executeRemoteCmd(ctx, i, hn, remoteCmd, rch); err != nil {
				totalErrors.Add(1)
				log.Error("error executing a remote command for", "host", hn, "error", err)
				return
			}
			totalSuccess.Add(1)
		}(idx, host, resultsChan)
	}

	poll(ctx, w, sem, allCompletions)

	// Prints a summary at the end of success vs failures as well as how long it took in seconds.
	summary := fmt.Sprintf("Finished: successes: %d, failures: %d, elapsed_sec: %.2f",
		totalSuccess.Load(),
		totalErrors.Load(),
		time.Since(startTime).Seconds())

	if _, err := fmt.Fprintln(w, summary); err != nil {
		log.Error("error on `Fprintln` when writing elapsed time", "error", err)
	}
}

func poll(ctx context.Context, w io.Writer, sem chan struct{}, allCompletions []*chanCompletions) {
	var totalCompleted int

	// Loop indefinitely until all totalCompleted are accounted for, then bail.
	for {
		if totalCompleted == len(allCompletions) {
			// We're completely done, so return.
			return
		}

		// Continually iterate through all completions and check that the following conditions:
		// 1. If we're completed for the first time, coordinate shutdown of this completion.
		// 2. If we're already completed, just skip
		// 3. Otherwise, if we're not completed consume up to maxLinesToProcess lines before moving onto the next
		// channel. We do this to best-effort cluster the output.

		for _, comp := range allCompletions {
		nextCompletion:
			// Attempt to read n lines from the channel before trying the next one, assuming they're ready.
			// This is intended to minimize the interleaving of lines across hosts.
			for i := 0; i < maxLinesToProcess; i++ {
				select {
				case stream, isOpen := <-comp.ch:
					if !isOpen && !comp.completed {
						// Mark this completion as done!
						comp.completed = true

						// Track how many completions are done.
						totalCompleted++

						// Mark sem for letting more work come in.
						<-sem

						// This completion is complete, move on to the next non-closed completion.
						break nextCompletion
					} else if comp.completed {
						// We've already drained this to completion, our work is done so skip.
						break nextCompletion
					}

					RenderLogLine(ctx, w, stream.idx, stream.hostname, stream.line)
				case <-time.After(maxCompletionTimeout):
					// We've waited long enough maybe another completion is ready.
					break nextCompletion
				}
			}
		}
	}
}

func executeRemoteCmd(ctx context.Context, idx int, host string, remoteCmd string, outputChan chan<- hostLine) error {
	// Construct the SSH command
	sshCmd := exec.Command(sshBin, host, remoteCmd)
	defer close(outputChan)

	// Get the output pipe
	stdout, err := sshCmd.StdoutPipe()
	if err != nil {
		return err
	}

	// Start the command
	if err := sshCmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		outputChan <- hostLine{
			idx:      idx,
			hostname: host,
			line:     scanner.Text(),
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Wait for the command to finish
	if err := sshCmd.Wait(); err != nil {
		return err
	}

	return nil
}
