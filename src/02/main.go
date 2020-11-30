package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type TokenRecursiveMutex struct {
	sync.Mutex
	token int64
	recursion int32
}

func (m *TokenRecursiveMutex) Lock(token int64) {
	//如果传入的token和持有锁的token一致，
	//说明是递归调用
	if atomic.LoadInt64(&m.token) == token {
		m.recursion++
		return
	}
	m.Mutex.Lock() //传入的token不一致，说明不是递归调用
	//抢到这个锁后记录这个token
	atomic.StoreInt64(&m.token, token)
	m.recursion = 1
}

func (m *TokenRecursiveMutex) Unlock(token int64) {
	if atomic.LoadInt64(&m.token) != token {
		//释放其他token持有的锁
		panic(fmt.Sprintf("wrong owner(%d): %d!", m.token, token))
	}
	m.recursion-- //当前持有这个锁的token释放锁
	if m.recursion != 0 { //还没有回退到最初的递归调用
		return
	}
	atomic.StoreInt64(&m.token, 0) //没有递归调用了，释放锁
	m.Mutex.Unlock()
}