package app

import (
	"get-tweet-image/handlers/tweet"
	"get-tweet-image/image"
	img "get-tweet-image/modules/image"
	"get-tweet-image/modules/twitter"
	"get-tweet-image/twitter/client"
	"get-tweet-image/twitter/media"
	"get-tweet-image/twitter/stream"
	twitterDriver "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type App struct {
	strm stream.IStream
}

func LoadApp(cfg Config) (*App, error) {
	var application App
	config := oauth1.NewConfig(cfg.ConsumerKey, cfg.ConsumerSecret)
	token := oauth1.NewToken(cfg.AccessToken, cfg.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	ttClient := twitterDriver.NewClient(httpClient)
	twitterModule := twitter.NewModule(client.NewClient(ttClient), cfg.Filter)
	imageModule := img.NewModule(image.NewGenerator(), media.NewMediaUploader(httpClient))
	tweetHandler := tweet.NewHandler(imageModule, twitterModule)
	application.strm = stream.NewStream(ttClient)
	application.strm.SetTweetHandler(tweetHandler.Handler)
	filterParams := &twitterDriver.StreamFilterParams{
		Track:         []string{cfg.Filter},
		StallWarnings: twitterDriver.Bool(true),
	}
	if err := application.strm.SetParams(filterParams); err != nil {
		return nil, err
	}
	return &application, nil
}

func (a *App) Run() {
	a.strm.Start()
}

func (a *App) Close() {
	a.strm.Stop()
}
