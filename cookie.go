package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"net/http"
	"time"
	"encoding/json"
	"log"
)

type CookieStore struct {
	SCookies		[]SingleCookie	`json:"cookies"`
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

func convertCookiesToGolang(cs *CookieStore) []*http.Cookie{
	goCookies := []*http.Cookie{}
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
		goCookies = append(goCookies, cookie)
	}
	return goCookies
}

func SetCookiesToBow(cs *CookieStore) {

	//cookie := &http.Cookie{
	//	Domain: "localhost",
	//	Name:   "PHPSESSID",
	//	Value:  "j5mdm88v2ougl2kjrrinsfcsd1",
	//}
	//cookies = append(cookies, cookie)
	//
	//cookie = &http.Cookie{
	//	Domain: "localhost",
	//	Name:   "security",
	//	Value:  "impossible",
	//}
	//cookies = append(cookies, cookie)



	goCookies := convertCookiesToGolang(cs)

	jar.SetCookies(urlParsed, goCookies)
	bow.SetCookieJar(jar)

	fmt.Println("Cookies are set successfully :")
	//for _, c := range bow.CookieJar().Cookies(urlParsed) {
	//	res, _ := json.Marshal(c)
	//	fmt.Println(string(res))
	//}

}

func readCookiesFromFile(filePath string) (*CookieStore,error){
	// change cookie.json with filePath
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	// parse cookies to cookiestore which is only defined to read from json
	var cs CookieStore
	e = json.Unmarshal(file,&cs)
	if e != nil{
		log.Println(e)
		return nil,e
	}
	return &cs,e
}