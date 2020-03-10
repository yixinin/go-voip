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
	// TimeStamp uint32
	Data []byte
}

func NewPacket(dataType byte, data []byte, uid int64) *Packet {
	return &Packet{
		isAudio: dataType == AudioType,
		isVideo: dataType == VideoType,
		Data:    data,
		Uid:     uid,
	}
}

func (p *Packet) IsVideo() bool {
	return p.isVideo
}

func (p *Packet) IsAudio() bool {
	return p.isAudio
}
