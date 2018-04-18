package source

import (
	"io/ioutil"
	"os"
	"net/http"
	"time"
	"encoding/json"
)

type CookieStore struct {
	SCookies		[]SingleCookie	`json:"Cookies"`
}

type SingleCookie struct {
	Name  		string	`json:"name"`
	Value 		string	`json:"value"`
	Path       	string	`json:"path"`
	Domain    	string	`json:"domain"`
	Expires    	int64	`json:"expires"`
	Secure   	bool	`json:"secure"`
	HttpOnly 	bool	`json:"httpOnly"`
}

func ConvertCookiesToGolang(cs *CookieStore) []*http.Cookie {
	//GoCookies = []*http.Cookie{}
	// copying my cookie type to goLang http.cookie
	for _,sCookie := range cs.SCookies{
		expiration := time.Unix(sCookie.Expires,0)
		cookie := &http.Cookie{
			Name:sCookie.Name,
			Value:sCookie.Value,
			Path:sCookie.Path,
			Domain:sCookie.Domain,
			Expires:expiration,
			Secure:sCookie.Secure,
			HttpOnly:sCookie.HttpOnly,
		}
		GoCookies = append(GoCookies, cookie)
	}
	return GoCookies
}

func SetCookiesToBow(cs *CookieStore) {

	// FIXME

	////cookie := &http.Cookie{
	////	Domain: "localhost",
	////	Name:   "PHPSESSID",
	////	Value:  "1l2fb11illqvns8ucu0fku1684",
	////}
	////Cookies = append(Cookies, cookie)
	////
	////cookie = &http.Cookie{
	////	Domain: "localhost",
	////	Name:   "security",
	////	Value:  "impossible",
	////}
	////Cookies = append(Cookies, cookie)
	//
	//
	//
	//Jar.SetCookies(UrlParsed, ConvertCookiesToGolang(cs))
	//Bow.SetCookieJar(Jar)
	//
	//
	//
	//
	////
	////
	////GoCookies := ConvertCookiesToGolang(cs)
	////Jar.SetCookies(UrlParsed,GoCookies)
	////Bow.SetCookieJar(Jar)
	////
	////
	////fmt.Println(Bow.SiteCookies())
	//
	//
	//fmt.Println("Cookies are set successfully :")
	////for _, c := range Bow.CookieJar().Cookies(UrlParsed) {
	////	res, _ := json.Marshal(c)
	////	fmt.Println(string(res))
	////}

}

func ReadCookiesFromFile(filePath string) (*CookieStore,error){
	// change cookie.json with filePath
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		LogError(e)
		os.Exit(1)
	}

	// parse Cookies to cookiestore which is only defined to read from json
	var cs CookieStore
	e = json.Unmarshal(file,&cs)
	if e != nil{
		LogError(e)
		return nil,e
	}
	return &cs,e
}