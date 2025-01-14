package main

import (
	"sync"
	"time"

	"github.com/gosuri/uiprogress"
)

func main() {
	waitTime := time.Millisecond * 100
	uiprogress.Start()

	var wg sync.WaitGroup

	bar1 := uiprogress.AddBar(20).AppendCompleted().PrependElapsed() // total is 20
	wg.Add(1)
	go func() {
		defer wg.Done()
		for bar1.Incr() {
			time.Sleep(waitTime) // something you do
		}
	}()

	bar2 := uiprogress.AddBar(40).AppendCompleted().PrependCompleted() // total is 40
	wg.Add(1)
	go func() {
		defer wg.Done()
		for bar2.Incr() {
			time.Sleep(waitTime)
		}
	}()

	time.Sleep(time.Second)
	bar3 := uiprogress.AddBar(20).PrependElapsed().AppendCompleted()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for bar3.Incr() {
			time.Sleep(waitTime)
		}
	}()

	wg.Wait()
}
