package stream

import "github.com/dghubble/go-twitter/twitter"

type IStream interface {
	SetTweetHandler(handler HandlerFunc)
	SetParams(params *twitter.StreamFilterParams) error
	Start()
	Stop()
}
