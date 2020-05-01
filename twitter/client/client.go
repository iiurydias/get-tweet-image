package client

import (
	"github.com/dghubble/go-twitter/twitter"
	"net/http"
)

type Client struct {
	client *twitter.Client
}

func NewClient(client *twitter.Client) IClient {
	return &Client{
		client: client,
	}
}

func (c *Client) GetTweet(id int64, params *twitter.StatusShowParams) (*twitter.Tweet, *http.Response, error) {
	return c.client.Statuses.Show(id, params)
}
