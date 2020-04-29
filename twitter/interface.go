package twitter

type IMediaUploader interface {
	MediaInit(media []byte) (*MediaInitResponse, error)
	MediaAppend(mediaId uint64, media []byte) error
	MediaFinalize(mediaId uint64) error
	UpdateStatusWithMedia(name string, inReplyToStatusId, mediaId int64) error
}
