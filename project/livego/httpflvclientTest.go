package main

import (
	"github.com/livego/protocol/rtmp/rtmprelay"
	"log"
	"time"
)

/*
type FlvRcvHandle struct {
	version string
}

func (self *FlvRcvHandle) HandleFlvData(data []byte, length int) error {
	return httpflv.WriteFlvFile(data, length)
}
*/
func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	/*
		client := httpflv.NewHttpFlvClient("http://pull99.a8.com/live/1499336947298690.flv")

		handle := &FlvRcvHandle{
			version: "flv1.0",
		}
		client.Start(handle)
	*/
	flvurl := "http://pull2.a8.com/live/1500359051219357.flv"
	rtmpurl := "rtmp://127.0.0.1/live/trans/inke/mlinkm/shiwei123"
	flvPull := rtmprelay.NewFlvPull(&flvurl, &rtmpurl)
	err := flvPull.Start()
	if err != nil {
		return
	}

	time.Sleep(time.Second * 60)
	flvPull.Stop()
	done := make(chan int)

	<-done
}
