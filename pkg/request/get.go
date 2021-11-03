package request

import (
	"net/http"
)

func Get(url string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// 解析Response
	return ParseResponse(response)

}
