package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"encoding/json"
)

func crawlURL(u string) {

	// Crawler travers as BFS tree

	err := bow.Open(u)
	if err != nil {
		panic(err)
	}

	// return all links in the current page
	links := bow.Links()


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