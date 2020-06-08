package crud

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
)

// Commodity is representation of commodity on Trading212 platform.
type Commodity struct {
	Name     string  `json:"name,omitempty"`  // full name of commodity
	Quantity float64 `json:"quantity"`        // if empty, then use default amount
	Price    float64 `json:"price"`           // if empty, then use market price
	Order    string  `json:"order,omitempty"` // type of order, buy or sells
}

// Buy performs a buy order in market, with specified order,
// it returns output of alert that comes out, eg: Order Executed.
func (c *Commodity) Buy(wd selenium.WebDriver) (string, error) {
	op := "crud.Buy"

	var (
		err  error
		elem selenium.WebElement
	)

	// Select search box
	elem, err = wd.FindElement(selenium.ByCSSSelector, ".search-button")
	if err != nil { // search field possibly minimized
		return "", errors.WithMessagef(err, "(%s): ", op)
	}

	// Find specified commodity
	err = elem.SendKeys(c.Name)
	if err != nil { // search field possibly minimized
		elem, err = wd.FindElement(selenium.ByXPATH, `//*[@id="navigation-search-button"]`)
		if err != nil {
			return "", errors.WithMessagef(err, "(%s): ", op)
		}

		// click on minimized search button
		err = elem.Click()
		if err != nil {
			return "", errors.WithMessagef(err, "(%s): ", op)
		}

		// select search box
		elem, err = wd.FindElement(selenium.ByXPATH, "/html/body/div[7]/div[1]/div/div[2]/input")
		if err != nil {
			return "", errors.WithMessagef(err, "(%s): ", op)
		}

		// Find specified commodity
		err = elem.SendKeys(c.Name)
		if err != nil {
			return "", errors.WithMessagef(err, "(%s): ", op)
		}
	}

	// Click on it, so buy order is possible
	elem, err = wd.FindElement(selenium.ByCSSSelector, ".search-results-instrument")
	if err != nil {
		return "", errors.WithMessagef(err, "(%s): ", op)
	}

	err = elem.Click()
	if err != nil {
		return "", errors.WithMessagef(err, "(%s): ", op)
	}

	// if quantity not set, then set default quantity
	if c.Quantity == 0 {
		c.Quantity = 1
	}

	// Convert quantity to string
	quantityString := strconv.FormatFloat(c.Quantity, 'f', 3, 64)
	if quantityString == "" {
		return "", errors.Wrapf(err, "(%s): wrong commodity quantity", op)
	}

	time.Sleep(time.Second * 1)

	// Find detailed trade box
	elem, err = wd.FindElement(selenium.ByCSSSelector, ".search-instrument-details .tradebox .tradebox-open-order-dialog-icon")
	if err != nil {
		return "", errors.WithMessagef(err, "(%s): ", op)
	}

	// Click it
	err = elem.Click()
	if err != nil {
		return "", errors.WithMessagef(err, "(%s): ", op)
	}

	// Buy at market price
	if c.Price == 0 {
		// Find market order button
		elem, err = wd.FindElement(selenium.ByCSSSelector, `span[data-tab="market-order"]`)
		if err != nil {
			return "", errors.WithMessagef(err, "(%s): ", op)
		}

		// Click it
		err = elem.Click()
		if err != nil {
			return "", errors.WithMessagef(err, "(%s): ", op)
		}

		time.Sleep(time.Second * 1)

		// Find quantity box
		elem, err = wd.FindElement(selenium.ByCSSSelector, ".visible-input input")
		if err != nil {
			return "", errors.WithMessagef(err, "(%s): ", op)
		}

		// Clear default quantity
		err = elem.Clear()
		if err != nil {
			return "", errors.WithMessagef(err, "(%s): ", op)
		}

		// Input quantity
		err = elem.SendKeys(quantityString)
		if err != nil {
			return "", errors.WithMessagef(err, "(%s): ", op)
		}
	} else { // Buy at specified price
		// Convert price to string
		priceString := strconv.FormatFloat(c.Price, 'f', 3, 64)
		if priceString == "" {
			return "", errors.Wrapf(err, "(%s): wrong commodity price", op)
		}

		// Find limist/stop order button
		elem, err = wd.FindElement(selenium.ByCSSSelector, `span[data-tab="limit_stop-order"]`)
		if err != nil {
			return "", errors.WithMessagef(err, "(%s): ", op)
		}

		// Click it
		err = elem.Click()
		if err != nil {
			return "", errors.WithMessagef(err, "(%s): ", op)
		}

		time.Sleep(time.Second * 1)

		// TODO: fix issue with element not interactable

		// Find price input box
		elem, err = wd.FindElement(selenium.ByCSSSelector, ".spinner input")
		if err != nil {
			return "", errors.WithMessagef(err, "(%s): ", op)
		}

		// Clear default price
		err = elem.Clear()
		if err != nil {
			return "", errors.WithMessagef(err, "(%s): ", op)
		}

		// Input price
		err = elem.SendKeys(priceString)
		if err != nil {
			return "", errors.WithMessagef(err, "(%s): ", op)
		}

		// Find quantity box
		elem, err = wd.FindElement(selenium.ByCSSSelector, ".visible-input input")
		if err != nil {
			return "", errors.WithMessagef(err, "(%s): here", op)
		}

		// Clear default quantity
		err = elem.Clear()
		if err != nil {
			return "", errors.WithMessagef(err, "(%s): ", op)
		}

		// Input quantity
		err = elem.SendKeys(quantityString)
		if err != nil {
			return "", errors.WithMessagef(err, "(%s): ", op)
		}
	}

	// Find confirm button
	elem, err = wd.FindElement(selenium.ByCSSSelector, ".confirm-button")
	if err != nil {
		return "", errors.WithMessagef(err, "(%s): ", op)
	}

	// Click it
	err = elem.Click()
	if err != nil {
		return "", errors.WithMessagef(err, "(%s): ", op)
	}

	time.Sleep(time.Millisecond * 200)

	// Output alert, eg: "Order Executed"
	elem, err = wd.FindElement(selenium.ByCSSSelector, ".alert")
	if err != nil {
		return "", errors.WithMessagef(err, "(%s): ", op)
	}

	output, err := elem.Text()
	if err != nil {
		return "", errors.WithMessagef(err, "(%s): order not executed", op)
	}

	return output, nil
}

// Sell performs a sell order in market, with specified order,
// it returns output of alert that comes out, eg: Order Executed.
func (c *Commodity) Sell(wd selenium.WebDriver) (string, error) {
	return "", nil
}
