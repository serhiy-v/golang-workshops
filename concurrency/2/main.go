// Hint 1: time.Ticker can be used to cancel function
// Hint 2: to calculate time-diff for Advanced lvl use:
//  start := time.Now()
//	// your work
//	t := time.Now()
//	elapsed := t.Sub(start) // 1s or whatever time has passed

package main

import (
	"math"
	"sync/atomic"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool  // can be used for 2nd level task. Premium users won't have 10 seconds limit.
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false if process had to be killed
func HandleRequest(process func(), u *User) bool {
	done := make(chan bool)
	go func() {
		process()
		done <- true
	}()

	start := time.Now()
	for {
		select {
		case <-done:
			return true
		case <-time.Tick(time.Second * 1):
			elapsed := time.Now().Sub(start)
			seconds := int64(math.Round(elapsed.Seconds()))
			atomic.AddInt64(&u.TimeUsed, seconds)
			start = time.Now()

			if !u.IsPremium && u.TimeUsed >= 10 {
				return false
			}
		}
	}
}

func main() {
	RunMockServer()
}
