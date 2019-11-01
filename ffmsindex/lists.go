package ffmsindex

// GetKeyframes gets a list of visible frame numbers which are
// keyframes. This frame number will differ from the frame number
// in the Frames slice, since some codecs, such as VP8, have
// hidden frames. This function takes those into account when
// calculating the frame numbers.
func (t *Track) GetKeyframes() []uint {
	// frame/24 is an arbitrary number to start allocating for
	ret := make([]uint, 0, len(t.Frames)/24)

	for k, v := range t.visibleFrames {
		if t.Frames[v].KeyFrame {
			ret = append(ret, uint(k))
		}
	}

	return ret
}

// GetKeyframeIndicies gets a list of indicies that can be used in
// a given Track's frame slice. This differs from GetKeyframes in
// that it takes into account all frames, and not only visible frames.
// GetKeyframes would be more suitable, to, for example, use for
// seeking with FFMS2's API, or measuring keyframe distances.
func (t *Track) GetKeyframeIndices() []uint {
	// frame/24 is an arbitrary number to start allocating for
	ret := make([]uint, 0, len(t.Frames)/24)

	for _, v := range t.visibleFrames {
		if t.Frames[v].KeyFrame {
			ret = append(ret, uint(v))
		}
	}

	return ret
}

// GetTimestamps returns a list of timestamps for visible frames.
// The number of timestamps may differ from the length of the Frames
// slice, since some codecs, such as VP8, have hidden frames. This
// function takes those into account.
func (t *Track) GetTimestamps() []int64 {
	ret := make([]int64, len(t.visibleFrames))

	for k, n := range t.visibleFrames {
		ret[k] = t.Frames[n].PTS
	}

	return ret
}
