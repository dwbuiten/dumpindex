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
	}
}

func TestReadIndex(t *testing.T) {
	data := "eF7tzb0JwmAQxvHnYlAhhU0MuIGtjRIHyABa273Yi907hpUBa53CAdzCIrWQwtp8CGpGEP4PHBz3" +
		"u+P2w8PalAQmk9pq4qKF3DiVi+tSqPPq1Gvm/c18UD5zf/GTx9WVt+2sKDTNRu3RcalvzJs6sSx5" +
		"N6+ffB7afRd01lEURVEURVEURVEURVEURVEURVEURVEURVEURVEURVEU/XOtALph804="

	r := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer([]byte(data)))

	idx, err := ReadIndex(r)
	if err != nil {
		t.Errorf("Failed to read index: %s.", err)
	}

	if idx.Header.Tracks != 1 || len(idx.Tracks) != 1 {
		t.Errorf("Wrong number of tracks.")
	}

	if len(idx.Tracks[0].Frames) != 378 {
		t.Errorf("Wrong number of frames.")
	}

	for i, v := range idx.Tracks[0].Frames {
		if v.PTS != int64(i) {
			t.Errorf("Frame %d was wrong PTS (%d vs %d).", i, i, v.PTS)
		}

		if v.Hidden {
			t.Errorf("Frame %d should not be hidden.", i)
		}

		if !v.KeyFrame {
			t.Errorf("Frame %d should be a keyframe.", i)
		}
	}
}
