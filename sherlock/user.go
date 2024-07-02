package sherlock

import "strings"

const WildcardPlaceholder string = "{*}"

var WildcardReplacements []string = []string{"_", "-", " ", ".", ""}

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
		CheckingNSFW:      false,
	}
	ur.AddKnownUsername(usernames...)
	return ur
}

// adds usernames, as well as smart generating associated usernames
func (ur *UserRecordings) AddKnownUsername(names ...string) {
	for _, name := range names {
		if strings.Contains(name, WildcardPlaceholder) {
			for _, replacement := range WildcardReplacements {
				likelyName := strings.Replace(name, WildcardPlaceholder, replacement, 1)
				ur.LikelyUsernames = append(ur.LikelyUsernames, likelyName)
				// we add the likely name to usernames to allow for multiple wildcards
				names = append(names, likelyName)
			}
		} else {
			ur.KnownUsernames = append(ur.KnownUsernames, name)
		}
	}
}

// this is going to be the return type for all of our data testers
type DataTestResults struct {
	Name      string //name of where it was found
	InfoFound map[string][]string
	/*
		should be map[InfoName][InfoItem1,InfoItem2...]. Basically, the way to store all info found from a data source.
		EG.
		Name: "TestSite123"
		InfoFound:{
			"OtherSocialMedias": ["Linkedin","Facebook"...],
			"Password":["123"],
			"Username": ["HelloWorld"],
			...
		}
	*/
}

func NewDataTestResults(source string) *DataTestResults {
	return &DataTestResults{
		Name:      source,
		InfoFound: map[string][]string{},
	}
}

// add the new, unique data! All keys are stored as lowercase only. To prevent confusion between site formats.
func (dtr *DataTestResults) Add(key string, data ...string) {
	key = strings.ToLower(key)
	if dtr.InfoFound[key] == nil {
		// it's nil, so we need a new array!
		dtr.InfoFound[key] = make([]string, 0, len(data))
	}
	// we check for duplicate data, and disregard that if it's an exact duplicate.
	// first we'll go over and check that all our data is new
	newData := make([]string, 0, len(data))
	for _, d := range data {
		dataIsNew := true
		for _, t := range dtr.InfoFound[key] {
			if strings.Compare(t, d) == 0 {
				// they're the same!
				dataIsNew = false
				break
			}
			if dataIsNew {
				newData = append(newData, d)
			}
		}
	}
	dtr.InfoFound[key] = append(dtr.InfoFound[key], newData...)
}

// returns self if values have been found, else, empty
func (dtr *DataTestResults) NilIfEmpty() *DataTestResults {
	if len(dtr.InfoFound) == 0 {
		return nil
	}
	return dtr
}

// tries its best to combine two data test results. Returns nill if anything goes wrong
func CombineDataTestResults(a *DataTestResults, b *DataTestResults) *DataTestResults {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}

	if a.Name != b.Name {
		return nil
	}

	ans := NewDataTestResults(a.Name)
	for key, vals := range a.InfoFound {
		ans.Add(key, vals...)
	}
	for key, vals := range b.InfoFound {
		ans.Add(key, vals...)
	}
	return ans
}

// attempts to condense down all of the data down to a smaller array
func CombineAllDataTestResults(in []*DataTestResults) []*DataTestResults {

	//we'll make a sorted form of the inputs, combining all the names together first.
	sortedByIns := map[string]*DataTestResults{}
	for _, test := range in {
		sortedByIns[test.Name] = CombineDataTestResults(test, sortedByIns[test.Name])
	}
	// now we put it back together!
	out := make([]*DataTestResults, 0, len(sortedByIns))
	for _, v := range sortedByIns {
		out = append(out, v)
	}

	return out
}
