package rtmp

import (
	"errors"
	"flag"
	"fmt"
	"github.com/livego/av"
	"github.com/livego/configure"
	"github.com/livego/container/flv"
	log "github.com/livego/logging"
	"github.com/livego/protocol/rtmp/core"
	"github.com/livego/record"
	"github.com/livego/utils/uid"
	"net"
	"net/url"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

const (
	maxQueueNum = 1024
)

var (
	readTimeout  = flag.Int("readTimeout", 10, "read time out")
	writeTimeout = flag.Int("writeTimeout", 10, "write time out")
)

type Client struct {
	handler av.Handler
	getter  av.GetWriter
}

func NewRtmpClient(h av.Handler, getter av.GetWriter) *Client {
	return &Client{
		handler: h,
		getter:  getter,
	}
}

func (c *Client) Dial(url string, method string) error {
	connClient := core.NewConnClient()
	if err := connClient.Start(url, method); err != nil {
		return err
	}
	if method == av.PUBLISH {
		writer := NewVirWriter(connClient)
		log.Infof("client Dial call NewVirWriter url=%s, method=%s", url, method)
		c.handler.HandleWriter(writer)
	} else if method == av.PLAY {
		reader := NewVirReader(connClient)
		log.Infof("client Dial call NewVirReader url=%s, method=%s", url, method)
		c.handler.HandleReader(reader)
		if c.getter != nil {
			writer := c.getter.GetWriter(reader.Info())
			c.handler.HandleWriter(writer)
		}
	}
	return nil
}

func (c *Client) GetHandle() av.Handler {
	return c.handler
}

type Server struct {
	handler av.Handler
	getter  av.GetWriter
}

func NewRtmpServer(h av.Handler, getter av.GetWriter) *Server {
	return &Server{
		handler: h,
		getter:  getter,
	}
}

func (s *Server) Serve(listener net.Listener) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("rtmp serve panic: ", r)
		}
	}()

	for {
		var netconn net.Conn
		netconn, err = listener.Accept()
		if err != nil {
			return
		}

		conn := core.NewConn(netconn, 4*1024)
		//log.Info("new client, connect remote:", conn.RemoteAddr().String(),
		//	"local:", conn.LocalAddr().String())
		go s.handleConn(conn)
	}
}

func (s *Server) ExecPush(key string) {
	execList := configure.GetExecPush()

	for _, execItem := range execList {
		cmdString := fmt.Sprintf("%s -k %s", execItem, key)
		go func(cmdString string) {
			log.Info("ExecPush:", cmdString)
			cmd := exec.Command("/bin/sh", "-c", cmdString)
			_, err := cmd.Output()
			if err != nil {
				log.Info("Excute error:", err)
			}
		}(cmdString)
	}
}

func (s *Server) handleConn(conn *core.Conn) error {
	if err := conn.HandshakeServer(); err != nil {
		conn.Close()
		//log.Error("handleConn HandshakeServer err:", err)
		return err
	}
	connServer := core.NewConnServer(conn)

	if err := connServer.ReadMsg(); err != nil {
		conn.Close()
		log.Error("handleConn read msg err:", err)
		return err
	}

	appname, name, url, remoteconn := connServer.GetInfo()
	log.Infof("handleConn: appname=%s, name=%s, url=%s, peerIP=%s", appname, name, url, remoteconn.RemoteAddr().String())

	log.Infof("handleConn: IsPublisher=%v", connServer.IsPublisher())
	if connServer.IsPublisher() {
		reader := NewVirReader(connServer)
		/*
			//从宽处理，新来的覆盖旧的流。
				if rtmpstream, ok := s.handler.(*RtmpStream); ok {
					if rtmpstream.IsExist(reader) {
						connServer.Close(nil)
						return nil
					}
				}
		*/
		s.handler.HandleReader(reader)
		log.Infof("new publisher: %v", reader.Info())

		if s.getter != nil {
			writeType := reflect.TypeOf(s.getter)
			log.Infof("handleConn:writeType=%v", writeType)
			writer := s.getter.GetWriter(reader.Info())
			s.handler.HandleWriter(writer)
		}
		if ret, recCfg := configure.IsRecordEnable(url); ret {
			info := reader.Info()

			recWriter := record.NewRecordWriter(info, recCfg)
			log.Infof("new publish add record: %v", info)
			s.handler.HandleWriter(recWriter)
		}
		s.ExecPush(reader.Info().Key)
	} else {
		writer := NewVirWriter(connServer)
		log.Infof("new player: %+v", writer.Info())
		s.handler.HandleWriter(writer)
	}

	return nil
}

type GetInFo interface {
	GetInfo() (string, string, string, *core.Conn) //app string, name string, url string, *core.Conn
}

type StreamReadWriteCloser interface {
	GetInFo
	Close(error)
	Write(core.ChunkStream) error
	Read(c *core.ChunkStream) error
}

type VirWriter struct {
	Uid    string
	closed bool
	av.RWBaser
	conn        StreamReadWriteCloser
	packetQueue chan *av.Packet
	WriteBWInfo av.StaticsBW
}

func NewVirWriter(conn StreamReadWriteCloser) *VirWriter {
	ret := &VirWriter{
		Uid:         uid.NewId(),
		conn:        conn,
		RWBaser:     av.NewRWBaser(time.Second * time.Duration(*writeTimeout)),
		packetQueue: make(chan *av.Packet, maxQueueNum),
		WriteBWInfo: av.StaticsBW{0, "", 0, 0, 0, 0, 0, 0, 0},
	}

	go ret.Check()
	go func() {
		err := ret.SendPacket()
		if err != nil {
			log.Error(err)
		}
	}()
	return ret
}

func (v *VirWriter) SaveStatics(streamid uint32, length uint64, isVideoFlag bool) {
	nowInMS := int64(time.Now().UnixNano() / 1e6)

	_, _, _, conn := v.conn.GetInfo()
	v.WriteBWInfo.PeerIP = conn.RemoteAddr().String()
	v.WriteBWInfo.StreamId = streamid
	if isVideoFlag {
		v.WriteBWInfo.VideoDatainBytes = v.WriteBWInfo.VideoDatainBytes + length
	} else {
		v.WriteBWInfo.AudioDatainBytes = v.WriteBWInfo.AudioDatainBytes + length
	}

	if v.WriteBWInfo.LastTimestamp == 0 {
		v.WriteBWInfo.LastTimestamp = nowInMS
	} else if (nowInMS - v.WriteBWInfo.LastTimestamp) >= av.SAVE_STATICS_INTERVAL {
		diffTimestamp := (nowInMS - v.WriteBWInfo.LastTimestamp) / 1000

		v.WriteBWInfo.VideoSpeedInBytesperMS = (v.WriteBWInfo.VideoDatainBytes - v.WriteBWInfo.LastVideoDatainBytes) * 8 / uint64(diffTimestamp) / 1000
		v.WriteBWInfo.AudioSpeedInBytesperMS = (v.WriteBWInfo.AudioDatainBytes - v.WriteBWInfo.LastAudioDatainBytes) * 8 / uint64(diffTimestamp) / 1000

		v.WriteBWInfo.LastVideoDatainBytes = v.WriteBWInfo.VideoDatainBytes
		v.WriteBWInfo.LastAudioDatainBytes = v.WriteBWInfo.AudioDatainBytes
		v.WriteBWInfo.LastTimestamp = nowInMS
	}
}

func (v *VirWriter) Check() {
	var c core.ChunkStream
	for {
		if err := v.conn.Read(&c); err != nil {
			v.Close(err)
			return
		}
	}
}

func (v *VirWriter) DropPacket(pktQue chan *av.Packet, info av.Info) {
	log.Infof("[%v] packet queue max!!!", info)
	for i := 0; i < maxQueueNum-84; i++ {
		tmpPkt, ok := <-pktQue
		// try to don't drop audio
		if ok && tmpPkt.IsAudio {
			if len(pktQue) > maxQueueNum-2 {
				log.Info("drop audio pkt")
				<-pktQue
			} else {
				pktQue <- tmpPkt
			}

		}

		if ok && tmpPkt.IsVideo {
			videoPkt, ok := tmpPkt.Header.(av.VideoPacketHeader)
			// dont't drop sps config and dont't drop key frame
			if ok && (videoPkt.IsSeq() || videoPkt.IsKeyFrame()) {
				pktQue <- tmpPkt
			}
			if len(pktQue) > maxQueueNum-10 {
				log.Info("drop video pkt")
				<-pktQue
			}
		}

	}
	log.Info("packet queue len: ", len(pktQue))
}

//
func (v *VirWriter) Write(p *av.Packet) (err error) {
	err = nil

	if v.closed {
		err = errors.New("VirWriter closed")
		return
	}
	defer func() {
		if e := recover(); e != nil {
			errString := fmt.Sprintf("VirWriter has already been closed:%v", e)
			err = errors.New(errString)
		}
	}()
	if len(v.packetQueue) >= maxQueueNum-24 {
		v.DropPacket(v.packetQueue, v.Info())
	} else {
		v.packetQueue <- p
	}

	return
}

func (v *VirWriter) SendPacket() error {
	var cs core.ChunkStream
	for {
		p, ok := <-v.packetQueue
		if ok {
			cs.Data = p.Data
			cs.Length = uint32(len(p.Data))
			cs.StreamID = p.StreamID
			cs.Timestamp = p.TimeStamp
			cs.Timestamp += v.BaseTimeStamp()

			if p.IsVideo {
				cs.TypeID = av.TAG_VIDEO
			} else {
				if p.IsMetadata {
					cs.TypeID = av.TAG_SCRIPTDATAAMF0
				} else {
					cs.TypeID = av.TAG_AUDIO
				}
			}

			//log.Printf("TypeID=%d, Length=%d", cs.TypeID, cs.Length)
			v.SaveStatics(p.StreamID, uint64(cs.Length), p.IsVideo)
			v.SetPreTime()
			v.RecTimeStamp(cs.Timestamp, cs.TypeID)
			err := v.conn.Write(cs)
			if err != nil {
				v.closed = true
				return err
			}

		} else {
			return errors.New("closed")
		}

	}
	return nil
}

func (v *VirWriter) Info() (ret av.Info) {
	ret.UID = v.Uid
	_, _, URL, _ := v.conn.GetInfo()
	ret.URL = URL
	_url, err := url.Parse(URL)
	if err != nil {
		log.Info(err)
	}
	ret.Key = strings.TrimLeft(_url.Path, "/")
	ret.Inter = true
	return
}

func (v *VirWriter) Close(err error) {
	log.Info("VirWriter.player ", v.Info(), "closed: "+err.Error())
	if !v.closed {
		close(v.packetQueue)
	}
	v.closed = true
	v.conn.Close(err)
}

type VirReader struct {
	Uid string
	av.RWBaser
	demuxer    *flv.Demuxer
	conn       StreamReadWriteCloser
	ReadBWInfo av.StaticsBW
}

func NewVirReader(conn StreamReadWriteCloser) *VirReader {
	return &VirReader{
		Uid:        uid.NewId(),
		conn:       conn,
		RWBaser:    av.NewRWBaser(time.Second * time.Duration(*writeTimeout)),
		demuxer:    flv.NewDemuxer(),
		ReadBWInfo: av.StaticsBW{0, "", 0, 0, 0, 0, 0, 0, 0},
	}
}

func (v *VirReader) SaveStatics(streamid uint32, length uint64, isVideoFlag bool) {
	nowInMS := int64(time.Now().UnixNano() / 1e6)

	_, _, _, conn := v.conn.GetInfo()
	v.ReadBWInfo.StreamId = streamid
	v.ReadBWInfo.PeerIP = conn.RemoteAddr().String()
	if isVideoFlag {
		v.ReadBWInfo.VideoDatainBytes = v.ReadBWInfo.VideoDatainBytes + length
	} else {
		v.ReadBWInfo.AudioDatainBytes = v.ReadBWInfo.AudioDatainBytes + length
	}

	if v.ReadBWInfo.LastTimestamp == 0 {
		v.ReadBWInfo.LastTimestamp = nowInMS
	} else if (nowInMS - v.ReadBWInfo.LastTimestamp) >= av.SAVE_STATICS_INTERVAL {
		diffTimestamp := (nowInMS - v.ReadBWInfo.LastTimestamp) / 1000

		//log.Printf("now=%d, last=%d, diff=%d", nowInMS, v.ReadBWInfo.LastTimestamp, diffTimestamp)
		v.ReadBWInfo.VideoSpeedInBytesperMS = (v.ReadBWInfo.VideoDatainBytes - v.ReadBWInfo.LastVideoDatainBytes) * 8 / uint64(diffTimestamp) / 1000
		v.ReadBWInfo.AudioSpeedInBytesperMS = (v.ReadBWInfo.AudioDatainBytes - v.ReadBWInfo.LastAudioDatainBytes) * 8 / uint64(diffTimestamp) / 1000

		v.ReadBWInfo.LastVideoDatainBytes = v.ReadBWInfo.VideoDatainBytes
		v.ReadBWInfo.LastAudioDatainBytes = v.ReadBWInfo.AudioDatainBytes
		v.ReadBWInfo.LastTimestamp = nowInMS
	}
}

func (v *VirReader) Read(p *av.Packet) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("rtmp read packet panic: ", r)
		}
	}()

	v.SetPreTime()
	cs := &core.ChunkStream{}

	for {
		err = v.conn.Read(cs)
		if err != nil {
			return err
		}
		if cs.TypeID == av.TAG_AUDIO ||
			cs.TypeID == av.TAG_VIDEO ||
			cs.TypeID == av.TAG_SCRIPTDATAAMF0 ||
			cs.TypeID == av.TAG_SCRIPTDATAAMF3 {
			break
		}
	}

	p.IsAudio = cs.TypeID == av.TAG_AUDIO
	p.IsVideo = cs.TypeID == av.TAG_VIDEO
	p.IsMetadata = cs.TypeID == av.TAG_SCRIPTDATAAMF0 || cs.TypeID == av.TAG_SCRIPTDATAAMF3
	p.StreamID = cs.StreamID
	p.Data = cs.Data
	p.TimeStamp = cs.Timestamp

	v.SaveStatics(p.StreamID, uint64(len(p.Data)), p.IsVideo)
	v.demuxer.DemuxH(p)
	return err
}

func (v *VirReader) Info() (ret av.Info) {
	ret.UID = v.Uid
	_, _, URL, _ := v.conn.GetInfo()
	ret.URL = URL
	_url, err := url.Parse(URL)
	if err != nil {
		log.Error(err)
	}
	ret.Key = strings.TrimLeft(_url.Path, "/")
	return
}

func (v *VirReader) Close(err error) {
	log.Info("VirReader.publisher ", v.Info(), "closed: "+err.Error())
	v.conn.Close(err)
}
