package main

import (
	"fmt"
	"net/http"
	"flag"
	"encoding/json"
)

func get_subdomains(domain string, api_server_host string, api_server_port string) []string {
	url := fmt.Sprintf("http://%s:%s/subdomains/%s", api_server_host, api_server_port, domain)
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	var subdomains []string
	json.NewDecoder(resp.Body).Decode(&subdomains)
	return subdomains
}

func get_tlds(domain string, api_server_host string, api_server_port string) []string {
	url := fmt.Sprintf("http://%s:%s/tlds/%s", api_server_host, api_server_port, domain)
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	var tlds []string
	json.NewDecoder(resp.Body).Decode(&tlds)
	return tlds
}

func get_all(domain string, api_server_host string, api_server_port string) []map[string]string {
	url := fmt.Sprintf("http://%s:%s/all/%s", api_server_host, api_server_port, domain)
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	var data []map[string]string
	json.NewDecoder(resp.Body).Decode(&data)
	return data
}


func main() {
	domain_sub := flag.String("s", "", "Get subdomains for this value")
	domain_tld := flag.String("t", "", "Get tlds for this value")
	domain_all := flag.String("all", "", "Get all data for this query")

	flag.Parse()

	if (*domain_sub != "") {
		str, _ := json.MarshalIndent(get_subdomains(*domain_sub, "192.168.42.56", "1337"), "", "    ")
		fmt.Println(string(str))
	} else if (*domain_tld != "") {
		str, _ := json.MarshalIndent(get_tlds(*domain_tld, "192.168.42.56", "1337"), "", "    ")
		fmt.Println(string(str))
	} else if (*domain_all != "") {
		str, _ := json.MarshalIndent(get_all(*domain_all, "192.168.42.56", "1337"), "", "    ")
		fmt.Println(string(str))
	}

}











