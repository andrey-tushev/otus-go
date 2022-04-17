package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, workersNum, maxErrors int) error {
	batch := NewBatch(tasks, workersNum, maxErrors)
	return batch.Run()
}

// Запускалка пакета задач.
type Batch struct {
	queue      chan Task
	workersNum int
	maxErrors  int
	errCount   int
	errMutex   *sync.Mutex
}

// Конструктор.
func NewBatch(tasks []Task, workersNum, maxErrors int) Batch {
	b := Batch{
		queue:      make(chan Task, len(tasks)),
		workersNum: workersNum,
		maxErrors:  maxErrors,
		errCount:   0,
		errMutex:   &sync.Mutex{},
	}

	// Заполним очередь задач
	for _, task := range tasks {
		b.queue <- task
	}
	close(b.queue)

	return b
}

// Запускалка пакета задач.
func (b *Batch) Run() error {
	// Запустим воркеры
	wg := sync.WaitGroup{}
	for w := 0; w < b.workersNum; w++ {
		wg.Add(1)
		go func() {
			b.worker()
			wg.Done()
		}()
	}

	// Дождемся завершения всех воркеров
	wg.Wait()

	// Возможно превышен лимит
	if b.IsTooManyErr() {
		return ErrErrorsLimitExceeded
	}

	return nil
}

// Воркер.
func (b *Batch) worker() {
	// Вытаскиваем задачи из очереди, пока они незакончатся или не достигнется лимит ошибок
	for task := range b.queue {
		// Если лимит ошибок превышен, то больше задачи не берем и завершаем воркер
		if b.IsTooManyErr() {
			return
		}

		// Выполняем извлеченную задачу
		err := task()
		// Если были ошибки, то увеличиваем счетчик ошибок
		if err != nil {
			b.errMutex.Lock()
			b.errCount++
			b.errMutex.Unlock()
		}
	}
}

func (b *Batch) IsTooManyErr() bool {
	b.errMutex.Lock() // Вопрос: Нужны ли тут мьютексы или сравнение с int и так атомарно?
	tooManyErr := b.errCount >= b.maxErrors
	b.errMutex.Unlock()
	return tooManyErr
}
