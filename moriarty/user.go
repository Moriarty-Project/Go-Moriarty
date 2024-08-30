package moriarty

import "strings"

const WildcardPlaceholder string = "{*}"

var WildcardReplacements []string = []string{"_", "-", " ", ".", ""}

// The recordings of a single user.
type UserRecordings struct {
	AccountName      string
	KnownNamesets    []*NameSet
	KnownSites       []string
	PossibleNamesets []*NameSet
	PossibleSites    []string

	CheckingNSFW bool
}

// generate a user based on known items.
func NewUserRecordings(name string, usernames, emails, sites []string) *UserRecordings {
	// TODO: we should go through and remove empty strings.
	if emails == nil {
		emails = []string{}
	}
	if sites == nil {
		sites = []string{}
	}
	// for usernames, we check if there is a wildcard in there. If there is, we replace it with all equivalents.
	ur := &UserRecordings{
		KnownSites:       sites,
		KnownNamesets:    []*NameSet{},
		PossibleSites:    []string{},
		PossibleNamesets: []*NameSet{},
		CheckingNSFW:     false,
	}
	ur.AddKnownUsername(usernames...)
	return ur
}

// adds usernames, as well as smart generating associated usernames
func (ur *UserRecordings) AddKnownUsername(names ...string) {
	for _, name := range names {
		ur.KnownNamesets = append(ur.KnownNamesets, NewNameSet(name))
	}
}
func (ur *UserRecordings) GetAllNames(separators ...string) []string {
	return append(ur.GetAllKnownNames(separators...), ur.GetAllPossibleNames(separators...)...)
}
func (ur *UserRecordings) GetAllKnownNames(separators ...string) []string {
	ans := []string{}
	for _, nameset := range ur.KnownNamesets {
		ans = append(ans, nameset.GenerateNames(separators)...)
	}
	return ans
}
func (ur *UserRecordings) GetAllPossibleNames(separators ...string) []string {
	ans := []string{}
	for _, nameset := range ur.PossibleNamesets {
		ans = append(ans, nameset.GenerateNames(separators)...)
	}
	return ans
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
	for _, d := range data {
		dataIsNew := true
		for _, t := range dtr.InfoFound[key] {
			if strings.Compare(t, d) == 0 {
				// they're the same!
				dataIsNew = false
				break
			}
		}
		if dataIsNew {
			dtr.InfoFound[key] = append(dtr.InfoFound[key], d)
		}
	}
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
