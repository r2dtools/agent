# R2DTools agent
Simplify the maintenance of your websites and servers

## Secure your website with SSL/TLS certificate

R2DTools makes it possible to issue a Let`s Encrypt certificate for a website in a few clicks via a user-friendly interface.

If you have an already issued certificate ( with .pem extension ) you can just upload it and R2DTools will secure your website with the uploaded certificate.

## Supported web servers

* Nginx

## How to install

* Connect to the server via ssh
* Download the latest version of the agent installer:
  ```bash 
  wget https://github.com/r2dtools/installer/releases/download/v1.0.0/installer
  ```
* Make the installer executable:
  ```bash
  chmod +x /tmp/installer
  ```
* Install the agent:
  ```bash
  /tmp/installer install
  ```
* The agent will be installed to the <strong>/opt/r2dtools</strong> directory
* Add generated token to the agent configuration file <strong>/opt/r2dtools/config/params.yaml</strong>: Token: <em>token</em>
* Restart the agent:
  ```bash
  systemctl restart r2dtools.service
  ```
