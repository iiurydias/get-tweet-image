package twitter

import (
	"fmt"
	"get-tweet-image/twitter/client"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/pkg/errors"
	"net/http"
)

type Module struct {
	twitterClient client.IClient
	filter        string
}

func NewModule(twitterClient client.IClient, filter string) IModule {
	return &Module{
		twitterClient: twitterClient,
		filter:        filter,
	}
}

func (m *Module) GetTweet(id int64) (*twitter.Tweet, error) {
	tweet, res, err := m.twitterClient.GetTweet(id, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, errors.Wrap(err, "failed to get tweet")
	}
	fmt.Println("tweet was got")
	return tweet, nil
}

func (m *Module) CheckTweet(tweet *twitter.Tweet) error {
	if tweet == nil {
		return errors.New("tweet is empty")
	}
	if tweet.User.Protected {
		return errors.New("tweet owner is protected")
	}
	if len(tweet.Text) > 100 || len(tweet.Text) == 0 {
		return errors.New("not enough tweet text length")
	}
	if "@"+tweet.User.ScreenName == m.filter {
		return errors.New("tweet is from bot")
	}
	fmt.Println("tweet was checked")
	return nil
}
