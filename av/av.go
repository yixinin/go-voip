package av

import "voip/utils"

const (
	_ byte = iota
	AudioType
	VideoType
)

type Packet struct {
	isAudio    bool
	isVideo    bool
	isKeyFrame bool
	Uid        string
	TimeStamp  uint64
	Data       []byte
}

func NewPacket(data []byte, uid string) *Packet {
	var p = &Packet{
		isAudio:   data[1] == AudioType,
		isVideo:   data[1] == VideoType,
		TimeStamp: utils.BytesToUint64(data[6 : 8+6]),
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
