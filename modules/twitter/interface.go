package twitter

import "github.com/dghubble/go-twitter/twitter"

type IModule interface {
	GetTweet(id int64) (*twitter.Tweet, error)
	PostTweet(text string, params *twitter.StatusUpdateParams) (*twitter.Tweet, error)
	CheckTweet(tweet *twitter.Tweet) error
	ErrTooLong() error
}
