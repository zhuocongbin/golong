package record

import (
	"errors"
	"fmt"
	"github.com/livego/av"
	"github.com/livego/configure"
	log "github.com/livego/logging"
	"os"
	"strings"
	"time"
)

const REC_CHANNEL_MAX = 1000
const REC_TIMEOUT = 60

type RecordWriter struct {
	info        av.Info
	recCfg      configure.RecordConfig
	dir         string
	path        string
	closedFlag  bool
	flvWriter   *FlvFileWriter
	packetChann chan *av.Packet
	av.RWBaser
}

func NewRecordWriter(info av.Info, recCfg configure.RecordConfig) *RecordWriter {
	subUrlArray := strings.Split(info.URL, "/")

	liveid := subUrlArray[len(subUrlArray)-1]

	dir := fmt.Sprintf("%s/%s", recCfg.Path, liveid)

	if !checkFileIsExist(dir) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			log.Errorf("create dir(%s) error", dir)
			return nil
		}
	}

	var path string
	index := 0

	for {
		if index == 0 {
			path = fmt.Sprintf("%s/%s.flv", dir, liveid)
		} else {
			path = fmt.Sprintf("%s/%s-%d.flv", dir, liveid, index)
		}
		index++
		if !checkFileIsExist(path) {
			break
		}
	}
	log.Infof("NewRecordWriter:%s", path)
	recW := &RecordWriter{
		info:        info,
		recCfg:      recCfg,
		dir:         dir,
		path:        path,
		closedFlag:  false,
		flvWriter:   NewFlvFileWriter(path),
		packetChann: make(chan *av.Packet, REC_CHANNEL_MAX),
		RWBaser:     av.NewRWBaser(time.Second * 10),
	}

	go recW.onWork()
	//go recW.onCheck()

	return recW
}

/*
func (self *RecordWriter) onCheck() {
	for {
		<-time.After(1000 * time.Millisecond)
		ret := self.Alive()
		if !ret {
			break
		}
	}
	self.Close(errors.New("record write alive timeout"))
}
*/
func (self *RecordWriter) onWork() {
	self.flvWriter.WriterHeader()
	for {
		packet, isOK := <-self.packetChann
		if !isOK {
			//log.Error("record packet channel closed")
			break
		}
		//log.Infof("onWork: length=%d", len(packet.Data))
		self.flvWriter.WriterPacket(packet)
	}
}

func (self *RecordWriter) Write(packet *av.Packet) error {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("record packet channel closed, path=%s", self.path)
		}
	}()

	if self.closedFlag {
		return errors.New(fmt.Sprintf("record(%s) is over", self.path))
	}
	self.SetPreTime()
	if len(self.packetChann) >= (REC_CHANNEL_MAX - 1) {
		log.Errorf("record writer channel is over, info[%v]", self.info)
		return errors.New(fmt.Sprintf("record writer channel is over, info[%v]", self.info))
	}

	//log.Infof("record write data length=%d, info=%v, type=%s", len(packet.Data), self.info, self.recCfg.Recordtype)
	self.packetChann <- packet
	return nil
}

func (self *RecordWriter) Info() av.Info {
	info := self.info

	return info
}

func (self *RecordWriter) Close(error) {
	log.Info("record writer close:", self.info)
	self.closedFlag = true
	close(self.packetChann)
	return
}
