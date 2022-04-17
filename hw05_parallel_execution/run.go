package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, workersNum, maxErrors int) error {
	batch := NewBatch(tasks, workersNum, maxErrors)
	return batch.Run()
}

type Batch struct {
	queue      chan Task
	workersNum int
	maxErrors  int
	errCount   int
	errMutex   sync.Mutex
}

func NewBatch(tasks []Task, workersNum, maxErrors int) Batch {
	//fmt.Printf("Tasks number %d \n", len(tasks))

	b := Batch{
		queue:      make(chan Task, len(tasks)),
		workersNum: workersNum,
		maxErrors:  maxErrors,
		errCount:   0,
	}

	// Заполним очередь задач
	for _, task := range tasks {
		b.queue <- task
	}
	close(b.queue)

	return b
}

func (b *Batch) Run() error {
	// Запустим воркеры
	wg := sync.WaitGroup{}
	for w := 0; w < b.workersNum; w++ {
		wg.Add(1)
		go func(id int) {
			b.worker(id)
			wg.Done()
		}(w)
	}

	// Дождемся завершения всех воркеров
	wg.Wait()

	if b.IsTooManyErr() {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func (b *Batch) worker(id int) {
	//fmt.Printf("Worker %d start \n", workerId)
	for task := range b.queue {
		// Если лимит ошибок превышен, то больше задачи не берем и завершаем воркер
		if b.IsTooManyErr() {
			fmt.Println("TOO MANY")
			return
		}

		//fmt.Printf("Worker %d processing \n", workerId)
		err := task()
		if err != nil {
			b.errMutex.Lock()
			b.errCount++
			b.errMutex.Unlock()
		}

		//fmt.Printf("Worker success %t \n", err != nil)
	}
	//fmt.Printf("Worker %d finish \n", workerId)
}

func (b *Batch) IsTooManyErr() bool {
	b.errMutex.Lock() // Вопрос: Нужны ли тут мьютексы или сравнение с int и так атомарно?
	tooManyErr := b.errCount >= b.maxErrors
	b.errMutex.Unlock()
	return tooManyErr
}
