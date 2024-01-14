package multiplelock

var defaultMultipleLock = New()

func Lock(key interface{}) {
	defaultMultipleLock.Lock(key)
}

func Unlock(key interface{}) {
	defaultMultipleLock.Unlock(key)
}

func RLock(key interface{}) {
	defaultMultipleLock.RLock(key)
}

func RUnlock(key interface{}) {
	defaultMultipleLock.RUnlock(key)
}
