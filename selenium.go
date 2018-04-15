package main

import (
	"fmt"
	"github.com/serge1peshcoff/selenium-go-conditions"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"strings"
	"sync"
	"time"
	"net/http"
	"math"
)

type ChromeDriver struct {
	service   *selenium.Service
	err       error
	webDriver selenium.WebDriver
	caps      selenium.Capabilities
	mutex     *sync.Mutex
}

func (cd *ChromeDriver) initDriver(browserArgs []string) {
	cd.mutex = &sync.Mutex{}
	port := 1234

	fmt.Println("Chrome driver is going to be opened..")
	cd.service, cd.err = selenium.NewChromeDriverService("/home/cemal/chromedriver", port)
	if cd.err != nil {
		fmt.Println("couldnt open service")
		return
	}

	// Connect to the WebDriver instance running locally.
	cd.caps = selenium.Capabilities{}
	cd.caps.AddChrome(chrome.Capabilities{Args: browserArgs}) // ,"--headless", "--disable-gpu" ,"--proxy-server=http://localhost:8080"


	cd.webDriver, cd.err = selenium.NewRemote(cd.caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if cd.err != nil {
		panic(cd.err)
	}

	//cd.webDriver.SetPageLoadTimeout(800*time.Millisecond)
	//err := cd.webDriver.SetImplicitWaitTimeout(80*time.Millisecond)
	//if err != nil{
	//	fmt.Println(err)
	//}
}

func (cd *ChromeDriver) stopDriver() {
	cd.service.Stop()
	cd.webDriver.Quit()
}

func (cd *ChromeDriver) login(url string) {
	// Navigate to the simple playground interface.
	if err := cd.webDriver.Get(url); err != nil {
		log.Println("PANIC:", err)
	}

	// Find all Forms
	forms, err := cd.webDriver.FindElements(selenium.ByXPATH, "//form")
	if err != nil {
		panic(err)
	}

	var text string
	for _, form := range forms {
		inputs, err := form.FindElements(selenium.ByXPATH, "//input")
		if err != nil {
			log.Println("input could not be found", err)
		}
		for _, input := range inputs {
			InputName, err := input.GetAttribute("name")
			if err != nil {
				log.Println("input name could not be found", err)
			}
			InputType, err := input.GetAttribute("type")
			if err != nil {
				log.Println("input type could not be found", err)
			}
			if InputType == "text" || InputType == "password" {
				fmt.Print(InputName, ":")
				fmt.Scanln(&text)
				input.SendKeys(text)
			}
		}

		button, err := form.FindElement(selenium.ByTagName, "button")

		if button == nil {
			button, err = form.FindElement(selenium.ByXPATH, "//input[@type='submit']")
			if err != nil {
				log.Println(err)
			}
			err = button.Click()
			if err != nil {
				log.Println(err)
			}
		} else {
			s, _ := button.Text()
			fmt.Println(s)

			//err = button.Click()
			err = button.Submit()
			if err != nil {
				log.Println(err)
			}
		}

		// TODO: Find a way to submit buttons
		//button,err := form.FindElement(selenium.ByXPATH,"//input[@type='submit']")
		//if err != nil{
		//	// then it is not <input type="submit"> but a <button>
		//	button,err = form.FindElement(selenium.ByXPATH,"//button")
		//	if err != nil {
		//		log.Println("button could not be found ERROR:",err)
		//	}
		//
		//	fmt.Println(button.Text())
		//	err = button.Click()
		//	if err != nil {
		//		log.Println(err)
		//	}
		//
		//}else{
		//	err = button.Click()
		//	if err != nil {
		//		log.Println(err)
		//	}
		//
		//}

		//time.Sleep(10*time.Second)

		cookies, _ := cd.webDriver.GetCookies()
		for _, c := range cookies {
			fmt.Println(c)
		}
	}
}

func (cd *ChromeDriver) formTest(url string, payloads []string) (bool, string) {

	cd.mutex.Lock()
	defer cd.mutex.Unlock()

	// Navigate to the simple playground interface.
	if err := cd.webDriver.Get(url); err != nil {
		log.Println("PANIC:", err)
	}

	
	alertText, err := cd.webDriver.AlertText()
	if err != nil {
		//log.Println(err)
	} else {
		if err = cd.webDriver.AcceptAlert(); err != nil{
			fmt.Println(err)
		}
		fmt.Println("stored xss -------------->", url)
		return true, alertText
	}

	// for each form all payloads will be tested
	for _, payload := range payloads {

		//time.Sleep(100*time.Millisecond)

		// Find all forms
		forms, err := cd.webDriver.FindElements(selenium.ByXPATH, "//form")
		if err != nil {
			log.Println("form could not be found. ERROR:", err)
			return false, ""
		}

		// try each payload
		for _, form := range forms {

			//find all inputs inside the form
			inputs, err := form.FindElements(selenium.ByXPATH, "//input")
			if err != nil {
				log.Println("input could not be found. ERROR:", err)
				return false, ""
			}

			// for each input
			for _, input := range inputs {

				InputType, err := input.GetAttribute("type")
				if err != nil {
					log.Println("input type could not be found. ERROR:", err)
				}
				if InputType == "text" || InputType == "password" {
					input.SendKeys(payload)
				}
			}

			textareas, err := form.FindElements(selenium.ByXPATH, "//textarea")
			if err != nil {
				log.Println("input could not be found. ERROR:", err)
			} else {
				// for each textarea
				for _, textarea := range textareas {
					textarea.SendKeys(payload)
				}
			}

			//source,_ := cd.webDriver.PageSource()

			//log.Println(source)

			//selection, err := form.FindElement(selenium.ByXPATH, "//select")
			//if err != nil {
			//	log.Println("no selection in form")
			//}else{
			//	// if there is selection tag in form
			//	option, err := selection.FindElement(selenium.ByXPATH, "//option")
			//
			//	if err = option.SendKeys("Spanish2"); err != nil{
			//		log.Println("error setting value to selection")
			//	}
			//	option.Click()
			//}

			//time.Sleep(3*time.Second)


			button, err := form.FindElement(selenium.ByXPATH, "//input[@type='submit']")
			button.Click()

			 time.Sleep(100 * time.Millisecond)
			alertText, err2 := cd.webDriver.AlertText()
			if err2 != nil {
				//debug				//fmt.Println("error waiting alert box:",err2)
			} else if err2 == nil && alertText != "" && strings.Contains(payload, alertText) {
				fmt.Println("xss found. alert text:", alertText, " URL:", url)
				err = cd.webDriver.DismissAlert()
				if err != nil {
					//debug					fmt.Println(err)
				}
				cd.webDriver.DismissAlert()
				return true, payload
			} else {
				fmt.Println("xss not found. URL:", url)
			}

			//cd.webDriver.Back()
		}

	}

	return false, ""
}

func (cd *ChromeDriver) pageTest(url string) bool {

	cd.mutex.Lock()
	defer cd.mutex.Unlock()

	var err error

	// Navigate to the url
	if err = cd.webDriver.Get(url); err != nil {
		log.Println(err)
	}

	// TODO : page may have its own pop-up. Detect it and store in storedXssUrl dictionary.
	err = cd.webDriver.AcceptAlert()
	if err != nil {
		//log.Println(err)
	} else {
		fmt.Println("stored xss ------------> ", url)
		return true
	}

	// first try to load the page and see <html> tag
	err = cd.webDriver.WaitWithTimeout(conditions.ElementIsLocated(selenium.ByTagName, "html"), 150*time.Millisecond)
	if err != nil {
		// in case of time out, there are 2 different possibilities
		// either page could not be reached
		// or a pop-up alert is waiting

		// try to get alert text
		alertText, err2 := cd.webDriver.AlertText()
		if err2 != nil {
			fmt.Println("error waiting alert box:", err2)
		} else if err2 == nil && alertText != "" {
			fmt.Println("xss found. alert text:", alertText, " URL:", url)
			err = cd.webDriver.DismissAlert()
			if err != nil {
				fmt.Println(err)
			}
			return true
		} else {
			fmt.Println("xss not found. URL:", url)
		}
	} else {
		// if there is no timeout it means page is loaded.
		fmt.Println("xss not found. URL:", url)
	}

	return false

}

func (cd *ChromeDriver) loginToChromeAuto(url string,loginInformation map[string]string) {

	// Navigate to the url
	if err := cd.webDriver.Get(url); err != nil {
		log.Println(err)
	}

	form, err := cd.webDriver.FindElement(selenium.ByTagName, "form")
	if err != nil {
		panic("no login form")
	}

	for inputName, value := range loginInformation {
		inputStringList := []string{"//input[@name='", inputName, "']"}
		inputString := strings.Join(inputStringList, "")
		inputField, err := form.FindElement(selenium.ByXPATH, inputString)
		if err != nil {
			panic("input field could not found")
		}
		inputField.SendKeys(value)
	}
	time.Sleep(5*time.Second)
	button, err := cd.webDriver.FindElement(selenium.ByXPATH, "//input[@type='submit']")
	if err != nil {
		panic("could not send form")
	}
	err = button.Click()
	if err != nil {
		panic("could not log in to chrome")
	}
	fmt.Println("Succesfully loged in to chrome browser..")

}

func (cd *ChromeDriver) trySelection(url string) {
	// Navigate to the url
	if err := cd.webDriver.Get(url); err != nil {
		log.Println(err)
	}

	selection, err := cd.webDriver.FindElement(selenium.ByXPATH, "//select")
	if err != nil {
		panic("no login form")
	}

	option, err := selection.FindElement(selenium.ByXPATH, "//option")

	if err = option.SendKeys("Spanish"); err != nil{
		log.Println("error setting value to selection")
	}
	option.Click()

	button,err := cd.webDriver.FindElement(selenium.ByXPATH,"//input[@type='submit']")
	if err != nil {
		panic("no login form")
	}
	button.Click()

	c,_:=cd.webDriver.GetCookies()
	fmt.Println(c)

	time.Sleep(5*time.Second)
}

func (cd *ChromeDriver) setCookiesToChrome(url string,goCookies []*http.Cookie){

	// Navigate to the url
	if err := cd.webDriver.Get(url); err != nil {
		log.Println(err)
	}


	// WARNING: First You need to open the page before setting cookies
	cd.webDriver.DeleteAllCookies()

	for _,goCookie := range goCookies{
		c := selenium.Cookie{
			Name:goCookie.Name,
			Value:goCookie.Value,
			Path:goCookie.Path,
			Domain:goCookie.Domain,
			Secure:goCookie.Secure,
			Expiry:math.MaxUint32,
			//Expiry:uint(time.Now().Add(24*time.Hour).Unix()),
		}
		e:=cd.webDriver.AddCookie(&c)
		if e != nil {
			log.Println(e)
		}
	}


	c,e := cd.webDriver.GetCookies()
	if e != nil {
		log.Println(e)
	}

	log.Println(c)
}

//func (cd *chromeDriver) loginW(url string) {
//	uname := "admin"
//	pwd := "password"
//
//	// Navigate to the url
//	if err := cd.webDriver.Get(url); err != nil {
//		log.Println(err)
//	}
//
//	form, err := cd.webDriver.FindElement(selenium.ByTagName, "form")
//	input1, err := form.FindElement(selenium.ByXPATH, "//input[@name='username']")
//	input1.SendKeys(uname)
//	input2, err := form.FindElement(selenium.ByXPATH, "//input[@name='password']")
//	input2.SendKeys(pwd)
//
//	b, err := form.FindElement(selenium.ByXPATH, "//input[@type='submit']")
//
//	err = b.Click()
//
//	fmt.Println(err)
//
//	//time.Sleep(5*time.Second)
//
//}

