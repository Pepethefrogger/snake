package main

import (
	"fmt"
	"os"
	"regexp"
	"snake/argparse"
	"snake/crawler"
	"snake/formatter"
)

func main(){
	args := os.Args;
	options := argparse.ParseArgs(args);
	if options.Error!=nil||options.FilterCode==nil {
		argparse.Usage();
		return;
	}
	url := options.Url;
	r, err := regexp.Compile("https?://(.*?)/.*$");
	if err!=nil {
		fmt.Println("Invalid url.")
		return;
	}
	domain := r.FindStringSubmatch(url)[1];
	fmt.Printf("[*] Searching in domain %s.\n",domain)

	channels := crawler.NewChannels();
	go func() {
		channels.ToRequestChannel<-url;
	}()
	go crawler.RequestHandler(channels,options);
	validUrls := crawler.ResponseParser(url,domain,channels,options);
	formatter.Format(validUrls);
	channels.Close();
}
