package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	myurl := os.Getenv("URL")
	if myurl == "" {
		myurl = "google.com"
	}
	curl := os.Getenv("CURL")
	if curl == "" {
		curl = "true"
	}
	curlBool, err := strconv.ParseBool(curl)
	if err != nil {
		panic(err)
	}
	endpoints := []string{"https://google.com", "https://google.se", "https://gp.se", "https://dn.se", "https://di.se"}
	fmt.Printf("Endpoints: %v\n", endpoints)
	fmt.Printf("DNS URL: %v\n", myurl)
	for _, endpoint := range endpoints {
		_, err := url.Parse(endpoint)
		if err != nil {
			panic(err)
		}
	}
	run(myurl, curlBool, endpoints)
}

func run(myurl string, curlBool bool, endpoints []string) {
	c := &http.Client{Timeout: time.Duration(2) * time.Second}

	tracer.Start()
	defer tracer.Stop()
	for {
		err := lookup(myurl)
		if err != nil {
			break
		}
		if curlBool {
			for _, endpoint := range endpoints {
				err := getHTTP(endpoint, c)
				if err != nil {
					break
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
	fmt.Println("something broke")
	os.Exit(1)
}

func lookup(myurl string) error {
	span := tracer.StartSpan("dns.request", tracer.ResourceName("dns"))
	defer span.Finish()
	start := time.Now()
	ips, err := net.LookupIP(myurl)
	if err != nil {
		fmt.Printf("Could not get IPs: %v\n", err)
		return err
	}
	duration := time.Since(start)
	fmt.Printf("This was the duration for dns: %v\n", duration)
	fmt.Printf("This is the IP:s %v\n", ips)
	// Set tag
	span.SetTag("url", myurl)
	return nil
}

func getHTTP(endpoint string, c *http.Client) error {
	span := tracer.StartSpan("http.request", tracer.ResourceName("http"))
	defer span.Finish()
	start := time.Now()
	resp, err := c.Get(endpoint)
	if err != nil {
		fmt.Printf("Unable to get endpoint: %v, got error: %s", endpoint, err)
		return err
	}
	duration := time.Since(start)
	fmt.Printf("This was the duration for get: %v\n", duration)
	fmt.Printf("Endpoint: %v status code: %v\n", endpoint, resp.StatusCode)
	span.SetTag("endpoint", endpoint)
	return nil
}
