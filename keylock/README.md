# keylock

Напишите примитив синхронизации, позволяющий "лочить" строки из множества.

```go
package keylock

type KeyLock interface {
    // LockKeys locks all keys from provided set.
    // 
    // Upon successful completion, function guarantees that no other call with intersecting set of keys
    // will finish, until unlock() is called.
    //
    // If cancel channel is closed, function returns immediately.
    LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func())
}
```