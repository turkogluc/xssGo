package source

import (
	"net/url"
	"net/http/cookiejar"
	"net/http"
	"github.com/headzoo/surf/browser"
	"time"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"os"
)

type Empty struct {
	Payload string
}

var TargetURLsORIGINAL map[string]Empty
var TargetURLsALTERED map[string]Empty
var VulnerableURLs map[string]Empty
var LoginInformation map[string]string
var UrlSTR, Host, LoginURL string
var Usage, Description string
var UrlParsed *url.URL
var Bow *browser.Browser
var BadUrls []string
var Jar *cookiejar.Jar
var Cookies []*http.Cookie
var GoCookies []*http.Cookie
var Payloads []string
var Level int
var CD *ChromeDriver
var StartTime time.Time
var LogFile *os.File

func InitEntities(){
	Payloads = append(Payloads, "<script>alert('XSS');</script>")

	Bow = surf.NewBrowser()
	CD = &ChromeDriver{}

	//bow.SetAttribute(browser.SendReferer, false)
	//bow.SetAttribute(browser.MetaRefreshHandling, false)
	Bow.SetAttribute(browser.FollowRedirects, true)
	Bow.SetUserAgent(agent.Firefox())

	// we create a cookie jar that will provide being statefull
	// after attaching a cookie jar to http.Client, http.client will
	// add cookies and update in every request and response
	Jar, _ = cookiejar.New(nil)
	Bow.SetCookieJar(Jar)

	TargetURLsORIGINAL = map[string]Empty{}
	TargetURLsALTERED = map[string]Empty{}
	VulnerableURLs = map[string]Empty{}
	LoginInformation = make(map[string]string)

	BadUrls = append(BadUrls, []string{"%C3%A7%C4%B1k%C4%B1%C5%9F", "logout", ".png", ".jpg", ".jpeg", ".mp3", ".mp4", ".avi", ".gif", ".svg", "setup", "csrf","rar"}...)
	BadUrls = append(BadUrls, []string{"brute","exec","reset", "user_extra", "password_change","country_result.jsp"}...)
}