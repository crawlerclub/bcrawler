package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"path/filepath"
	"sync"

	"crawler.club/dl"
	"crawler.club/et"
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

var pool = &Parsers{items: make(map[string]*et.Parser)}

func ParseTask(task *et.UrlTask) (
	[]*et.UrlTask, []map[string]interface{}, error) {
	return Parse(task.ParserName, task.Url)
}

func Parse(name, url string) (
	[]*et.UrlTask, []map[string]interface{}, error) {
	p, err := pool.GetParser(name, false)
	if err != nil {
		return nil, nil, err
	}
	ret := dl.Download(&dl.HttpRequest{Url: url, Retry: 3, Platform: "pc"})
	if ret.Error != nil {
		return nil, nil, ret.Error
	}
	return p.Parse(ret.Text, url)
}
