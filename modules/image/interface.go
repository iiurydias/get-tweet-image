package image

type IModule interface {
	GenerateLocalImage(text, username string) (string, error)
	PublishTweetImage(username, imagePath string, tweetId int64) error
	DeleteLocalImage(path string) error
}
