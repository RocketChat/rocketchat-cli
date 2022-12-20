package traefik

const configDefaults = `entryPoints:
  web:
    address: ":80"
  web-secure:
    address: ":443"

api:
  dashboard: true
  insecure: true

providers:
  file:
    directory: "/config"
    watch: true
  docker:
    endpoint: "unix:///var/run/docker.sock"
    network: "federation"
    watch: true
    exposedByDefault: false

certificatesResolvers:
  letsencrypt:
    acme:
      caServer: https://acme-v02.api.letsencrypt.org/directory
      email: "email@address.com"
      storage: "/acme.json"
      httpChallenge:
        entryPoint: "web"`
