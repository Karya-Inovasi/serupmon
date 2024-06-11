# Serupmon

> A simple monitoring service that checks the status of various services and sends alerts.

### Installlation

Download the latest release from [here](https://github.com/karyainovasiab/serupmon/releases/latest) or use the installation script below.

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/karyainovasiab/serupmon/main/install.sh)"
```

TIP: You can also use the `install.sh` script to update/uninstall serupmon. Just pass the `update` or `uninstall` argument.

```bash
# update
sh -c "$(curl -fsSL https://raw.githubusercontent.com/karyainovasiab/serupmon/main/install.sh)" -- uninstall

# uninstall
sh -c "$(curl -fsSL https://raw.githubusercontent.com/karyainovasiab/serupmon/main/install.sh)" -- update
```

### Docker Image

The Docker image is available on Docker Hub at [fathurrohman26/serupmon](https://hub.docker.com/r/fathurrohman26/serupmon). You can pull the image using the following command:

```bash
docker pull fathurrohman26/serupmon

# run the container
docker run -d --name serupmon -v /path/to/config:/etc/serupmon -v /path/to/prefix:/tmp/serupmon fathurrohman26/serupmon
```

### Usage

Serupmon can be started using the serupmon binary itself or using the service manager (systemd on Linux).

#### Using the binary

```bash
serupmon -h # for help
serupmon start --config /path/to/serupmon.hcl --prefix /path/to/prefix

# example
serupmon start --config /etc/serupmon/serupmon.hcl --prefix /tmp/serupmon
serupmon start --config /home/user/serupmon.hcl --prefix /home/user/serupmon
```

#### Using systemd

Serupmon can be managed using systemd on Linux. The service file is included in the release.

```bash
sudo systemctl <start|stop|restart|status> serupmon

# example
sudo systemctl start serupmon

# enable on boot
sudo systemctl enable serupmon

# view logs with journalctl
journalctl -u serupmon
```

> The service file uses the default configuration file path `/etc/serupmon/serupmon.hcl` and the default prefix `/etc/serupmon`.

> Prefix is the directory where the logs and pid files are stored.

### Configuration

Serupmon uses a configuration file in HCL format. See the example below.

> Documentation for the configuration file is coming soon.

```bash
# example configuration file on Linux
sudo nano /etc/serupmon/serupmon.hcl
```

```hcl
# Global Configuration

global {
	interval = 60
	timeout = 15
	threshold = 3

	log {
		enabled = true
		path    = "/var/log/serupmon.log"
		format  = "json"
		maxsize = 1 # in MB
		mode    = "append"
	}

	alert {
		telegram {
			enabled = true
			config {
				token   = "" # Set Bot Token HERE
				chat_id = "" # Set Id Here
			}
		}

		# Unimplemented
		email {
			enabled = false
			config {
				host     = "smtp.gmail.com"
				port     = 587
				username = "amsms"
				password = "password"
				from     = "notif@kutt.app"
				to       = "allcwf@kutt.app"
				cc       = "a1@me.com,a2@me.com"
			}
		}
	}
}

# Monitors

monitor "server-1" {
	service "http" {
		interval = 5
		upstream = "http://localhost:3000"

		add_header "X-Forwarded-For" {
			value = "Serupmon"
		}

		alert {
			telegram {
				enabled = true
			}
		}
	}
}

# add other monitors here...
```

### Development

Serupmon is written in Go and uses Go modules for dependency management. You can download the source code and build it yourself.

```bash
git clone https://github.com/karyainovasiab/serupmon.git
cd serupmon
go build
```

### License

Serupmon is licensed under the MIT License. See the [LICENSE](https://github.com/karyainovasiab/serupmon/blob/main/LICENSE) file for more information.

### Contributing

This project is intended to be internal use only, but if you have any suggestions or improvements, feel free to open an issue or a pull request. We appreciate your feedback!
