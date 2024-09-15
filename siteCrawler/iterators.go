package siteCrawler

import "fmt"

type SiteIterator interface {
	Next() (val string, exists bool) //get the next value and increment the next iteration
	Peak(n int) (string, bool)       //peak ahead, but doesn't actually skip ahead.
	Skip(n int)                      //skip ahead by n points in the iteration.
}
type DefaultIterator struct {
	current  int
	max      int
	stepSize int
}

func NewDefaultIterator(start, stop, step int) *DefaultIterator {
	return &DefaultIterator{
		current:  start,
		max:      stop,
		stepSize: step,
	}
}
func (di *DefaultIterator) Next() (found string, exists bool) {
	if di.current > di.max {
		return "", false
	}
	found = fmt.Sprintf("%v", di.current)
	di.current += di.stepSize
	return found, true
}
func (di *DefaultIterator) Peak(n int) (found string, exists bool) {
	val := di.current + n
	if n > di.max {
		return "", false
	}
	return fmt.Sprintf("%v", val), true
}
func (di *DefaultIterator) Skip(n int) {
	di.current += n
}

type ArrayIterator struct {
	data    []any
	current int
}

func NewArrayIterator(data []any) *ArrayIterator {
	return &ArrayIterator{
		data:    data,
		current: 0,
	}
}
func (ai *ArrayIterator) Next() (found string, exists bool) {
	if ai.current < len(ai.data) {
		found = fmt.Sprintf("%v", ai.data[ai.current])
		ai.current += 1
		return found, true
	}
	return "", false
}
func (ai *ArrayIterator) Peak(n int) (found string, exists bool) {
	val := ai.current + n
	if n >= len(ai.data) {
		return "", false
	}
	return fmt.Sprintf("%v", ai.data[val]), true
}
func (ai *ArrayIterator) Skip(n int) {
	ai.current += n
}

type ChanIterator struct {
	buffer []any
	data   chan any
}

func NewChanIterator(input chan any) *ChanIterator {
	return &ChanIterator{
		buffer: []any{},
		data:   input,
	}
}
func (ci *ChanIterator) Next() (found string, exists bool) {
	if len(ci.buffer) > 0 {
		// we have a buffer going.
		found := fmt.Sprintf("%v", ci.buffer[0])
		ci.buffer = ci.buffer[1:]
		return found, true
	}
	// we have to see if we can get any data from data...
	val, ok := <-ci.data

	return fmt.Sprintf("%v", val), ok
}
func (ci *ChanIterator) Peak(n int) (found string, exists bool) {
	exists = true
	var val any
	for len(ci.buffer) <= n && exists {
		val, exists = <-ci.data
		ci.buffer = append(ci.buffer, val)
	}
	if exists {
		return fmt.Sprintf("%v", ci.buffer[n]), exists
	} else {
		return "", false
	}

}
func (ci *ChanIterator) Skip(n int) {
	for i := 0; i < n; i++ {
		ci.Next()
	}
}
