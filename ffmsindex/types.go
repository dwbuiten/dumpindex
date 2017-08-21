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

// The type of the track
type TrackType uint8

func (this TrackType) String() string {
	switch this {
	case TypeAudio:
		return "Audio"
	case TypeVideo:
		return "Video"
	default:
		return "Unknown"
	}
}

func (this TrackType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", this)), nil
}

// The FFIndex header. Contains info on what versions of various
// libraries the index was generated with, and some other misc ifo
// about how it was generated.
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
	Decoder       uint32
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

// Contains all info about a particular frame or sample.
type Frame struct {
	PTS         int64
	FilePos     int64
	SampleStart int64
	SampleCount uint32
	OriginalPos uint64
	FrameType   int
	RepeatPict  int32
	KeyFrame    bool
	Hidden      bool
}

// Contains all info about a particular track, and all of its frames.
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

// Contains all info from the parsed index.
type Index struct {
	Header *Header
	Tracks []*Track
}
