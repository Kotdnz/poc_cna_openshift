FROM nginx:alpine

#!/bin/sh

COPY ./nginx/nginx.conf /etc/nginx/nginx.conf

## Remove default nginx index page
RUN rm -rf /usr/share/nginx/html/*

COPY ./html/ /usr/share/nginx/html

RUN chgrp -R root /var/cache/nginx /var/run /var/log/nginx && \
    chmod -R 770 /var/cache/nginx /var/run /var/log/nginx

EXPOSE 3003 8083

ENTRYPOINT ["nginx", "-g", "daemon off;"]

# docker build -t cna-prod-fe:0.1 .
# docker run --publish 3003:8083 --detach cna-prod-fe:0.1