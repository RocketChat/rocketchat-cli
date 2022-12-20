package docker

const composeDefaults = `services:
  traefik:
    image: "traefik:<traefik_tag>"
    restart: "unless-stopped"
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock:ro"
      - "./data/traefik/traefik.yml:/etc/traefik/traefik.yml:ro"
      - "./data/traefik/config:/config:ro"
      - "./data/traefik/acme.json:/acme.json"
    labels:
      - "traefik.enable=true"
      - "traefik.http.services.traefik.loadbalancer.server.port=8080"
      - "traefik.http.routers.traefik.rule=Host(` + "`traefik.$SERVER_HOSTNAME`" + `)"
      - "traefik.http.routers.traefik.entrypoints=web-secure"
      - "traefik.http.routers.traefik.tls.certresolver=letsencrypt"
    networks:
      - internal

  postgres:
    image: "postgres:14"
    restart: "unless-stopped"
    environment:
      POSTGRES_PASSWORD: itsasecret
      POSTGRES_USER: synapse
      POSTGRES_DB: synapse
      POSTGRES_INITDB_ARGS: "--encoding='UTF8' --lc-collate='C' --lc-ctype='C'"
    volumes:
      - "./data/postgres/data:/var/lib/postgresql/data"
    networks:
      - internal

  redis:
    image: "redis:<redis_tag>"
    restart: "unless-stopped"
    networks:
      - internal

  synapse:
    image: "matrixdotorg/synapse:<synapse_tag>"
    restart: "unless-stopped"
    environment:
      SYNAPSE_CONFIG_DIR: "/data"
      SYNAPSE_CONFIG_PATH: "/data/homeserver.yaml"
      UID: "1000"
      GID: "1000"
      TZ: "America/New_York"
    volumes:
      - "./data/matrix/synapse:/data"
    ports:
      - 8008:8008
      - 8448:8448
    labels:
      - "traefik.enable=true"
      - "traefik.http.services.synapse.loadbalancer.server.port=8008"
      - "traefik.http.routers.synapse.rule=Host(` + "`synapse.$SERVER_HOSTNAME`" + `)"
      - "traefik.http.routers.synapse.entrypoints=web-secure"
      - "traefik.http.routers.synapse.tls=true"
      - "traefik.http.routers.synapse.tls.certresolver=letsencrypt"
    networks:
      - internal

  nginx:
    image: "nginx:<nginx_tag>"
    restart: "unless-stopped"
    volumes:
      - "./data/matrix/nginx/matrix.conf:/etc/nginx/conf.d/matrix.conf"
    labels:
      - "traefik.enable=true"
      - "traefik.http.services.nginx.loadbalancer.server.port=80"
      - "traefik.http.routers.nginx.rule=Host(` + "`$SERVER_HOSTNAME`" + `)"
      - "traefik.http.routers.nginx.entrypoints=web-secure"
      - "traefik.http.routers.nginx.tls=true"
      - "traefik.http.routers.nginx.tls.certresolver=letsencrypt"
    networks:
      - internal

  element:
    image: vectorim/element-web:<element_tag>
    restart: unless-stopped
    volumes:
      - ./data/element/config.json:/app/config.json
    labels:
      - "traefik.enable=true"
      - "traefik.http.services.element.loadbalancer.server.port=80"
      - "traefik.http.routers.element.rule=Host(` + "`element.$SERVER_HOSTNAME`" + `)"
      - "traefik.http.routers.element.entrypoints=web-secure"
      - "traefik.http.routers.element.tls=true"
      - "traefik.http.routers.element.tls.certresolver=letsencrypt"
    networks:
      - internal

  rocketchat:
    image: "rocketchat:<rocketchat_tag>"
    command: >
      bash -c
        "for i in ` + "`seq 1 30`" + `; do
          node main.js &&
          s=$$? && break || s=$$?;
          echo \"Tried $$i times. Waiting 5 secs...\";
          sleep 5;
        done; (exit $$s)"
    restart: unless-stopped
    volumes:
      - ./uploads:/app/uploads
      - ./data/matrix/synapse/registration.yaml:/app/matrix-federation-config/registration.yaml
    environment:
      - PORT=3000
      - ROOT_URL=https://localhost:3000
      - MONGO_URL=mongodb://mongodb:27017/rocketchat?replicaSet=rs0&directConnection=true
      - MONGO_OPLOG_URL=mongodb://mongodb:27017/local?replicaSet=rs0&directConnection=true
      - NODE_ENV=production
      # - REG_TOKEN=${REG_TOKEN}
    depends_on:
      - mongodb
    ports:
      - 3000:3000
      - 3300:3300
    networks:
      - internal
    # labels:
    #   - "traefik.enable=true"
    #   - "traefik.http.services.rc.loadbalancer.server.port=3000"
    #   - "traefik.http.routers.rc.rule=Host(` + "`$SERVER_HOSTNAME`" + `)"
    #   - "traefik.http.routers.rc.entrypoints=web-secure"
    #   - "traefik.http.routers.rc.tls=true"
    #   - "traefik.http.routers.rc.tls.certresolver=letsencrypt"
    #   - "traefik.http.services.bridge.loadbalancer.server.port=$RC_BRIDGE_PORT"
    #   - "traefik.http.routers.bridge.rule=Host(` + "`$SERVER_HOSTNAME`" + `)"
    #   - "traefik.http.routers.bridge.entrypoints=web-secure"
    #   - "traefik.http.routers.bridge.tls=true"
    #   - "traefik.http.routers.bridge.tls.certresolver=letsencrypt"

  mongodb:
    image: mongo:5.0
    restart: unless-stopped
    volumes:
     - ./data/mongo/db:/data/db
    command: mongod --oplogSize 128 --replSet rs0
    labels:
      - "traefik.enable=false"
    networks:
      - internal

  # this container's job is just run the command to initialize the replica set.
  # it will run the command and remove himself (it will not stay running)
  mongo-init-replica:
    image: mongo:5.0
    command: >
      bash -c
        "for i in ` + "`seq 1 30`" + `; do
          mongo mongodb/rocketchat --eval \"
            rs.initiate({
              _id: 'rs0',
              members: [ { _id: 0, host: 'localhost:27017' } ]})\" &&
          s=$$? && break || s=$$?;
          echo \"Tried $$i times. Waiting 5 secs...\";
          sleep 5;
        done; (exit $$s)"
    depends_on:
      - mongodb
    networks:
      - internal

networks:
  internal:
    attachable: true`
