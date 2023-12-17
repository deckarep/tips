package app

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

const (
	sshBin = "/usr/bin/ssh"
)

type hostLine struct {
	hostname string
	idx      int
	line     string
}

func ExecuteClusterRemoteCmd(ctx context.Context, hosts []string, remoteCmd string) {
	const (
		chanBuffer        = 10
		maxLinesToProcess = 5
		parallelism       = 5
	)

	type chanCompletions struct {
		host      string
		idx       int
		completed bool
		ch        chan hostLine
	}

	var allCompletions []*chanCompletions
	sem := make(chan struct{}, parallelism)
	var wg sync.WaitGroup
	wg.Add(len(hosts))

	for idx, host := range hosts {
		resultsChan := make(chan hostLine, chanBuffer)
		allCompletions = append(allCompletions, &chanCompletions{ch: resultsChan, host: host, idx: idx, completed: false})
		go func(i int, hn string, rch chan hostLine) {
			sem <- struct{}{}
			ExecuteRemoteCmd(ctx, i, hn, remoteCmd, rch)
		}(idx, host, resultsChan)
	}

	go func() {
		var totalCompleted int
		for {
			if totalCompleted == len(allCompletions) {
				// We're completely done.
				//log.Print("are we done?", "totalCompleted=", totalCompleted, "len(allCompletions)", len(allCompletions))
				return
			}

			for _, compl := range allCompletions {
			nextCompletion:
				// Attempt to read n lines from the channel before trying the next one.
				for i := 0; i < maxLinesToProcess; i++ {
					select {
					case stream, isOpen := <-compl.ch:
						if !isOpen && !compl.completed {
							// Mark this completion as done!
							//fmt.Printf("host: %s (%d) is completed!\n", compl.host, compl.idx)
							compl.completed = true

							// Track how many completions are done.
							totalCompleted++

							// Done fully consuming.
							wg.Done()
							<-sem

							// Move on to the next channel.
							break nextCompletion
						} else if compl.completed {
							// We've already drained this to completion, our work is done so skip.
							break nextCompletion
						}

						fmt.Println(fmt.Sprintf("%s (%d): %s ", stream.hostname, stream.idx, stream.line))
					case <-time.After(10 * time.Millisecond):
						// We've waited long enough maybe another chan is ready.
						break nextCompletion
					}
				}
			}
		}
	}()
	wg.Wait()
}

func ExecuteRemoteCmd(ctx context.Context, idx int, host string, remoteCmd string, outputChan chan<- hostLine) error {

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
	//log.Warn("waiting on Wait()")
	if err := sshCmd.Wait(); err != nil {
		return err
	}

	//log.Warn("Finished because we can return.")
	return nil
}
