package semaphore

type Semaphore struct {
	semaphoreCh chan struct{}
}

func New(maxTask int) *Semaphore {
	return &Semaphore{
		semaphoreCh: make(chan struct{}, maxTask),
	}
}

func (s *Semaphore) Acquire() {
	s.semaphoreCh <- struct{}{}
}
func (s *Semaphore) Release() {
	<-s.semaphoreCh
}
