package record

import (
	"github.com/livego/av"
	log "github.com/livego/logging"
	//"github.com/livego/protocol/amf"
	"errors"
	"fmt"
	"github.com/livego/utils/pio"
	"os"
)

const ITEM_HEADER_LEN = 11

type FlvFileWriter struct {
	path   string
	header []byte
}

func checkFileIsExist(filename string) bool {
	exist := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func NewFlvFileWriter(path string) *FlvFileWriter {
	return &FlvFileWriter{
		path:   path,
		header: make([]byte, ITEM_HEADER_LEN),
	}
}

func (self *FlvFileWriter) WriterHeader() error {
	var f *os.File
	var err error

	FLV_HEADER := [13]byte{0x46, 0x4c, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00}
	log.Info("WriterHeader:", self.path)

	if checkFileIsExist(self.path) { //如果文件存在
		log.Errorf("file(%s) exist, WriterHeader error", self.path)
		return errors.New(fmt.Sprintf("file(%s) exist, WriterHeader error", self.path))
		/*
			f, err = os.OpenFile(self.path, os.O_WRONLY|os.O_APPEND, 0666) //打开文件
			if err != nil {
				log.Errorf("openfile(%s) error:%v", self.path, err)
				return err
			}
		*/
	} else {
		f, err = os.Create(self.path) //创建文件
		if err != nil {
			log.Errorf("createfile(%s) error:%v", self.path, err)
			return err
		}
	}
	f.Write(FLV_HEADER[0:])

	f.Close()
	return nil
}

func (self *FlvFileWriter) WriterPacket(p *av.Packet) error {
	var f *os.File
	var err error

	f, err = os.OpenFile(self.path, os.O_WRONLY|os.O_APPEND, 0666) //打开文件
	if err != nil {
		log.Errorf("openfile(%s) error:%v", self.path, err)
		return err
	}
	defer f.Close()
	typeID := av.TAG_VIDEO
	if !p.IsVideo {
		if p.IsMetadata {
			log.Infof("record flv drop metadata:%02x %02x", p.Data[0], p.Data[1])
			return nil
		} else {
			aHdr, ok := p.Header.(av.AudioPacketHeader)
			if !ok {
				log.Error("flv file write audio header error")
				return errors.New(fmt.Sprintf("flv file write audio header error"))
			}
			if (aHdr.AACPacketType() == av.SOUND_AAC) && (aHdr.SoundFormat() == av.AAC_SEQHDR) {
				log.Infof("record flv audio header:%02x %02x", p.Data[0], p.Data[1])
			}
			typeID = av.TAG_AUDIO
		}
	} else {
		vHdr, ok := p.Header.(av.VideoPacketHeader)
		if !ok {
			log.Error("flv file write video header error")
			return errors.New(fmt.Sprintf("flv file write video header error"))
		}
		if vHdr.IsSeq() {
			log.Infof("video header info:%02x %02x %02x %02x", p.Data[0], p.Data[1], p.Data[2], p.Data[3])
		}
	}

	dataLen := len(p.Data)
	preDataLen := dataLen + ITEM_HEADER_LEN

	timestamp := p.TimeStamp
	timestampbase := timestamp & 0xffffff
	timestampExt := timestamp >> 24 & 0xff
	pio.PutU8(self.header[0:1], uint8(typeID))
	pio.PutI24BE(self.header[1:4], int32(dataLen))
	pio.PutI24BE(self.header[4:7], int32(timestampbase))
	pio.PutU8(self.header[7:8], uint8(timestampExt))

	f.Write(self.header[0:])
	f.Write(p.Data[:])

	pio.PutI32BE(self.header[:4], int32(preDataLen))
	f.Write(self.header[:4])

	return nil
}
