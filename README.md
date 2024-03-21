# go-main-autoconf - Simple mail client configuration

This is a simple mail client configuration server that provides auto-configuration for mail clients. It is written in Go and uses the `gofiber` package to serve the configuration files.

It provides mailbox configuration for Apple's Mobileconfig, Microsoft's Autodiscover & Mozilla's Thunderbird.

> **Note:** This is a work in progress and is not yet ready for production use.  
> And current is it only tested with Thunderbird.

## Installation
Fetch latest binaries for the proper OS from the [releases](https://github.com/kim-groenborg-nielsen/go-mail-autoconf/releases) page.

### On Ubuntu 22.04
Copy the binary to `/usr/local/bin`.

Create a user for the service:
```bash
adduser --no-create-home --no-shell --disabled-login --group go-mail-autoconf go-mail-autoconf
```

Create a systemd service file `/etc/systemd/system/go-mail-autoconf.service`:
```ini
[Unit]
Description=Go instance to server mail auto configuration
After=network.target

[Service]
Type=simple
EnvironmentFile=/etc/default/go-mail-autoconf
User=go-mail-autoconf
Group=nogroup
ExecStart=/usr/local/bin/go-mail-autoconf
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

Create a configuration file `/etc/default/go-mail-autoconf` eg.:
```bash
IMAP_SERVER=mail.dummy.org
SMTP_SERVER=mail.dummy.org
SERVER_URL=127.0.0.1:30090
```

Reload systemd and start the service:
```bash
systemctl daemon-reload
systemctl enable go-mail-autoconf
systemctl start go-mail-autoconf
```

Nginx site configuration '/etc/nginx/sites-available/go-mail-autoconf' eg.:
```nginx
server {
        listen 80;
        server_name autoconfig.<FQDN> autoconfig.* maildiscovery.* autodiscovery.* autodiscover.*;
        location ~* ^/(autodiscover/autodiscover|mail/config|email\.mobileconfig) {
                proxy_pass http://127.0.0.1:30090;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header Host $http_host;
                proxy_set_header X-Forwarded-Proto $scheme;
                proxy_buffering off;
        }

        location / {
                return 404;
        }
}

server {
        server_name autoconfig.network-it.dk autoconfig.* maildiscovery.* autodiscovery.* autodiscover.*;
        location ~* ^/(autodiscover/autodiscover|mail/config|email\.mobileconfig) {
                proxy_pass http://127.0.0.1:30090;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header Host $http_host;
                proxy_set_header X-Forwarded-Proto $scheme;
                proxy_buffering off;
        }

        location / {
                return 404;
        }

    listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/autoconfig.<FQDN>/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/autoconfig.<FQDN>/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot
}
```
The `<FQDN>` should be replaced with the fully qualified domain. And certbot is used to manage the SSL certificates.

## Environment variables
The following environment variables are used to configure the server:

| Variable              | Description                                        | Default        |
|-----------------------|----------------------------------------------------|----------------|
| IMAP_SERVER           | Hostname of the IMAP server                        | imap.dummy.org |
| IMAP_PORT             | Port of the IMAP server                            | 993            |
| IMAP_SSL              | Use SSL for IMAP                                   | true           |
| POP3_SERVER           | Hostname of the POP3 server                        |                |
| POP3_PORT             | Port of the POP3 server                            | 995            |
| POP3_SSL              | Use SSL for POP3                                   | true           |
| SMTP_SERVER           | Hostname of the SMTP server                        | smtp.dummy.org |
| SMTP_PORT             | Port of the SMTP server                            | 587            |
| SMTP_SSL              | Use SSL for SMTP                                   | false          |
| SMTP_GLOBAL_PREFERRED | Use the SMTP server as the global preferred server | false          |
| EMAIL_PROVIDER        | Domain of the email provider                       | dummy.org      |
| SERVER_URL            | URL of the server                                  | 127.0.0.1:3000 |


## References
This code is inspired by the following article https://pipo.blog/articles/20210826-email-autoconf and GitLab repository https://gitlab.com/onlime/email-autoconf by Philip Iezzi.
