package main

import (
	"fmt"
	"net/http"
	"flag"
	"encoding/json"
	"os"
	"os/user"
	"io/ioutil"
)

func get_subdomains(domain string) []string {
	config := load_config()
	url := fmt.Sprintf("http://%s:%s/subdomains/%s", config["host"], config["port"], domain)
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	var subdomains []string
	json.NewDecoder(resp.Body).Decode(&subdomains)
	return subdomains
}

func get_tlds(domain string) []string {
	config := load_config()
	url := fmt.Sprintf("http://%s:%s/tlds/%s", config["host"], config["port"], domain)
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	var tlds []string
	json.NewDecoder(resp.Body).Decode(&tlds)
	return tlds
}

func get_all(domain string) []map[string]string {
	config := load_config()
	url := fmt.Sprintf("http://%s:%s/all/%s", config["host"], config["port"], domain)
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	var data []map[string]string
	json.NewDecoder(resp.Body).Decode(&data)
	return data
}

func load_config() map[string]interface{} {
	usr, err := user.Current()

	path := fmt.Sprintf("%s/.crobatrc", usr.HomeDir)
	jsonFile, err := os.Open(path)

	if err != nil {
    	fmt.Println("Unable to load connection details from ~/.crobatrc, did you use --init?")
    	fmt.Println("Error:", err)
		os.Exit(1)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	return result
}

func main() {
	initialize := flag.Bool("init", false, "Initialize config and auth file")
	domain_sub := flag.String("s", "", "Get subdomains for this value")
	domain_tld := flag.String("t", "", "Get tlds for this value")
	domain_all := flag.String("all", "", "Get all data for this query")
	format := flag.String("f", "plain", "Set output format (json/plain)")

	flag.Parse()

	if (*initialize) {
		fmt.Println("Initializing ~/.crobatrc")
		fmt.Println("Warnining: this will overwrite existing data in .crobatrc, use ctrl+c to abort.")
		usr, _ := user.Current()
		path := fmt.Sprintf("%s/.crobatrc", usr.HomeDir)
		config := make(map[string]string)
		var host string
		fmt.Printf("Host: ")
		fmt.Scan(&host)
		var port string
		fmt.Printf("Port: ")
		fmt.Scan(&port)
		// var key string
		// fmt.Printf("Port: ")
		// fmt.Scan(&key)
		config["host"] = host
		config["port"] = port
		// config["key"] = key

		str, _ := json.MarshalIndent(config, "", "  ")
		fmt.Println(string(str))
		ioutil.WriteFile(path, str, 0644)
	}

	if (*domain_sub != "") {
		if (*format == "json") {
			str, _ := json.MarshalIndent(get_subdomains(*domain_sub), "", "    ")
			fmt.Println(string(str))
		} else if (*format == "plain") {
			subdomains := get_subdomains(*domain_sub)
			for i := 0; i < len(subdomains); i++ {
				fmt.Println(subdomains[i])
			}
		}	
	} else if (*domain_tld != "") {
		if (*format == "json") {
			str, _ := json.MarshalIndent(get_tlds(*domain_tld), "", "    ")
			fmt.Println(string(str))
		} else if (*format == "plain") {
			tlds := get_tlds(*domain_tld)
			for i := 0; i < len(tlds); i++ {
				fmt.Println(tlds[i])
			}
		}
	} else if (*domain_all != "") {
		if (*format == "json") {
			str, _ := json.MarshalIndent(get_all(*domain_all), "", "    ")
			fmt.Println(string(str))
		} else if (*format == "plain") {
		
			for _, i := range get_all(*domain_all) {
				fmt.Println(i["name"])
			}
		}
	}
}











