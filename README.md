# R2DTools SSLBot
Simplify the maintenance of your websites

## Secure your website with SSL/TLS certificate

R2DTools makes it possible to issue and deploy Let`s Encrypt certificate for a website in a few clicks via CLI.

If you have an already issued certificate ( with .pem extension ) you can just upload it and R2DTools will secure your website with the uploaded certificate.

## Supported web servers

* Nginx

## How to install

* Connect to the server via ssh
* Download the latest version of the SSLBot installer:
  ```bash 
  wget https://github.com/r2dtools/installer/releases/latest/download/installer
  ```
* Make the installer executable:
  ```bash
  chmod +x /tmp/installer
  ```
* Install SSLBot:
  ```bash
  /tmp/installer install
  ```
* The SSLBot will be installed in the <strong>/opt/r2dtools</strong> directory

## How to use

* Secure domain with Let`s Encrypt certificate
  ```bash
  ./r2dtools issue-cert --email example@gamil.com --domain example.com --alias www.example.com
  ```
  Use help command for more information
  ```bash
  ./r2dtools --help
  ```