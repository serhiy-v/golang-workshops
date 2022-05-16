package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type DynamicWP struct {
	// number of workers
	mu                 sync.Mutex
	min, max           int
	currentWorkerCount *int32
}

// fill the correct arguments
func (w *DynamicWP) work(ctx context.Context, workerTasks chan func()) {
	//atomic.AddInt32(w.currentWorkerCount, 1)
	//defer wg.Done()
	//start := time.Now()
	for {
		select {
		case task := <-workerTasks:
			//atomic.AddInt32(w.currentWorkerCount, 1)
			fmt.Println(atomic.LoadInt32(w.currentWorkerCount))
			if len(workerTasks) != 0 && int(*w.currentWorkerCount) < 20 {
				//w.mu.Lock()
				//defer w.mu.Unlock()
				atomic.AddInt32(w.currentWorkerCount, 1)
				go w.work(ctx, workerTasks)
			} else if len(workerTasks) == 0 && int(*w.currentWorkerCount) > 3 {
				atomic.AddInt32(w.currentWorkerCount, -1)
				return
			}
			//elapsed := time.Now().Sub(start)
			//seconds := int64(math.Round(elapsed.Seconds()))
			//if seconds > int64(2) {
			//	atomic.AddInt32(w.currentWorkerCount, -1)
			//	return
			//}
			task()
		case <-time.Tick(time.Millisecond * 50):
			if len(workerTasks) != 0 {
				atomic.AddInt32(w.currentWorkerCount, 1)
				go w.work(ctx, workerTasks)
			} else if len(workerTasks) == 0 && int(*w.currentWorkerCount) > 3 {
				atomic.AddInt32(w.currentWorkerCount, -1)
				return
			}
		case <-ctx.Done():
			atomic.AddInt32(w.currentWorkerCount, -1)
			return
		}
	}
	// work should call the task function from task ch
}

// Start starts dynamic worker pull logic
func (w *DynamicWP) Start(ctx context.Context, tasksCh chan func()) {

	for i := 0; i < w.min; i++ {
		atomic.AddInt32(w.currentWorkerCount, 1)
		go w.work(ctx, tasksCh)
	}

}

func NewDynamicWorkerPool(min, max int) *DynamicWP {
	return &DynamicWP{
		min:                min,
		max:                max,
		currentWorkerCount: new(int32),
	}
}

func main() {
	var (
		ctx, _          = context.WithCancel(context.TODO())
		maxWorkersCount = 20
		minWorkersCount = 3
		tasksCh         = make(chan func(), 100)
		taskFunc        = func() { time.Sleep(time.Second * 1) }
		wp              = NewDynamicWorkerPool(minWorkersCount, maxWorkersCount)
	)

	go wp.Start(ctx, tasksCh)
	for i := 0; i < 100; i++ {
		tasksCh <- taskFunc
	}
	fmt.Println("asd")
	time.Sleep(30 * time.Second)

}
