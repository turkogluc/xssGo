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
	"github.com/tebeka/selenium"
	"github.com/serge1peshcoff/selenium-go-conditions"
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

	flagURL := flag.String("url","http://192.168.56.101","Target Url to be scanned")
	flagLevel := flag.Int("level",3,"Scan Depth Level")
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
	TargetURLsORIGINAL[UrlSTR] = Empty{}
	TargetURLsALTERED[UrlSTR] = Empty{}

	Level = *flagLevel

	ReadPayloads(*flagPayload)

	if *flagBlackList != "" {
		blackList := *flagBlackList
		arr := strings.Split(blackList,",")
		BadUrls = append(BadUrls,arr...)
	}
	browserArgs := []string{"--disable-xss-auditor"}
	browserArgs = append(browserArgs, "--proxy-server=localhost:8888")
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

	//CrawlURL(UrlSTR)
	if Level > 1 {
		for i := 2; i <= Level; i++ {
			for tempURL, _ := range TargetURLsORIGINAL {
				CrawlURL(tempURL)
			}
		}
	}
	time.Sleep(250*time.Millisecond)

	//
	//TargetURLsALTERED["http://localhost/dvwa/vulnerabilities/xss_d/?default=English"] = Empty{}
	//TargetURLsALTERED["http://localhost/dvwa/vulnerabilities/xss_r/?name=cemal"] = Empty{}
	//TargetURLsALTERED["http://localhost/dvwa/vulnerabilities/xss_s/"] = Empty{}

	// List of target Urls
	i := 1
	fmt.Println("List of Target Urls:")
	log.Println("List of Target Urls:")
	for u, _ := range TargetURLsALTERED {
		fmt.Println(i, u)
		i++
	}
	time.Sleep(250*time.Millisecond)

	// Add one extra url to list
	//targetURLs["http://localhost/dvwa/vulnerabilities/xss_d/?default=English"] = empty{}

	for u := range TargetURLsALTERED {
		CD.Mutex.Lock()
		SendToProxy(u)
		CD.Mutex.Unlock()
		//time.Sleep(10*time.Second)
	}

	//// Query Parameters
	//log.Println("Query parameters are going to be tested")
	//fmt.Println("Query parameters are going to be tested")
	//for u := range TargetURLsALTERED {
	//	ControlQueryParameters(u)
	//}
	//time.Sleep(250*time.Millisecond)
	//
	//// Form input variables
	//log.Println("Form inputs are going to be tested")
	//fmt.Println("Form inputs are going to be tested")
	//for u := range TargetURLsORIGINAL {
	//	ControlFormInputs(u)
	//}

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

func SendToProxy(u string) {

	// TODO : ilk acilista stored xss yakalamaya calis
	// Navigate to the url
	//tempUrl := u + "$DONOTTOUCH"
	//if err := CD.WebDriver.Get(tempUrl); err != nil {
	//	LogError(err)
	//}
	//// TODO : page may have its own pop-up. Detect it and store in storedXssUrl dictionary.
	//err := CD.WebDriver.WaitWithTimeout(conditions.ElementIsLocated(selenium.ByTagName, "html"), 150*time.Millisecond)
	////if err != nil {
	////	err := CD.WebDriver.AcceptAlert()
	////	if err != nil {
	////		LogDebug(err)
	////	} else {
	////		VulnerableURLs[u] = Empty{Payload:"stored"}
	////		fmt.Println("stored xss " , u)
	////		LogInfo("stored xss ------------> ", u)
	////		return
	////	}
	////}
	//if err != nil {
	//	alertText, err2 := CD.WebDriver.AlertText()
	//	if err2 != nil {
	//		LogDebug("error waiting alert box:", err2)
	//	} else if err2 == nil {
	//		LogInfo("xss found. alert text:", alertText, " URL:", u)
	//		err = CD.WebDriver.DismissAlert()
	//		if err != nil {
	//			LogDebug(err)
	//		}
	//		err = CD.WebDriver.AcceptAlert()
	//		if err != nil {
	//			LogDebug(err)
	//		}
	//
	//		fmt.Println("bulundu-possible stored")
	//		VulnerableURLs[u] = Empty{Payload: "stored"}
	//		return
	//	}
	//}


	originalURL, _ := url.Parse(u)

	// testing query parameters
	q := originalURL.Query()

	if len(q) > 0 {
		//CD.Mutex.Lock()
		for _, payload := range Payloads {

			// SET PAYLOAD
			payloadSetter := "http://localhost/payload$" + payload
			fmt.Println("payload: " + payload)
			if err := CD.WebDriver.Get(payloadSetter); err != nil {
				LogError(err)
			}

			//fmt.Println(tempUrl)

			// Navigate to the url
			if err := CD.WebDriver.Get(u); err != nil {
				LogError(err)
			}

			//time.Sleep(5*time.Second)
			// first try to load the page and see <html> tag
			err := CD.WebDriver.WaitWithTimeout(conditions.ElementIsLocated(selenium.ByTagName, "html"), 150*time.Millisecond)
			if err != nil {
				alertText, err2 := CD.WebDriver.AlertText()
				if err2 != nil {
					LogDebug("error waiting alert box:", err2)
				} else if err2 == nil {
					LogInfo("xss found. alert text:", alertText, " URL:", u)
					err = CD.WebDriver.DismissAlert()
					if err != nil {
						LogDebug(err)
					}
					err = CD.WebDriver.AcceptAlert()
					if err != nil {
						LogDebug(err)
					}

					fmt.Println("bulundu")
					VulnerableURLs[u] = Empty{Payload: payload}
					break

				} else {
					LogInfo("xss not found. URL:", u)
				}
			} else {
				// if there is no timeout it means page is loaded.
				LogInfo("xss not found. URL:", u)
			}

			//fmt.Println("sending to test query parameter "+ u)
			//if CD.PageTest(u) {
			//	fmt.Println("Bulundu..")
			//	if _, contains := VulnerableURLs[u]; !contains {
			//		VulnerableURLs[u] = Empty{Payload:payload}
			//	}
			//	break
			//}
		}
		//CD.Mutex.Unlock()
	}

	//time.Sleep(5*time.Second)

	// testing form parameters
	// navigate to page

	found := false
	for _,payload := range Payloads {
		//CD.Mutex.Lock()

		if found {
			//break
		}

		// SET PAYLOAD
		payloadSetter := "http://localhost/payload$" + payload
		if err := CD.WebDriver.Get(payloadSetter); err != nil {
			LogError(err)
		}

		if err := CD.WebDriver.Get(u); err != nil {
			LogError(err)
		}



		forms, err := CD.WebDriver.FindElements(selenium.ByXPATH, "//form")
		if err != nil {
			LogDebug("form could not be found.")
			LogDebug(err)
		}

		if len(forms) >= 1 {
			fmt.Println("form inputs are going to be tested =>"+ u)
			for _, form := range forms {

				//find all inputs inside the form
				inputs, err := form.FindElements(selenium.ByXPATH, "//input")
				if err != nil {
					LogDebug("input could not be found. ERROR:")
					LogDebug(err)
					continue
				}

				// for each input
				for _, input := range inputs {

					InputType, err := input.GetAttribute("type")
					if err != nil {
						LogDebug("input type could not be found. ERROR:")
						LogDebug(err)
					}
					if InputType == "text" || InputType == "password" {
						input.Clear()
						// send random inputs
						input.SendKeys("123")
					}
				}

				textareas, err := form.FindElements(selenium.ByXPATH, "//textarea")
				if err != nil {
					LogDebug("input could not be found. ERROR:")
					LogDebug(err)
				} else {
					// for each textarea
					for _, textarea := range textareas {
						// send random inputs
						textarea.SendKeys("123")
					}
				}

				button, err := form.FindElement(selenium.ByXPATH, "//input[@type='submit']")
				button.Click()

				err = CD.WebDriver.WaitWithTimeout(conditions.ElementIsLocated(selenium.ByTagName, "html"), 150*time.Millisecond)
				if err != nil {
					//time.Sleep(100 * time.Millisecond)
					_, err2 := CD.WebDriver.AlertText()
					if err2 != nil {
						LogDebug(err2)
						//debug				//fmt.Println("error waiting alert box:",err2)
					} else if err2 == nil  { //&& && alertText != "" strings.Contains(payload, alertText) {
						LogInfo("xss found. payload:", payload, " URL:", u)
						fmt.Println("xss found. payload:", payload, " URL:", u)
						err = CD.WebDriver.DismissAlert()
						if err != nil {
							LogError(err)
							//debug					fmt.Println(Err)
						}
						err = CD.WebDriver.AcceptAlert()
						if err != nil {
							LogError(err)
							//debug					fmt.Println(Err)
						}
						VulnerableURLs[u] = Empty{Payload:payload}
						//CD.WebDriver.DismissAlert()

						found = true
						return

					} else {
						LogInfo("xss not found. URL:", u)
					}
				}



			}
		}else{
			break
		}

		//CD.Mutex.Unlock()
	}



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