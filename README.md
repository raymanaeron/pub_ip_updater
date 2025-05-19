# Public IP Updater for DigitalOcean DNS

A lightweight Go application that automatically updates DigitalOcean DNS records with your current public IP address.

## Overview

This utility solves the problem of maintaining DNS records when your home or office has a dynamic public IP address that changes periodically. The application runs in the background, periodically checking your public IP address and updating a specified DNS record in your DigitalOcean account whenever changes are detected.

## Features

- Automatic detection of your public IP address
- Configurable DNS update for any domain managed in DigitalOcean
- Support for any DNS record type (A, AAAA, etc.)
- Customizable TTL (Time-To-Live) settings
- Error handling with detailed logging
- Simple configuration through a `.env` file

## Requirements

- Go 1.13 or higher
- DigitalOcean account with DNS management
- A domain managed through DigitalOcean's nameservers 

## Configuration

When first run, the application will create a `.env` file with default values:

```
digital_ocean_token=YOUR_DIGITAL_OCEAN_TOKEN
domain_name=aeronlab.net
sub_domain_name=rayman
record_type=A
ttl=3600
update_interval=5
```

Edit this file and replace the values:

- `digital_ocean_token`: Your DigitalOcean API token (create one in the API section of your DigitalOcean account)
- `domain_name`: The domain name you want to update (e.g., `example.com`)
- `sub_domain_name`: The subdomain prefix (e.g., `www` for `www.example.com`)
- `record_type`: The DNS record type (typically `A` for IPv4 addresses)
- `ttl`: Time-To-Live in seconds (3600 = 1 hour)
- `update_interval`: How often to check for IP changes (in minutes)

## Building the Application

### Windows
Run the included batch file:
```
b.bat
```

### macOS and Linux
Run the included shell script:
```
chmod +x b.sh
./b.sh
```

Or build manually:
```
mkdir -p build
go build -o build/ip_updater main.go
```

## Usage

Run the application:

```
# On Windows
.\build\ip_updater.exe

# On macOS and Linux
./build/ip_updater
```

The application will:
1. Check for a `.env` file and create one if it doesn't exist
2. Read the configuration
3. Get your current public IP address
4. Update the specified DNS record in DigitalOcean
5. Wait for the configured update interval and repeat

For continuous operation, you might want to set it up as a service/daemon using systemd or similar.

## Running as a Background Service

### Using systemd (Linux)

Create a systemd service file:

```bash
sudo nano /etc/systemd/system/ip_updater.service
```

Add the following content (adjust paths as needed):

```
[Unit]
Description=Public IP Updater for DigitalOcean DNS
After=network.target

[Service]
ExecStart=/path/to/pub_ip_updater/build/ip_updater
WorkingDirectory=/path/to/pub_ip_updater
Restart=always
User=yourusername

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl enable ip_updater
sudo systemctl start ip_updater
```
