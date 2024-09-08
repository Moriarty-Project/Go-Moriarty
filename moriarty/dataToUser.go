package moriarty

import (
	"fmt"
	"sync"
)

type DataSearcher interface {
	GetData(string) (*DataTestResults, error) //get the data through this searcher. Nil if nothing is found
	GetName() string                          //get the name of the data searcher
	IsNsfw() bool
}

// this represents the element that actually checks each data clustering.
type DataToUserTester struct {
	dataElement    DataSearcher
	cache          map[string]*DataTestResults //a cache of the responses from a username
	nameFoundCache map[string]bool
	lock           *sync.Mutex
}

func NewSiteUserTester(ubse *StringBasedSiteElement) *DataToUserTester {
	return &DataToUserTester{
		dataElement:    ubse,
		cache:          map[string]*DataTestResults{},
		nameFoundCache: map[string]bool{},
		lock:           &sync.Mutex{},
	}
}
func (sut *DataToUserTester) GetSiteName() string {
	sut.lock.Lock()
	defer sut.lock.Unlock()
	return sut.dataElement.GetName()
}

// test all of these with ourselves.
func (sut *DataToUserTester) TestSiteWith(user *UserRecordings) (fromKnow, fromPossible []*DataTestResults, err error) {
	if user == nil {
		err = fmt.Errorf("no user reference given")
		return
	}
	if sut.dataElement == nil {
		err = fmt.Errorf("no site reference found")
		return
	}
	fromKnow = []*DataTestResults{}
	fromPossible = []*DataTestResults{}
	sut.lock.Lock()
	defer sut.lock.Unlock()
	if sut.dataElement.IsNsfw() && !user.CheckingNSFW {
		// this user doesn't need to be checked on NSFW sites
		return
	}
	for _, name := range user.GetAllKnownNames() {
		if sut.unsafeTestSiteHas(name) {
			fromKnow = append(fromKnow, sut.cache[name])
		}
	}

	for _, name := range user.GetAllPossibleNames() {
		if sut.unsafeTestSiteHas(name) {
			fromPossible = append(fromPossible, sut.cache[name])
		}
	}
	user.AddKnownFindings(fromKnow...)
	user.AddPossibleFindings(fromPossible...)
	return
}
func (sut *DataToUserTester) TestSiteHas(names ...string) []bool {
	sut.lock.Lock()
	defer sut.lock.Unlock()
	ans := make([]bool, 0, len(names))

	for _, name := range names {
		ans = append(ans, sut.unsafeTestSiteHas(name))
	}
	return ans
}

func (sut *DataToUserTester) IsNSFW() bool {
	return sut.dataElement.IsNsfw()
}

// returns true if any of these usernames are found.
func (sut *DataToUserTester) TestSiteHasAny(names ...string) bool {
	sut.lock.Lock()
	defer sut.lock.Unlock()
	for _, name := range names {
		if sut.unsafeTestSiteHas(name) {
			return true
		}
	}
	return false
}

// gets the data test results from
func (sut *DataToUserTester) GetSiteResults(names ...string) *DataTestResults {
	sut.lock.Lock()
	defer sut.lock.Unlock()
	var ans *DataTestResults
	for _, name := range names {
		if sut.unsafeTestSiteHas(name) {
			ans = CombineDataTestResults(ans, sut.cache[name])
		}
	}
	return ans
}

// WARNING: this is not thread safe! use with caution!
func (sut *DataToUserTester) unsafeTestSiteHas(name string) bool {
	if val, hasVal := sut.nameFoundCache[name]; hasVal {
		return val
	}
	found, err := sut.dataElement.GetData(name)
	if err != nil {
		return false
	}
	sut.nameFoundCache[name] = found != nil
	if found != nil {
		sut.cache[name] = found
	}
	return found != nil
}
