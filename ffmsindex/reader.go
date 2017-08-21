package ffmsindex

import (
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	ffmsindex    = 0x53920873
	indexVersion = 4
)

func read(r io.Reader, dst interface{}) error {
	return binary.Read(r, binary.LittleEndian, dst)
}

func readVersion(r io.Reader) (struct {
	Major uint8
	Minor uint8
	Micro uint8
}, error) {

	var ret struct {
		Major uint8
		Minor uint8
		Micro uint8
	}

	var version uint32
	err := read(r, &version)
	if err != nil {
		return ret, err
	}

	ret.Major = uint8((version & 0xFF0000) >> 16)
	ret.Minor = uint8((version & 0x00FF00) >> 8)
	ret.Micro = uint8((version & 0x0000FF))

	return ret, nil
}

func readHeader(r io.Reader) (*Header, error) {
	ret := new(Header)

	err := read(r, &ret.ID)
	if err != nil {
		return nil, err
	} else if ret.ID != ffmsindex {
		return nil, fmt.Errorf("corrupted ffindex")
	}

	err = read(r, &ret.Version.Bump)
	if err != nil {
		return nil, err
	}

	err = read(r, &ret.Version.Micro)
	if err != nil {
		return nil, err
	}

	err = read(r, &ret.Version.Minor)
	if err != nil {
		return nil, err
	}

	err = read(r, &ret.Version.Major)
	if err != nil {
		return nil, err
	}

	// Nothing before this version had the index version field.
	if ret.Version.Major <= 2 && ret.Version.Minor <= 22 &&
		ret.Version.Micro == 0 && ret.Version.Bump < 1 {
		return nil, fmt.Errorf("FFMS2 version used to create index is too old")
	}

	err = read(r, &ret.IndexVersion)
	if err != nil {
		return nil, err
	} else if ret.IndexVersion != indexVersion {
		return nil, fmt.Errorf("unsupported index version (%d)", ret.IndexVersion)
	}

	err = read(r, &ret.Tracks)
	if err != nil {
		return nil, err
	} else if ret.Tracks < 1 {
		return nil, fmt.Errorf("invlaid number of tracks")
	}

	err = read(r, &ret.ErrorHandling)
	if err != nil {
		return nil, err
	}

	ret.AVUtilVersion, err = readVersion(r)
	if err != nil {
		return nil, err
	}

	ret.AVFormatVersion, err = readVersion(r)
	if err != nil {
		return nil, err
	}

	ret.AVCodecVersion, err = readVersion(r)
	if err != nil {
		return nil, err
	}

	ret.SWScaleVersion, err = readVersion(r)
	if err != nil {
		return nil, err
	}

	err = read(r, &ret.FileSize)
	if err != nil {
		return nil, err
	}

	n, err := io.ReadFull(r, ret.Digest[:])
	if err != nil {
		return nil, err
	} else if n != len(ret.Digest) {
		return nil, fmt.Errorf("digest too short")
	}

	return ret, nil
}

func readFrames(r io.Reader, frames uint64, typ TrackType) ([]Frame, error) {
	ret := make([]Frame, frames)

	oldPTS := int64(0)
	oldOriginalPTS := int64(0)
	oldPos := int64(0)
	oldSamp := int64(0)
	oldCount := uint32(0)
	oldOrigPos := uint64(0)

	for i := uint64(0); i < frames; i++ {
		err := read(r, &ret[i].PTS)
		if err != nil {
			return nil, err
		}
		ret[i].PTS += oldPTS
		oldPTS = ret[i].PTS

		err = read(r, &ret[i].OriginalPTS)
		if err != nil {
			return nil, err
		}
		ret[i].OriginalPTS += oldOriginalPTS
		oldOriginalPTS = ret[i].OriginalPTS

		var tmp uint8
		err = read(r, &tmp)
		if err != nil {
			return nil, err
		}
		ret[i].KeyFrame = (tmp != 0)

		err = read(r, &ret[i].FilePos)
		if err != nil {
			return nil, err
		}
		ret[i].FilePos += oldPos
		oldPos = ret[i].FilePos

		err = read(r, &tmp)
		if err != nil {
			return nil, err
		}
		ret[i].Hidden = (tmp != 0)

		if typ == TypeVideo {
			err = read(r, &ret[i].OriginalPos)
			if err != nil {
				return nil, err
			}
			ret[i].OriginalPos += oldOrigPos + 1
			oldOrigPos = ret[i].OriginalPos

			err = read(r, &ret[i].RepeatPict)
			if err != nil {
				return nil, err
			}
		} else if typ == TypeAudio {
			ret[i].SampleStart = oldSamp + int64(oldCount)
			oldSamp = ret[i].SampleStart

			err = read(r, &ret[i].SampleCount)
			if err != nil {
				return nil, err
			}
			ret[i].SampleCount += oldCount
			oldCount = ret[i].SampleCount
		} else {
			return nil, fmt.Errorf("unknown track type")
		}
	}

	return ret, nil
}

func readTrack(r io.Reader) (*Track, error) {
	ret := new(Track)

	err := read(r, &ret.TrackType)
	if err != nil {
		return nil, err
	}

	err = read(r, &ret.TimeBase.Num)
	if err != nil {
		return nil, err
	}

	err = read(r, &ret.TimeBase.Den)
	if err != nil {
		return nil, err
	}

	err = read(r, &ret.MaxBFrames)
	if err != nil {
		return nil, err
	}

	var tmp uint8
	err = read(r, &tmp)
	if err != nil {
		return nil, err
	}
	ret.UseDTS = (tmp != 0)

	err = read(r, &tmp)
	if err != nil {
		return nil, err
	}
	ret.HasTS = (tmp != 0)

	var frames uint64
	err = read(r, &frames)
	if err != nil {
		return nil, err
	}

	ret.Frames, err = readFrames(r, frames, ret.TrackType)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (t *Index) populateVisibleFrames() {
	for i := uint32(0); i < t.Header.Tracks; i++ {
		t.Tracks[i].visibleFrames = make([]int, 0, len(t.Tracks[i].Frames))

		index := 0
		for k, f := range t.Tracks[i].Frames {
			if !f.Hidden {
				t.Tracks[i].visibleFrames = append(t.Tracks[i].visibleFrames, k)
				index++
			}
		}
	}
}

// ReadIndex parses the ffindex from the reader, and returns
// information on all tracks and headers. Unlike FFMS2 itself,
// it will not fail if the libav* library versions do not match,
// since it is intended to extract information only.
func ReadIndex(r io.Reader) (*Index, error) {
	ret := new(Index)

	zr, err := zlib.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	ret.Header, err = readHeader(zr)
	if err != nil {
		return nil, err
	}

	ret.Tracks = make([]*Track, ret.Header.Tracks)

	for i := uint32(0); i < ret.Header.Tracks; i++ {
		ret.Tracks[i], err = readTrack(zr)
		if err != nil {
			return nil, err
		}
	}

	// Some formats, such as VP8, have invisible frames, which shou;d
	// not be counted in the publically visible info.
	ret.populateVisibleFrames()

	return ret, nil
}
