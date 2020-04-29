package twitter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

const StatusUpdate string = "https://api.twitter.com/1.1/statuses/update.json"
const MediaUpload string = "https://upload.twitter.com/1.1/media/upload.json"
const FormContentType string = "application/x-www-form-urlencoded"

type MediaUploader struct {
	client *http.Client
}

type MediaInitResponse struct {
	MediaId          uint64 `json:"media_id"`
	MediaIdString    string `json:"media_id_string"`
	ExpiresAfterSecs uint64 `json:"expires_after_secs"`
}

func NewMediaUploader(client *http.Client) IMediaUploader {
	self := &MediaUploader{}
	self.client = client
	return self
}

func (self *MediaUploader) MediaInit(media []byte) (*MediaInitResponse, error) {
	form := url.Values{
		"command":     []string{"INIT"},
		"media_type":  []string{"image/png"},
		"total_bytes": []string{fmt.Sprint(len(media))},
	}
	content := strings.NewReader(form.Encode())
	res, err := self.makeRequest(http.MethodPost, MediaUpload, FormContentType, content)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusAccepted {
		return nil, errors.New("request was not accepted on media init")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var mediaInitResponse MediaInitResponse
	if err = json.Unmarshal(body, &mediaInitResponse); err != nil {
		return nil, err
	}
	fmt.Println(fmt.Sprintf("Initialized upload of media number %d", mediaInitResponse.MediaId))
	return &mediaInitResponse, nil
}

func (self *MediaUploader) MediaAppend(mediaId uint64, media []byte) error {
	step := 500 * 1024
	for s := 0; s*step < len(media); s++ {
		var body bytes.Buffer
		rangeBeginning := s * step
		rangeEnd := (s + 1) * step
		if rangeEnd > len(media) {
			rangeEnd = len(media)
		}
		contentType, err := createContentForm(&body, s, mediaId, media[rangeBeginning:rangeEnd])
		if err != nil {
			return err
		}
		res, err := self.makeRequest(http.MethodPost, MediaUpload, contentType, &body)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusNoContent {
			return errors.New("request was not accepted on media append")
		}
	}
	fmt.Println(fmt.Sprintf("Media number %d was upploaded", mediaId))
	return nil
}

func (self *MediaUploader) MediaFinalize(mediaId uint64) error {
	form := url.Values{
		"command":  []string{"FINALIZE"},
		"media_id": []string{fmt.Sprint(mediaId)},
	}
	content := strings.NewReader(form.Encode())
	res, err := self.makeRequest(http.MethodPost, MediaUpload, FormContentType, content)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusCreated {
		return errors.New("request was not accepted on media finalized")
	}
	fmt.Println(fmt.Sprintf("Media number %d was finalized", mediaId))
	return nil
}

func (self *MediaUploader) UpdateStatusWithMedia(inReplyToStatusId, mediaId int64) error {
	form := url.Values{
		"status":                []string{""},
		"in_reply_to_status_id": []string{fmt.Sprint(inReplyToStatusId)},
		"media_ids":             []string{fmt.Sprint(mediaId)},
	}
	content := strings.NewReader(form.Encode())
	res, err := self.makeRequest(http.MethodPost, StatusUpdate, FormContentType, content)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("request was not accepted on media append")
	}
	fmt.Println(fmt.Sprintf("Media number %d was published", mediaId))
	return nil
}

func (self *MediaUploader) makeRequest(method, url, contentType string, content io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, content)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", contentType)
	res, err := self.client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func createContentForm(buffer *bytes.Buffer, segmentNumber int, mediaId uint64, content []byte) (string, error) {
	w := multipart.NewWriter(buffer)
	if err := w.WriteField("command", "APPEND"); err != nil {
		return "", err
	}
	if err := w.WriteField("media_id", fmt.Sprint(mediaId)); err != nil {
		return "", err
	}
	if err := w.WriteField("segment_index", fmt.Sprint(segmentNumber)); err != nil {
		return "", err
	}
	fw, err := w.CreateFormFile("media", "orkutMeGenerated.jpg")
	if err != nil {
		return "", err
	}
	if _, err = fw.Write(content); err != nil {
		return "", err
	}
	if err = w.Close(); err != nil {
		return "", err
	}
	return w.FormDataContentType(), nil
}
