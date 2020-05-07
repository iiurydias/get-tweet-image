package twitter

import (
	"fmt"
	"get-tweet-image/twitter/client"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
	"strings"
)

const TooLong = tooLong("tweet: too long")

type tooLong string

func (e tooLong) Error() string { return string(e) }

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

func (m *Module) PostTweet(text string, params *twitter.StatusUpdateParams) (*twitter.Tweet, error) {
	tweet, res, err := m.twitterClient.PostTweet(text, params)
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, errors.Wrap(err, "failed to post tweet")
	}
	fmt.Println("tweet was posted")
	return tweet, nil
}

func (m *Module) CheckTweet(tweet *twitter.Tweet) error {
	if tweet == nil {
		return errors.New("tweet is empty")
	}
	if tweet.User.Protected {
		return errors.New("tweet owner is protected")
	}
	if strings.Contains(tweet.Text, m.filter) {
		return errors.New("invalid tweet")
	}
	tweet.Text = removeLink(tweet.Text)
	if len(tweet.Text) == 0 {
		return errors.New("not enough tweet text length")
	}
	if len(tweet.Text) > 100 {
		return TooLong
	}
	if "@"+tweet.User.ScreenName == m.filter {
		return errors.New("tweet is from bot")
	}
	fmt.Println("tweet was checked")
	return nil
}

func (m *Module) ErrTooLong() error {
	return TooLong
}

func removeLink(text string) string {
	valid := regexp.MustCompile(`https?://(www\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_+.~#?&//=]*)`)
	if valid.MatchString(text) {
		newText := valid.ReplaceAllString(text, "$1W")
		if len(newText) == 0 {
			return newText
		}
		return newText[:len(newText)-1]
	}
	return text
}
