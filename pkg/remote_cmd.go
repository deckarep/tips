/*
Open Source Initiative OSI - The MIT License (MIT):Licensing

The MIT License (MIT)
Copyright (c) 2023 - 2024 Ralph Caraveo (deckarep@gmail.com)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package pkg

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/log"
)

const (
	maxLinesToCluster    = 10
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
		wg           sync.WaitGroup
	)

	wg.Add(len(hosts))

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
			defer wg.Done()
			if err := executeRemoteCmd(ctx, i, hn, remoteCmd, rch); err != nil {
				totalErrors.Add(1)
				log.Error("error executing remote command for", "host", hn, "cmd", remoteCmd, "error", err)
				return
			}
			totalSuccess.Add(1)
		}(idx, host, resultsChan)
	}

	// This blocks until all completions have shutdown.
	// However, upon an early ssh connection error this polling will immediately fallthrough.
	poll(ctx, w, sem, allCompletions)

	// But we still want to wait for all goroutines executed above to run to completion.
	wg.Wait()

	// Prints a summary at the end of success vs failures as well as how long it took in seconds.
	summary := fmt.Sprintf("Finished: successes: %d, failures: %d, elapsed_sec: %.2f",
		totalSuccess.Load(),
		totalErrors.Load(),
		time.Since(startTime).Seconds())

	if _, err := fmt.Fprintln(w, summary); err != nil {
		log.Error("error on `Fprintln` when writing elapsed time", "error", err)
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
			for i := 0; i < maxLinesToCluster; i++ {
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
