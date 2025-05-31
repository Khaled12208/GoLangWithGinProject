package service

import (
	"fmt"
	"golangwithgin/internal/domain"
	"sync"
	"time"
)

// TaskProcessor implements the TaskProcessor interface
type TaskProcessor struct {
	tasks      chan *domain.Task
	workers    int
	wg         sync.WaitGroup
	stopChan   chan struct{}
	resultPool *sync.Pool
}

// NewTaskProcessor creates a new task processor with the specified number of workers
func NewTaskProcessor() domain.TaskProcessor {
	processor := &TaskProcessor{
		tasks:    make(chan *domain.Task, 100),
		workers:  5,
		stopChan: make(chan struct{}),
		resultPool: &sync.Pool{
			New: func() interface{} {
				s := ""
				return &s
			},
		},
	}

	processor.start()
	return processor
}

func (p *TaskProcessor) start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

func (p *TaskProcessor) worker() {
	defer p.wg.Done()

	for {
		select {
		case task := <-p.tasks:
			// Get a result string from the pool
			result := p.resultPool.Get().(*string)

			// Process the task
			*result = fmt.Sprintf("Processed task '%s' by worker %d at %v", task.Title, p.workers, time.Now())
			
			// Simulate some work
			time.Sleep(time.Second)

			// Update task status
			task.Status = "completed"
			task.UpdatedAt = time.Now()

			// Put the result string back in the pool
			p.resultPool.Put(result)
		case <-p.stopChan:
			return
		}
	}
}

func (p *TaskProcessor) Process(task *domain.Task) error {
	task.Status = "processing"
	task.UpdatedAt = time.Now()
	p.tasks <- task
	return nil
}

// Shutdown gracefully shuts down the processor
func (p *TaskProcessor) Shutdown() {
	close(p.stopChan)
	p.wg.Wait()
	close(p.tasks)
} 