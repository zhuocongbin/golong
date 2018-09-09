package httpflv

import (
	"errors"
	"fmt"
	"github.com/livego/av"
	log "github.com/livego/logging"
	"github.com/livego/protocol/amf"
	"github.com/livego/utils/pio"
	"github.com/livego/utils/uid"
	"net/http"
	"time"
)

const (
	headerLen   = 11
	maxQueueNum = 1024
)

type FLVWriter struct {
	Uid string
	av.RWBaser
	app, title, url string
	buf             []byte
	closed          bool
	closedChan      chan struct{}
	ctx             http.ResponseWriter
	req             *http.Request
	packetQueue     chan *av.Packet
	WriteBWInfo     av.StaticsBW
}

func NewFLVWriter(app, title, url string, req *http.Request, ctx http.ResponseWriter) *FLVWriter {
	ret := &FLVWriter{
		Uid:         uid.NewId(),
		app:         app,
		title:       title,
		url:         url,
		req:         req,
		ctx:         ctx,
		RWBaser:     av.NewRWBaser(time.Second * 10),
		closedChan:  make(chan struct{}),
		buf:         make([]byte, headerLen),
		packetQueue: make(chan *av.Packet, maxQueueNum),
		WriteBWInfo: av.StaticsBW{0, "", 0, 0, 0, 0, 0, 0, 0},
	}

	ret.ctx.Write([]byte{0x46, 0x4c, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x09})
	pio.PutI32BE(ret.buf[:4], 0)
	ret.ctx.Write(ret.buf[:4])
	go func() {
		err := ret.SendPacket()
		if err != nil {
			log.Error("SendPacket error:", err)
			ret.closed = true
		}
	}()
	return ret
}

func (flvWriter *FLVWriter) DropPacket(pktQue chan *av.Packet, info av.Info) {
	log.Info("[%v] packet queue max!!!", info)
	for i := 0; i < maxQueueNum-84; i++ {
		tmpPkt, ok := <-pktQue
		if ok && tmpPkt.IsVideo {
			videoPkt, ok := tmpPkt.Header.(av.VideoPacketHeader)
			// dont't drop sps config and dont't drop key frame
			if ok && (videoPkt.IsSeq() || videoPkt.IsKeyFrame()) {
				log.Info("insert keyframe to queue")
				pktQue <- tmpPkt
			}

			if len(pktQue) > maxQueueNum-10 {
				<-pktQue
			}
			// drop other packet
			<-pktQue
		}
		// try to don't drop audio
		if ok && tmpPkt.IsAudio {
			log.Info("insert audio to queue")
			pktQue <- tmpPkt
		}
	}
	log.Info("packet queue len: ", len(pktQue))
}

func (flvWriter *FLVWriter) Write(p *av.Packet) (err error) {
	err = nil
	if flvWriter.closed {
		err = errors.New("flvwrite source closed")
		return
	}
	defer func() {
		if e := recover(); e != nil {
			errString := fmt.Sprintf("FLVWriter has already been closed:%v", e)
			err = errors.New(errString)
		}
	}()
	if len(flvWriter.packetQueue) >= maxQueueNum-24 {
		flvWriter.DropPacket(flvWriter.packetQueue, flvWriter.Info())
	} else {
		flvWriter.packetQueue <- p
	}

	return
}

func (flvWriter *FLVWriter) SaveStatics(streamid uint32, length uint64, isVideoFlag bool) {
	nowInMS := int64(time.Now().UnixNano() / 1e6)

	flvWriter.WriteBWInfo.PeerIP = flvWriter.req.RemoteAddr
	flvWriter.WriteBWInfo.StreamId = streamid
	if isVideoFlag {
		flvWriter.WriteBWInfo.VideoDatainBytes = flvWriter.WriteBWInfo.VideoDatainBytes + length
	} else {
		flvWriter.WriteBWInfo.AudioDatainBytes = flvWriter.WriteBWInfo.AudioDatainBytes + length
	}

	if flvWriter.WriteBWInfo.LastTimestamp == 0 {
		flvWriter.WriteBWInfo.LastTimestamp = nowInMS
	} else if (nowInMS - flvWriter.WriteBWInfo.LastTimestamp) >= av.SAVE_STATICS_INTERVAL {
		diffTimestamp := (nowInMS - flvWriter.WriteBWInfo.LastTimestamp) / 1000

		flvWriter.WriteBWInfo.VideoSpeedInBytesperMS = (flvWriter.WriteBWInfo.VideoDatainBytes - flvWriter.WriteBWInfo.LastVideoDatainBytes) * 8 / uint64(diffTimestamp) / 1000
		flvWriter.WriteBWInfo.AudioSpeedInBytesperMS = (flvWriter.WriteBWInfo.AudioDatainBytes - flvWriter.WriteBWInfo.LastAudioDatainBytes) * 8 / uint64(diffTimestamp) / 1000

		flvWriter.WriteBWInfo.LastVideoDatainBytes = flvWriter.WriteBWInfo.VideoDatainBytes
		flvWriter.WriteBWInfo.LastAudioDatainBytes = flvWriter.WriteBWInfo.AudioDatainBytes
		flvWriter.WriteBWInfo.LastTimestamp = nowInMS
	}
}

func (flvWriter *FLVWriter) SendPacket() error {
	for {
		p, ok := <-flvWriter.packetQueue
		if ok {
			flvWriter.RWBaser.SetPreTime()
			h := flvWriter.buf[:headerLen]
			typeID := av.TAG_VIDEO
			if !p.IsVideo {
				if p.IsMetadata {
					var err error
					typeID = av.TAG_SCRIPTDATAAMF0
					p.Data, err = amf.MetaDataReform(p.Data, amf.DEL)
					if err != nil {
						return err
					}
				} else {
					typeID = av.TAG_AUDIO
					packetLen := len(p.Data) + 12
					flvWriter.SaveStatics(p.StreamID, uint64(packetLen), false)
				}
			} else {
				packetLen := len(p.Data) + 12
				flvWriter.SaveStatics(p.StreamID, uint64(packetLen), true)
			}

			dataLen := len(p.Data)
			timestamp := p.TimeStamp
			timestamp += flvWriter.BaseTimeStamp()
			flvWriter.RWBaser.RecTimeStamp(timestamp, uint32(typeID))

			preDataLen := dataLen + headerLen
			timestampbase := timestamp & 0xffffff
			timestampExt := timestamp >> 24 & 0xff

			pio.PutU8(h[0:1], uint8(typeID))
			pio.PutI24BE(h[1:4], int32(dataLen))
			pio.PutI24BE(h[4:7], int32(timestampbase))
			pio.PutU8(h[7:8], uint8(timestampExt))

			if _, err := flvWriter.ctx.Write(h); err != nil {
				return err
			}

			if _, err := flvWriter.ctx.Write(p.Data); err != nil {
				return err
			}

			pio.PutI32BE(h[:4], int32(preDataLen))
			if _, err := flvWriter.ctx.Write(h[:4]); err != nil {
				return err
			}
		} else {
			return errors.New("closed")
		}

	}

	return nil
}

func (flvWriter *FLVWriter) Wait() {
	select {
	case <-flvWriter.closedChan:
		return
	}
}

func (flvWriter *FLVWriter) Close(error) {
	log.Info("http flv closed")
	if !flvWriter.closed {
		close(flvWriter.packetQueue)
		close(flvWriter.closedChan)
	}
	flvWriter.closed = true
}

func (flvWriter *FLVWriter) Info() (ret av.Info) {
	ret.UID = flvWriter.Uid
	ret.URL = flvWriter.url
	ret.Key = flvWriter.app + "/" + flvWriter.title
	ret.Inter = true
	return
}
