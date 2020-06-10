# www http2 main website
server{
listen 443 ssl http2;
listen [::]:443 ssl http2;
server_name h2.testmyprotocol.com;
access_log /var/www/testmyprotocol.com/logs/access.log;
error_log /var/www/testmyprotocol.com/logs/error.log;
    ssl_certificate /etc/letsencrypt/live/h2.testmyprotocol.com/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/h2.testmyprotocol.com/privkey.pem; # managed by Certbot

location / {
add_header Timing-Allow-Origin *;
add_header Content-Security-Policy "frame-ancestors testmyprotocol.com";
root /var/www/testmyprotocol.com/public/;
index index.html;

#kill cache
add_header Last-Modified $date_gmt;
add_header Cache-Control 'no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0';
if_modified_since off;
expires off;
etag off;

#kill cors
add_header Access-Control-Allow-Origin '*';
}

location ~* \.sh$ {
    deny all;
}

location ~* \.pem$ {
    deny all;
}

location ~* \.txt {
    deny all;
}


location  /php/ {
    root /var/www/testmyprotocol.com/public/;
    include snippets/fastcgi-php.conf;
    fastcgi_pass unix:/run/php/php7.2-fpm.sock;
    include fastcgi_params;
    fastcgi_buffers 8 16k;
    fastcgi_buffer_size 32k;

    client_max_body_size 24M;
    client_body_buffer_size 128k;
}

location  /cloned_sites/php/ {
    root /var/www/testmyprotocol.com/public/;
    include snippets/fastcgi-php.conf;
    fastcgi_pass unix:/run/php/php7.2-fpm.sock;
    include fastcgi_params;
    fastcgi_buffers 8 16k;
    fastcgi_buffer_size 32k;

    client_max_body_size 24M;
    client_body_buffer_size 128k;
}



}


server{
    if ($host = h2.testmyprotocol.com) {
        return 301 https://$host$request_uri;
    } # managed by Certbot



listen 80;
listen [::]:80;
server_name h2.testmyprotocol.com;
    return 404; # managed by Certbot


}
