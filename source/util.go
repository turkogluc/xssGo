package source

import (
	"os"
	"bufio"
	"fmt"
)

func ReadPayloads(path string){
	fileHandle, err := os.Open(path)
	defer fileHandle.Close()
	if err != nil {
		LogError(err)
		panic(err)
	}
	fileScanner := bufio.NewScanner(fileHandle)

	for fileScanner.Scan() {
		text := fileScanner.Text()
		Payloads = append(Payloads,text)
	}

}



func PrintUsage(){
	Usage = `Usage : ./xssGo	[options[=value]]

Options:
	  --Help
			Display this page

	  --BlackList
			Forbid some links to be visited via giving comma seperated list

	  --CookieFile
			Authenticate via Cookies. Give path of cookie file as json

	  --Headless
			Browser can to run as headless

	  --Level int
			Scan Depth Level (default 3)

	  --LoginPage string
			Authenticate via Login page. (default "http://localhost/dvwa/Login.php")

	  --PayloadFile string
			Set Payloads from file (default "Payload.txt")

	  --URL string
			Target Url to be scanned (default "http://localhost/dvwa/")`

	Description = `		xssGo - Automated Xss Scanner

"xssGo" is designed to automate scanning cross site scripting vulnerabilities.
Most important functionality of the program is that because of testing with a
real browser, not producing wrong results, aka false positives.

It is possible to scan with or without authentication. Authentication can be
supplied via Login credentials or Cookies.

Phases of program:
1) Scanning
	xssGo scans the webpage and creates kind of site-map. Scanning phase is
	consist of parsing the html page, excluding all links refers to same host.
	Scanning is done as Breadt-First-Search (BFS) algorithm. Leves is determined
	by user. Output of the phase is a list consist of target URLs.

2) Testing Query Parameters.
	All target URLs are examined and tested if there is paramater that can be
	tested. All parameters like /page.php?name=user&date=10#last are replaced
 	with XSS Payloads and new link sent to real browser to be tested.
	Real Browser runs the URL and examines the response.If there is alert
	URL is included to vulnerable url list.

3) Testing Form Inputs.
	All target URLs are examined and the forms in page are scanned.All input
	tags in every form is filled with XSS Payload and sent. Result is
	observed by real browser.`


	fmt.Print(Description)
	fmt.Println(Usage)
}
