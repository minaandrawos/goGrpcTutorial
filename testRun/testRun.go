package main

import (
	"log"
	"os/exec"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	for i := 0; i <= 20; i++ {
		wg.Add(1)
		go func() {
			cmd := exec.Command("client.exe")
			if err := cmd.Start(); err != nil {
				log.Fatal(err)
			}
			cmd.Wait()
			wg.Done()
		}()
	}
	wg.Wait()
}
