package crawler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"snake/argparse"
	"time"
)

type Channels struct {
	ResponseChannel chan []string;
	ToRequestChannel chan string;
	DoneRequesting chan struct{};
}

func NewChannels() Channels {
	var channels Channels;
	channels.ToRequestChannel = make(chan string);
	channels.ResponseChannel = make(chan []string);
	channels.DoneRequesting = make(chan struct{});
	return channels 
}

func (channels Channels) Close() {
	close(channels.DoneRequesting);
	close(channels.ToRequestChannel);
}

func request(url string,channels Channels,options argparse.Options) { //Makes one request
	defer func() {
		channels.DoneRequesting<-struct{}{};
	}();
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(url);
	if err!=nil {
		return;
	}
	if (!options.FilterCode.Match([]byte(resp.Status))) {
		body, err := io.ReadAll(resp.Body);
		if err!=nil {
			body=[]byte("");
		}
		channels.ResponseChannel<-[]string{url,string(body),resp.Status}
		if options.Verbose&&options.Pretty {
			fmt.Fprint(os.Stderr,url+"  ("+resp.Status+")\n");
		} else if options.Verbose {
			fmt.Println(url+"   ("+resp.Status+")");
		}
	}
}

func RequestHandler(channels Channels,options argparse.Options) { //Controls launching and ending requests
	var requestedcount int;
	var parsedcount int;
	closed := false;
	for {
		select {
		case <-channels.DoneRequesting:
			requestedcount++;
		case url := <-channels.ToRequestChannel:
			parsedcount++;
			go request(url,channels,options);
		default:
			if (requestedcount>1&&requestedcount>=parsedcount) {
				lastcount := requestedcount;
				time.Sleep(time.Millisecond*100);
				if ((requestedcount == lastcount)&&!closed) {
					closed=true;
					close(channels.ResponseChannel);
				}
			}
		}
	}
}


func ResponseParser(url string,domain string,channels Channels,options argparse.Options) [][]string {
	var visited []string;
	var validUrls [][]string;
	visited = append(visited, url)
	for array := range(channels.ResponseChannel) {
		url := array[0];
		if options.Pretty&&options.StatusCode {
			sc := array[2];
			validUrls = append(validUrls, []string{url,sc});
		} else {
			validUrls = append(validUrls, []string{url,""});
		}
		parse := array[1];
		r1, _ := regexp.Compile(fmt.Sprintf("(https?://.*?%s.*?)[\" ']",domain));
		r2, _ := regexp.Compile("href=\"(/?[^(http)].*?)[\" ']|src=\"(/?[^(http)].*?)[\" ']");
		match1 := r1.FindAllStringSubmatch(parse,-1);
		match2 := r2.FindAllSubmatch([]byte(parse),-1);
		var is_filtered bool;
		for _,m := range(match1) {
			option := m[1];
			if (options.FilterString==nil) {
				is_filtered = false;
			} else {
				is_filtered = options.FilterString.Match([]byte(option));
			}
			if !is_filtered {
				if !slices.Contains(visited,option) {
					visited = append(visited, option);
					channels.ToRequestChannel<-option;
				}
			}
		}
		for _,m := range(match2) {
			for i,s := range(m) {
				if (i==1||i==2) {
					if (options.FilterString==nil) {
						is_filtered = false;
					} else {
						is_filtered = options.FilterString.Match([]byte(m[1]));
					}
					if !is_filtered {
						str := url + string(s);
						if !slices.Contains(visited,str) {
							visited = append(visited, str);
							channels.ToRequestChannel<-str;
						}
					}
				}
			}
		}
	}
	return validUrls;
}