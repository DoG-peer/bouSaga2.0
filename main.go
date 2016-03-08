package main

import (
	"errors"
	"fmt"
	"github.com/DoG-peer/gobou/utils"
	"log"
	"time"
)

// "github.com/dullgiulio/pingo"

// Task is interface with gobou
type Task struct {
	bbs      BBS
	voiceMng VoiceManager
	config   Config
}

// ErrConfig is error about config
type ErrConfig error

// ErrBBS is error about bbs
type ErrBBS error

var (
	errConfig = ErrConfig(errors.New("Error on Config"))
	errBBS    = ErrBBS(errors.New("Error on BBS"))
)

// Configure is
// load config
// make wav files
func (p *Task) Configure(configFile string, e *error) error {
	// load config file
	var conf Config
	err := conf.Load(configFile)
	if err != nil {
		*e = err
		fmt.Println(errConfig)
		return *e
	}

	// needs mkdir

	p.config = conf
	p.bbs.baseurl, err = GetRawURL(conf.URL)
	p.bbs.MoveTo(conf.Res)

	// voice
	p.voiceMng = makeVoiceManager(conf.Cache, conf.Jtalk)
	// add voice
	for _, v := range conf.Voice {
		err := p.voiceMng.add(v)
		if err != nil {
			log.Println(err)
		}
	}
	time.Sleep(500 * time.Millisecond)

	return err
}

// Main task
func (p *Task) Main(configFile string, s *[]gobou.Message) error {
	data, err := p.bbs.Read()
	if err != nil {
		return err
	}
	*s = []gobou.Message{}
	// find voice and  play the voice
	for _, res := range data {
		p.voiceMng.playAllMatch(res.body, 5*time.Second)
		*s = append(*s, gobou.Print(res.body))
	}
	nextNum := NextID(data, p.bbs.num)
	p.config.Res = nextNum
	p.bbs.MoveTo(nextNum)
	if len(*s) != 0 {
		*s = append(*s, gobou.Notify("レスがありました"))
	}
	return nil
}

// SaveData is loaded by gobou
func (p *Task) SaveData(configFile string, e *error) error {
	return nil
}

// SaveConfig is loaded by gobou
func (p *Task) SaveConfig(configFile string, e *error) error {
	p.config.Save(configFile)
	return nil
}

// Interval is loaded by gobou
func (p *Task) Interval(a string, d *time.Duration) error {
	*d = 15 * time.Second
	return nil
}

/*
// End is loaded by gobou
func (p *Task) End() error {
	return nil
}
*/
func makeVoiceManager(dir string, jtalk Jtalk) VoiceManager {
	return VoiceManager{
		voice: map[string]Voice{},
		dir:   dir,
		jtalk: jtalk,
	}
}

func main() {
	gobou.Register(&Task{})
}
