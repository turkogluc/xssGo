package source

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
	Service   *selenium.Service
	Err       error
	WebDriver selenium.WebDriver
	Caps      selenium.Capabilities
	Mutex     *sync.Mutex
}

func (cd *ChromeDriver) InitDriver(browserArgs []string) {
	cd.Mutex = &sync.Mutex{}
	port := 1234

	fmt.Println("Chrome driver is going to be opened..")
	cd.Service, cd.Err = selenium.NewChromeDriverService("/home/cemal/chromedriver", port)
	if cd.Err != nil {
		LogError(cd.Err)
		return
	}

	// Connect to the WebDriver instance running locally.
	cd.Caps = selenium.Capabilities{}
	cd.Caps.AddChrome(chrome.Capabilities{Args: browserArgs}) // ,"--headless", "--disable-gpu" ,"--proxy-server=http://localhost:8080"


	cd.WebDriver, cd.Err = selenium.NewRemote(cd.Caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if cd.Err != nil {
		LogError(cd.Err)
		panic(cd.Err)
	}

	//CD.WebDriver.SetPageLoadTimeout(800*time.Millisecond)
	//Err := CD.WebDriver.SetImplicitWaitTimeout(80*time.Millisecond)
	//if Err != nil{
	//	fmt.Println(Err)
	//}
}

func (cd *ChromeDriver) StopDriver() {
	cd.Service.Stop()
	cd.WebDriver.Quit()
}

func (cd *ChromeDriver) Login(url string) {
	// Navigate to the simple playground interface.
	if err := cd.WebDriver.Get(url); err != nil {
		LogError(err)
	}

	// Find all Forms
	forms, err := cd.WebDriver.FindElements(selenium.ByXPATH, "//form")
	if err != nil {
		LogError(err)
		panic(err)
	}

	var text string
	for _, form := range forms {
		inputs, err := form.FindElements(selenium.ByXPATH, "//input")
		if err != nil {
			LogDebug("input could not be found: ") //err)
			LogDebug(err)
		}
		for _, input := range inputs {
			InputName, err := input.GetAttribute("name")
			if err != nil {
				LogDebug("input name could not be found")
				LogDebug(err)
			}
			InputType, err := input.GetAttribute("type")
			if err != nil {
				LogDebug("input type could not be found")
				LogDebug(err)
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
				LogDebug(err)
			}
			err = button.Click()
			if err != nil {
				LogDebug(err)
			}
		} else {
			s, _ := button.Text()
			fmt.Println(s)

			//Err = button.Click()
			err = button.Submit()
			if err != nil {
				LogDebug(err)
			}
		}

		// TODO: Find a way to submit buttons
		//button,Err := form.FindElement(selenium.ByXPATH,"//input[@type='submit']")
		//if Err != nil{
		//	// then it is not <input type="submit"> but a <button>
		//	button,Err = form.FindElement(selenium.ByXPATH,"//button")
		//	if Err != nil {
		//		log.Println("button could not be found ERROR:",Err)
		//	}
		//
		//	fmt.Println(button.Text())
		//	Err = button.Click()
		//	if Err != nil {
		//		log.Println(Err)
		//	}
		//
		//}else{
		//	Err = button.Click()
		//	if Err != nil {
		//		log.Println(Err)
		//	}
		//
		//}

		//time.Sleep(10*time.Second)

		cookies, _ := cd.WebDriver.GetCookies()
		for _, c := range cookies {
			fmt.Println(c)
		}
	}
}

func (cd *ChromeDriver) FormTest(url string, payloads []string) (bool, string) {

	cd.Mutex.Lock()
	defer cd.Mutex.Unlock()

	// Navigate to the simple playground interface.
	if err := cd.WebDriver.Get(url); err != nil {
		LogError(err)
	}

	
	alertText, err := cd.WebDriver.AlertText()
	if err != nil {
		//log.Println(Err)
	} else {
		if err = cd.WebDriver.AcceptAlert(); err != nil{
			LogDebug(err)
		}
		LogInfo("stored xss -------------->", url)
		return true, alertText
	}

	// for each form all Payloads will be tested
	for _, payload := range payloads {

		//time.Sleep(100*time.Millisecond)

		// Find all forms
		forms, err := cd.WebDriver.FindElements(selenium.ByXPATH, "//form")
		if err != nil {
			LogDebug("form could not be found. ERROR:")
			LogDebug(err)
			return false, ""
		}

		// try each Payload
		for _, form := range forms {

			//find all inputs inside the form
			inputs, err := form.FindElements(selenium.ByXPATH, "//input")
			if err != nil {
				LogDebug("input could not be found. ERROR:")
				LogDebug(err)
				return false, ""
			}

			// for each input
			for _, input := range inputs {

				InputType, err := input.GetAttribute("type")
				if err != nil {
					LogDebug("input type could not be found. ERROR:")
					LogDebug(err)
				}
				if InputType == "text" || InputType == "password" {
					input.SendKeys(payload)
				}
			}

			textareas, err := form.FindElements(selenium.ByXPATH, "//textarea")
			if err != nil {
				LogDebug("input could not be found. ERROR:")
				LogDebug(err)
			} else {
				// for each textarea
				for _, textarea := range textareas {
					textarea.SendKeys(payload)
				}
			}

			//source,_ := CD.WebDriver.PageSource()

			//log.Println(source)

			//selection, Err := form.FindElement(selenium.ByXPATH, "//select")
			//if Err != nil {
			//	log.Println("no selection in form")
			//}else{
			//	// if there is selection tag in form
			//	option, Err := selection.FindElement(selenium.ByXPATH, "//option")
			//
			//	if Err = option.SendKeys("Spanish2"); Err != nil{
			//		log.Println("error setting value to selection")
			//	}
			//	option.Click()
			//}

			//time.Sleep(3*time.Second)


			button, err := form.FindElement(selenium.ByXPATH, "//input[@type='submit']")
			button.Click()

			//time.Sleep(100 * time.Millisecond)
			alertText, err2 := cd.WebDriver.AlertText()
			if err2 != nil {
				LogDebug(err2)
				//debug				//fmt.Println("error waiting alert box:",err2)
			} else if err2 == nil && alertText != "" && strings.Contains(payload, alertText) {
				LogInfo("xss found. alert text:", alertText, " URL:", url)
				err = cd.WebDriver.DismissAlert()
				if err != nil {
					LogError(err)
					//debug					fmt.Println(Err)
				}
				cd.WebDriver.DismissAlert()
				return true, payload
			} else {
				LogInfo("xss not found. URL:", url)
			}

			//CD.WebDriver.Back()
		}

	}

	return false, ""
}

func (cd *ChromeDriver) PageTest(url string) bool {

	cd.Mutex.Lock()
	defer cd.Mutex.Unlock()

	var err error

	// Navigate to the url
	if err = cd.WebDriver.Get(url); err != nil {
		LogError(err)
	}

	// TODO : page may have its own pop-up. Detect it and store in storedXssUrl dictionary.
	err = cd.WebDriver.AcceptAlert()
	if err != nil {
		LogDebug(err)
	} else {
		LogInfo("stored xss ------------> ", url)
		return true
	}

	// first try to load the page and see <html> tag
	err = cd.WebDriver.WaitWithTimeout(conditions.ElementIsLocated(selenium.ByTagName, "html"), 150*time.Millisecond)
	if err != nil {
		// in case of time out, there are 2 different possibilities
		// either page could not be reached
		// or a pop-up alert is waiting

		// try to get alert text
		alertText, err2 := cd.WebDriver.AlertText()
		if err2 != nil {
			LogDebug("error waiting alert box:", err2)
		} else if err2 == nil && alertText != "" {
			LogInfo("xss found. alert text:", alertText, " URL:", url)
			err = cd.WebDriver.DismissAlert()
			if err != nil {
				LogDebug(err)
			}
			return true
		} else {
			LogInfo("xss not found. URL:", url)
		}
	} else {
		// if there is no timeout it means page is loaded.
		LogInfo("xss not found. URL:", url)
	}

	return false

}

func (cd *ChromeDriver) LoginToChromeAuto(url string,loginInformation map[string]string) {

	// Navigate to the url
	if err := cd.WebDriver.Get(url); err != nil {
		log.Println(err)
	}

	form, err := cd.WebDriver.FindElement(selenium.ByTagName, "form")
	if err != nil {
		LogError(err)
	}

	for inputName, value := range loginInformation {
		inputStringList := []string{"//input[@name='", inputName, "']"}
		inputString := strings.Join(inputStringList, "")
		inputField, err := form.FindElement(selenium.ByXPATH, inputString)
		if err != nil {
			LogError("input field could not found")
		}
		inputField.SendKeys(value)
	}
	//time.Sleep(5*time.Second)
	button, err := cd.WebDriver.FindElement(selenium.ByXPATH, "//input[@type='submit']")
	if err != nil {
		LogError("could not send form")
	}
	err = button.Click()
	if err != nil {
		LogError("could not log in to chrome")
	}
	fmt.Println("Succesfully loged in to chrome browser..")

}

func (cd *ChromeDriver) TrySelection(url string) {
	// Navigate to the url
	if err := cd.WebDriver.Get(url); err != nil {
		log.Println(err)
	}

	selection, err := cd.WebDriver.FindElement(selenium.ByXPATH, "//select")
	if err != nil {
		panic("no Login form")
	}

	option, err := selection.FindElement(selenium.ByXPATH, "//option")

	if err = option.SendKeys("Spanish"); err != nil{
		log.Println("error setting value to selection")
	}
	option.Click()

	button,err := cd.WebDriver.FindElement(selenium.ByXPATH,"//input[@type='submit']")
	if err != nil {
		panic("no Login form")
	}
	button.Click()

	c,_:=cd.WebDriver.GetCookies()
	fmt.Println(c)

	//time.Sleep(5*time.Second)
}

func (cd *ChromeDriver) SetCookiesToChrome(url string,goCookies []*http.Cookie){

	// Navigate to the url
	if err := cd.WebDriver.Get(url); err != nil {
		LogError(err)
	}


	// WARNING: First You need to open the page before setting Cookies
	cd.WebDriver.DeleteAllCookies()

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
		e:=cd.WebDriver.AddCookie(&c)
		if e != nil {
			LogError(e)
		}
	}


	c,e := cd.WebDriver.GetCookies()
	if e != nil {
		LogError(e)
	}

	LogInfo(c)
}

//func (CD *chromeDriver) loginW(url string) {
//	uname := "admin"
//	pwd := "password"
//
//	// Navigate to the url
//	if Err := CD.WebDriver.Get(url); Err != nil {
//		log.Println(Err)
//	}
//
//	form, Err := CD.WebDriver.FindElement(selenium.ByTagName, "form")
//	input1, Err := form.FindElement(selenium.ByXPATH, "//input[@name='username']")
//	input1.SendKeys(uname)
//	input2, Err := form.FindElement(selenium.ByXPATH, "//input[@name='password']")
//	input2.SendKeys(pwd)
//
//	b, Err := form.FindElement(selenium.ByXPATH, "//input[@type='submit']")
//
//	Err = b.Click()
//
//	fmt.Println(Err)
//
//	//time.Sleep(5*time.Second)
//
//}

