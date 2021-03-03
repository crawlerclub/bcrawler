package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"crawler.club/dl"
	"crawler.club/et"
	"github.com/crawlerclub/ce"
)

var (
	conf = flag.String("conf", "./conf", "dir for parsers conf")
)

type Parsers struct {
	sync.Mutex
	items map[string]*et.Parser
}

func (p *Parsers) GetParser(name string, refresh bool) (*et.Parser, error) {
	p.Lock()
	defer p.Unlock()
	if !refresh && p.items[name] != nil {
		return p.items[name], nil
	}
	file := filepath.Join(*conf, "parsers", name+".json")
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	parser := new(et.Parser)
	if err := json.Unmarshal(content, parser); err != nil {
		return nil, err
	}
	p.items[name] = parser
	return parser, nil
}

func GetSeeds() ([]*et.UrlTask, error) {
	seedsFile := filepath.Join(*conf, "seeds.json")
	content, err := ioutil.ReadFile(seedsFile)
	if err != nil {
		return nil, err
	}
	var seeds []*et.UrlTask
	if err = json.Unmarshal(content, &seeds); err != nil {
		return nil, err
	}
	return seeds, nil
}

var pool = &Parsers{items: make(map[string]*et.Parser)}

func ParseTask(task *et.UrlTask) (
	[]*et.UrlTask, []map[string]interface{}, error) {
	return Parse(task.ParserName, task.Url)
}

func Parse(name, url string) (
	[]*et.UrlTask, []map[string]interface{}, error) {
	ret := dl.Download(&dl.HttpRequest{Url: url, Retry: 3, Platform: "pc"})
	if ret.Error != nil {
		return nil, nil, ret.Error
	}
	page := ret.Text
	switch strings.ToLower(name) {
	case "content_":
		items := strings.Split(ret.RemoteAddr, ":")
		ip := ""
		if len(items) > 0 {
			ip = items[0]
		}
		doc := ce.ParsePro(url, page, ip, false)
		return nil, []map[string]interface{}{
			map[string]interface{}{"doc": doc}}, nil
	case "link_":
		links, err := et.ParseNewLinks(page, url)
		if err != nil {
			return nil, nil, err
		}
		var tasks []*et.UrlTask
		for _, link := range links {
			tasks = append(tasks, &et.UrlTask{
				ParserName: "content_", Url: link})
		}
		return tasks, nil, nil
	case "raw_":
		return nil, []map[string]interface{}{
			map[string]interface{}{
				"base64_content": base64.StdEncoding.EncodeToString(ret.Content)}}, nil
	default:
		p, err := pool.GetParser(name, false)
		if err != nil {
			return nil, nil, err
		}
		return p.Parse(page, url)
	}
}
