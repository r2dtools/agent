# SSLBot – Server Agent for SSLPanel

**SSLBot** is a lightweight server agent developed by [R2DTools](https://github.com/r2dtools) that works seamlessly with [SSLPanel](https://github.com/r2dtools/sslpanel). It automates issuing, installing, and renewing SSL/TLS certificates — making it simple to secure your web domains via a user-friendly UI.

## 🔒 Features

- One-click SSL/TLS issuance and renewal
- Let`sEncrypt certificate automation
- Integration with Nginx and Apache
- Lightweight agent with CLI interface
- Works with the SSLPanel UI

---

## 🖥 Supported Web Servers

- **Nginx**
- **Apache**

---

## 🚀 Installation

1. **Connect to your server via SSH.**

2. **Download the latest SSLBot installer:**
   ```bash
   wget https://github.com/r2dtools/installer/releases/latest/download/installer -O /tmp/installer
   ```

3. **Make the installer executable:**
   ```bash
   chmod +x /tmp/installer
   ```

4. **Run the installer:**
   ```bash
   /tmp/installer install
   ```

5. **SSLBot will be installed in:**
   ```
   /opt/r2dtools
   ```

6. **Check if the SSLBot service is running:**
   ```bash
   systemctl status sslbot.service
   ```

7. **Ensure port `60150` is open (default):**
   - This is required for communication with SSLPanel.
   - You can change the port in:
     ```
     /opt/r2dtools/config/params.yaml
     ```
   - Restart the service after changing the port:
     ```bash
     systemctl restart sslbot.service
     ```

---

## 🔑 Connecting SSLBot to SSLPanel

Generate a connection token:
```bash
/opt/r2dtools/sslbot generate-token
```

To view the token:
```bash
/opt/r2dtools/sslbot show-token
```

> 🔍 The token is also stored in:
> ```
> /opt/r2dtools/config/params.yaml
> ```

---

## ⚙️ SSLBot CLI Usage

| Task | Command |
|------|---------|
| **Issue a Let's Encrypt certificate** | <pre>/opt/r2dtools/sslbot issue-cert \<br>  --email your@email.com \<br>  --domain example.com \<br>  --alias www.example.com \<br>  --webserver nginx</pre> |
| **Generate SSLPanel token** | ```/opt/r2dtools/sslbot generate-token``` |
| **Show existing token** | ```/opt/r2dtools/sslbot show-token``` |
| **Deploy an existing certificate** | <pre>/opt/r2dtools/sslbot deploy-cert \<br>  --domain example.com \<br>  --cert /path/to/cert.pem \<br>  --key /path/to/key.pem \<br>  --webserver nginx</pre> |
| **List configured domains** | ```/opt/r2dtools/sslbot hosts``` |
| **Manage ACME challenge directory** | <pre>/opt/r2dtools/sslbot common-dir \<br>  --domain example.com \<br>  --enable \<br>  --webserver apache</pre> |
| **Run SSLBot service manually** | ```/opt/r2dtools/sslbot serve``` |
| **Show help for all commands** | ```/opt/r2dtools/sslbot --help``` |

---

## 📂 Default Installation Path

```
/opt/r2dtools
```

---

## 🛠 Troubleshooting

- Ensure `systemctl status sslbot.service` shows the service is **active**.
- Make sure port `60150` is **open** and **not blocked by firewall rules**.
- If you change the port or any config, remember to restart:
  ```bash
  systemctl restart sslbot.service
  ```

---

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

---

## 👥 Community & Support

- Join the project on [GitHub](https://github.com/r2dtools/sslbot)
- Report issues or request features via [GitHub Issues](https://github.com/r2dtools/sslbot/issues)

---

Secure your web server today with SSLBot + SSLPanel. Easy. Automated. Free.
