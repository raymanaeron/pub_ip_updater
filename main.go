package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Config stores the application configuration settings
// loaded from the .env file
type Config struct {
	DigitalOceanToken string `json:"digital_ocean_token"` // API token for Digital Ocean authentication
	DomainName       string `json:"domain_name"`         // Domain name to update (e.g. example.com)
	SubDomainName    string `json:"sub_domain_name"`     // Subdomain prefix (e.g. www)
	RecordType       string `json:"record_type"`         // DNS record type (e.g. A, CNAME)
	TTL              int    `json:"ttl"`                 // Time To Live value in seconds
}

// PublicIP represents the structure of the response from the ipify.org API
// containing the public IP address
type PublicIP struct {
	IP string `json:"ip"` // Public IP address as a string
}

// init_config creates a default configuration file (.env) if it doesn't exist
// with placeholder values for Digital Ocean DNS updates
func init_config() {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		file, err := os.Create(".env")
		if err != nil {
			fmt.Println("Error creating .env file:", err)
			return
		}
		defer file.Close()
		file.WriteString("digital_ocean_token=YOUR_DIGITAL_OCEAN_TOKEN\n")
		file.WriteString("domain_name=aeronlab.net\n")
		file.WriteString("sub_domain_name=rayman\n")
		file.WriteString("record_type=A\n")
		file.WriteString("ttl=3600\n")
		fmt.Println(".env file created with placeholder values.")
	}
}

// read_config reads and parses the configuration file (.env)
// and returns a populated Config struct
func read_config() Config {
	file, err := os.Open(".env")
	if err != nil {
		fmt.Println("Error opening .env file:", err)
		return Config{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	config := Config{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || len(strings.TrimSpace(line)) == 0 {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "digital_ocean_token":
			config.DigitalOceanToken = value
		case "domain_name":
			config.DomainName = value
		case "sub_domain_name":
			config.SubDomainName = value
		case "record_type":
			config.RecordType = value
		case "ttl":
			// Parse TTL value from string to integer
			fmt.Sscanf(value, "%d", &config.TTL)
		}
	}
	return config
}

// get_public_ip retrieves the current public IP address of the system
// by making a request to ipify.org API and returns it in a structured format
func get_public_ip() PublicIP {
	// Make HTTP request to ipify.org to get public IP
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		fmt.Println("Error getting public IP:", err)
		return PublicIP{}
	}
	defer resp.Body.Close()
	
	// Read and parse the response body
	body, _ := ioutil.ReadAll(resp.Body)
	var ip PublicIP
	json.Unmarshal(body, &ip)
	return ip
}

// update_do_dns updates a DNS record in Digital Ocean's DNS service
// Parameters:
//   token: Digital Ocean API token for authentication
//   domain_name: The domain to update (e.g., example.com)
//   sub_domain_name: The subdomain to update (e.g., www)
//   record_type: The DNS record type (e.g., A, CNAME)
//   value: The value to set for the DNS record (e.g., IP address)
//   ttl: Time To Live value in seconds
// Returns:
//   bool: true if update was successful, false otherwise
func update_do_dns(token, domain_name, sub_domain_name, record_type, value string, ttl int) bool {
	// Only check if token is completely empty
	if token == "" {
		fmt.Println("Error: Digital Ocean API token is empty.")
		return false
	}
	
	client := &http.Client{}
	
	// Step 1: Get all DNS records to find the specific record ID
	url := fmt.Sprintf("https://api.digitalocean.com/v2/domains/%s/records", domain_name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}
	req.Header.Add("Authorization", "Bearer "+token)
	
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error connecting to Digital Ocean API:", err)
		return false
	}
	defer resp.Body.Close()
	
	// Check for unauthorized or other error responses
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Digital Ocean API error (status %d): %s\n", resp.StatusCode, body)
		if resp.StatusCode == http.StatusUnauthorized {
			fmt.Println("Error: Invalid Digital Ocean API token. Please check your configuration.")
		}
		return false
	}
	
	// Parse the response to extract record information
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return false
	}
	
	var records struct {
		DomainRecords []struct {
			ID   int    `json:"id"`
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"domain_records"`
	}
	
	if err := json.Unmarshal(body, &records); err != nil {
		fmt.Println("Error parsing response:", err)
		return false
	}
	
	// Find the ID of the record that matches our criteria
	var recordID int
	for _, rec := range records.DomainRecords {
		if rec.Type == record_type && rec.Name == sub_domain_name {
			recordID = rec.ID
			break
		}
	}
	
	if recordID == 0 {
		fmt.Println("DNS record not found. Make sure the domain and subdomain exist in your Digital Ocean account.")
		return false
	}
	
	// Update the record
	updateURL := fmt.Sprintf("https://api.digitalocean.com/v2/domains/%s/records/%d", domain_name, recordID)
	updateBody := fmt.Sprintf(`{"data":"%s","ttl":%d}`, value, ttl)
	req, err = http.NewRequest("PUT", updateURL, strings.NewReader(updateBody))
	if err != nil {
		fmt.Println("Error creating update request:", err)
		return false
	}
	
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error during update request:", err)
		return false
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		fmt.Println("DNS record updated successfully.")
		return true
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Failed to update DNS record. Status: %s, Response: %s\n", resp.Status, body)
		return false
	}
}

// main is the entry point of the application
// It runs in an infinite loop, checking and updating the DNS record
// with the current public IP address every 5 minutes
func main() {
	init_config()
	
	for {
		config := read_config()
		
		 // Don't pre-validate the token - let the API tell us if it's invalid
		// Just check if it's completely empty
		if config.DigitalOceanToken == "" {
			fmt.Println("Error: Digital Ocean API token not found in config. Please update your .env file.")
			fmt.Println("Waiting 1 minute before retrying...")
			time.Sleep(1 * time.Minute)
			continue
		}
		
		ip := get_public_ip()
		if ip.IP == "" {
			fmt.Println("Failed to get current public IP. Will retry in 5 minutes.")
			time.Sleep(5 * time.Minute)
			continue
		}
		
		fmt.Println("Current public IP:", ip.IP)
		success := update_do_dns(config.DigitalOceanToken, config.DomainName, config.SubDomainName, config.RecordType, ip.IP, config.TTL)
		
		if success {
			fmt.Println("DNS update successful. Next check in 5 minutes.")
		} else {
			fmt.Println("DNS update failed. Will retry in 5 minutes.")
		}
		
		time.Sleep(5 * time.Minute)
	}
}
