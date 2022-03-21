package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	url := os.Getenv("URL")
	if url == "" {
		url = "google.com"
	}
	fmt.Println(url)
	tracer.Start()
	defer tracer.Stop()
	for {
		err := lookup(url)
		if err != nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	os.Exit(1)
}

func lookup(url string) error {
	span := tracer.StartSpan("dns.request", tracer.ResourceName("dns"))
	defer span.Finish()
	ips, err := net.LookupIP(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		return err
	}
	fmt.Printf("This is the IP:s %v", ips)
	// Set tag
	span.SetTag("url", url)
	return nil
}
