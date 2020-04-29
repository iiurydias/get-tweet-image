package twitter

type IMediaUploader interface {
	MediaInit(media []byte) (*MediaInitResponse, error)
	MediaAppend(mediaId uint64, media []byte) error
	MediaFinalize(mediaId uint64) error
	UpdateStatusWithMedia(inReplyToStatusId, mediaId int64) error
}
