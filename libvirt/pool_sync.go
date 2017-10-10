package libvirt

import (
	"log"
	"sync"
)

// LibVirtPoolSync makes possible to synchronize operations
// against libvirt pools.
// Doing pool.Refresh() operations while uploading or removing
// a volume into the pool causes errors inside of libvirtd
type LibVirtPoolSync struct {
	PoolLocks     map[string]*sync.Mutex
	internalMutex *sync.Mutex
}

// Allocate a new instance of LibVirtPoolSync
func NewLibVirtPoolSync() *LibVirtPoolSync {
	pool := new(LibVirtPoolSync)
	pool.internalMutex = new(sync.Mutex)
	pool.PoolLocks = make(map[string]*sync.Mutex)

	return pool
}

// Acquire a lock for the specified pool
func (ps *LibVirtPoolSync) GetLock(pool string) *sync.Mutex {
	var lock *sync.Mutex
	log.Printf("[DEBUG] going to acquire lock for pool: '%s'", pool)
	ps.internalMutex.Lock()
	defer ps.internalMutex.Unlock()

	lock, exists := ps.PoolLocks[pool]
	if !exists {
		lock = new(sync.Mutex)
		log.Printf("[DEBUG] acquire lock (1); pool: '%s', ps: %+v, lock: %+v\n", pool, ps, lock)
		ps.PoolLocks[pool] = lock
	}

	return lock
}

//// Release the look for the specified pool
//func (ps *LibVirtPoolSync) ReleaseLock(pool string) {
//  log.Printf("[DEBUG] going to release lock for pool: '%s'", pool)
//  ps.internalMutex.Lock()

//  lock, exists := ps.PoolLocks[pool]
//  ps.internalMutex.Unlock()

//  if !exists {
//    return
//  }

//  log.Printf("[DEBUG] release lock; pool: '%s', ps: %+v, lock: %+v\n", pool, ps, lock)
//  lock.Unlock()
//}
