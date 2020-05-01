package twitter

import "github.com/dghubble/go-twitter/twitter"

type IModule interface {
	GetTweet(id int64) (*twitter.Tweet, error)
	CheckTweet(tweet *twitter.Tweet) error
}
