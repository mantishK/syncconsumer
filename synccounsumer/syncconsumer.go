package syncconsumer

import (
	"fmt"
	"sync"
)

type syncconsumer[T any] struct {
	bridge   chan T
	wg       sync.WaitGroup
	consumer []func(a T)
}

func New[T any](fn ...func(a T)) *syncconsumer[T] {
	s := &syncconsumer[T]{
		bridge:   make(chan T),
		wg:       sync.WaitGroup{},
		consumer: fn,
	}
	go func() {
		for {
			data, ok := <-s.bridge
			if !ok {
				return
			}
			for _, consumer := range s.consumer {
				consumer(data)
			}
			s.wg.Done()
		}
	}()
	return s
}

func (s *syncconsumer[T]) Close() {
	s.wg.Wait()
	close(s.bridge)
}

func (s *syncconsumer[T]) Publish(dataSlice ...T) (err error) {
	defer func() {
		if r := recover(); r != nil {
			recoveryErr := r.(error)
			err = fmt.Errorf("can't publish to a closed syncconsumer: %w", recoveryErr)
		}
	}()
	s.wg.Add(len(dataSlice))
	for _, data := range dataSlice {
		s.bridge <- data
	}
	return nil
}
