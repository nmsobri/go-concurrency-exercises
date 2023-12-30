//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID                  int
	IsPremium           bool
	TimeUsed            int64 // in seconds
	mutex               sync.Mutex
	IsTenSecondsElapsed bool
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(chan<- bool, *User), u *User) bool {
	procChan := make(chan bool)
	done := make(chan bool)

	tenSecondsElapsed := false

	go func() {
		<-time.After(time.Second * 10)
		done <- true
	}()

	go process(procChan, u)

	for {
		select {
		case <-done:
			tenSecondsElapsed = true
		case <-procChan:
			if tenSecondsElapsed {
				if u.IsPremium {
					return true
				} else {
					return false
				}
			}

			if !u.IsPremium {
				return !u.IsTenSecondsElapsed
			}

			return true
		}
	}
}

func main() {
	RunMockServer()
}
