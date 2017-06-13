package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	myMutex := sync.RWMutex{}
	for i := 0; i < 10; i++ {
		go acquiringFunc(&myMutex, i)
	}

	fmt.Println("main will now go to sleep for 10 seconds")

	time.Sleep(1 * time.Second)
	nowTime := time.Now()
	myMutex.Lock()
	timeSince := time.Since(nowTime)
	fmt.Println(timeSince)
	fmt.Println("main got the lock")
	time.Sleep(10 * time.Second)
	myMutex.Unlock()
	time.Sleep(10 * time.Second)

	fmt.Println("waking up after second sleep")
}

func acquiringFunc(myMutex *sync.RWMutex, i int) {
	myMutex.RLock()
	fmt.Printf("I am goroutine no. %d, and I got the read lock!\n", i)
	myMutex.RUnlock()
}
