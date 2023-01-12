package av

import "github.com/yixinin/go-voip/bi"

const (
	_ byte = iota
	AudioType
	VideoType
)

type Packet struct {
	isAudio    bool
	isVideo    bool
	isKeyFrame bool
	Uid        int64
	TimeStamp  uint64
	Data       []byte
}

func NewPacket(data []byte, uid int64) *Packet {
	var p = &Packet{
		isAudio:   data[1] == AudioType,
		isVideo:   data[1] == VideoType,
		TimeStamp: bi.BytesToInt[uint64](data[6 : 6+8]),
		// isKeyFrame: data[6+3] == 0x65,
		Data: data,
		Uid:  uid,
	}

	if p.IsVideo() {
		var flags byte
		if data[8] == 0 && data[9] == 1 {
			flags = data[10]
		} else {
			flags = data[9]
		}

		p.isKeyFrame = flags>>4 == 1
	}

	return p
}

func (p *Packet) IsVideo() bool {
	return p.isVideo
}

func (p *Packet) IsAudio() bool {
	return p.isAudio
}

func (p *Packet) IsKeyFrame() bool {
	return p.isKeyFrame
}
