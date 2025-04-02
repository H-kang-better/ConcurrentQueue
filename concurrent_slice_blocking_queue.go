package main

import (
	"context"
	"golang.org/x/sync/semaphore"
	"sync"
)

// ConcurrentBlockingQueue 基于semaphore.Weighted实现，支持超时控制、阻塞等待。
type ConcurrentBlockingQueue[T any] struct {
	data  []T // 循环队列
	mutex *sync.RWMutex
	head  int // 对头元素下标
	tail  int // 对尾元素下标
	count int // 元素数量

	enqueueCap *semaphore.Weighted
	dequeueCap *semaphore.Weighted

	zero T // 不可作为返回值返回，防止用户修改
}

func NewConcurrentBlockingQueue[T any](capacity int) *ConcurrentBlockingQueue[T] {
	mutex := &sync.RWMutex{}

	semaEnqueue := semaphore.NewWeighted(int64(capacity))
	semaDequeue := semaphore.NewWeighted(int64(capacity))

	_ = semaDequeue.Acquire(context.TODO(), int64(capacity))

	return &ConcurrentBlockingQueue[T]{
		data:       make([]T, capacity),
		mutex:      mutex,
		enqueueCap: semaEnqueue,
		dequeueCap: semaDequeue,
	}
}

// DeQueue 出队
func (c *ConcurrentBlockingQueue[T]) DeQueue(ctx context.Context) (any, error) {
	err := c.dequeueCap.Acquire(ctx, int64(c.count))
	var res T
	if err != nil {
		return res, err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if ctx.Err() != nil {
		c.dequeueCap.Release(1)
		return res, ctx.Err()
	}

	res = c.data[c.head]
	c.data[c.head] = c.zero
	c.count++
	c.head = (c.head + 1) % len(c.data)

	c.enqueueCap.Release(1)
	return res, nil
}

// EnQueue 入队
func (c *ConcurrentBlockingQueue[T]) EnQueue(ctx context.Context, value T) error {
	err := c.enqueueCap.Acquire(ctx, 1)
	if err != nil {
		return err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 先判断是否超时，防止在抢锁的过程中超时
	if ctx.Err() != nil {
		c.dequeueCap.Release(1) // 归还拿到的信号量
		return ctx.Err()
	}

	c.data[c.tail] = value
	c.tail = (c.tail + 1) % len(c.data)
	c.count++

	c.dequeueCap.Release(1)

	return nil
}

func (c *ConcurrentBlockingQueue[T]) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.count
}

func (c *ConcurrentBlockingQueue[T]) AsSlice() []T {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	res := make([]T, 0, c.count)
	count := 0
	for count < c.count {
		index := (count + c.head) % len(c.data)
		res = append(res, c.data[index])
		count++
	}
	return res
}
