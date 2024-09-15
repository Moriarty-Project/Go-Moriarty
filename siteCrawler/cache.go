package siteCrawler

/**
cache acts as the cached reference of all saved tests from the crawler. All positive and negative results are stored in here.
*/
import (
	"GoMoriarty/utils"
	"encoding/json"
	"io"
	"os"
	"path"
)

type siteCrawlerCache struct {
	savePath string //where this file is saved.
	Found    map[string]bool
}

func LoadOrNewSiteCrawlerCache(savePath string) *siteCrawlerCache {
	absPath, _ := utils.GetAbsolutePath(path.Dir(savePath))
	scc := &siteCrawlerCache{
		savePath: absPath,
		Found:    map[string]bool{},
	}
	file, err := os.Open(path.Join(scc.savePath, "cache.json"))
	if err == nil {
		// file exists
		defer file.Close()
		data, err := io.ReadAll(file)
		if err == nil {
			json.Unmarshal(data, scc)
		}

	}
	return scc
}
func (scc *siteCrawlerCache) Save() error {
	file, err := os.Create(path.Join(scc.savePath, "cache.json"))
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(scc)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}
func (scc *siteCrawlerCache) IsNew(val string) bool {
	_, exists := scc.Found[val]
	return !exists
}
func (scc *siteCrawlerCache) IsFound(val string) bool {
	return scc.Found[val]
}
func (scc *siteCrawlerCache) AddFound(val string, correctly bool) {
	scc.Found[val] = correctly
}
