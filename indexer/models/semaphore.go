package models

// Semaphore creates a semaphore with n slots
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore creates a semaphore with n slots
func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, n),
	}
}

// Acquire acquires a slot in the semaphore (blocks if the semaphore is full)
func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

// Release releases a slot in the semaphore
func (s *Semaphore) Release() {
	<-s.ch
}
