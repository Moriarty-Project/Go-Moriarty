package sherlock

import (
	"sync"
)

type Sherlock struct {
	siteTesters  map[string]*DataToUserTester //the site structures we test against.
	trackingUser *UserRecordings              // the user profile.
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
		go func(sut *DataToUserTester) {
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

func LoadAllSiteUserTesters(filePath string) (map[string]*DataToUserTester, error) {
	// we just create a mapping of each name, to it's siteUserTester.
	ubses, err := LoadAllSiteElements(filePath)
	if err != nil {
		return nil, err
	}
	// go through all of the ubses, and add them!
	ans := map[string]*DataToUserTester{}
	for _, ubse := range ubses {
		ans[ubse.Name] = NewSiteUserTester(ubse)
	}

	return ans, nil
}
