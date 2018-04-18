package source

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func CrawlURL(u string) {

	// Crawler travers as BFS tree

	err := Bow.Open(u)
	if err != nil {
		LogError(err)
		panic(err)
	}

	// return all links in the current page
	links := Bow.Links()


	//i:=0
	for _, link := range links {

		_, ifContains := TargetURLs[link.Url().String()]
		whitelisted := true
		for _, b := range BadUrls {
			if strings.Contains(link.Url().String(), b) {
				whitelisted = false
				break
			}

		}

		// search for same domain
		// do not allow unwanted urls like: png,jpg,pdf etc.
		// add if the URL not already in slice
		if strings.Compare(Host, link.Url().Host) == 0 && !ifContains && whitelisted {

			TargetURLs[link.Url().String()] = Empty{}
			//i++
			//fmt.Println(i, link.Url())

		}

	}
}

func LoginByCredentials(loginURL string, user string, pass string) {

	fm, _ := Bow.Form("form")
	fm.Input("username", "admin")
	fm.Input("password", "password")
	//fmt.Println(Bow.Dom().Html())
	fm.Submit()
	//fmt.Println(Bow.Dom().Html())

}

func LoginToBow(loginURL string) {

	err := Bow.Open(loginURL)
	if err != nil {
		LogError(err)
		panic(err)
	}

	allForms := Bow.Forms()
	var text string

	for _, fm := range allForms {
		if fm != nil {
			fmt.Println("Login form is found.. : ")
			fm.Dom().Find("input").Each(func(i int, s *goquery.Selection) {
				// For each item found, get the band and title
				if inputName, ok := s.Attr("name"); ok {
					if inputType, ok2 := s.Attr("type"); ok2 {
						if inputType == "text" || inputType == "password" {
							fmt.Print("Enter ", inputName, ":")
							fmt.Scanln(&text)
							fm.Input(inputName, text)
							LoginInformation[inputName] = text
						}
					}
				}

			})

			err = fm.Submit()
			if err != nil {
				LogError(err)
				panic(err)
			}

			fmt.Println("")
			fmt.Println("Succesfully loged in to Bow browser..")

			//s,_ := Bow.Dom().Html()
			//fmt.Println(s)
			//time.Sleep(12*time.Second)
			//fmt.Println("Cookies are set as:")
			//
			//for _, c := range Bow.CookieJar().Cookies(UrlParsed) {
			//	res, _ := json.Marshal(c)
			//	fmt.Println(string(res))
			//}
		}
	}
}