// +build !solution

package keylock

type KeyLock struct{}

func New() *KeyLock {
	panic("implement me")
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	panic("implement me")
}
