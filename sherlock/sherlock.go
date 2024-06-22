package sherlock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

type Sherlock struct {
	siteTesters  map[string]*SiteUserTester //the site structures we test against.
	trackingUser *UserRecordings            // the user profile.
}

func NewSherlock(filePath string) (*Sherlock, error) {
	testers, err := LoadAllSiteUserTesters(filePath)
	if err != nil {
		return nil, err
	}
	s := &Sherlock{
		siteTesters: testers,
	}
	return s, nil
}
func (s *Sherlock) AssignNewUser(user *UserRecordings) {
	s.trackingUser = user
}

// attempts to
func (s *Sherlock) TrackUser(caching bool) (sitesFoundByKnown, sitesFoundByLikely, SitesFoundByPossible []string, err error) {
	sitesFoundByKnown = []string{}
	sitesFoundByLikely = []string{}
	SitesFoundByPossible = []string{}
	wg := sync.WaitGroup{}
	wg.Add(len(s.siteTesters))
	doneChan := make(chan bool, 1)
	// add a little waiter function here
	go func() {
		wg.Wait()
		doneChan <- true
	}()

	bufferSize := len(s.siteTesters) / 5
	knownChan := make(chan string, bufferSize)
	likelyChan := make(chan string, bufferSize)
	possibleChan := make(chan string, bufferSize)
	errChan := make(chan error)

	user := s.trackingUser
	// now, we go tracking.
	for _, sut := range s.siteTesters {
		go func(sut *SiteUserTester) {
			defer wg.Done()
			knowns := append(user.KnownEmails, user.KnownUsernames...)
			if sut.TestSiteHasAny(caching, knowns...) {
				knownChan <- sut.GetSiteName()
			}
			likelys := append(user.LikelyEmails, user.LikelyUsernames...)
			if sut.TestSiteHasAny(caching, likelys...) {
				likelyChan <- sut.GetSiteName()
			}
			possibles := append(user.PossibleEmails, user.PossibleUsernames...)
			if sut.TestSiteHasAny(caching, possibles...) {
				possibleChan <- sut.GetSiteName()
			}
		}(sut)
	}

	// next we do need to get this data into the answer arrays...
	for {
		select {
		case errResult := <-errChan:
			err = errResult
			return
		case know := <-knownChan:
			sitesFoundByKnown = append(sitesFoundByKnown, know)
		case likely := <-likelyChan:
			sitesFoundByLikely = append(sitesFoundByLikely, likely)
		case possible := <-possibleChan:
			SitesFoundByPossible = append(SitesFoundByPossible, possible)
		default:
			if len(knownChan) == 0 && len(likelyChan) == 0 && len(possibleChan) == 0 && len(doneChan) == 1 {
				return
			}
		}
	}
}

// The recordings of a single user.
type UserRecordings struct {
	KnownUsernames    []string
	KnownEmails       []string
	KnownSites        []string
	LikelyUsernames   []string
	LikelyEmails      []string
	LikelySites       []string
	PossibleUsernames []string
	PossibleEmails    []string
	PossibleSites     []string

	CheckingNSFW bool
}

const WildcardPlaceholder string = "{*}"

var WildcardReplacements []string = []string{"_", "-", " ", ".", ""}

// generate a user based on known items.
func NewUserRecordings(usernames, emails, sites []string) *UserRecordings {
	if emails == nil {
		emails = []string{}
	}
	if sites == nil {
		sites = []string{}
	}
	// for usernames, we check if there is a wildcard in there. If there is, we replace it with all equivalents.
	ur := &UserRecordings{
		KnownUsernames:    []string{},
		KnownEmails:       emails,
		KnownSites:        sites,
		LikelyUsernames:   []string{},
		LikelyEmails:      []string{},
		LikelySites:       []string{},
		PossibleUsernames: []string{},
		PossibleEmails:    []string{},
		PossibleSites:     []string{},
	}

	for _, name := range usernames {
		if strings.Contains(name, WildcardPlaceholder) {
			for _, replacement := range WildcardReplacements {
				likelyName := strings.Replace(name, WildcardPlaceholder, replacement, 1)
				ur.LikelyUsernames = append(ur.LikelyUsernames, likelyName)
				// we add the likely name to usernames to allow for multiple wildcards
				usernames = append(usernames, likelyName)
			}
		} else {
			ur.KnownUsernames = append(ur.KnownUsernames, name)
		}
	}
	return ur
}
func LoadAllSiteElements(filePath string) ([]*StringBasedSiteElement, error) {
	// first, check if the path leads to a large data.json file...
	if !strings.HasSuffix(filePath, ".json") {
		// at some point we might want to deal with this type, eg, multiple files within a folder, but for now, we just throw an error
		return nil, fmt.Errorf("path is not to a JSON file")
	}
	if !path.IsAbs(filePath) {
		// if it's not an absolute path, add the CWD.
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		filePath = path.Join(cwd, filePath)
	}
	// now we try to find the file.
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	fileData, err := io.ReadAll(file)
	file.Close()
	if err != nil {
		return nil, err
	}
	ubses := make([]*StringBasedSiteElement, 0, 100)
	err = json.Unmarshal(fileData, &ubses)
	if err != nil {
		return ubses, err
	}
	// I've split this up so we can check if we want to do any further sanitization here!

	return ubses, nil
}

// this is the base way that sites are saved.
type StringBasedSiteElement struct {
	Name                   string //the name of the site
	UrlHome                string //the home URL for the site
	UrlUsername            string //the url that has a way to format in the username
	UnclaimedIfResponseHas string //if the response text has this, then it is unclaimed.
	IsNSFW                 bool   //is this site NSFW
}

func LoadAllSiteUserTesters(filePath string) (map[string]*SiteUserTester, error) {
	// we just create a mapping of each name, to it's siteUserTester.
	ubses, err := LoadAllSiteElements(filePath)
	if err != nil {
		return nil, err
	}
	// go through all of the ubses, and add them!
	ans := map[string]*SiteUserTester{}
	for _, ubse := range ubses {
		ans[ubse.Name] = NewSiteUserTester(ubse)
	}

	return ans, nil
}

// this represents the element that actually checks each site.
// this may be turned into an abstract system for different ways sites need to be tested.
type SiteUserTester struct {
	SiteElement    *StringBasedSiteElement
	cache          map[string]string //a cache of the responses from a username
	nameFoundCache map[string]bool
	lock           sync.Mutex
}

func NewSiteUserTester(ubse *StringBasedSiteElement) *SiteUserTester {
	return &SiteUserTester{
		SiteElement:    ubse,
		cache:          map[string]string{},
		nameFoundCache: map[string]bool{},
		lock:           sync.Mutex{},
	}
}
func (sut *SiteUserTester) GetSiteName() string {
	sut.lock.Lock()
	defer sut.lock.Unlock()
	return sut.SiteElement.Name
}

// test all of these with ourselves.
func (sut *SiteUserTester) TestSiteWith(user *UserRecordings, usingCache bool) (fromKnow, fromLikely, fromPossible []*StringBasedSiteElement, err error) {
	if user == nil {
		err = fmt.Errorf("no user reference given")
		return
	}
	if sut.SiteElement == nil {
		err = fmt.Errorf("no site reference found")
		return
	}
	fromKnow = []*StringBasedSiteElement{}
	fromLikely = []*StringBasedSiteElement{}
	fromPossible = []*StringBasedSiteElement{}
	sut.lock.Lock()
	defer sut.lock.Unlock()
	if sut.SiteElement.IsNSFW && !user.CheckingNSFW {
		// this user doesn't need to be checked on NSFW sites
		return
	}
	for _, name := range append(user.KnownUsernames, user.KnownEmails...) {
		if sut.unsafeTestSiteHas(name, usingCache) {
			fromKnow = append(fromKnow, sut.SiteElement)
		}
	}
	for _, name := range append(user.LikelyUsernames, user.LikelyEmails...) {
		if sut.unsafeTestSiteHas(name, usingCache) {
			fromLikely = append(fromLikely, sut.SiteElement)
		}
	}
	for _, name := range append(user.PossibleUsernames, user.PossibleEmails...) {
		if sut.unsafeTestSiteHas(name, usingCache) {
			fromPossible = append(fromPossible, sut.SiteElement)
		}
	}

	return
}
func (sut *SiteUserTester) TestSiteHas(usingCache bool, names ...string) []bool {
	sut.lock.Lock()
	defer sut.lock.Unlock()
	ans := make([]bool, 0, len(names))

	for _, name := range names {
		ans = append(ans, sut.unsafeTestSiteHas(name, usingCache))
	}
	return ans
}

// returns true if any of these usernames are found.
func (sut *SiteUserTester) TestSiteHasAny(usingCache bool, names ...string) bool {
	sut.lock.Lock()
	defer sut.lock.Unlock()

	for _, name := range names {
		if sut.unsafeTestSiteHas(name, usingCache) {
			return true
		}
	}
	return false
}

// WARNING: this is not thread safe! use with caution!
func (sut *SiteUserTester) unsafeTestSiteHas(name string, usingCache bool) bool {
	if usingCache && sut.nameFoundCache[name] {
		return true
	}
	response, err := http.Get(fmt.Sprintf(sut.SiteElement.UrlUsername, name))
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
	if sut.SiteElement.UnclaimedIfResponseHas != "" {
		found = !strings.Contains(string(responseData), sut.SiteElement.UnclaimedIfResponseHas)
	}
	if usingCache {
		sut.nameFoundCache[name] = found
	}
	return found
}
