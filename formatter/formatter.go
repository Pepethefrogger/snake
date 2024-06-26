package formatter

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type webpage struct {
	directory string;
	statusCode string;
	children []webpage;
}

func Format(array [][]string) {
	r,err := regexp.Compile("https?://(.*)$");
	if err!=nil {
		panic("Error");
	}

	var lines []string;
	var statusCodes []string;
	for _,lineStatusPair := range(array) {
		url := r.FindSubmatch([]byte(lineStatusPair[0]))[1]
		trimmed := strings.TrimSpace(string(url));
		lines = append(lines,string(trimmed));
		statusCodes = append(statusCodes, lineStatusPair[1]);
	}

	slices.Sort(lines);

	wp := webpage{directory: "",children: nil}
	for _,line := range(lines) {
		array := strings.Split(line,"/");
		if array[len(array)-1]=="" {
			array = array[:len(array)-1]
		}
		arrayToWebpage(array,statusCodes,&wp)
	}
	
	printWebpage(&wp,"")
}

func arrayToWebpage(slice []string,statusCodes []string,wp *webpage) {
	if len(slice)==0 {
		return
	}
	if wp.children==nil {
		wp.children=[]webpage{{directory: slice[0],statusCode: statusCodes[0],children: nil}};
	} else {
		found := false;
		for index := range wp.children {
			if wp.children[index].directory==slice[0] {
				arrayToWebpage(slice[1:],statusCodes[1:],&wp.children[index]);
				found = true;
			}
		}

		if !found {
			wp.children = append(wp.children, webpage{directory: slice[0],statusCode: statusCodes[0],children: nil});
		}
	}
}

func printWebpage(wp *webpage,arrows string) {
	for index := range(wp.children) {
		node := "";
		if (wp.children[index].children!=nil) {
			node="──┐ "
		}
		
		var length int;
		if (wp.children[index].statusCode!="") {
			fmt.Println(arrows+wp.children[index].directory+" ("+wp.children[index].statusCode+")"+node);
			length = len(wp.children[index].directory)+len(wp.children[index].statusCode)+3;
		} else {
			fmt.Println(arrows+wp.children[index].directory+node);
			length = len(wp.children[index].directory);
		}

		spaces := strings.Replace(arrows,"─"," ",-1)
		spaces = strings.Replace(spaces,"├","│",1)
		newspaces := spaces+strings.Repeat(" ",length+2)+"├───";
		printWebpage(&wp.children[index],newspaces);
	}
}