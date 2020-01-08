# bcrawler
A standalone crawling  tool for extracting data from the web.

## Usage

```sh
Usage of ./bcrawler:
  -alsologtostderr
    	log to standard error as well as files
  -conf string
    	dir for parsers conf (default "./conf")
  -dir string
    	the data dir (default "data")
  -log_dir string
    	If non-empty, write log files in this directory
  -logtostderr
    	log to standard error instead of files
  -q string
    	the queue dir (default "q")
  -sleep int
    	in seconds (default -1)
  -start string
    	the parser name for the start url (default "addr_year")
```

There should be a `parsers` directory under `./conf`, which stores all the parser configurations in json format.

