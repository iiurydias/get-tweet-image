package stream

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
)

type HandlerFunc func(tweet *twitter.Tweet)

type Stream struct {
	client *twitter.Client
	demux  twitter.SwitchDemux
	stream *twitter.Stream
}

func NewStream(client *twitter.Client) IStream {
	return &Stream{
		client: client,
		demux:  twitter.NewSwitchDemux(),
	}
}

func (s *Stream) SetTweetHandler(handler HandlerFunc) {
	s.demux.Tweet = handler
}

func (s *Stream) SetParams(params *twitter.StreamFilterParams) error {
	stream, err := s.client.Streams.Filter(params)
	if err != nil {
		return err
	}
	s.stream = stream
	return nil
}

func (s *Stream) Start() {
	fmt.Println("Starting Stream...")
	go s.demux.HandleChan(s.stream.Messages)
}

func (s *Stream) Stop() {
	fmt.Println("Stopping Stream...")
	s.stream.Stop()
}
