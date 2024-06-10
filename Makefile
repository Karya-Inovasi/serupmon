# Variabel untuk nama aplikasi, versi, dan direktori kerja
APP_NAME = serupmon
VERSION = 1.1.0
WORK_DIR = $(APP_NAME)-$(VERSION)

# Direktori aplikasi dan sumber konfigurasi
APP_SRC = /build/bin/serupmon
CONFIG_SRC = /build/stubs/etc/serupmon/serupmon.hcl
SERVICE_SRC = /build/stubs/etc/systemd/system/serupmon.service

# Direktori target untuk file biner, layanan, dan paket Debian
BIN_DIR = $(WORK_DIR)/usr/local/bin
SERVICE_DIR = $(WORK_DIR)/etc/systemd/system
CONFIG_DIR = $(WORK_DIR)/etc/serupmon
DEBIAN_DIR = $(WORK_DIR)/DEBIAN

# Path ke file kontrol dan file pasang (post-install script)
CONTROL_FILE = $(DEBIAN_DIR)/control
POSTINST_FILE = $(DEBIAN_DIR)/postinst

# Target default
all: build

# Target untuk membangun paket Debian
build: $(WORK_DIR).deb

# Target untuk membersihkan hasil build
clean:
	rm -rf $(WORK_DIR) $(WORK_DIR).deb

# Target untuk menginstal paket Debian
install: $(WORK_DIR).deb
	sudo dpkg -i $(WORK_DIR).deb

# Target untuk membuat paket Debian
$(WORK_DIR).deb: $(BIN_DIR)/$(APP_NAME) $(SERVICE_DIR)/$(APP_NAME).service $(CONTROL_FILE) $(POSTINST_FILE) $(CONFIG_DIR)/serupmon.hcl
	dpkg-deb --build $(WORK_DIR)


# Target untuk menyalin file konfigurasi
$(CONFIG_DIR)/serupmon.hcl: $(CONFIG_SRC)
	mkdir -p $(CONFIG_DIR)
	cp $(CONFIG_SRC) $(CONFIG_DIR)/serupmon.hcl

# Target untuk menyalin file biner aplikasi
$(BIN_DIR)/$(APP_NAME): $(APP_SRC)
	mkdir -p $(BIN_DIR)
	cp $(APP_SRC) $(BIN_DIR)/$(APP_NAME)

# Target untuk menyalin file layanan systemd
$(SERVICE_DIR)/$(APP_NAME).service: $(SERVICE_SRC)
	mkdir -p $(SERVICE_DIR)
	cp $(SERVICE_SRC) $(SERVICE_DIR)/$(APP_NAME).service

# Target file kontrol Debian
$(CONTROL_FILE):
	mkdir -p $(DEBIAN_DIR)
	echo "Package: $(APP_NAME)" > $(CONTROL_FILE)
	echo "Version: $(VERSION)" >> $(CONTROL_FILE)
	echo "Section: utils" >> $(CONTROL_FILE)
	echo "Priority: optional" >> $(CONTROL_FILE)
	echo "Architecture: amd64" >> $(CONTROL_FILE)
	echo "Depends: systemd" >> $(CONTROL_FILE)
	echo "Maintainer: Fathurrohman <fathurrohmanrosyadi@gmail.com>" >> $(CONTROL_FILE)
	echo "Description: Serupmon Monitoring Service" >> $(CONTROL_FILE)
	echo " A simple monitoring service that checks the status of various services and sends alerts." >> $(CONTROL_FILE)

# Target (post-install script)
$(POSTINST_FILE):
	mkdir -p $(DEBIAN_DIR)
	echo "#!/bin/bash" > $(POSTINST_FILE)
	echo "systemctl daemon-reload" >> $(POSTINST_FILE)
	echo "systemctl enable $(APP_NAME)" >> $(POSTINST_FILE)
	echo "systemctl start $(APP_NAME)" >> $(POSTINST_FILE)
	chmod +x $(POSTINST_FILE)
