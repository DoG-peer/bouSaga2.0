package main

import (
	"errors"
	"fmt"
	"github.com/dullgiulio/pingo"
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
	p.voiceMng = makeVoiceManager(conf.Cache)
	// add voice
	for _, v := range conf.Voice {
		p.voiceMng.add(v)
	}
	time.Sleep(500 * time.Millisecond)

	return err
}

// Main task
func (p *Task) Main(configFile string, e *error) error {
	data, err := p.bbs.Read()
	if err != nil {
		return err
	}

	// find voice and  play the voice
	for _, res := range data {
		p.voiceMng.playAllMatch(res.body, 5*time.Second)
		fmt.Println(res.body)
	}
	nextNum := NextID(data, p.bbs.num)
	p.config.Res = nextNum
	p.bbs.MoveTo(nextNum)
	return nil
}

// SaveData is loaded by gobou
func (p *Task) SaveData(configFile string, e *error) error {
	return nil
}

// SaveConfig is loaded by gobou
func (p *Task) SaveConfig(configFile string, e *error) error {
	//p.config.Save(configFile)
	return nil
}

// Interval is loaded by gobou
func (p *Task) Interval(a string, d *time.Duration) error {
	*d = 30 * time.Second
	return nil
}

// End is loaded by gobou
func (p *Task) End() error {
	return nil
}

func makeVoiceManager(dir string) VoiceManager {
	return VoiceManager{
		voice: map[string]Voice{},
		dir:   dir,
	}
}

func main() {
	task := &Task{}
	pingo.Register(task)
	pingo.Run()
}
