package workerpool

type Task func() error

type WorkerPool struct {
	pool chan Task
}

func NewWorkerPool(tasksCount uint) *WorkerPool {
	return &WorkerPool{
		pool: make(chan Task, tasksCount),
	}
}

func (wp *WorkerPool) Run(task Task) {

}
