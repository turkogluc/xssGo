package main

import (
	"fmt"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	//"github.com/tebeka/selenium"
	"time"
)

type empty struct {
	payload string
}

var targetURLs map[string]empty
var vulnerableURLs map[string]empty
var loginInformation map[string]string
var urlSTR, host, loginURL string
var urlParsed *url.URL
var bow *browser.Browser
var badUrls []string
var jar *cookiejar.Jar
var cookies []*http.Cookie
var payloads []string
var level int
var cd *chromeDriver
var startTime time.Time

func init() {
	startTime = time.Now()
	p1 := []string{"<script>alert('XSS');</script>", "<BODY ONLOAD=alert('XSS')>", "\";alert(1);//", "';alert(1);//"}
	p2 := []string{"><script>alert(0)</script>", "\" onfocus=\"alert(1);", "javascript:alert(1)", "\"><img src=\"x:x\" onerror=\"alert(0)\">"}

	payloads = append(payloads, p1...)
	payloads = append(payloads, p2...)

	bow = surf.NewBrowser()
	cd = &chromeDriver{}

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

	//urlSTR = "http://192.168.56.101"
	urlSTR = "http://localhost/dvwa/"
	//urlSTR = "http://localhost/bwapp/"
	//urlSTR = "https://www.seslisozluk.net/"
	//urlSTR = "https://xss-game.appspot.com"
	urlParsed, _ = url.Parse(urlSTR)

	host = urlParsed.Host
	level = 3

	badUrls = append(badUrls, []string{"%C3%A7%C4%B1k%C4%B1%C5%9F", "logout", ".png", ".jpg", ".jpeg", ".mp3", ".mp4", ".avi", ".gif", ".svg", "setup", "csrf"}...)
	badUrls = append(badUrls, []string{"reset", "user_extra", "password_change"}...)

	//fmt.Println(time.Now().Unix())
}

// TODO ## List ##

// TODO Accept Cookies from JSON File
// 		Inserting to both modules

// TODO Command line arguments with flag
//		Welcome message
//		Usage
// 		Get Black-list (jpg,png,pdf..)
//		Payloads from file

// TODO Better Logging and stdout mechanizm (seperated)
//

// TODO Implement Go routines and sync , maybe workers
//		To increase the speed run on all cores of CPU





func main() {

	defer finishTime()
	//err := bow.Open(urlSTR)
	//if err != nil {
	//	log.Println("PANIC:", err)
	//}
	//
	//cs,err := readCookiesFromFile("cookies.json")
	//if err != nil{
	//	panic("couldnt read cookies from file")
	//}
	//SetCookies(cs)
	//
	//x := bow.SiteCookies()
	//log.Println(x)
	//err = bow.Open(urlSTR)
	//if err != nil {
	//	log.Println("PANIC:", err)
	//}
	//
	//s,_:= bow.Dom().Html()
	//fmt.Println(s)
	//
	//// add the first target to list
	//targetURLs[urlSTR] = empty{}
	//
	////LoginByCredentials(urlSTR,"admin","password")
	//Login(urlSTR)
	//
	////crawling
	//crawlURL(urlSTR)
	//
	//if level > 1 {
	//	for i := 2; i <= level; i++ {
	//		for tempURL, _ := range targetURLs {
	//			crawlURL(tempURL)
	//		}
	//	}
	//}
	////
	//i := 1
	//for u, _ := range targetURLs {
	//	fmt.Println(i, u)
	//	i++
	//}





	//cd.initDriver()
	//defer cd.stopDriver()
	////cd.loginAuto(urlSTR)		   // let the chrome-driver log in with credentials taken from user
	//
	//
	//
	//cd.webDriver.Get("http://localhost/dvwa/vulnerabilities/xss_d/?default=English")
	//
	//
	//cs,err := readCookiesFromFile("cookies.json")
	//if err != nil{
	//	panic("couldnt read cookies from file")
	//}
	//
	//goCookie := convertCookiesToGolang(cs)
	//cd.setCookiesToChrome(goCookie)
	//
	//cd.webDriver.Get("http://localhost/dvwa/vulnerabilities/xss_d/?default=English")
	//
	//
	//targetURLs["http://localhost/dvwa/vulnerabilities/xss_d/?default=English"] = empty{}
	//
	//log.Println("Query parameters are going to be tested")
	//for u := range targetURLs {
	//	controlQueryParameters(u)
	//}
	//
	//log.Println("Form inputs are going to be tested")
	//for u := range targetURLs {
	//	controlFormInputs(u)
	//}

	// ###### tests ######

	//controlFormInputs("http://localhost/dvwa/vulnerabilities/xss_d/")

	//cd.trySelection("http://localhost/dvwa/security.php")

	// test
	//controlFormInputs("http://localhost/dvwa/vulnerabilities/xss_s/")
	// controlFormInputs("http://localhost/dvwa/vulnerabilities/xss_s/")

	//cd.formTest("http://192.168.56.101/xss/example8.php",payloads)

	//controlQueryParameters("http://192.168.56.101/xss/example9.php#hacker")

	// test
	//cd.pageTest("http://localhost/dvwa/vulnerabilities/xss_d/?default=%3Cscript%3Ealert(1);%3C/script%3E")
	//time.Sleep(5*time.Second)
	//cd.pageTest("http://192.168.56.101/xss/example3.php?name=%3CBODY+ONLOAD%3Dalert(%27cemal%27)%3E")	//xss
	//cd.pageTest("http://192.168.56.101/xss/example1.php?name=%3Cbody%20onload=alert(1)%3E")		// xss
	//cd.pageTest("http://192.168.56.101/codeexec/example1.php?name=hacker")
	//cd.pageTest("http://192.168.56.101/xss/example3.php?name=%3CBODY+ONLOAD%3Dalert(%27cemal%27)%3E")	//xss
	//cd.pageTest("http://192.168.56.101/xss/example1.php?name=%3Cbody%20onload=alert(1)%3E")		// xss
	//cd.pageTest("http://192.168.56.101/codeexec/example1.php?name=hacker")
	//cd.pageTest("http://192.168.56.101/xss/example2.php?name=%3CBODY+ONLOAD%3Dalert(%27XSS%27)%3E") //xss

	// test
	//cd.pageTest("http://localhost/dvwa/vulnerabilities/xss_r/?name=%3Cscript%3Ealert(1)%3B%3C%2Fscript%3E")
	//cd.pageTest("http://localhost/dvwa/vulnerabilities/weak_id/")
	//cd.pageTest("http://localhost/dvwa/vulnerabilities/brute/")
	//cd.pageTest("http://localhost/dvwa/instructions.php")
	//cd.pageTest("http://localhost/dvwa/vulnerabilities/xss_r/?name=<BODY+ONLOAD%3Dalert('XSS')>")

	//
	////for u := range targetURLs{
	////	controlFormInputs(u,bow.SiteCookies())
	////}
	//

	//Output
	fmt.Println(" ")
	fmt.Println("vulnerable urls:")
	j := 1
	for i, v := range vulnerableURLs {
		fmt.Println(j, ":", i, " payload:", v.payload)
		j++
	}
}





func finishTime() {
	fmt.Println("Time Elapsed:", time.Since(startTime))
}
