package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenNameFrom(t *testing.T) {
	ns := NewNameSet("a", "b")
	names := ns.generateNamesFrom(".")
	assert.ElementsMatch(t,
		[]string{"a", "b", "a.b", "b.a"},
		names,
		"names dont match expected.",
	)
	ns = NewNameSet("a", "b", "c")
	names = ns.generateNamesFrom(".")
	assert.ElementsMatch(t,
		[]string{"a", "b", "c", "a.b", "a.c", "b.a", "b.c", "c.a", "c.b", "a.b.c", "a.c.b", "b.a.c", "b.c.a", "c.a.b", "c.b.a"},
		names,
		"names dont match expected.",
	)
}
func TestGenNames(t *testing.T) {
	ns := NewNameSet("a", "b")
	names := ns.GenerateNames([]string{".", ","})
	assert.ElementsMatch(t,
		[]string{"a", "b", "a.b", "b.a", "a", "b", "a,b", "b,a"},
		names,
	)
}

func BenchmarkByNames(b *testing.B) {
	names := make([]string, 0, b.N)
	for i := 0; i < b.N; i++ {
		names = append(names, fmt.Sprintf("name%v", i))
	}
	ns := NewNameSet(names...)
	b.ResetTimer()
	ns.GenerateNames([]string{"."})
}

func BenchmarkBySeparators(b *testing.B) {
	ns := NewNameSet("a", "b", "c", "d")
	seps := make([]string, 0, b.N)
	for i := 0; i < b.N; i++ {
		seps = append(seps, fmt.Sprintf(":%v:", i))
	}
	b.ResetTimer()
	ns.GenerateNames(seps)
}
