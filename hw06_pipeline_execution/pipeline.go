package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Этап-прерыватель
	breaker := func(in In) Out {
		out := make(Bi)

		go func() {
			for {
				select {
				case v, exists := <-in:
					if exists { // Получен элемент
						out <- v
					} else { // Канал закрылся
						close(out)
						return
					}

				case <-done: // Сигнал на завершение
					close(out)
					return
				}
			}
		}()

		return out
	}

	// Вырожденный случай
	if len(stages) == 0 {
		return in
	}

	// Собираем в цепочку этапы пайплайна, прокладывая между ними прерыватель
	var out Out
	for _, stage := range stages {
		br := breaker(in)
		out = stage(br)
		in = out
	}
	return out
}
