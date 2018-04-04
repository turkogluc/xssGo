package main

import (
	"github.com/sclevine/agouti"
	"log"
	"sync"
	//"time"
)

var driver *agouti.WebDriver
var page *agouti.Page

var mutex *sync.Mutex

func DriverInit()(err error){
	driver = agouti.ChromeDriver(agouti.ChromeOptions("args", []string{"--headless", "--disable-gpu", "--no-sandbox","--disable-xss-auditor"}))

	//command := []string{"java", "-jar", "/home/cemal/selenium-server.jar", "-port", ""}
	//driver := agouti.NewWebDriver("http://localhost/wd/hub", command)

	//driver := agouti.GeckoDriver()

	mutex = &sync.Mutex{}

	err = nil
	if err = driver.Start(); err != nil {
		log.Fatal("Failed to start driver:", err)
	}

	page, err = driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Fatal("Failed to open page:", err)
	}
	return err
}

func BrowserTest(url string,payload string)(bool) {


	mutex.Lock()
	//time.Sleep(time.Millisecond)

	if err := page.Navigate(url); err != nil {
		log.Fatal("Failed to navigate:", err)
	}


	//sectionTitle, err := page.HTML()
	//log.Println("file:", sectionTitle)
	//
	////log.Println(page.Find("#alert").Text())
	//
	//u,_ := page.URL()
	//log.Print("url",u)
	state := false
	popup,err := page.PopupText()
	if err != nil {
		log.Print("No pop up:", err)
		state = false
	}else{
		page.ConfirmPopup()
		log.Println("yes ", popup)
		state = true
	}
	mutex.Unlock()
	return state
}

func DriverStop(){
	if err := driver.Stop(); err != nil {
		log.Fatal("Failed to close pages and stop WebDriver:", err)
	}
}