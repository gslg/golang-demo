package main

import "log"

//rate limit

type ConnLimiter struct {
	concurrentConn int
	bucket         chan int
}

func NewConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter{
		concurrentConn: cc,
		bucket:         make(chan int, cc),
	}
}

func (cl ConnLimiter) GetConn() bool{
	if len(cl.bucket) >= cl.concurrentConn {
		log.Printf("Reached the rate limitation: %d",cl.concurrentConn)
		return false
	}
	cl.bucket <- 1
	return true
}

func (cl ConnLimiter) ReleaseConn(){
	c := <- cl.bucket

	log.Printf("new Connection released: %d",c)
}
