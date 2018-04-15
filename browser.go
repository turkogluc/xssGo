package main

import (
	"github.com/sclevine/agouti"
	"log"
	"sync"
	//"time"
	"fmt"
	"net/http"
	//"time"
	"net/url"
	"time"
)

type browserDriver struct {
	driver *agouti.WebDriver
	page   *agouti.Page
	mutex  *sync.Mutex
	cookie []*http.Cookie
}

func (bd *browserDriver) DriverInit(url string) (err error) {
	bd.driver = agouti.ChromeDriver(agouti.ChromeOptions("args", []string{"--disable-xss-auditor"}))

	// "--headless", "--disable-gpu", "--no-sandbox",

	//command := []string{"java", "-jar", "/home/cemal/selenium-server.jar", "-port", ""}
	//driver := agouti.NewWebDriver("http://localhost/wd/hub", command)

	//driver := agouti.GeckoDriver()

	bd.mutex = &sync.Mutex{}

	err = nil
	if err = bd.driver.Start(); err != nil {
		log.Fatal("Failed to start driver:", err)
	}

	bd.page, err = bd.driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Fatal("Failed to open page:", err)
	} else {
		bd.page.SetPageLoad(3)
	}

	if err := bd.page.Navigate(url); err != nil {
		log.Fatal("Failed to navigate:", err)
	}

	return err
}

func (bd *browserDriver) DriverSetCookies(jar http.CookieJar, url *url.URL) {

}

func (bd *browserDriver) BrowserQueryTest(url string, payload string) bool {

	bd.mutex.Lock()
	//time.Sleep(time.Millisecond)

	if err := bd.page.Navigate(url); err != nil {
		log.Fatal("Failed to navigate:", err)
	}

	time.Sleep(100 * time.Millisecond)

	sectionTitle, err := bd.page.Title()
	if err != nil {
		fmt.Println(err)
	}
	u := ""
	u, err = bd.page.URL()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("title:", sectionTitle, " url:", u)
	//
	////log.Println(page.Find("#alert").Text())
	//
	//u,_ := page.URL()
	//log.Print("url",u)
	time.Sleep(100 * time.Millisecond)
	state := false
	popup, err := bd.page.PopupText()
	if err != nil {
		log.Print("No pop up:", err)
		state = false
	} else {
		bd.page.ConfirmPopup()
		log.Println("yes ", popup)
		state = true
	}
	bd.mutex.Unlock()
	return state
}

func (bd *browserDriver) DriverStop() {
	if err := bd.driver.Stop(); err != nil {
		log.Fatal("Failed to close pages and stop WebDriver:", err)
	}
}

func (bd *browserDriver) login() {
	url := "http://localhost/dvwa/login.php"
	if err := bd.page.Navigate(url); err != nil {
		log.Fatal("Failed to navigate:", err)
	}

	//current,err := bd.page.URL()
	//if current != url {
	//	log.Fatal("Something is wrong:", err)
	//}

	var err error

	if err = bd.page.FindByLabel("Username").Fill("admin"); err != nil {
		log.Fatal("Something is wrong:", err)
	}
	if err = bd.page.FindByLabel("password").Fill("password"); err != nil {
		log.Fatal("Something is wrong:", err)
	}
	if err = bd.page.FindByButton("login").Submit(); err != nil {
		log.Fatal("Something is wrong:", err)
	}

}

//func controlFormInputs(url string,cookie []*http.Cookie){
//
//	url = urlSTR
//	// open the page
//	if err := page.Navigate(url); err != nil {
//		log.Fatal("Failed to navigate:", err)
//	}
//
//	//// clear cookies
//	//if err := page.ClearCookies(); err!=nil{
//	//	log.Print(err)
//	//}
//	//
//	//// set new cookies from beginning
//	//for _,c := range cookie{
//	//	page.SetCookie(c)
//	//}
//
//
//	forms := page.FindByXPath("//form")
//	user := forms.FindByName("username")
//	pass := forms.FindByName("password")
//	if err := user.Fill("admin"); err != nil {
//		fmt.Println(err)
//	}
//
//	pass.Fill("password")
//
//	button := forms.FindByName("submit")
//	button.Submit()
//
//
//
//
//
//	fmt.Println(page.URL())
//
//
//
//}
