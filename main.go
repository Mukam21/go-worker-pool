package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Worker представляет собой структуру с ID и функцией отмены context.
// Context используется для управления завершением работы горутины.
type Worker struct {
	ID     int
	Cancel context.CancelFunc
}

// Pool реализует структуру worker-pool.
// Включает мьютекс для синхронизации, список воркеров, канал заданий и счётчик активных горутин.
type Pool struct {
	mu      sync.Mutex
	workers map[int]Worker
	jobs    chan string
	nextID  int
	wg      sync.WaitGroup
}

// NewPool создаёт новый пул с буфером для заданий
func NewPool(bufferSize int) *Pool {
	return &Pool{
		workers: make(map[int]Worker),
		jobs:    make(chan string, bufferSize),
	}
}

// AddWorker запускает нового воркера в виде горутины.
// Каждому воркеру присваивается уникальный ID и создаётся свой context.
func (p *Pool) AddWorker() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	id := p.nextID
	p.nextID++

	worker := Worker{
		ID:     id,
		Cancel: cancel,
	}
	p.workers[id] = worker
	p.wg.Add(1)

	// Запускаем горутину — сам воркер
	go func(id int, ctx context.Context) {
		defer func() {
			// При завершении удаляем воркера из пула и помечаем завершение wg
			p.mu.Lock()
			delete(p.workers, id)
			p.mu.Unlock()
			p.wg.Done()
			fmt.Printf("Worker %d stopped\n", id)
		}()

		fmt.Printf("Worker %d started\n", id)
		for {
			select {
			case <-ctx.Done():
				// Контекст отменён — завершение воркера
				return
			case job, ok := <-p.jobs:
				if !ok {
					// Канал закрыт — завершение воркера
					return
				}
				// Обработка задания
				fmt.Printf("Worker %d processing job: %s\n", id, job)
				time.Sleep(500 * time.Millisecond) // имитация обработки
			}
		}
	}(id, ctx)

	return id
}

// RemoveWorker отключает конкретного воркера по ID.
// Контекст воркера будет отменён, и тот завершит выполнение.
func (p *Pool) RemoveWorker(id int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if worker, exists := p.workers[id]; exists {
		worker.Cancel() // воркер удалится сам при завершении горутины
	}
}

// SendJob помещает задание в очередь.
// Если очередь заполнена, возвращается ошибка.
func (p *Pool) SendJob(job string) error {
	select {
	case p.jobs <- job:
		return nil
	default:
		return fmt.Errorf("job queue is full")
	}
}

// Shutdown завершает работу всех воркеров и очищает ресурсы.
// Важен порядок: сначала закрываем канал, затем отменяем контексты, потом ждём завершения всех горутин.
func (p *Pool) Shutdown() {
	// Сигнализируем воркерам, что больше не будет заданий
	close(p.jobs)

	// Отменяем контексты всех активных воркеров
	p.mu.Lock()
	for _, worker := range p.workers {
		worker.Cancel()
	}
	p.mu.Unlock()

	// Ждём завершения всех воркеров
	p.wg.Wait()
}

func main() {
	pool := NewPool(10) // создаём пул с буфером на 10 заданий

	// Добавляем двух воркеров
	worker1 := pool.AddWorker()
	pool.AddWorker()

	// Отправляем 5 заданий
	for i := 0; i < 5; i++ {
		if err := pool.SendJob(fmt.Sprintf("Task %d", i)); err != nil {
			fmt.Println("Failed to send job:", err)
		}
	}

	time.Sleep(2 * time.Second)

	// Добавляем третьего воркера
	pool.AddWorker()

	// Отправляем ещё 3 задания
	for i := 5; i < 8; i++ {
		if err := pool.SendJob(fmt.Sprintf("Task %d", i)); err != nil {
			fmt.Println("Failed to send job:", err)
		}
	}

	time.Sleep(2 * time.Second)

	// Удаляем одного воркера
	pool.RemoveWorker(worker1)

	time.Sleep(2 * time.Second)

	// Завершаем пул
	fmt.Println("Shutting down pool...")
	pool.Shutdown()
	fmt.Println("Pool shutdown complete.")
}
