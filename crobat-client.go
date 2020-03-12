package main

import (
	"log"
	"fmt"
	"net/http"
	"flag"
	"encoding/json"
	"os"
	"os/user"
	"io/ioutil"
	"time"
)

func get_subdomains(domain string) ([]string, error) {
	config := load_config()
	url := fmt.Sprintf("http://%s:%s/subdomains/%s", config["host"], config["port"], domain)
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	// defer resp.Body.Close()

	if err != nil {
        log.Fatal("Unable to connect to crobat server.", err)
	}

	var subdomains []string
	json.NewDecoder(resp.Body).Decode(&subdomains)
	return subdomains, err
}

func get_tlds(domain string) ([]string,  error) {
	config := load_config()
	url := fmt.Sprintf("http://%s:%s/tlds/%s", config["host"], config["port"], domain)
	client := http.Client{
    	Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	defer resp.Body.Close()

	if err != nil {
        log.Fatal("Unable to connect to crobat server.", err)
	}

	var tlds []string
	json.NewDecoder(resp.Body).Decode(&tlds)
	return tlds, err
}

func get_all(domain string) ([]map[string]string, error)  {
	config := load_config()
	url := fmt.Sprintf("http://%s:%s/all/%s", config["host"], config["port"], domain)
	client := http.Client{
    	Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	defer resp.Body.Close()

	if err != nil {
        log.Fatal("Unable to connect to crobat server.", err)
	}

	var data []map[string]string
	json.NewDecoder(resp.Body).Decode(&data)
	return data, err
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

		// fmt.Println(string(str))
		err := ioutil.WriteFile(path, str, 0644)
		if (err == nil) {
			fmt.Println("Saved to ~/.crobatrc successfully")
		} else {
			fmt.Println("Saving ~/.crobatrc failed")
			fmt.Println("Error:", err)
		}
	}

	if (*domain_sub != "") {

		data, err := get_subdomains(*domain_sub)
		if err != nil {
        	log.Fatal("Unable to connect to crobat server.", err)
		}

		if (*format == "json") {
			str, _ := json.MarshalIndent(data, "", "    ")
			fmt.Println(string(str))
		} else if (*format == "plain") {
			for i:=0; i < len(data); i++ {
				fmt.Println(data[i])
			}
		}	

	} else if (*domain_tld != "") {

		data, err := get_tlds(*domain_tld)
		if err != nil {
        	log.Fatal("Unable to connect to crobat server.", err)
		}	

		if (*format == "json") {
			str, _ := json.MarshalIndent(data, "", "    ")
			fmt.Println(string(str))
		} else if (*format == "plain") {
			for i:=0; i < len(data); i++ {
				fmt.Println(data[i])
			}
		}

	} else if (*domain_all != "") {

		data, err := get_all(*domain_all)
		if err != nil {
        	log.Fatal("Unable to connect to crobat server.", err)
		}	

		if (*format == "json") {
			str, _ := json.MarshalIndent(data, "", "    ")
			fmt.Println(string(str))
		} else if (*format == "plain") {
			for _, i := range data {
				fmt.Println(i["name"])
			}
		}
	}
}











