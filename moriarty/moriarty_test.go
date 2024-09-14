package moriarty

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

const jsonFile = "resources/data.json"

func TestMoriartyUsability(t *testing.T) {
	arty, err := NewMoriarty(jsonFile)
	if err != nil {
		t.Fatal(err)
	}
	user := NewUserRecordings(
		"jonah",
	)
	arty.AssignNewUser(user)
	arty.trackingUser.AddNames("jpwilmsmeyer@gmail.com")
	assert.Empty(t, user.KnownFindings)
	assert.Empty(t, user.PossibleFindings)
	arty.TrackUserAcrossSites()

	// now check we found anything
	assert.NotEmpty(t, user.KnownFindings)
	// assert.NotEmpty(t, user.PossibleFindings)

}
