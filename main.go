package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"
)

type options struct {
	url string;
	filter_code *regexp.Regexp;
	filter_string *regexp.Regexp;
	error error;
}

type channels struct {
	responseChannel chan []string;
	toRequestChannel chan string;
	doneRequesting chan struct{};
}

func NewChannels() channels {
	var channels channels;
	channels.toRequestChannel = make(chan string);
	channels.responseChannel = make(chan []string);
	channels.doneRequesting = make(chan struct{});
	return channels 
}

func (channels channels) close() {
	close(channels.doneRequesting);
	close(channels.toRequestChannel);
}

func main(){
	args := os.Args;
	options := parse_args(args);
	if options.error!=nil||options.filter_code==nil {
		usage();
		return;
	}
	url := options.url;
	r, err := regexp.Compile("https?://(.*?)/.*$");
	if err!=nil {
		fmt.Println("Invalid url.")
		return;
	}
	domain := r.FindStringSubmatch(url)[1];
	fmt.Printf("[*] Searching in domain %s.\n",domain)

	channels := NewChannels();
	go func() {
		channels.toRequestChannel<-url;
	}()
	go big_request(channels,options);
	retrieve_pages(url,domain,channels,options);
	channels.close();
}

func request(url string,channels channels,options options) {
	defer func() {
		channels.doneRequesting<-struct{}{};
	}();
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(url);
	if err!=nil {
		return;
	}
	if (!options.filter_code.Match([]byte(resp.Status))) {
		fmt.Println(url,resp.Status);
		body, err := io.ReadAll(resp.Body);
		if err!=nil {
			return;
		}
		channels.responseChannel<-[]string{url,string(body)}
	}
}

func big_request(channels channels,options options) {
	var requestedcount int;
	var parsedcount int;
	closed := false;
	for {
		select {
		case <-channels.doneRequesting:
			requestedcount++;
		case url := <-channels.toRequestChannel:
			parsedcount++;
			go request(url,channels,options);
		default:
			if (requestedcount>1&&requestedcount>=parsedcount) {
				lastcount := requestedcount;
				time.Sleep(time.Millisecond*100);
				if ((requestedcount == lastcount)&&!closed) {
					closed=true;
					close(channels.responseChannel);
				}
			}
		}
	}
}


func retrieve_pages(url string,domain string,channels channels,options options) {
	var visited []string;
	visited = append(visited, url)
	for array := range(channels.responseChannel) {
		url := array[0];
		parse := array[1];
		r1, _ := regexp.Compile(fmt.Sprintf("(https?://.*?%s.*?)[\" ']",domain));
		r2, _ := regexp.Compile("href=\"(/?[^(http)].*?)[\" ']|src=\"(/?[^(http)].*?)[\" ']");
		match1 := r1.FindAllStringSubmatch(parse,-1);
		match2 := r2.FindAllSubmatch([]byte(parse),-1);
		var is_filtered bool;
		for _,m := range(match1) {
			option := m[1];
			if (options.filter_string==nil) {
				is_filtered = false;
			} else {
				is_filtered = options.filter_string.Match([]byte(option));
			}
			if !is_filtered {
				if !slices.Contains(visited,option) {
					visited = append(visited, option);
					channels.toRequestChannel<-option;
				}
			}
		}
		for _,m := range(match2) {
			for i,s := range(m) {
				if (i==1||i==2) {
					if (options.filter_string==nil) {
						is_filtered = false;
					} else {
						is_filtered = options.filter_string.Match([]byte(m[1]));
					}
					if !is_filtered {
						str := url + string(s);
						if !slices.Contains(visited,str) {
							visited = append(visited, str);
							channels.toRequestChannel<-str;
						}
					}
				}
			}
		}
	}
}

func parse_args(args []string) options {
	var opt options;
	for i,flag := range(args) {
		if (i!=len(args)) {
			switch string(flag) {
			case "-u","--url":
				opt.url = args[i+1];
			case "-fs","--filter-string":
				array := strings.Split(args[i+1], ",");
				filter := strings.Join(array,"|");
				reg,err := regexp.Compile(filter);
				if err!=nil {
					return options{error: err};
				}
				opt.filter_string = reg;
			case "-fc","--filter-code":
				array := strings.Split(args[i+1],",");
				filter := strings.Join(array,"|");
				reg,err := regexp.Compile(filter);
				if err!=nil {
					return options{error: err};
				}
				opt.filter_code = reg;
			}
		}
	}
	if (opt.filter_code==nil) {
		return options{};
	}
	return opt;
}

func usage() {
	fmt.Println("Usage: ./snake -u http://crawl.website -fs js,css -fc 404,500,403");
	fmt.Println("-fs|--filter-string : Filters url's containing this string.");
	fmt.Println("-fc|--filter-code   : Filters url's that return these status codes.");
	fmt.Println("-u |--url           : Url to use espaider against.")
}