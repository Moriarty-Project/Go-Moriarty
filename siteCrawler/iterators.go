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
