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
	data := "eF6t2L1LW1EYx/ETo00HF5eCdO9apw4lIIpIQZFSi4gU09zeG4KEGPIiIlKKlFqKg2ToUDp0Fh26" +
		"9R8p/gsO/gEFB5vmRYpk+AznQt7u/d5zzj3P8/ud56T1sLsecrl8IeTC4Einnod0shiyB8WQThTC" +
		"ZudJ//zx71+Pb37+uVyZf7r6bitcHVxc30z0zk/3Xlm9nNSyUtrMKq1/7czN9N46raxUTlq7tU47" +
		"KzXK7Wr/SrjKD/qZHfY36nf442sYf+SePRp+u/3vGN1897mxcO++sdQiUTWiGkRZW3WiKkTtEXVA" +
		"lD3jFlHbRB0S9ZqoKlEWIXvGnYiURegLUctEWY8xlZYQtUqUxbFLlOnRlGbjsow2qk1UzLYsJz4R" +
		"ZTlhXvieqM9E7RP1jSjzVZv7lKiPRJlj2ky8ikiZHt9EpMyjTdtGnRJl69ALoizaRjWJMqWZOl4S" +
		"ZZljPZqbmDPZrJ4RlRAVc+WzqsOUZnNvs2qj/0GUua/5hI1+jSgbl+WXqfaIqA9E2Tr0nSgbl1W1" +
		"RlnFZ3WO1eSmbVsVzL+WiDIHiLm6W9Vh1ZCp1pRmuxMbl9XkVm2b52wSZbtye0br8YQo8y/Le2sr" +
		"JmUuF9PJLUJGvSXKom01QMy9lbmv/Wdl83VOlCnNom37jiOirNo2jzb3tfwaG8e/cv/jCw=="

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

		if v.Hidden {
			t.Errorf("Frame %d should not be hidden.", i)
			return
		}

		if !v.KeyFrame {
			t.Errorf("Frame %d should be a keyframe.", i)
			return
		}
	}
}
