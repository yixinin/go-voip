package av

const (
	_ byte = iota
	AudioType
	VideoType
)

type Packet struct {
	isAudio bool
	isVideo bool
	Uid     int64
	// TimeStamp uint64
	Data []byte
}

func NewPacket(data []byte, uid int64) *Packet {
	return &Packet{
		isAudio: data[1] == AudioType,
		isVideo: data[1] == VideoType,
		// TimeStamp: utils.BytesToUint64(data[6 : 8+6]),
		Data: data,
		Uid:  uid,
	}
}

func (p *Packet) IsVideo() bool {
	return p.isVideo
}

func (p *Packet) IsAudio() bool {
	return p.isAudio
}
