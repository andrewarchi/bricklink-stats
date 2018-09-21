package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"

	"golang.org/x/net/publicsuffix"
)

// https://gist.github.com/varver/f327ef9087ebf76aa4c4
// https://stackoverflow.com/questions/16784419/in-golang-how-to-determine-the-final-url-after-a-series-of-redirects

func main() {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{Jar: jar}
	fmt.Println("Logging in")
	_, err = client.PostForm("https://www.bricklink.com/ajax/renovate/loginandout.ajax", url.Values{
		"userid":          {"USERNAME/EMAIL"},
		"password":        {"PASSWORD"},
		"keepme_loggedin": {"true"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Logged in")

	count, err := getOrderCount(client, 9940500, 10000001)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d total orders\n", count)
}

func getOrderCount(client http.Client, min, max int) (int, error) {
	var id int
	for {
		id = (max-min)/2 + min
		exists, err := checkOrderExist(client, id)
		if err != nil {
			return 0, err
		}
		fmt.Println(id, exists)
		if exists {
			min = id + 1
			if min > max {
				return id, nil
			}
		} else {
			max = id - 1
			if max < min {
				return id - 1, nil
			}
		}
	}
}

func checkOrderExist(client http.Client, id int) (bool, error) {
	url := "https://www.bricklink.com/orderDetail.asp?ID=" + strconv.Itoa(id)
	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}

	finalURL := resp.Request.URL.String()
	if finalURL == url || finalURL == "https://www.bricklink.com/oops.asp?err=403" {
		return true, nil
	} else if finalURL == "https://www.bricklink.com/notFound.asp?nf=order&mFolder=o&mSub=o" {
		return false, nil
	} else {
		return false, fmt.Errorf("unexpected URL: %v", finalURL)
	}
}
