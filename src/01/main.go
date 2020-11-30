package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)
import "github.com/petermattis/goid"

type RecursiveMutex struct {
	sync.Mutex
	owner     int64 //当前持有锁的goroutine id
	recursion int32 //这个goroutine重入的次数
}

func getGoId() int64 {
	return 0
}

func (m *RecursiveMutex) Lock() {
	gid := getGoId()
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	atomic.StoreInt64(&m.owner, gid)
	m.recursion = 1
}

func (m *RecursiveMutex) Unlock() {
	gid := getGoId()
	if atomic.LoadInt64(&m.owner) != gid {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.owner, gid))
	}
	m.recursion--
	if m.recursion != 0 {
		return
	}
	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()
}