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

func TestCheck(t *testing.T) {
	defer goleak.VerifyNone(t)

	tasksCount := 10
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

	workersCount := 3
	maxErrorsCount := 23
	_ = Run(tasks, workersCount, maxErrorsCount)

}

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
		var count int32

		qty := 100
		successRatio := 0.7

		successCnt := 0
		failCnt := 0
		count = 0

		// Generate tasks
		tasks := make([]Task, qty)
		for i := 0; i < qty; i++ {
			if rand.Float64() < successRatio {
				// Make success task
				successCnt++
				tasks[i] = func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					atomic.AddInt32(&count, 1)
					return nil
				}
			} else {
				// Make failed task
				failCnt++
				tasks[i] = func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					return fmt.Errorf("error from task")
				}
			}
		}

		Run(tasks, 10, 10)
		fmt.Printf("Success tasks: %d, Failed tasks: %d, Count %d \n", successCnt, failCnt, count)

		//require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})
}
