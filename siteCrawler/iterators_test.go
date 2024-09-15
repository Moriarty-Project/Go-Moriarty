package siteCrawler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultIterator(t *testing.T) {
	di := NewDefaultIterator(0, 100, 1)
	for i := 0; i <= 100; i++ {
		ansString, exists := di.Next()
		assert.True(t, exists)
		assert.Equal(t, fmt.Sprintf("%v", i), ansString)
	}
	_, exists := di.Next()
	assert.False(t, exists)

	di = NewDefaultIterator(0, 100, 1)
	ans, exists := di.Peak(50)
	assert.True(t, exists)
	assert.Equal(t, "50", ans)
	di.Skip(50)
	ans, exists = di.Next()
	assert.True(t, exists)
	assert.Equal(t, "50", ans)
}
func TestArrayIterator(t *testing.T) {
	arr := make([]any, 0, 101)
	for i := 0; i < cap(arr); i++ {
		arr = append(arr, fmt.Sprintf("%v", i))
	}
	ai := NewArrayIterator(arr)
	for i := 0; i <= 100; i++ {
		ansString, exists := ai.Next()
		assert.True(t, exists)
		assert.Equal(t, fmt.Sprintf("%v", i), ansString)
	}
	_, exists := ai.Next()
	assert.False(t, exists)

	ai = NewArrayIterator(arr)
	ans, exists := ai.Peak(50)
	assert.True(t, exists)
	assert.Equal(t, "50", ans)
	ai.Skip(50)
	ans, exists = ai.Next()
	assert.True(t, exists)
	assert.Equal(t, "50", ans)
}
func TestChanIterator(t *testing.T) {
	input := make(chan any, 101)
	for i := 0; i < cap(input); i++ {
		input <- fmt.Sprintf("%v", i)
	}
	ci := NewChanIterator(input)
	for i := 0; i <= 100; i++ {
		ansString, exists := ci.Next()
		assert.True(t, exists)
		assert.Equal(t, fmt.Sprintf("%v", i), ansString)
	}
	close(input)
	_, exists := ci.Next()
	assert.False(t, exists)
	input = make(chan any, 101)
	for i := 0; i < cap(input); i++ {
		input <- fmt.Sprintf("%v", i)
	}
	ci = NewChanIterator(input)
	ans, exists := ci.Peak(50)
	assert.True(t, exists)
	assert.Equal(t, "50", ans)
	ci.Skip(50)
	ans, exists = ci.Next()
	assert.True(t, exists)
	assert.Equal(t, "50", ans)
}
