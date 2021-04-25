package physical

const DefaultParallelOperations = 125

type PermitPool struct {
	sem chan int
}

func NewPermitPool(permits int) *PermitPool {
	if permits < 1 {
		permits = DefaultParallelOperations
	}
	return &PermitPool{sem: make(chan int, permits)}
}

func (c *PermitPool) Acquire() {
	c.sem <- 1
}

func (c *PermitPool) Release() {
	<-c.sem
}

func (c *PermitPool) CurrentPermits() int {
	return len(c.sem)
}
