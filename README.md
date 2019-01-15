# Contacts Book Using Go

## Installation

Build the source for the target machine; e.g.; for linux use: `env GOOS=linux GOARCH=arm go build`. If your build machine is similar to the target machine, you can just use `go build` instead. 

The binary will be named `contact-book-api`.
 
### First time setup (with nginx on CentOS)
1. Move binary to `/usr/share/contact-book-api/` and copy the config.json file to `/usr/share/contact-book-api/config/` folder.
2. Setup a reverse proxy in nginx to point to `http://localhost:8888`. This is the default port, use whatever port you set in `config.json`
3. Remember to add `proxy_set_header Authorization $http_authorization;` in nginx config for this service so that the `Authorization` header is sent to the API for every request.
4. Restart nginx service using `systemctl restart nginx`
5. We will call this service with `contactbookapi`. To create a service configuration, copy the below config to `/etc/systemd/system/contactbookapi.service`
```
[Unit]
Description=Contacts Boook Service
ConditionPathExists=/usr/share/contact-book-api/contact-book-api
After=network.target

[Service]
Type=simple
User=<linux user>
Group=<linux group>
LimitNOFILE=1024

Restart=on-failure
RestartSec=10
startLimitIntervalSec=60

WorkingDirectory=/usr/share/contact-book-api/
ExecStart=/usr/share/contact-book-api/contact-book-api

# make sure log directory exists and owned by the user specified
PermissionsStartOnly=true
ExecStartPre=/bin/mkdir -p /var/log/contactbookapi
ExecStartPre=/bin/chown <linux user>:<linux user> /var/log/contactbookapi
ExecStartPre=/bin/chmod 755 /var/log/contactbookapi
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=contactbookapi

[Install]
WantedBy=multi-user.target
```
6. Enable this service using `systemctl enable contactbookapi`
7. We want the logs to be stored at `/var/log/contactbookapi/contactbookapi.log`, so make contactbookapi dir under /var/log
8. For logging, create a file at `/etc/rsyslog.d/30-contactbookapi.conf` with the following contents:
```
if $programname == 'contactbookapi' or $syslogtag == 'contactbookapi' then /var/log/contactbookapi/contactbookapi.log
& stop
```
9. Start contactbookapi by `systemctl start contactbookapi`.

You will start seeing the logs in the log file.

### Updates
To push an update, just move contact-book-api binary to `/usr/share/contact-book-api/` and restart the service by `systemctl restart contactbookapi`.
