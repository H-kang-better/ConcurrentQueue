package main

import "context"

// Queue 范型
type Queue[T any] interface {
	DeQueue(ctx context.Context) (any, error)
	EnQueue(ctx context.Context, value T) error

	IsEmpty() bool
	IsFull() bool
	Len() uint64
}
