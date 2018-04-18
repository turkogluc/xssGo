package main

import (
	"fmt"
	"net/url"
	"time"
	"log"
	"flag"
	"strings"
	"os"
	. "xssGo/source"
)

func init() {
	InitEntities()
	// Init logging
	var err error
	LogFile, err = os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(LogFile)
	InitLogger(LogFile, LogFile, LogFile, true)
}

// TODO : bad urls için map yap, kontrolu daha kolay, ok := map[sth]

// FIXME: There is a bug for crawling, urls need to be parsed and analyzed.
//5 https://www.numbeo.com/cost-of-living/country_result.jsp?country=Croatia
//6 https://www.numbeo.com/cost-of-living/country_result.jsp?country=Gambia

func main() {

	//LogDebug("debug message hello 1")
	//LogInfo("info message")
	//LogError("error message")

	flagURL := flag.String("url","http://localhost/dvwa/","Target Url to be scanned")
	flagLevel := flag.Int("level",2,"Scan Depth Level")
	flagLogin := flag.String("loginPage","","Authenticate via login page.")
	flagCookie := flag.String("cookieFile","","Authenticate via cookies. Give path of cookie file as json")
	flagPayload := flag.String("payloadFile","payloads.txt","Set Payloads from file")
	flagHeadless := flag.Bool("headless",false,"Browser can run as headless")
	flagBlackList := flag.String("blackList","","Forbid some links to be visited via giving comma seperated list")
	flagHelp  := flag.Bool("help",false,"Print usage")

	flag.Parse()
	if *flagHelp {
		PrintUsage()
		os.Exit(1)
	}

	UrlSTR = *flagURL
	//UrlSTR = "https://www.numbeo.com"
	UrlParsed, _ = url.Parse(UrlSTR)
	Host = UrlParsed.Host

	// add the first target to list
	TargetURLs[UrlSTR] = Empty{}

	Level = *flagLevel

	ReadPayloads(*flagPayload)

	if *flagBlackList != "" {
		blackList := *flagBlackList
		arr := strings.Split(blackList,",")
		BadUrls = append(BadUrls,arr...)
	}
	browserArgs := []string{"--disable-xss-auditor"}
	if *flagHeadless {
		browserArgs = append(browserArgs, "--headless")
		browserArgs = append(browserArgs, "--disable-gpu")
	}
	// start chrome driver
	CD.InitDriver(browserArgs)
	fmt.Println("")
	defer CD.StopDriver()

	// open the initial url in both browsers
	err := Bow.Open(UrlSTR)
	if err != nil {
		LogError(err)
	}
	if err := CD.WebDriver.Get(UrlSTR); err != nil {
		LogError(err)
	}

	if *flagLogin != "" {
		LoginURL = *flagLogin
		LoginToBow(LoginURL)
		CD.LoginToChromeAuto(LoginURL,LoginInformation)

	}else if *flagCookie != "" {
		// reading from file
		cs,err := ReadCookiesFromFile(*flagCookie)
		if err != nil {
			LogError(err)
			panic(err)
		}
		// setting cookies to bow browser
		SetCookiesToBow(cs)
		// setting cookies to chrome
		cookies := ConvertCookiesToGolang(cs)
		CD.SetCookiesToChrome(*flagURL,cookies)
	}

	StartTime = time.Now()
	defer finishTime()

	fmt.Println("Page traversal operation is started:")
	log.Println("Page traversal operation is started:")
	//crawling
	CrawlURL(UrlSTR)
	if Level > 1 {
		for i := 2; i <= Level; i++ {
			for tempURL, _ := range TargetURLs {
				CrawlURL(tempURL)
			}
		}
	}
	time.Sleep(250*time.Millisecond)

	// List of target Urls
	i := 1
	fmt.Println("List of Target Urls:")
	log.Println("List of Target Urls:")
	for u, _ := range TargetURLs {
		fmt.Println(i, u)
		i++
	}
	time.Sleep(250*time.Millisecond)

	// Add one extra url to list
	//targetURLs["http://localhost/dvwa/vulnerabilities/xss_d/?default=English"] = empty{}

	// Query Parameters
	log.Println("Query parameters are going to be tested")
	fmt.Println("Query parameters are going to be tested")
	for u := range TargetURLs {
		ControlQueryParameters(u)
	}
	time.Sleep(250*time.Millisecond)

	// Form input variables
	log.Println("Form inputs are going to be tested")
	fmt.Println("Form inputs are going to be tested")
	for u := range TargetURLs {
		ControlFormInputs(u)
	}

	// Output, Vulnerable URLs
	fmt.Println("")
	fmt.Println("vulnerable urls:")
	j := 1
	for i, v := range VulnerableURLs {
		fmt.Println(j, ":", i, " payload:", v.Payload)
		j++
	}




	// ###### tests ######

	//cd.trySelection("http://localhost/dvwa/security.php")

}

func finishTime() {
	fmt.Println("Time Elapsed:", time.Since(StartTime))
}


// TODO ## List ##

// TODO Better Logging and stdout mechanizm (seperated)

// TODO Implement Go routines and sync , maybe workers
//		To increase the speed run on all cores of CPU


// TODO click all buttons for the possibility of having another target link

// FIXME Implement cookie insertion for bow

// DONE Command line arguments with flag
//		Welcome message
//		Usage
// 		Get Black-list (jpg,png,pdf..)
//		Payloads from file