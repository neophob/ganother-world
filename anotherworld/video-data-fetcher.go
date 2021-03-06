package anotherworld

//VideoDataFetcher is the DTO for polygon drawing
type VideoDataFetcher struct {
	asset      *[]uint8
	readOffset int
}

func (reader *VideoDataFetcher) fetchByte() uint8 {
	result := (*reader.asset)[reader.readOffset]
	reader.readOffset++
	return result
}

func (reader *VideoDataFetcher) fetchWord() uint16 {
	b1 := (*reader.asset)[reader.readOffset]
	b2 := (*reader.asset)[reader.readOffset+1]
	reader.readOffset += 2
	return ToUint16BE(b1, b2)
}

func (reader *VideoDataFetcher) cloneWithUpdatedOffset(readOffset int) VideoDataFetcher {
	return VideoDataFetcher{reader.asset, readOffset}
}

func (video *Video) BuildReader(useVideo2Source bool, readOffset int) VideoDataFetcher {
	videoAssets := video.videoAssets
	if useVideo2Source == true {
		return VideoDataFetcher{&videoAssets.Video2, readOffset}
	}
	return VideoDataFetcher{&videoAssets.Cinematic, readOffset}
}
