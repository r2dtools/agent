FROM golang:1.23

RUN apt update && apt install -y \
    nginx

COPY test/nginx/nginxconfig.io /etc/nginx/nginxconfig.io

COPY test/nginx/sites-available /etc/nginx/sites-available

COPY test/nginx/fastcgi_params /etc/nginx/
COPY test/nginx/mime.types /etc/nginx/
COPY test/nginx/nginx.conf /etc/nginx/

RUN mkdir /opt/r2dtools
VOLUME /opt/r2dtols
WORKDIR  /opt/r2dtools

RUN ln -s /etc/nginx/sites-available/example.com.conf /etc/nginx/sites-enabled/example.com.conf
RUN ln -s /etc/nginx/sites-available/example2.com.conf /etc/nginx/sites-enabled/example2.com.conf
RUN ln -s /etc/nginx/sites-available/example3.com.conf /etc/nginx/sites-enabled/example3.com.conf

ENTRYPOINT ["/bin/sh", "./script/testcmd.sh"]
