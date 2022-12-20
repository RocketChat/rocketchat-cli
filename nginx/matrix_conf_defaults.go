package nginx

const matrixConfDefaults = `server {
    Listen 80;
    Listen 3300;
    server_name <server_name>;
    location /.well-known/Matrix/server {
        access_log off;
        add_header Access-Control-Allow-Origin *;
        default_type application/json;
        return 200 '{"m.homeserver": {"base_url": "https://synapse.<server_name>"}}';
    }
    location /.well-known/Matrix/client {
        access_log off;
        add_header Access-Control-Allow-Origin *;
        default_type application/json;
        return 200 '{"m.homeserver": {"base_url": "https://synapse.<server_name>"}}';
    }
    location / {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_pass http://rocketchat:3000;
        proxy_read_timeout 90;
    }
}`
