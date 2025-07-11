FROM golang:1.24

RUN apt update && apt install -y \
    nginx

COPY test/nginx/nginxconfig.io /etc/nginx/nginxconfig.io
COPY test/nginx/sites-available /etc/nginx/sites-available

COPY test/nginx/fastcgi_params /etc/nginx/
COPY test/nginx/mime.types /etc/nginx/
COPY test/nginx/nginx.conf /etc/nginx/

RUN mkdir -p /usr/local/r2dtools/var/certificates
COPY test/certificate /usr/local/r2dtools/var/certificates

RUN mkdir /opt/r2dtools
VOLUME /opt/r2dtools
WORKDIR  /opt/r2dtools

RUN ln -s /etc/nginx/sites-available/example.com.conf /etc/nginx/sites-enabled/example.com.conf && \
    ln -s /etc/nginx/sites-available/example2.com.conf /etc/nginx/sites-enabled/example2.com.conf && \
    ln -s /etc/nginx/sites-available/example3.com.conf /etc/nginx/sites-enabled/example3.com.conf && \
    ln -s /etc/nginx/sites-available/example4.com.conf /etc/nginx/sites-enabled/example4.com.conf && \
    ln -s /etc/nginx/sites-available/webmail.conf /etc/nginx/sites-enabled/webmail.conf

ENTRYPOINT ["/bin/sh", "./script/testcmd.sh"]
