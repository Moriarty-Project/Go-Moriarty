package utils

import (
	"strings"
	"sync"
)

// the name set is a way to have multiple parts of a name stored together
type NameSet struct {
	nameParts  []string
	nameCombos []string
	lock       *sync.Mutex
}

func NewNameSet(names ...string) *NameSet {
	if names == nil {
		names = []string{}
	}
	ns := &NameSet{
		nameParts: names,
		lock:      &sync.Mutex{},
	}
	return ns
}
func (ns *NameSet) AddName(names ...string) {
	ns.lock.Lock()
	defer ns.lock.Unlock()
	ns.nameParts = append(ns.nameParts, names...)
	ns.nameCombos = nil
}

/**
 * if separators is nil, a default list of separators is used.
 * generates all possible combinations of names
**/
func (ns *NameSet) GenerateNames(separators []string) []string {
	ns.lock.Lock()
	defer ns.lock.Unlock()
	if ns.nameCombos != nil {
		return ns.nameCombos
	}
	if len(ns.nameParts) <= 1 {
		return ns.nameParts
	}
	wg := &sync.WaitGroup{}
	waitCh := make(chan struct{})
	ansChan := make(chan string, 255)
	for _, sep := range separators {
		wg.Add(1)
		go ns.generateNameChan(sep, ansChan, wg)
		// ans = append(ans, ns.generateNamesFrom(sep)...)
	}
	go func() {
		wg.Wait()
		close(waitCh)
	}()
	ns.nameCombos = []string{}
	for {
		select {
		case <-waitCh:
			// check if ansChan has anything?
			for len(ansChan) != 0 {
				name := <-ansChan
				ns.nameCombos = append(ns.nameCombos, name)
			}
			return ns.nameCombos
		case name := <-ansChan:
			ns.nameCombos = append(ns.nameCombos, name)
		}
	}
}

// generate every possible combination of names based on a separator
// not thread safe, only to be used internally
func (ns *NameSet) generateNamesFrom(separator string) []string {
	ans := []string{}
	comb(len(ns.nameParts)+1, len(ns.nameParts), func(c []int) {
		name := &strings.Builder{}
		for index, i := range c {
			if i == 0 {
				continue
			}
			name.WriteString(ns.nameParts[i-1])
			if index != len(c)-1 {
				name.WriteString(separator)
			}
		}
		ans = append(ans, name.String())
	})
	return ans
}

// generate every possible combination of names based on a separator
func (ns *NameSet) generateNameChan(separator string, ans chan string, wg *sync.WaitGroup) {
	comb(len(ns.nameParts)+1, len(ns.nameParts), func(c []int) {
		name := &strings.Builder{}
		for index, i := range c {
			if i == 0 {
				continue
			}
			name.WriteString(ns.nameParts[i-1])
			if index != len(c)-1 {
				name.WriteString(separator)
			}
		}
		ans <- name.String()
	})
	wg.Done()
}

// zeros may repeat, but all n>0 will be unique
func comb(max, size int, emit func([]int)) {
	// lets imagine it as wheels of big endian numbers
	arr := make([]int, size)
	last := size - 1
	arr[last] += 1
	// a non recursive based combine method
	for arr[0] < max {
		// check everything in arr is unique, or 0
		seen := map[int]bool{}
		duplicatesFound := false
		for i := last; i >= 0; i-- {
			if arr[i] != 0 && seen[arr[i]] {
				duplicatesFound = true
				break
			}
			seen[arr[i]] = true
		}
		if !duplicatesFound {
			emit(arr)
		}
		// increment and round.
		arr[last] += 1
		for i := last; i > 0; i-- {
			if arr[i] == max {
				arr[i] = 1
				arr[i-1] += 1
			}
		}
	}

}
