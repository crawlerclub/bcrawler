package main

import (
	"encoding/json"
	"flag"
	"time"

	"crawler.club/et"
	"github.com/golang/glog"
	"zliu.org/ds"
	"zliu.org/filestore"
	"zliu.org/goutil"
)

var (
	start = flag.String("start", "", "the parser name for the start url, if omitted, will start from seeds.json")
	dir   = flag.String("dir", "data", "the data dir")
	qDir  = flag.String("q", "q", "the queue dir")
	sleep = flag.Int("sleep", -1, "in seconds")
)

func main() {
	flag.Parse()
	defer glog.Flush()

	fs, err := filestore.NewFileStore(*dir)
	if err != nil {
		glog.Fatal(err)
	}
	defer fs.Close()

	q, err := ds.OpenQueue(*qDir)
	if err != nil {
		glog.Fatal(err)
	}
	defer q.Close()

	if *start != "" {
		p, err := pool.GetParser(*start, false)
		if err != nil {
			glog.Fatal(err)
		}

		glog.Infof("start crawling from %s", p.ExampleUrl)

		if goutil.FileGuard("first.lock") {
			q.EnqueueObject(&et.UrlTask{ParserName: *start, Url: p.ExampleUrl})
		}
	} else {
		seeds, err := GetSeeds()
		if err != nil {
			glog.Fatal(err)
		}
		if len(seeds) == 0 {
			glog.Fatal("no seeds found!")
		}
		glog.Infof("start crawling from %d seeds in seeds.json", len(seeds))
		for _, seed := range seeds {
			q.EnqueueObject(seed)
		}

	}

	var task = new(et.UrlTask)
	for {
		if q.Length() == 0 {
			break
		}

		item, err := q.Dequeue()
		if err != nil {
			glog.Fatal(err)
		}
		if err = item.ToObject(task); err != nil {
			glog.Fatal(err)
		}

		glog.Info(task.Url)
		new_tasks, items, err := ParseTask(task)
		if err != nil {
			glog.Error(err)
			continue
		}
		for _, t := range new_tasks {
			q.EnqueueObject(t)
		}
		for _, item := range items {
			b, _ := json.Marshal(item)
			fs.WriteLine(b)
		}
		if *sleep > 0 {
			time.Sleep(time.Duration(*sleep) * time.Second)
		}
	}
}
