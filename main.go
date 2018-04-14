package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	//"github.com/tebeka/selenium"
	"log"
	"time"
)

type empty struct{
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

func init(){
	startTime = time.Now()
	p1 := []string{"<script>alert('XSS');</script>","<BODY ONLOAD=alert('XSS')>","\";alert(1);//","';alert(1);//"}
	p2 := []string{"><script>alert(0)</script>", "\" onfocus=\"alert(1);", "javascript:alert(1)","\"><img src=\"x:x\" onerror=\"alert(0)\">"}

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

	badUrls = append(badUrls, []string{"%C3%A7%C4%B1k%C4%B1%C5%9F", "logout", ".png", ".jpg", ".jpeg", ".mp3", ".mp4", ".avi", ".gif", ".svg","setup","csrf"}...)
	badUrls = append(badUrls,[]string{"reset","user_extra","password_change"}...)
}

func main() {
	defer finishTime()
	err := bow.Open(urlSTR)
	if err != nil {
		log.Println("PANIC:",err)
	}

	// add the first target to list
	targetURLs[urlSTR] = empty{}

	//LoginByCredentials(urlSTR,"admin","password")
	Login(urlSTR)
	//cookies := bow.CookieJar().Cookies(urlParsed)
	//cookies := bow.SiteCookies()


	//c := bow.SiteCookies()  // get cookies
	//SetCookie()

	//crawling
	crawlURL(urlSTR)

	if level > 1 {
		for i := 2; i <= level; i++ {
			for tempURL, _ := range targetURLs {
				crawlURL(tempURL)
			}
		}
	}
	//
	i := 1
	for u, _ := range targetURLs {
		fmt.Println(i, u)
		i++
	}

	cd.initDriver()
	defer cd.stopDriver()
	//cd.login(urlSTR)
	//cd.loginW(urlSTR)
	cd.loginAuto(urlSTR)

	log.Println("Query parameters are going to be tested")
	for u := range targetURLs{
		controlQueryParameters(u)
	}

	log.Println("Form inputs are going to be tested")
	for u := range targetURLs{
		controlFormInputs(u)
	}

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
	j:=1
	for i,v := range vulnerableURLs{
		fmt.Println(j,":",i," payload:",v.payload)
		j++
	}


}

func controlFormInputs(u string){

	isVul,p := cd.formTest(u,payloads)
	if isVul {
		if _,contains := vulnerableURLs[u]; !contains{
			vulnerableURLs[u]=empty{payload:p}
		}
	}

}

func controlQueryParameters(u string){

	// convert string to url.URL
	originalURL,_ := url.Parse(u)
	modifiedURL := originalURL


	// testing for DOM XSS

	// browser does not allow , encodes url
	// TODO: find a way to bypass

	if modifiedURL.Fragment != "" {
		for _, payload := range payloads {

			modifiedURL.Fragment = payload

			if cd.pageTest(modifiedURL.String()) {
				if _, contains := vulnerableURLs[modifiedURL.String()]; !contains {
					vulnerableURLs[modifiedURL.String()] = empty{}
				}
				break
			}

		}

	}


	// get query parameters and values as map[string][]string
	q := originalURL.Query()

	// checking query parameters Ex: username=xx&email=yy
	// only one parameter is changed at once
	if len(q)>0{
		for parameter := range q{
			for _,payload := range payloads{

				modifiedURL = originalURL
				q.Set(parameter,payload)
				modifiedURL.RawQuery = q.Encode()



				if cd.pageTest(modifiedURL.String()){
					if _,contains := vulnerableURLs[modifiedURL.String()]; !contains{
								vulnerableURLs[modifiedURL.String()]=empty{}
					}
					break
				}



				//bow.Open(modifiedURL.String())

				//if BrowserQueryTest(modifiedURL.String(),payload){
				//	if _,contains := vulnerableURLs[modifiedURL.String()]; !contains{
				//		vulnerableURLs[modifiedURL.String()]=empty{}
				//	}
				//	break
				//}
			}
		}
	}

}


func crawlURL(u string) {

	// Crawler travers as BFS tree

	err := bow.Open(u)
	if err != nil {
		panic(err)
	}

	// return all links in the current page
	links := bow.Links()

	// TODO: dont get non-html files (jpg,png,pdf..)
	//i:=0
	for _, link := range links {

		//fmt.Println(link)

		_, ifContains := targetURLs[link.Url().String()]
		whitelisted := true
		for _, b := range badUrls {
			if strings.Contains(link.Url().String(), b) {
				whitelisted = false
				break
			}

		}

		// search for same domain
		// do not allow unwanted urls like: png,jpg,pdf etc.
		// add if the URL not already in slice
		if strings.Compare(host, link.Url().Host) == 0 && !ifContains && whitelisted {

			//fmt.Println(link)
			targetURLs[link.Url().String()] = empty{}
			//i++
			//fmt.Println(i, link.Url())

		}

	}
}

func LoginByCredentials(loginURL string, user string, pass string) {

	fm, _ := bow.Form("form")
	fm.Input("username", "admin")
	fm.Input("password", "password")
	//fmt.Println(bow.Dom().Html())
	fm.Submit()
	//fmt.Println(bow.Dom().Html())


}

func Login(loginURL string) {

	err := bow.Open(urlSTR)
	if err != nil {
		panic(err)
	}

	allForms := bow.Forms()
	var text string

	for _, fm := range allForms {
		if fm != nil {
			fmt.Println("Form found.. : ")
			fm.Dom().Find("input").Each(func(i int, s *goquery.Selection) {
				// For each item found, get the band and title
				if inputName, ok := s.Attr("name"); ok {
					if inputType, ok2 := s.Attr("type"); ok2 {
						if inputType == "text" || inputType == "password" {
							fmt.Print("Enter ", inputName, ":")
							fmt.Scanln(&text)
							fm.Input(inputName, text)
							loginInformation[inputName] = text
						}
					}
				}

			})

			fmt.Println(fm.GetFields())

			err = fm.Submit()
			if err != nil {
				panic(err)
			}

			fmt.Println("Succesfully loged in..")
			fmt.Println("Cookies are set as:")

			for _, c := range bow.CookieJar().Cookies(urlParsed) {
				res, _ := json.Marshal(c)
				fmt.Println(string(res))
			}
		}
	}
}

func SetCookie() {

	cookie := &http.Cookie{
		Domain: "localhost",
		Name:   "PHPSESSID",
		Value:  "j5mdm88v2ougl2kjrrinsfcsd1",
	}
	cookies = append(cookies, cookie)

	cookie = &http.Cookie{
		Domain: "localhost",
		Name:   "security",
		Value:  "impossible",
	}
	cookies = append(cookies, cookie)

	jar.SetCookies(urlParsed, cookies)
	bow.SetCookieJar(jar)

	fmt.Println("afterr:")
	for _, c := range bow.CookieJar().Cookies(urlParsed) {
		res, _ := json.Marshal(c)
		fmt.Println(string(res))
	}

}

func finishTime(){
	fmt.Println("Time Elapsed:",time.Since(startTime))
}
