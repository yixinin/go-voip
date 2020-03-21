package cache

import "voip/av"

type Gop struct {
	Video []*av.Packet
	Audio []*av.Packet
	ts    uint64
}

type Cache struct {
	gop   []*Gop
	ttl   int
	index int
	cap   int
}

func NewCache() *Cache {
	return &Cache{
		gop: make([]*Gop, 0, 30),
		// ttl:   10,
		cap:   30,
		index: 0,
	}
}

func (c *Cache) Put(p *av.Packet) {
	// c.gop = append(c.gop, g)
	if p.IsVideo() {
		if p.IsKeyFrame() { //新起一个gop
			var g = &Gop{
				Video: make([]*av.Packet, 0, 5),
				Audio: make([]*av.Packet, 0, 5),
				ts:    p.TimeStamp,
			}
			if c.index >= c.cap { //满了 移除第一个
				c.gop = append(c.gop[1:], g)
				return
			} else {
				c.index++
			}
			c.gop = append(c.gop, g)
			return
		}
		if c.index == 0 {
			return
		}
		c.gop[c.index-1].Video = append(c.gop[c.index-1].Video, p)
		return
	}
	if c.index == 0 {
		return
	}
	c.gop[c.index-1].Audio = append(c.gop[c.index-1].Audio, p)
}

//Get 传入最后一次获取时间戳
func (c *Cache) Get(ts uint64) *Gop {
	for _, v := range c.gop {
		if ts > v.ts {
			return v
		}
	}
	return nil
}
