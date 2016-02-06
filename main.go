package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

/*
	"strings"
	"github.com/dullgiulio/pingo"
*/

// Edit your task
type Task struct {
	bbsURL   string
	resNum   int
	cacheDir string
	url      string
	vinfo    VoiceInfo
}

type Config struct {
	bbsURL    string        `json:"url"`
	resNum    int           `json:"res"`
	cacheDir  string        `json:"cache"`
	voiceList []VoiceConfig `json:"voice"`
}
type VoiceConfig struct {
	name      string `json:"name"`
	words     string `json:"words"`
	condition string `json:"condition"`
}

func (c Config) String() string {
	return c.bbsURL + c.cacheDir
}

func loadJSONFile(fname string) (*Config, error) {
	file, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(file))
	var conf Config
	// err = json.Unmarshal(file, &conf)
	var c map[string]interface{}
	err = json.Unmarshal(file, &c)
	if err != nil {
		return nil, err
	}
	fmt.Println(c)
	fmt.Println(c["cache"])
	//if c["cache"] != "" {
	//		conf.cacheDir = string(c["cache"])
	//	}
	return &conf, nil
}

func getRawURL(url string) (string, error) {
	r := regexp.MustCompile(`/(\d+)/(\d+)/$`)
	if !r.MatchString(url) {
		return "", fmt.Errorf("%s should be like ***/1234/5678/", url)
	}
	rawmodeURLFormat := "http://jbbs.shitaraba.net/bbs/rawmode.cgi/game/($1)/($2)"
	return r.ReplaceAllString(url, rawmodeURLFormat), nil
}

func (p *Task) Configure(configFile string, e *error) error {
	// load config file
	conf, err := loadJSONFile(configFile)
	if err != nil {
		config := map[string]interface{}{}
		config["cache"] = "/home/your_user_name/.cache/gobou/bouSaga2.0"
		config["url"] = "**Write your bbsURL**"
		config["res"] = 500
		v, _ := json.Marshal(config)
		fmt.Println(configFile, "is not found")
		fmt.Println(string(v))
		*e = err
		fmt.Println(err)
		return *e
	}
	fmt.Println(json.Marshal(*conf))
	// needs mkdir

	p.cacheDir = conf.cacheDir
	p.bbsURL = conf.bbsURL
	p.resNum = conf.resNum
	// p.cacheDir = "/home/user/.cache/gobou/bouSaga2.0"
	// p.bbsURL = "http://jbbs.shitaraba.net/bbs/rawmode.cgi/game/57358/1389905050/"
	// p.resNum = 500
	p.url = p.bbsURL + strconv.Itoa(p.resNum) + "-"

	// voice
	p.vinfo = makeVoiceInfo(p.cacheDir)
	p.vinfo.addVoice("hello", "ぷええーーーー", "sta")
	p.vinfo.addVoice("hello2", "すたすた", "sta")
	p.vinfo.addVoice("tada", "ただの香車を", "タダ")
	time.Sleep(500 * time.Millisecond)

	return nil
}

func (p *Task) Main(configFile string, e *error) error {
	res, err := http.Get(p.url)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	arr, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		return err2
	}

	data := string(arr)
	fmt.Println(data)
	// find voice
	// play the voice
	return *e
}
func (p *Task) SaveData(configFile string, e *error) error {
	return nil
}
func (p *Task) SaveConfig(configFile string, e *error) error {
	return nil
}
func (p *Task) Interval(a string, d *time.Duration) error {
	*d = 30 * time.Second
	return nil
}
func (p *Task) End() error {
	return nil
}

// VoiceInfo is
type VoiceInfo struct {
	voice map[string]Voice
	dir   string
}

func makeVoiceInfo(dir string) VoiceInfo {
	vinfo := VoiceInfo{}
	vinfo.voice = map[string]Voice{}
	vinfo.dir = dir
	return vinfo
}

type Voice struct {
	path      string
	condition *regexp.Regexp
	name      string
	words     string
}

func makeVoice(name, words, condition, dir string) (Voice, error) {
	r, e := regexp.Compile(condition)
	v := Voice{}
	if e != nil {
		return v, e
	}

	v.name = name
	v.words = words
	v.path = filepath.Join(dir, name+".wav")
	v.condition = r
	return v, nil
}

func (v *Voice) makeWavFile() error {
	voice := "/usr/share/hts-voice/nitech-jp-atr503-m001/nitech_jp_atr503_m001.htsvoice"
	dict := "/var/lib/mecab/dic/open-jtalk/naist-jdic"

	cmd := exec.Command("open_jtalk", "-x", dict, "-m", voice, "-ow", v.path)
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

func (vinfo *VoiceInfo) addVoice(name, words, cond string) error {
	v, err := makeVoice(name, words, cond, vinfo.dir)
	if err != nil {
		return err
	}
	vinfo.voice[name] = v
	err = v.makeWavFile()
	return err
}

func (vinfo *VoiceInfo) play(name string) bool {
	v, found := vinfo.voice[name]
	if found {
		v.play()
	}
	return found
}

func (vinfo *VoiceInfo) playAll(t time.Duration) {
	for _, v := range vinfo.voice {
		v.play()
		time.Sleep(t)
	}
}

func main() {
	task := &Task{}
	var err error
	task.Configure("/home/user/.config/gobou/plugin_config/bouSaga2.0.json", &err)
	//vinfo.play("hello")
	//task.vinfo.playAll(1 * time.Second)
	//pingo.Register(task)
	//pingo.Run()
}
