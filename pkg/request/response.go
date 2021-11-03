package request

import (
	"encoding/json"
	cmap "github.com/orcaman/concurrent-map"
	"io/ioutil"
	"net/http"
)

func ParseResponse(response *http.Response) (map[string]interface{}, error) {
	result := cmap.New().Items()
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		err = json.Unmarshal(body, &result)
	}
	return result, err
}
