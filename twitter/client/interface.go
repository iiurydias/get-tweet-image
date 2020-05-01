package client

import (
	"github.com/dghubble/go-twitter/twitter"
	"net/http"
)

type IClient interface {
	GetTweet(id int64, params *twitter.StatusShowParams) (*twitter.Tweet, *http.Response, error)
}
