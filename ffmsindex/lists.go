package ffmsindex

// Gets a list of visible frame numbers which are keyframes.
// This frame number will differ from the frame number in the
// Frames slice, since some codecs, such as VP8, have hidden
// frames. This function takes those into account when
// calculating the frame numbers.
func (this *Track) GetKeyframes() []uint {
    // frame/24 is an arbitrary number to start allocating for
    ret := make([]uint, 0, len(this.Frames) / 24)

    for k, v := range this.visibleFrames {
        if this.Frames[v].KeyFrame {
            ret = append(ret, uint(k))
        }
    }

    return ret
}

// Returns a list of timestamps for visible frames. The number
// of timestamps may differ from the length of the Frames slice,
// since some codecs, such as VP8, have hidden frames. This
// function takes those into account.
func (this *Track) GetTimestamps() []int64 {
    ret := make([]int64, len(this.visibleFrames))

    for _, n := range this.visibleFrames {
        ret = append(ret, this.Frames[n].PTS)
    }

    return ret
}
