// Package ffmsindex provides a native Go library for parsing FFMS2 Indexes, for
// versions post-2.22.0.1.
package ffmsindex

import (
	"fmt"
)

// Track types
const (
	TypeAudio = TrackType(1)
	TypeVideo = TrackType(0)
)

// TrackType is the type of the track
type TrackType uint8

func (t TrackType) String() string {
	switch t {
	case TypeAudio:
		return "Audio"
	case TypeVideo:
		return "Video"
	default:
		return "Unknown"
	}
}

// MarshalJSON is a standard implementation of JSON marshalling for
// the TrackType type.
func (t TrackType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t)), nil
}

// Header is the FFIndex header. It contains info on what versions of
// various libraries the index was generated with, and some other misc
// info about how it was generated.
type Header struct {
	ID      uint32
	Version struct {
		Major uint8
		Minor uint8
		Micro uint8
		Bump  uint8
	}
	IndexVersion  uint16
	Tracks        uint32
	ErrorHandling uint32
	AVUtilVersion struct {
		Major uint8
		Minor uint8
		Micro uint8
	}
	AVFormatVersion struct {
		Major uint8
		Minor uint8
		Micro uint8
	}
	AVCodecVersion struct {
		Major uint8
		Minor uint8
		Micro uint8
	}
	SWScaleVersion struct {
		Major uint8
		Minor uint8
		Micro uint8
	}
	FileSize int64
	Digest   [20]byte
}

// Frame contains all info about a particular frame or sample.
type Frame struct {
	PTS         int64
	OriginalPTS int64
	FilePos     int64
	SampleStart int64
	SampleCount uint32
	OriginalPos uint64
	FrameType   int
	RepeatPict  int32
	KeyFrame    bool
	Hidden      bool
}

// Track contains all info about a particular track, and all of its frames.
type Track struct {
	TrackType TrackType
	TimeBase  struct {
		Num int64
		Den int64
	}
	MaxBFrames    int32
	UseDTS        bool
	HasTS         bool
	Frames        []Frame
	visibleFrames []int
}

// Index contains all info from the parsed index.
type Index struct {
	Header *Header
	Tracks []*Track
}
