package moriarty

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	file := getTestFile(t)
	assert.True(t, SearchBufferFor(bufio.NewReader(file), []byte("user_username")))
	file.Seek(0, 0)
	assert.True(t, SearchBufferFor(bufio.NewReader(file), []byte("user")))
	file.Seek(0, 0)
	assert.True(t, SearchBufferFor(bufio.NewReader(file), []byte("user_phone_number\": \"123-456-7980")))
	file.Seek(0, 0)
	assert.True(t, SearchBufferFor(bufio.NewReader(file), []byte("id")))
}
func TestSearchAll(t *testing.T) {
	file := getTestFile(t)
	assert.False(t, SearchBufferForAll(file, [][]byte{}), "failed while testing empty")
	testData := [][]byte{
		[]byte("Firstname"),
		[]byte("Sirname"),
		[]byte("testEmail"),
	}
	file.Seek(0, 0) //must be reset between searches.
	assert.True(t, SearchBufferForAll(file, testData), "failed while testing all in order")
	testData = append(testData,
		[]byte("id"),
		[]byte("user_gender"),
	)
	file.Seek(0, 0) //must be reset between searches.
	assert.True(t, SearchBufferForAll(file, testData), "failed while out of order")
}
func TestSearchAny(t *testing.T) {
	file := getTestFile(t)
	assert.False(t, SearchBufferForAny(file, [][]byte{}), "failed while testing empty")

	file.Seek(0, 0) //must be reset between searches.
	assert.True(t,
		SearchBufferForAny(file, [][]byte{
			[]byte("Firstname"),
			[]byte("Sirname"),
			[]byte("testEmail"),
		}),
		"failed while testing all true and in order",
	)

	file.Seek(0, 0) //must be reset between searches.
	assert.True(t,
		SearchBufferForAny(file, [][]byte{
			[]byte("id"),
			[]byte("user_gender"),
		}),
		"failed while all true and out of order",
	)

	file.Seek(0, 0) //must be reset between searches.
	assert.True(t,
		SearchBufferForAny(file, [][]byte{
			[]byte(strings.Repeat("a", 500)), //larger than the file.
			[]byte("user_gender"),
		}),
		"failed while with larger than file test",
	)

	file.Seek(0, 0) //must be reset between searches.
	assert.False(t,
		SearchBufferForAny(file, [][]byte{
			[]byte(strings.Repeat("a", 500)), //larger than the file.
			[]byte("user_genders"),
		}),
		"failed while with larger than file test and no positives",
	)

	file.Seek(0, 0) //must be reset between searches.
	assert.False(t,
		SearchBufferForAny(file, [][]byte{
			[]byte(strings.Repeat("a", 500)), //larger than the file.
			[]byte("user_genders"),
		}),
		"failed while with larger than file test and no positives",
	)

}

func getTestFile(t *testing.T) *os.File {
	filePath := "moriarty/resources/testingData/UserDataDemoFile.json"
	filePath, err := getAbsolutePath(filePath)
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	return file
}
