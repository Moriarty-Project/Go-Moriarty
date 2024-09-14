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
func (m *Moriarty) GetUserResultsFromSites() (doneSignal chan (bool)) {
	wg := &sync.WaitGroup{}
	wg.Add(len(m.siteTesters))
	doneSignal = make(chan bool, 1)

	// get the important user info setup
	user := m.trackingUser
	// now, we go tracking.
	fmt.Println("now starting the go routines")
	for _, sut := range m.siteTesters {
		if sut.IsNSFW() && !user.CheckingNSFW {
			wg.Done()
			continue
			// we aren't checking nsfw sites this time
		}
		go func(sut *DataToUserTester, user *UserRecordings) {
			sut.TestSourceWith(user) //automatically adds the logs to the user.
			wg.Done()
		}(sut, m.trackingUser)
	}
	fmt.Println("all routines have been started!")

	// add a little waiter function here. This is answers that we're done.
	go func() {
		wg.Wait()
		fmt.Println("all go routines have finished!\n ")
		doneSignal <- true
	}()
	fmt.Println("returning function")
	return doneSignal
}

func (m *Moriarty) GetKnownFromSites() (knownChan chan string, doneSignal chan bool) {
	return m.GetAllSitesFrom(m.trackingUser.GetAllKnownNames("")...)
}

func (m *Moriarty) GetAllSitesFrom(names ...string) (sitesWithNames chan string, doneSignal chan bool) {
	sitesWithNames = make(chan string, len(m.siteTesters))
	wg := &sync.WaitGroup{}
	wg.Add(len(m.siteTesters))
	doneSignal = make(chan bool, 1)
	for sutName, sut := range m.siteTesters {
		if sut.IsNSFW() && !m.trackingUser.CheckingNSFW {
			wg.Done()
			continue
			// we aren't checking nsfw sites this time
		}
		go func(sut *DataToUserTester, sutName string) {
			if sut.TestSourceHasAny(names...) {
				sitesWithNames <- sutName
			}
			wg.Done()
		}(sut, sutName)
	}
	go func() {
		wg.Wait()
		doneSignal <- true
	}()
	return
}

// attempts to
func (m *Moriarty) TrackUserAcrossSites() {
	done := m.GetUserResultsFromSites()
	<-done
	// we start it right away, then just wait till it's done.
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
		ans[ubse.Name] = NewDataToUserTester(ubse)
	}

	return ans, nil
}
