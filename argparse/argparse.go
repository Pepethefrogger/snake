package argparse

import (
	"fmt"
	"regexp"
	"strings"
)

type Options struct {
	Url string;
	FilterCode *regexp.Regexp;
	FilterString *regexp.Regexp;
	Verbose bool;
	StatusCode bool;
	Pretty bool;
	Error error;
}

func Usage() {
	fmt.Println("Usage: ./snake -u http://crawl.website -fs js,css -fc 404,500,403");
	fmt.Println("-fs|--filter-string : Filters url's containing this string.");
	fmt.Println("-fc|--filter-code   : Filters url's that return these status codes.");
	fmt.Println("-u |--url           : Url to use espaider against.")
	fmt.Println("-v |--verbose       : Outputs directories as they are found. Only makes sense with pretty.")
	fmt.Println("-p |--pretty        : Output is sepparated by directory. Less greppable though.")
	fmt.Println("-sc|--status-code   : Prints status code alongside with found directories.")
}

func ParseArgs(args []string) Options {
	var opt Options;
	for i,flag := range(args) {
		if (i!=len(args)) {
			switch string(flag) {
			case "-u","--url":
				opt.Url = args[i+1];
			case "-fs","--filter-string":
				array := strings.Split(args[i+1], ",");
				filter := strings.Join(array,"|");
				reg,err := regexp.Compile(filter);
				if err!=nil {
					return Options{Error: err};
				}
				opt.FilterString = reg;
			case "-fc","--filter-code":
				array := strings.Split(args[i+1],",");
				filter := strings.Join(array,"|");
				reg,err := regexp.Compile(filter);
				if err!=nil {
					return Options{Error: err};
				}
				opt.FilterCode = reg;
			case "-v","--verbose":
				opt.Verbose = true;
			case "-sc","--status-code":
				opt.StatusCode = true;
			case "-p","--pretty":
				opt.Pretty = true;
			}
		}
	}
	if (opt.FilterCode==nil) {
		return Options{};
	}
	return opt;
}
