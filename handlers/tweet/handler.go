package tweet

import (
	"fmt"
	"get-tweet-image/modules/image"
	"get-tweet-image/modules/twitter"
	twitterDriver "github.com/dghubble/go-twitter/twitter"
	"log"
)

type Handler struct {
	imageModule   image.IModule
	twitterModule twitter.IModule
}

func NewHandler(imageModule image.IModule, twitterModule twitter.IModule) *Handler {
	return &Handler{
		imageModule:   imageModule,
		twitterModule: twitterModule,
	}
}

func (h *Handler) Handler(tweet *twitterDriver.Tweet) {
	fmt.Println("a new mention was made")
	if tweet.InReplyToStatusID != 0 {
		repliedTweet, err := h.twitterModule.GetTweet(tweet.InReplyToStatusID)
		if err != nil {
			log.Println("failed to get replied twitter", err)
			return
		}
		if err = h.twitterModule.CheckTweet(repliedTweet); err != nil {
			if err == h.twitterModule.ErrTooLong() {
				if _, err = h.twitterModule.PostTweet("This tweet is too long, @"+tweet.User.ScreenName+" :(", &twitterDriver.StatusUpdateParams{
					InReplyToStatusID: tweet.ID,
				}); err != nil {
					log.Println(err.Error())
					return
				}
				return
			}
			log.Println(err.Error())
			return
		}
		imagePath, err := h.imageModule.GenerateLocalImage(repliedTweet.Text, repliedTweet.User.ScreenName)
		if err != nil {
			log.Println("failed to generate image", err.Error())
			return
		}
		if err = h.imageModule.PublishTweetImage("@"+tweet.User.ScreenName, imagePath, tweet.ID); err != nil {
			log.Println(err.Error())
			return
		}
		if err = h.imageModule.DeleteLocalImage(imagePath); err != nil {
			log.Println("failed to delete image", err.Error())
			return
		}
		return
	}
	fmt.Println("it is a invalid mention")
}
