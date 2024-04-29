package ffmsindex

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestReadOldIndex(t *testing.T) {
	data := "eF7tzbEJwlAQxvHvoqhgYaMBN7C1iSQDZACt7R72YvfGsFKwNlM4gFukSB1IYa0+BTUjCP8PDo77" +
		"3XG7wX4lxZFJChXihgu5SSo3fpa6Oi9PnTDvrZN+czv6wk/ri2uum3lVaZaPXkeHTN+YN7Viefxu" +
		"7j/5PLRyG7XWURRFURRFURRFURRFURRFURRFURRFURRFURRFURRFURT9c30AXiHzTA=="

	r := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer([]byte(data)))

	_, err := ReadIndex(r)
	if err == nil {
		t.Errorf("Shouldn't have been able to read this index.")
		return
	}
}

func TestReadIndex(t *testing.T) {
	data := "eF7tzD0OAVEYheEzMxmZUquwArQamUQzrcYW3MwCWAFRSxQ2Yi/WoLABCSExdwri6CXv09y/937L" +
		"Yj+XlBdK1AjdieqsVJ2XWqSFbuv4Mt6sRp2wrWaX6ylMq2N/cH/dn7PmvRcHtIPi4aDvkmG7e3yI" +
		"H/W+7lL9QO6QO+QOuUPukDvkDrlD7pA75A65Q+6QO+QOuUPukDvkDrlD7pA75A6588/5E6FMdbM="

	r := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer([]byte(data)))

	idx, err := ReadIndex(r)
	if err != nil {
		t.Errorf("Failed to read index: %s.", err)
		return
	}

	if idx.Header.Tracks != 1 || len(idx.Tracks) != 1 {
		t.Errorf("Wrong number of tracks.")
		return
	}

	if len(idx.Tracks[0].Frames) != 150 {
		t.Errorf("Wrong number of frames: %d vs 150.", len(idx.Tracks[0].Frames))
		return
	}

	for i, v := range idx.Tracks[0].Frames {
		if v.PTS != int64(i) {
			t.Errorf("Frame %d was wrong PTS (%d vs %d).", i, i, v.PTS)
			return
		}

		if v.MarkedHidden {
			t.Errorf("Frame %d should not be hidden.", i)
			return
		}

		if !v.KeyFrame {
			t.Errorf("Frame %d should be a keyframe.", i)
			return
		}
	}
}
