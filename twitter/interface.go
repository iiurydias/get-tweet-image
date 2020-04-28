package twitter

type IMediaUploader interface {
	MediaInit(media []byte) (*MediaInitResponse, error)
	MediaAppend(mediaId uint64, media []byte) error
	MediaFinilize(mediaId uint64) error
	UpdateStatusWithMedia(text string, mediaId uint64) error
}
