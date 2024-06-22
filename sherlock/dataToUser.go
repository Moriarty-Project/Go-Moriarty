package sherlock

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// this represents the element that actually checks each data clustering.
type DataToUserTester struct {
	dataElement    *StringBasedSiteElement
	cache          map[string]string //a cache of the responses from a username
	nameFoundCache map[string]bool
	lock           sync.Mutex
}

func NewSiteUserTester(ubse *StringBasedSiteElement) *DataToUserTester {
	return &DataToUserTester{
		dataElement:    ubse,
		cache:          map[string]string{},
		nameFoundCache: map[string]bool{},
		lock:           sync.Mutex{},
	}
}
func (sut *DataToUserTester) GetSiteName() string {
	sut.lock.Lock()
	defer sut.lock.Unlock()
	return sut.dataElement.Name
}

// test all of these with ourselves.
func (sut *DataToUserTester) TestSiteWith(user *UserRecordings, usingCache bool) (fromKnow, fromLikely, fromPossible []*StringBasedSiteElement, err error) {
	if user == nil {
		err = fmt.Errorf("no user reference given")
		return
	}
	if sut.dataElement == nil {
		err = fmt.Errorf("no site reference found")
		return
	}
	fromKnow = []*StringBasedSiteElement{}
	fromLikely = []*StringBasedSiteElement{}
	fromPossible = []*StringBasedSiteElement{}
	sut.lock.Lock()
	defer sut.lock.Unlock()
	if sut.dataElement.IsNSFW && !user.CheckingNSFW {
		// this user doesn't need to be checked on NSFW sites
		return
	}
	for _, name := range append(user.KnownUsernames, user.KnownEmails...) {
		if sut.unsafeTestSiteHas(name, usingCache) {
			fromKnow = append(fromKnow, sut.dataElement)
		}
	}
	for _, name := range append(user.LikelyUsernames, user.LikelyEmails...) {
		if sut.unsafeTestSiteHas(name, usingCache) {
			fromLikely = append(fromLikely, sut.dataElement)
		}
	}
	for _, name := range append(user.PossibleUsernames, user.PossibleEmails...) {
		if sut.unsafeTestSiteHas(name, usingCache) {
			fromPossible = append(fromPossible, sut.dataElement)
		}
	}

	return
}
func (sut *DataToUserTester) TestSiteHas(usingCache bool, names ...string) []bool {
	sut.lock.Lock()
	defer sut.lock.Unlock()
	ans := make([]bool, 0, len(names))

	for _, name := range names {
		ans = append(ans, sut.unsafeTestSiteHas(name, usingCache))
	}
	return ans
}

func (sut *DataToUserTester) IsNSFW() bool {
	return sut.dataElement.IsNSFW
}

// returns true if any of these usernames are found.
func (sut *DataToUserTester) TestSiteHasAny(usingCache bool, names ...string) bool {
	sut.lock.Lock()
	defer sut.lock.Unlock()
	// time to beat with these is 51.7s. Without is 52.4
	for _, name := range names {
		if sut.unsafeTestSiteHas(name, usingCache) {
			return true
		}
	}
	return false
}

// WARNING: this is not thread safe! use with caution!
// TODO: this should be split in half. With caching staying here, and active checking moving to the data handler.
func (sut *DataToUserTester) unsafeTestSiteHas(name string, usingCache bool) bool {
	if usingCache && sut.nameFoundCache[name] {
		return true
	}
	response, err := http.Get(fmt.Sprintf(sut.dataElement.UrlUsername, name))
	if err != nil {
		return false
	}
	// now we need to parse the response to see if it was found.

	// assume that a 404 is always a fail. So we check that it's all 200's
	if response.StatusCode != 200 {
		return false
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return false
	}
	// check if it was found properly. We do this by checking if it doesn't have the false case string in its body.
	found := true
	if sut.dataElement.UnclaimedIfResponseHas != "" {
		found = !strings.Contains(string(responseData), sut.dataElement.UnclaimedIfResponseHas)
	}
	if usingCache {
		sut.nameFoundCache[name] = found
	}
	return found
}
