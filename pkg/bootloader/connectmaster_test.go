package bootloader

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"testing"
)

func TestNetwork(t *testing.T) {
	ips, err := net.LookupIP("ghostnetroot.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ips)
	fmt.Println("-------------")
	url := "https://api.ipify.org?format=text"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	ip, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(ip))
	fmt.Println("-------------")
	addrs, _ := net.InterfaceAddrs()
	fmt.Printf("%v\n", addrs)
	for _, addr := range addrs {
		fmt.Println(addr.String())
	}
}
