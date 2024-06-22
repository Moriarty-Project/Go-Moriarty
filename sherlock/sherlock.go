package sherlock

import (
	"fmt"
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
	bufferSize := len(s.siteTesters) / 5
	sitesFoundByKnown = make([]string, 0, bufferSize)
	sitesFoundByLikely = make([]string, 0, bufferSize)
	SitesFoundByPossible = make([]string, 0, bufferSize)
	wg := sync.WaitGroup{}
	wg.Add(len(s.siteTesters))
	doneChan := make(chan bool, 1)
	// add a little waiter function here

	knownChan := make(chan string, bufferSize)
	likelyChan := make(chan string, bufferSize)
	possibleChan := make(chan string, bufferSize)
	errChan := make(chan error)

	user := s.trackingUser
	knownNames := append(user.KnownEmails, user.KnownUsernames...)
	likelyNames := append(user.LikelyEmails, user.LikelyUsernames...)
	possibleNames := append(user.PossibleEmails, user.PossibleUsernames...)
	// now, we go tracking.
	fmt.Println("now starting the go routines")
	for sutName, sut := range s.siteTesters {
		if sut.IsNSFW() && !user.CheckingNSFW {
			wg.Done()
			continue
			// we arent checking nsfw sites this time
		}
		go func(sut *DataToUserTester, sutName string) {
			if sut.TestSiteHasAny(caching, knownNames...) {
				knownChan <- sutName
			}
			if sut.TestSiteHasAny(caching, likelyNames...) {
				likelyChan <- sutName
			}
			if sut.TestSiteHasAny(caching, possibleNames...) {
				possibleChan <- sutName
			}
			wg.Done()
		}(sut, sutName)
	}
	fmt.Println("all routines have been started!")
	go func() {
		wg.Wait()
		fmt.Println("all go routines have finished!")
		doneChan <- true
	}()
	fmt.Println("waitgroup is waiting now!")
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
		}
		if len(knownChan) == 0 && len(likelyChan) == 0 && len(possibleChan) == 0 && len(doneChan) == 1 {
			fmt.Println("done!\n ")
			return
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
