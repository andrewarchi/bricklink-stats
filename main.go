package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/publicsuffix"
)

// https://gist.github.com/varver/f327ef9087ebf76aa4c4
// https://stackoverflow.com/questions/16784419/in-golang-how-to-determine-the-final-url-after-a-series-of-redirects

var timeFormat = "2006/01/02 15:04:05"

type order struct {
	id   int
	time time.Time
}

func main() {
	delayStr := "5s"
	if len(os.Args) >= 2 {
		delayStr = os.Args[1]
	}
	delay, err := time.ParseDuration(delayStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s query delay\n", delay.String())

	client := createClient("USERNAME", "PASSWORD")

	t := time.Now()
	estimate := int(0.0459291*float64(t.Unix()) - 60679590.20236)
	fmt.Printf("%s  %d  estimated\n", t.Format(timeFormat), estimate)

	o, err := getOrderRange(client, estimate-10000, estimate+10000)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s  %d  actual\n", o.time.Format(timeFormat), o.id)
	orders := []order{o}

	id := o.id + 1
	for {
		t = time.Now()
		exist, err := checkOrderExist(client, id)
		if err != nil {
			log.Fatal(err)
		}
		if exist {
			orders = addOrder(orders, order{id, t})
			id++
			continue
		}
		time.Sleep(delay)
	}
}

func createClient(username, password string) http.Client {
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
		"userid":          {username},
		"password":        {password},
		"keepme_loggedin": {"true"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Logged in")
	return client
}

func addOrder(orders []order, o order) []order {
	orders = append(orders, o)
	diff := o.time.Sub(orders[len(orders)-2].time)
	count := 25
	if len(orders) < count {
		count = len(orders)
	}
	avg := float64(o.time.Sub(orders[len(orders)-count].time).Nanoseconds()) / float64(count-1)
	fmt.Printf("%s  %d  %-10s  %s\n",
		o.time.Format(timeFormat),
		o.id,
		diff.Round(time.Millisecond),
		(time.Duration(avg) * time.Nanosecond).Round(time.Millisecond))
	return orders
}

func getOrderRange(client http.Client, minID, maxID int) (order, error) {
	min, max := minID, maxID
	var id int
	for {
		id = (max-min)/2 + min
		t := time.Now()
		exists, err := checkOrderExist(client, id)
		if err != nil {
			return order{}, err
		}
		//fmt.Println(id, exists)
		if exists {
			min = id + 1
			if min > max {
				if max == maxID {
					return getOrderRange(client, maxID+1, maxID+20000)
				}
				return order{id, t}, nil
			}
		} else {
			max = id - 1
			if max < min {
				if min == minID {
					return getOrderRange(client, minID-20000, minID-1)
				}
				return order{id - 1, t}, nil
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
