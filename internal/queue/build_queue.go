package queue

type BuildQueue struct {
	jobs chan string
}

func NewBuildQueue(buffer int) *BuildQueue {
	return &BuildQueue{
		jobs: make(chan string, buffer),
	}
}

func (q *BuildQueue) Enqueue(id string) {
	q.jobs <- id
}

func (q *BuildQueue) Dequeue() string {
	return <-q.jobs
}
