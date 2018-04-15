package main

import (
	"fmt"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
	"log"
	"flag"
	"strings"
	"os"
)

type empty struct {
	payload string
}

var targetURLs map[string]empty
var vulnerableURLs map[string]empty
var loginInformation map[string]string
var urlSTR, host, loginURL string
var usage,description string
var urlParsed *url.URL
var bow *browser.Browser
var badUrls []string
var jar *cookiejar.Jar
var cookies []*http.Cookie
var goCookies []*http.Cookie
var payloads []string
var level int
var cd *ChromeDriver
var startTime time.Time

func init() {

	payloads = append(payloads, "<script>alert('XSS');</script>")

	bow = surf.NewBrowser()
	cd = &ChromeDriver{}

	//bow.SetAttribute(browser.SendReferer, false)
	//bow.SetAttribute(browser.MetaRefreshHandling, false)
	bow.SetAttribute(browser.FollowRedirects, true)
	bow.SetUserAgent(agent.Firefox())

	// we create a cookie jar that will provide being statefull
	// after attaching a cookie jar to http.Client, http.client will
	// add cookies and update in every request and response
	jar, _ = cookiejar.New(nil)
	bow.SetCookieJar(jar)

	targetURLs = map[string]empty{}
	vulnerableURLs = map[string]empty{}
	loginInformation = make(map[string]string)

	badUrls = append(badUrls, []string{"%C3%A7%C4%B1k%C4%B1%C5%9F", "logout", ".png", ".jpg", ".jpeg", ".mp3", ".mp4", ".avi", ".gif", ".svg", "setup", "csrf"}...)
	badUrls = append(badUrls, []string{"reset", "user_extra", "password_change"}...)

}

func main() {

	flagURL := flag.String("URL","http://localhost/dvwa/","Target Url to be scanned")
	flagLevel := flag.Int("Level",3,"Scan Depth Level")
	flagLogin := flag.String("LoginPage","http://localhost/dvwa/login.php","Authenticate via login page.")
	flagCookie := flag.String("CookieFile","","Authenticate via cookies. Give path of cookie file as json")
	flagPayload := flag.String("PayloadFile","payloads.txt","Set Payloads from file")
	flagHeadless := flag.Bool("Headless",false,"Browser can run as headless")
	flagBlackList := flag.String("BlackList","","Forbid some links to be visited via giving comma seperated list")
	flagHelp  := flag.Bool("Help",false,"Print usage")

	flag.Parse()
	if *flagHelp {
		printUsage()
		os.Exit(1)
	}

	urlSTR = *flagURL
	urlParsed, _ = url.Parse(urlSTR)
	host = urlParsed.Host

	// add the first target to list
	targetURLs[urlSTR] = empty{}

	level = *flagLevel

	readPayloads(*flagPayload)

	if *flagBlackList != "" {
		blackList := *flagBlackList
		arr := strings.Split(blackList,",")
		badUrls = append(badUrls,arr...)
	}
	browserArgs := []string{"--disable-xss-auditor"}
	if *flagHeadless {
		browserArgs = append(browserArgs, "--headless")
		browserArgs = append(browserArgs, "--disable-gpu")
	}
	// start chrome driver
	cd.initDriver(browserArgs)
	fmt.Println("")
	defer cd.stopDriver()

	// open the initial url in both browsers
	err := bow.Open(urlSTR)
	if err != nil {
		log.Println("PANIC:", err)
	}
	if err := cd.webDriver.Get(urlSTR); err != nil {
		log.Println("PANIC:", err)
	}

	if *flagLogin != "" {
		loginURL = *flagLogin
		LoginToBow(loginURL)
		cd.loginToChromeAuto(loginURL,loginInformation)

	}else if *flagCookie != "" {
		// reading from file
		cs,err :=readCookiesFromFile(*flagCookie)
		if err != nil {
			panic(err)
		}
		// setting cookies to bow browser
		SetCookiesToBow(cs)
		// setting cookies to chrome
		cookies := convertCookiesToGolang(cs)
		cd.setCookiesToChrome(*flagURL,cookies)
	}

	startTime = time.Now()
	defer finishTime()

	//crawling
	crawlURL(urlSTR)
	if level > 1 {
		for i := 2; i <= level; i++ {
			for tempURL, _ := range targetURLs {
				crawlURL(tempURL)
			}
		}
	}
	time.Sleep(250*time.Millisecond)

	// List of target Urls
	i := 1
	log.Println("List of Target Urls:")
	for u, _ := range targetURLs {
		fmt.Println(i, u)
		i++
	}
	time.Sleep(250*time.Millisecond)

	// Add one extra url to list
	targetURLs["http://localhost/dvwa/vulnerabilities/xss_d/?default=English"] = empty{}

	// Query Parameters
	log.Println("Query parameters are going to be tested")
	for u := range targetURLs {
		controlQueryParameters(u)
	}
	time.Sleep(250*time.Millisecond)

	// Form input variables
	log.Println("Form inputs are going to be tested")
	for u := range targetURLs {
		controlFormInputs(u)
	}

	// Output, Vulnerable URLs
	fmt.Println("")
	fmt.Println("vulnerable urls:")
	j := 1
	for i, v := range vulnerableURLs {
		fmt.Println(j, ":", i, " payload:", v.payload)
		j++
	}




	// ###### tests ######

	//cd.trySelection("http://localhost/dvwa/security.php")

}

func finishTime() {
	fmt.Println("Time Elapsed:", time.Since(startTime))
}


// TODO ## List ##

// TODO Command line arguments with flag
//		Welcome message
//		Usage
// 		Get Black-list (jpg,png,pdf..)
//		Payloads from file

// TODO Better Logging and stdout mechanizm (seperated)
//

// TODO Implement Go routines and sync , maybe workers
//		To increase the speed run on all cores of CPU


// FIXME Implement cookie insertion for bow