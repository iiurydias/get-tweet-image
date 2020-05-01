package media

type Response struct {
	MediaId          uint64 `json:"media_id"`
	MediaIdString    string `json:"media_id_string"`
	ExpiresAfterSecs uint64 `json:"expires_after_secs"`
}
