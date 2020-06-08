package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
)

const (
	port = 4444 // default port of Selenium server instance

	username = "bejaneps@gmail.com"
	password = "Saburi23031999"
)

var (
	seleniumURL = fmt.Sprintf("http://localhost:%d/wd/hub", port)

	loginURL   = "https://www.trading212.com/en/login"
	tradingURL = "https://demo.trading212.com/"
)

func main() {
	log.SetFlags(log.Lshortfile) // debug error line

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		log.Fatal(err)
	}
	defer wd.Quit()

	// Navigate to the login page.
	if err := wd.Get(loginURL); err != nil {
		log.Fatal(err)
	}

	// Login account
	// Fill username
	usernameBtn, err := wd.FindElement(selenium.ByID, "username-real")
	if err != nil {
		log.Fatal(err)
	}

	err = usernameBtn.SendKeys(username)
	if err != nil {
		log.Fatal(err)
	}

	// Fill password
	passwordBtn, err := wd.FindElement(selenium.ByID, "pass-real")
	if err != nil {
		log.Fatal(err)
	}

	err = passwordBtn.SendKeys(password)
	if err != nil {
		log.Fatal(err)
	}

	// Submit form
	submitBtn, err := wd.FindElement(selenium.ByCSSSelector, `input[type="submit"]`)
	if err != nil {
		log.Fatal(err)
	}

	err = submitBtn.Click()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 6) // wait for login

	// Navigate to trading page.
	if err := wd.Get(tradingURL); err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 7) // wait for trading page loading

	// Find buy button
	btn, err := wd.FindElement(selenium.ByCSSSelector, ".tradebox-buy")
	if err != nil {
		log.Fatal(err)
	}

	// Click on buy button, with default amount
	if err := btn.Click(); err != nil {
		log.Fatal(err)
	}

	// Check how much is bought
	elem, err := wd.FindElement(selenium.ByXPATH, "/html/body/div[5]/div[3]/div/div[2]/div[4]/div/table/tbody/tr/td[3]")
	if err != nil {
		log.Fatal(err)
	}

	amount, err := elem.Text()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Amount: %s\n", amount)
}
