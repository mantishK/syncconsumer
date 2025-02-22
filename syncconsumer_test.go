package syncconsumer

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncconsumer(t *testing.T) {
	all := make([]int, 0)
	s := New(
		func(data int) {
			all = append(all, data)
		})

	wg := sync.WaitGroup{}

	// send data
	wg.Add(1)
	go func() {
		s.Publish(1, 2, 3)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		s.Publish(8, 9, 10)
		wg.Done()
	}()
	wg.Wait()
	s.WaitNClose()

	assert.ElementsMatch(t, all, []int{1, 2, 3, 8, 9, 10})
}

func TestSyncconsumerMultiple(t *testing.T) {
	all := make([]int, 0)
	square := make([]int, 0)
	s := New(
		func(data int) {
			all = append(all, data)
		},
		func(data int) {
			square = append(square, data*data)
		},
	)

	wg := sync.WaitGroup{}

	// send data
	wg.Add(1)
	go func() {
		s.Publish(1, 2, 3)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		s.Publish(8, 9, 10)
		wg.Done()
	}()
	wg.Wait()
	s.WaitNClose()

	assert.ElementsMatch(t, all, []int{1, 2, 3, 8, 9, 10})
	assert.ElementsMatch(t, square, []int{1, 4, 9, 64, 81, 100})
}

func TestIncorrectClose(t *testing.T) {
	s := New(func(data int) {})

	s.WaitNClose()

	err := s.Publish(1, 2, 3)
	assert.Error(t, err)
}
