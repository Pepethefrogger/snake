package argparse

import (
	"fmt"
	"os"
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
	Banner bool;
	Error error;
}

func Banner() {
fmt.Fprintf(os.Stderr,"                                                                           L                                   \n");
fmt.Fprintf(os.Stderr,"____   ____.__                                                     .___    J   .-\"\"\"\"-.               J        \n");
fmt.Fprintf(os.Stderr,"\\   \\ /   /|__|_____   ___________  ____________   ____   ____   __| _/     \\ /        \\   __    /    F        \n");
fmt.Fprintf(os.Stderr," \\   Y   / |  \\____ \\_/ __ \\_  __ \\/  ___/\\____ \\_/ __ \\_/ __ \\ / __ |  \\    (|)(|)_   .-'\".'  .'    /         \n");
fmt.Fprintf(os.Stderr,"  \\     /  |  |  |_> >  ___/|  | \\/\\___ \\ |  |_> >  ___/\\  ___// /_/ |   \\    \\   /_>-'  .<_.-'     /          \n");
fmt.Fprintf(os.Stderr,"   \\___/   |__|   __/ \\___  >__|  /____  >|   __/ \\___  >\\___  >____ |    `.   `-'     .'         .'           \n");
fmt.Fprintf(os.Stderr,"              |__|        \\/           \\/ |__|        \\/     \\/     \\/      `--.|___.-'`._    _.-'             \n");
fmt.Fprintf(os.Stderr,"                                                                                ^         \"\"\"\"                 \n");
}

func Usage() {
	fmt.Println("Usage: ./snake -u http://crawl.website -fs js,css -fc 404,500,403");
	fmt.Println("-fs|--filter-string : Filters url's containing this string.");
	fmt.Println("-fc|--filter-code   : Filters url's that return these status codes.");
	fmt.Println("-u |--url           : Url to use espaider against.")
	fmt.Println("-v |--verbose       : Outputs directories as they are found. Only makes sense with pretty.")
	fmt.Println("-p |--pretty        : Output is sepparated by directory. Less greppable though.")
	fmt.Println("-sc|--status-code   : Prints status code alongside with found directories.")
	fmt.Println("-b |--hide-banner   : Hides banner.")
}

func ParseArgs(args []string) Options {
	var opt Options;
	opt.Banner = true;
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
			case "-b","--hide-banner":
				opt.Banner = false;
			}
		}
	}
	if (!opt.Pretty) {
		opt.Verbose = true;
	}
	if (opt.FilterCode==nil) {
		return Options{};
	}
	return opt;
}
