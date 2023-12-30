//////////////////////////////////////////////////////////////////////
//
// DO NOT EDIT THIS PART
// Your task is to edit `main.go`
//

package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

// RunMockServer pretends to be a video processing service. It
// simulates user interacting with the Server.
func RunMockServer() {
	u1 := User{ID: 0, IsPremium: false, TimeUsed: 0, IsTenSecondsElapsed: false}
	u2 := User{ID: 1, IsPremium: true, TimeUsed: 0, IsTenSecondsElapsed: false}

	wg.Add(5)

	go createMockRequest(1, shortProcess, &u1)
	time.Sleep(1 * time.Second)

	go createMockRequest(2, longProcess, &u2)
	time.Sleep(2 * time.Second)

	go createMockRequest(3, shortProcess, &u1)
	time.Sleep(1 * time.Second)

	go createMockRequest(4, longProcess, &u1)
	go createMockRequest(5, shortProcess, &u1)

	wg.Wait()
}

func createMockRequest(pid int, fn func(chan<- bool, *User), u *User) {
	fmt.Println("UserID:", u.ID, "\tProcess", pid, "started.")
	res := HandleRequest(fn, u)

	if res {
		fmt.Println("UserID:", u.ID, "\tProcess", pid, "done.")
	} else {
		fmt.Println("UserID:", u.ID, "\tProcess", pid, "killed. (No quota left)")
	}

	wg.Done()
}

func shortProcess(proc chan<- bool, u *User) {
	time.Sleep(6 * time.Second)

	u.mutex.Lock()
	u.TimeUsed += 6

	if u.TimeUsed > 10 {
		u.IsTenSecondsElapsed = true
	}

	u.mutex.Unlock()
	proc <- true
}

func longProcess(proc chan<- bool, u *User) {
	time.Sleep(11 * time.Second)

	u.mutex.Lock()
	u.TimeUsed += 11

	if u.TimeUsed > 10 {
		u.IsTenSecondsElapsed = true
	}

	u.mutex.Unlock()

	proc <- true
}
