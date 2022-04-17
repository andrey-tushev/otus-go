package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type errCount struct {
	sync.Mutex
	Count int
}

func (e *errCount) Inc() {
	e.Mutex.Lock()
	e.Count++
	e.Mutex.Unlock()
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workersNum, maxErrors int) error {
	//fmt.Printf("Tasks number %d \n", len(tasks))

	// Создадим очередь задач
	queue := make(chan Task, len(tasks))
	for _, task := range tasks {
		queue <- task
	}
	close(queue)

	var ec errCount

	// Запустим воркеры
	wg := sync.WaitGroup{}
	for w := 0; w < workersNum; w++ {
		wg.Add(1)
		go func(id int) {
			worker(id, queue, &ec)
			wg.Done()
		}(w)
	}

	// Дождемся завершения всех воркеров
	wg.Wait()

	return nil
}

// Вопрос: Почему каналы можно передавать по значению и все работает?

func worker(workerId int, queue chan Task, ec *errCount) {
	//fmt.Printf("Worker %d start \n", workerId)
	for task := range queue {
		//fmt.Printf("Worker %d processing \n", workerId)
		err := task()
		if err != nil {
			ec.Inc()
		}

		fmt.Printf("Worker success %t \n", success)
	}
	//fmt.Printf("Worker %d finish \n", workerId)
}
