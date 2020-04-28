package main

import (
	tt "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	cfg "github.com/micro/go-micro/config"
	"log"
	"os"
	"os/signal"
	"syscall"
	"ttwitterProject/handlers/tweet"
	"ttwitterProject/image"
	"ttwitterProject/twitter"
)

func main() {
	configParams := getConfigParams()
	config := oauth1.NewConfig(configParams.ConsumerKey, configParams.ConsumerSecret)
	token := oauth1.NewToken(configParams.AccessToken, configParams.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	ttClient := tt.NewClient(httpClient)
	stream := twitter.NewStream(ttClient)
	tweetHandler := tweet.NewHandler(ttClient, image.NewGenerator(), twitter.NewMediaUploader(httpClient))
	stream.SetTweetHandler(tweetHandler.Handler)
	filterParams := &tt.StreamFilterParams{
		Track:         []string{configParams.Filter},
		StallWarnings: tt.Bool(true),
	}
	if err := stream.SetParams(filterParams); err != nil {
		log.Fatal(err)
	}
	stream.Start()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	stream.Stop()
}

func getConfigParams() Config {
	var configParams Config
	err := cfg.LoadFile("./config.json")
	if err != nil {
		log.Fatal("fail to load config file", err.Error())
	}
	err = cfg.Scan(&configParams)
	if err != nil {
		log.Fatal("failed to scan file", err.Error())
	}
	return configParams
}
