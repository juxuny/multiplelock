package main

import (
	"fmt"
	"sync"

	"github.com/juxuny/multiplelock"
)

func main() {
	k1 := "first"
	k2 := "second"
	wg := &sync.WaitGroup{}
	s1 := 0
	s2 := 0
	num := 10000
	wg.Add(num * 2)
	go func() {
		for i := 0; i < num; i++ {
			go func(value int) {
				defer wg.Done()
				multiplelock.Lock(k1)
				defer multiplelock.Unlock(k1)
				s1 += value
			}(i)
		}
	}()

	go func() {
		for i := 0; i < num; i++ {
			go func(value int) {
				defer wg.Done()
				multiplelock.Lock(k2)
				defer multiplelock.Unlock(k2)
				s2 += value
			}(i * 2)
		}
	}()

	wg.Wait()
	multiplelock.RLock(k1)
	fmt.Println("s1: ", s1)
	multiplelock.RUnlock(k1)

	multiplelock.RLock(k2)
	fmt.Println("s2: ", s2)
	multiplelock.RUnlock(k2)
}
