package retriever

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Provider interface {
	GetCount(url string) int
}

func Make(p string) (Provider, error) {
	switch p {
	case "facebook":
		return new(Facebook), nil
	case "twitter":
		return new(Twitter), nil
	default:
		return nil, errors.New("Provider not supported!")
	}
}

type Facebook struct {
}

func (retriever *Facebook) GetCount(urlString string) int {

	type Response struct {
		Count int `json:"like_count"`
	}

	providerResp := make([]Response, 0)

	endpoint := "https://api.facebook.com/method/links.getStats?urls=" + url.QueryEscape(urlString) + "&format=json"

	resp, err := http.Get(endpoint)
	checkErr(err, "Remote call failed")

	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err, "Reading body failed")

	unmarshalErr := json.Unmarshal(body, &providerResp)
	checkErr(unmarshalErr, "Unmarshal failed")

	return providerResp[0].Count
}

type Twitter struct {
}

func (retriever *Twitter) GetCount(urlString string) int {

	type Response struct {
		Count int `json:"count"`
	}

	var providerResp Response

	endpoint := "http://urls.api.twitter.com/1/urls/count.json?url=" + url.QueryEscape(urlString)

	resp, err := http.Get(endpoint)
	checkErr(err, "Remote call failed")

	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err, "Reading body failed")

	unmarshalErr := json.Unmarshal(body, &providerResp)
	checkErr(unmarshalErr, "Unmarshal failed")

	return providerResp.Count

}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
