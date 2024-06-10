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