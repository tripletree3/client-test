package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	resp, err := http.Get("http://172.20.3.34:8501/health")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	u, err := url.Parse("http://admin@172.20.3.34:8501/health?abc=2&ewwe=3")
	fmt.Println("scheme=", u.Scheme, "opaque=", u.Opaque, "user=", u.User, "host=", u.Host, "path=", u.Path,
		"rawpath=", u.RawPath, "forecequery=", u.ForceQuery, "rawquery=", u.RawQuery, "fragment=", u.Fragment)
}
