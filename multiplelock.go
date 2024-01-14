package multiplelock

import "sync"

type MultipleLock interface {
	Lock(key interface{})
	Unlock(key interface{})
	RLock(key interface{})
	RUnlock(key interface{})
}

type refCounter struct {
	counter int64
	mutex   *sync.RWMutex
}

type lock struct {
	lock      *sync.Mutex
	inUse     map[interface{}]*refCounter
	mutexPool *sync.Pool
}

func (self *lock) getLocker(key interface{}) *refCounter {
	ret, found := self.inUse[key]
	if found {
		return ret
	}
	ret = &refCounter{
		counter: 0,
		mutex:   self.mutexPool.Get().(*sync.RWMutex),
	}
	self.inUse[key] = ret
	return ret
}

func (self *lock) Lock(key interface{}) {
	self.lock.Lock()
	l := self.getLocker(key)
	l.counter += 1
	self.lock.Unlock()
	l.mutex.Lock()
}

func (self *lock) Unlock(key interface{}) {
	self.lock.Lock()
	l := self.getLocker(key)
	l.counter -= 1
	if l.counter == 0 {
		delete(self.inUse, key)
		self.mutexPool.Put(l.mutex)
	}
	self.lock.Unlock()
	l.mutex.Unlock()
}

func (self *lock) RLock(key interface{}) {
	self.lock.Lock()
	l := self.getLocker(key)
	l.counter += 1
	self.lock.Unlock()
	l.mutex.RLock()
}

func (self *lock) RUnlock(key interface{}) {
	self.lock.Lock()
	l := self.getLocker(key)
	l.counter -= 1
	if l.counter == 0 {
		delete(self.inUse, key)
		self.mutexPool.Put(l.mutex)
	}
	self.lock.Unlock()
	l.mutex.RUnlock()
}

func New() MultipleLock {
	return &lock{
		lock:  &sync.Mutex{},
		inUse: make(map[interface{}]*refCounter),
		mutexPool: &sync.Pool{
			New: func() interface{} {
				return &sync.RWMutex{}
			},
		},
	}
}
