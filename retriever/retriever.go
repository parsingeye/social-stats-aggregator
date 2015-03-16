package retriever

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Facebook struct {
}

func (retriever *Facebook) GetLikes(urlString string) int {

	endpoint := "https://api.facebook.com/method/links.getStats?urls=" + url.QueryEscape(urlString) + "&format=json"

	resp, err := http.Get(endpoint)
	checkErr(err, "Remote call failed")

	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err, "Reading body failed")

	type StatDetails struct {
		Likes int `json:"like_count"`
	}

	data := make([]StatDetails, 0)

	unmarshalErr := json.Unmarshal(body, &data)
	checkErr(unmarshalErr, "Unmarshal failed")

	return data[0].Likes
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
