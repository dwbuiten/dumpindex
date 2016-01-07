// Package ffmsindex provides a native Go library for parsing FFMS2 Major Version 2
// indexes.
package ffmsindex

import (
    "compress/zlib"
    "encoding/binary"
    "fmt"
    "io"
)

// Track types
const (
    ffmsindex     = 0x53920873
    requiredMajor = 2
    TypeAudio     = 1
    TypeVideo     = 0
)

// The FFIndex header. Contains info on what versions of various
// libraries the index was generated with, and some other misc ifo
// about how it was generated.
type Header struct {
    ID uint32
    Version struct {
        Major uint8
        Minor uint8
        Micro uint8
        Bump uint8
    }
    Tracks uint32
    Decoder uint32
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
    Digest [20]byte
}

// Contains all info about a particular frame or sample.
type Frame struct {
    PTS int64
    FilePos int64
    SampleStart int64
    SampleCount uint32
    OriginalPos uint64
    FrameType int
    RepeatPict int32
    KeyFrame bool
    Hidden bool
}

// Contains all info about a particular track, and all of its frames.
type Track struct {
    TrackType uint8
    TimeBase struct {
        Num int64
        Den int64
    }
    MaxBFrames int32
    UseDTS bool
    HasTS bool
    Frames []Frame
}

// Contains all info from the parsed index.
type Index struct {
    Header *Header
    Tracks []*Track
}

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
    ret.Minor = uint8((version & 0x00FF00) >> 8 )
    ret.Micro = uint8((version & 0x0000FF)      )

    return ret, nil
}

func readHeader(r io.Reader) (*Header, error) {
    ret := new(Header)

    err := read(r, &ret.ID)
    if err != nil {
        return nil, err
    } else if ret.ID != ffmsindex {
        return nil, fmt.Errorf("Corrupted FFINDEX.")
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

    if ret.Version.Major != requiredMajor {
        return nil, fmt.Errorf("Wrong FFMS major version.")
    }

    err = read(r, &ret.Tracks)
    if err != nil {
        return nil, err
    }

    err = read(r, &ret.Decoder)
    if err != nil {
        return nil, err
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
        return nil, fmt.Errorf("Digest too short.")
    }

    return ret, nil
}

func readFrames(r io.Reader, frames uint64, typ uint8) ([]Frame, error) {
    ret := make([]Frame, frames)

    oldPTS     := int64(0)
    oldPos     := int64(0)
    oldSamp    := int64(0)
    oldCount   := uint32(0)
    oldOrigPos := uint64(0)

    for i := uint64(0); i < frames; i++ {
        err := read(r, &ret[i].PTS)
        if err != nil {
            return nil, err
        }
        ret[i].PTS += oldPTS
        oldPTS      = ret[i].PTS

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
        oldPos          = ret[i].FilePos

        if typ == TypeVideo {
            err = read(r, &ret[i].OriginalPos)
            if err != nil {
                return nil, err
            }
            ret[i].OriginalPos += oldOrigPos
            oldOrigPos          = ret[i].OriginalPos

            err = read(r, &ret[i].RepeatPict)
            if err != nil {
                return nil, err
            }

            err = read(r, &tmp)
            if err != nil {
                return nil, err
            }
            ret[i].Hidden = (tmp != 0)
        } else if typ == TypeAudio {
            err = read(r, &ret[i].SampleStart)
            if err != nil {
                return nil, err
            }
            ret[i].SampleStart += oldSamp
            oldSamp             = ret[i].SampleStart

            err = read(r, &ret[i].SampleCount)
            if err != nil {
                return nil, err
            }
            ret[i].SampleCount += oldCount
            oldCount            = ret[i].SampleCount
        } else {
            return nil, fmt.Errorf("Unknown Track Type.")
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

// Parses the ffindex from the reader, and returns information
// on all tracks and headers. Unlike FFMS2 itself, it will not
// fail if the libav* library versions do not match, since it
// is intended to extract information only.
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

    return ret, nil
}
