package moriarty

import (
	"GoMoriarty/utils"
	"fmt"
	"sync"
)

type DataSearcher interface {
	GetData(string) (*utils.DataTestResults, error) //get the data through this searcher. Nil if nothing is found
	GetName() string                                //get the name of the data searcher
	IsNsfw() bool
}

// this represents the element that actually checks each data clustering.
// allows for caching of results, among other useful parts.
type DataToUserTester struct {
	dataElement    DataSearcher
	cache          map[string]*utils.DataTestResults //a cache of the responses from a username
	nameFoundCache map[string]bool
	lock           *sync.Mutex
}

func NewDataToUserTester(source DataSearcher) *DataToUserTester {
	return &DataToUserTester{
		dataElement:    source,
		cache:          map[string]*utils.DataTestResults{},
		nameFoundCache: map[string]bool{},
		lock:           &sync.Mutex{},
	}
}
func (sut *DataToUserTester) GetSourceName() string {
	sut.lock.Lock()
	defer sut.lock.Unlock()
	return sut.dataElement.GetName()
}

// test all of these with ourselves.
func (dt *DataToUserTester) TestSourceWith(user *utils.UserRecordings) (fromKnow, fromPossible []*utils.DataTestResults, err error) {
	if user == nil {
		err = fmt.Errorf("no user reference given")
		return
	}
	if dt.dataElement == nil {
		err = fmt.Errorf("no source reference found")
		return
	}
	fromKnow = []*utils.DataTestResults{}
	fromPossible = []*utils.DataTestResults{}
	dt.lock.Lock()
	defer dt.lock.Unlock()
	if dt.dataElement.IsNsfw() && !user.CheckingNSFW {
		// this user doesn't need to be checked on NSFW sources
		return
	}
	for _, name := range user.GetAllKnownNames() {
		if dt.unsafeTestSourceHas(name) {
			fromKnow = append(fromKnow, dt.cache[name])
		}
	}

	for _, name := range user.GetAllPossibleNames() {
		if dt.unsafeTestSourceHas(name) {
			fromPossible = append(fromPossible, dt.cache[name])
		}
	}
	user.AddKnownFindings(fromKnow...)
	user.AddPossibleFindings(fromPossible...)
	return
}
func (dt *DataToUserTester) TestSourceHas(names ...string) []bool {
	dt.lock.Lock()
	defer dt.lock.Unlock()
	ans := make([]bool, 0, len(names))

	for _, name := range names {
		ans = append(ans, dt.unsafeTestSourceHas(name))
	}
	return ans
}

func (dt *DataToUserTester) IsNSFW() bool {
	return dt.dataElement.IsNsfw()
}

// returns true if any of these usernames are found.
func (dt *DataToUserTester) TestSourceHasAny(names ...string) bool {
	dt.lock.Lock()
	defer dt.lock.Unlock()
	for _, name := range names {
		if dt.unsafeTestSourceHas(name) {
			return true
		}
	}
	return false
}

// gets the data test results from
func (sut *DataToUserTester) GetSourceResults(names ...string) *utils.DataTestResults {
	sut.lock.Lock()
	defer sut.lock.Unlock()
	var ans *utils.DataTestResults
	for _, name := range names {
		if sut.unsafeTestSourceHas(name) {
			ans = utils.CombineDataTestResults(ans, sut.cache[name])
		}
	}
	return ans
}

// WARNING: this is not thread safe! use with caution!
func (sut *DataToUserTester) unsafeTestSourceHas(name string) bool {
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
