package moriarty

import (
	"fmt"
	"sync"
)

type Moriarty struct {
	siteTesters  map[string]*DataToUserTester //the site structures we test against.
	trackingUser *UserRecordings              // the user profile.
}

func NewMoriarty(filePath string) (*Moriarty, error) {
	testers, err := LoadAllSiteUserTesters(filePath)
	if err != nil {
		return nil, err
	}
	s := &Moriarty{
		siteTesters: testers,
	}
	return s, nil
}
func (s *Moriarty) AssignNewUser(user *UserRecordings) {
	s.trackingUser = user
}

// get the results from the user as channels
func (m *Moriarty) GetUserResultsFromSites() (knownChan, likelyChan, possibleChan chan string, doneSignal chan bool) {
	bufferSize := len(m.siteTesters)

	knownChan = make(chan string, bufferSize)
	likelyChan = make(chan string, bufferSize)
	possibleChan = make(chan string, bufferSize)

	wg := &sync.WaitGroup{}
	wg.Add(len(m.siteTesters))
	doneSignal = make(chan bool, 1)

	// get the important user info setup
	user := m.trackingUser
	knownNames := append(user.KnownEmails, user.KnownUsernames...)
	likelyNames := append(user.LikelyEmails, user.LikelyUsernames...)
	possibleNames := append(user.PossibleEmails, user.PossibleUsernames...)

	// now, we go tracking.
	fmt.Println("now starting the go routines")
	for sutName, sut := range m.siteTesters {
		if sut.IsNSFW() && !user.CheckingNSFW {
			wg.Done()
			continue
			// we aren't checking nsfw sites this time
		}
		go func(sut *DataToUserTester, sutName string) {
			if sut.TestSiteHasAny(knownNames...) {
				knownChan <- sutName
			}
			if sut.TestSiteHasAny(likelyNames...) {
				likelyChan <- sutName
			}
			if sut.TestSiteHasAny(possibleNames...) {
				possibleChan <- sutName
			}
			wg.Done()
		}(sut, sutName)
	}
	fmt.Println("all routines have been started!")

	// add a little waiter function here. This is answers that we're done.
	go func() {
		wg.Wait()
		fmt.Println("all go routines have finished!\n ")
		doneSignal <- true
	}()
	fmt.Println("returning function")
	return knownChan, likelyChan, possibleChan, doneSignal
}

// attempts to
func (s *Moriarty) TrackUserAcrossSites() (sitesFoundByKnown, sitesFoundByLikely, sitesFoundByPossible []string) {
	known, likely, possible, done := s.GetUserResultsFromSites()
	<-done
	// we start it right away, then just wait till it's done.
	sitesFoundByKnown = make([]string, len(known))
	sitesFoundByLikely = make([]string, len(likely))
	sitesFoundByPossible = make([]string, len(possible))
	wg := &sync.WaitGroup{}
	wg.Add(3)
	// none of these interact with each other... so we could try paralleling them...
	WriteAll(known, sitesFoundByKnown, wg)
	WriteAll(likely, sitesFoundByLikely, wg)
	WriteAll(possible, sitesFoundByPossible, wg)
	wg.Wait()
	return sitesFoundByKnown, sitesFoundByLikely, sitesFoundByPossible
}

// writes all new data to the array!
func WriteAll(from chan string, to []string, wg *sync.WaitGroup) {
	for i := 0; i < len(to) && len(from) != 0; i++ {
		to[i] = <-from
	}
	wg.Done()
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
