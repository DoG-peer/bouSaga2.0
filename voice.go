package main

import (
	"io"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
)

// VoiceConfig is used in config file
type VoiceConfig struct {
	Name      string
	Words     string
	Condition string
}

// VoiceManager manages all voice
type VoiceManager struct {
	voice map[string]Voice
	dir   string
	jtalk Jtalk
}

// Jtalk is about open_jtalk
type Jtalk struct {
	Voice      string
	Dictionary string
}

// Voice is a data about a voice
type Voice struct {
	path      string
	condition *regexp.Regexp
	name      string
	words     string
}

// MakeWavFile makes .wav file
func (mng *VoiceManager) MakeWavFile(v Voice) error {
	cmd := mng.jtalk.Command(v.path)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()
	if err = cmd.Start(); err != nil {
		return err
	}
	io.WriteString(stdin, v.words)
	return nil
}

func (v *Voice) play() {
	exec.Command("aplay", "--quiet", v.path).Start()
}

func (mng *VoiceManager) add(vc VoiceConfig) error {
	return mng.addVoice(vc.Name, vc.Words, vc.Condition)
}

func (mng *VoiceManager) addVoice(name, words, cond string) error {
	v, err := makeVoice(name, words, cond, mng.dir)
	if err != nil {
		return err
	}
	mng.voice[name] = v
	return mng.MakeWavFile(v)
}

func (mng *VoiceManager) play(name string) bool {
	v, found := mng.voice[name]
	if found {
		v.play()
	}
	return found
}

func (mng *VoiceManager) playAll(t time.Duration) {
	for _, v := range mng.voice {
		v.play()
		time.Sleep(t)
	}
}

func (mng *VoiceManager) playAllMatch(s string, t time.Duration) {
	for _, v := range mng.voice {
		if v.condition.MatchString(s) {
			v.play()
			time.Sleep(t)
		}
	}
}

func (mng *VoiceManager) playOneMatch(s string) {
	for _, v := range mng.voice {
		if v.condition.MatchString(s) {
			v.play()
			break
		}
	}
}

func makeVoice(name, words, condition, dir string) (Voice, error) {
	r, e := regexp.Compile(condition)
	if e != nil {
		return Voice{}, e
	}
	return Voice{
		name:      name,
		words:     words,
		path:      filepath.Join(dir, name+".wav"),
		condition: r,
	}, nil
}

// Command returns open_jtalk command
func (j *Jtalk) Command(out string) *exec.Cmd {
	cmd := exec.Command("open_jtalk", "-x", j.Dictionary, "-m", j.Voice, "-ow", out)
	return cmd
}
