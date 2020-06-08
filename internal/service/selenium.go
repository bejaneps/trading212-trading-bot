package service

import (
	"fmt"
	"time"

	"github.com/bejaneps/trading212/internal/models"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"github.com/tebeka/selenium"
)

// WebDriver is a struct that derives functions from Selenium interface,
// it's created so functions can be created against it.
type WebDriver struct {
	selenium.WebDriver
}

// NewSelenium returns an instance of initialized web driver,
// it will listen on specified port and use specified browser.
//
// Don't forget to defer close opened web driver !
func NewSelenium(port int, browser string) (*WebDriver, error) {
	wd := &WebDriver{}
	var err error

	seleniumURL := fmt.Sprintf("http://localhost:%d/wd/hub", port)

	caps := selenium.Capabilities{"browserName": browser}
	wd.WebDriver, err = selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		return nil, err
	}

	return wd, nil
}

// Navigate goes to a specified url.
func (wd *WebDriver) Navigate(url string) error {
	op := "service.Go"

	if err := wd.WebDriver.Get(url); err != nil {
		return errors.Wrapf(err, "(%s): navigating to %s page", op, url)
	}

	return nil
}

// LoginTrading212 is a utility function, that basically passes login page.
// TODO: perform login again, when session timeout on platform.
func (wd *WebDriver) LoginTrading212(username, password string, loadTime int) error {
	op := "service.LoginTrading212"

	// Navigate to the login page.
	err := wd.Navigate(models.LoginURL)
	if err != nil {
		return errors.WithMessagef(err, "(%s): ", op)
	}

	// Login account
	// Fill username
	usernameBtn, err := wd.FindElement(selenium.ByID, "username-real")
	if err != nil {
		return errors.WithMessagef(err, "(%s): ", op)
	}

	err = usernameBtn.SendKeys(username)
	if err != nil {
		return errors.WithMessagef(err, "(%s): ", op)
	}

	// Fill password
	passwordBtn, err := wd.FindElement(selenium.ByID, "pass-real")
	if err != nil {
		return errors.WithMessagef(err, "(%s): ", op)
	}

	err = passwordBtn.SendKeys(password)
	if err != nil {
		return errors.WithMessagef(err, "(%s): ", op)
	}

	// Submit form
	submitBtn, err := wd.FindElement(selenium.ByCSSSelector, `input[type="submit"]`)
	if err != nil {
		return errors.WithMessagef(err, "(%s): ", op)
	}

	err = submitBtn.Click()
	if err != nil {
		return errors.WithMessagef(err, "(%s): ", op)
	}

	log.Infof("Logged to Trading212 successfully.")

	time.Sleep(time.Second * time.Duration(loadTime)) // wait for login

	return nil
}
