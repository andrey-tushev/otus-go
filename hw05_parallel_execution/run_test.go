package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("random", func(t *testing.T) {
		// Таблица тестов с разными параметрами
		tests := []struct {
			tasksQty     int     // Количество задач
			successRatio float64 // Процент успешных
			workers      int     // Количество воркеров
			maxErrors    int     // Лимит ошибок
		}{
			{tasksQty: 100, successRatio: 0.66, workers: 10, maxErrors: 20},
			{tasksQty: 50, successRatio: 0.66, workers: 2, maxErrors: 20},
			{tasksQty: 20, successRatio: 0.66, workers: 5, maxErrors: 1},
			{tasksQty: 100, successRatio: 0.66, workers: 50, maxErrors: 100},
			{tasksQty: 10, successRatio: 0.5, workers: 3, maxErrors: 1},
		}

		for _, test := range tests {
			var launchCnt, successCnt, failCnt int32

			// Генерируем набор задач
			tasks := make([]Task, test.tasksQty)
			for i := 0; i < test.tasksQty; i++ {
				if rand.Float64() < test.successRatio {
					// Успешная задача
					tasks[i] = func() error {
						time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
						atomic.AddInt32(&launchCnt, 1)
						atomic.AddInt32(&successCnt, 1)
						return nil
					}
				} else {
					// Задача с ошибкой
					tasks[i] = func() error {
						time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
						atomic.AddInt32(&launchCnt, 1)
						atomic.AddInt32(&failCnt, 1)
						return fmt.Errorf("error from task")
					}
				}
			}

			_ = Run(tasks, test.workers, test.maxErrors)
			// fmt.Printf("Launches: %d, Success: %d, Fail: %d \n", launchCnt, successCnt, failCnt)

			// Количество запусков с ошибкой не может быть больше лимита ошибок
			// плюс кол-во воркеров, т.к. некоторые задачи могут быть уже в работе, когда достигнется лимит ошибок
			require.LessOrEqual(t, failCnt, int32(test.maxErrors+test.workers), "extra tasks were started")
		}
	})
}
