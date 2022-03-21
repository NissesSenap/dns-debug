package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	myurl := os.Getenv("URL")
	if myurl == "" {
		myurl = "google.com"
	}
	endpoints := []string{"google.com", "google.se"}
	fmt.Printf("Endpoints: %v\n", endpoints)
	fmt.Printf("DNS URL: %v\n", myurl)
	for _, name := range endpoints {
		_, err := url.Parse(name)
		if err != nil {
			panic(err)
		}
	}

	tracer.Start()
	defer tracer.Stop()
	for {
		err := lookup(myurl)
		if err != nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	os.Exit(1)
}

func lookup(myurl string) error {
	span := tracer.StartSpan("dns.request", tracer.ResourceName("dns"))
	defer span.Finish()
	start := time.Now()
	ips, err := net.LookupIP(myurl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		return err
	}
	duration := time.Since(start)
	fmt.Printf("This was the duration: %v\n", duration)
	fmt.Printf("This is the IP:s %v\n", ips)
	// Set tag
	span.SetTag("url", myurl)
	return nil
}
