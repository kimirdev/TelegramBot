package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var memeApi = "https://meme-api.herokuapp.com/gimme"

type Meme struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func GetMeme() (Meme, error) {
	var meme Meme

	resp, err := http.Get(memeApi)

	if err != nil {
		return meme, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return meme, err
	}

	json.Unmarshal(body, &meme)

	return meme, nil
}
