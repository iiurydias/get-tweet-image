package tweet

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"ttwitterProject/image"
	tt "ttwitterProject/twitter"
)

type Handler struct {
	client         *twitter.Client
	imageGenerator image.IGenerator
	mediaUploader  tt.IMediaUploader
}

func NewHandler(client *twitter.Client, imageGenerator image.IGenerator, mediaUploader tt.IMediaUploader) *Handler {
	return &Handler{client: client, imageGenerator: imageGenerator, mediaUploader: mediaUploader}
}

func (h *Handler) Handler(tweet *twitter.Tweet) {
	if tweet.InReplyToStatusID != 0 {
		fmt.Println("reply tweet was received")
		repliedTweet, res, err := h.client.Statuses.Show(tweet.InReplyToStatusID, nil)
		if err != nil || res.StatusCode != http.StatusOK {
			log.Println("failed to get replied tweet", err)
			return
		}
		if err = h.checkRepliedTweet(repliedTweet); err != nil {
			log.Println(err.Error())
			return
		}
		fmt.Println("tweet replied was got")
		filePath, err := h.imageGenerator.Generate(repliedTweet.Text, repliedTweet.User.ScreenName)
		if err != nil {
			log.Println("failed to generate image", err.Error())
			return
		}
		if err = h.publishImage(filePath, tweet.ID); err != nil {
			log.Println(err.Error())
			return
		}
		if err = deleteFile(filePath); err != nil {
			log.Println("failed to delete image", err.Error())
			return
		}
	}
}

func (h *Handler) checkRepliedTweet(repliedTweet *twitter.Tweet) error {
	if len(repliedTweet.Text) > 100 {
		return errors.New("replied tweet text is too long")
	}
	if len(repliedTweet.Text) == 0 {
		return errors.New("replied tweet text is too short")
	}
	if repliedTweet == nil {
		return errors.New("replied tweet is empty")
	}
	if repliedTweet.User.Protected {
		return errors.New("replied tweet owner is protected")
	}
	return nil
}

func (h *Handler) publishImage(filePath string, inReplyToStatusId int64) error {
	data, _ := ioutil.ReadFile(filePath)
	media, err := h.mediaUploader.MediaInit(data)
	if err != nil {
		return errors.Wrap(err, "failed to init image upload")
	}
	if err = h.mediaUploader.MediaAppend(media.MediaId, data); err != nil {
		return errors.Wrap(err, "failed to upload image")
	}
	if err = h.mediaUploader.MediaFinalize(media.MediaId); err != nil {
		return errors.Wrap(err, "failed to finalize image upload")
	}
	if err = h.mediaUploader.UpdateStatusWithMedia(inReplyToStatusId, int64(media.MediaId)); err != nil {
		return errors.Wrap(err, "failed to finalize image upload")
	}
	return nil
}

func deleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}
