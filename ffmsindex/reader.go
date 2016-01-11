package ffmsindex

import (
    "compress/zlib"
    "encoding/binary"
    "fmt"
    "io"
)

const (
    ffmsindex    = 0x53920873
    indexVersion = 1
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

    // Nothing before this version had the index version field.
    if ret.Version.Major <= 2 && ret.Version.Minor <= 22 &&
       ret.Version.Micro == 0 && ret.Version.Bump < 1 {
        return nil, fmt.Errorf("FFMS2 version used to create index is too old.")
    }

    err = read(r, &ret.IndexVersion)
    if err != nil {
        return nil, err
    } else if ret.IndexVersion != indexVersion {
        return nil, fmt.Errorf("Unsupported index version (%d).", ret.IndexVersion)
    }

    err = read(r, &ret.Tracks)
    if err != nil {
        return nil, err
    } else if ret.Tracks < 1 {
        return nil, fmt.Errorf("Invlaid number of tracks.")
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

func readFrames(r io.Reader, frames uint64, typ TrackType) ([]Frame, error) {
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
    } else if frames < 1 {
        return nil, fmt.Errorf("Invalid number of frames in track.")
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
