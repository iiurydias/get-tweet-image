package image

import (
	"fmt"
	"get-tweet-image/image"
	"get-tweet-image/twitter/media"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

type Module struct {
	imageGenerator image.IGenerator
	mediaUploader  media.IMediaUploader
}

func NewModule(imageGenerator image.IGenerator, mediaUploader media.IMediaUploader) IModule {
	return &Module{
		imageGenerator: imageGenerator,
		mediaUploader:  mediaUploader,
	}
}

func (m *Module) GenerateLocalImage(text, username string) (string, error) {
	filePath, err := m.imageGenerator.Generate(text, username)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate image")
	}
	fmt.Println("image was generated")
	return filePath, nil
}

func (m *Module) PublishTweetImage(username, imagePath string, tweetId int64) error {
	data, _ := ioutil.ReadFile(imagePath)
	media, err := m.mediaUploader.MediaInit(data)
	if err != nil {
		return errors.Wrap(err, "failed to init image upload")
	}
	if err = m.mediaUploader.MediaAppend(media.MediaId, data); err != nil {
		return errors.Wrap(err, "failed to upload image")
	}
	if err = m.mediaUploader.MediaFinalize(media.MediaId); err != nil {
		return errors.Wrap(err, "failed to finalize image upload")
	}
	if err = m.mediaUploader.UpdateStatusWithMedia(username, tweetId, int64(media.MediaId)); err != nil {
		return errors.Wrap(err, "failed to finalize image upload")
	}
	return nil
}

func (m *Module) DeleteLocalImage(path string) error {
	if err := os.Remove(path); err != nil {
		return err
	}
	fmt.Println("image was deleted locally")
	return nil
}
