package traefik

const routersDefaults = `http:
  routers:
    redirecttohttps:
      entryPoints:
        - "web"
      middlewares:
        - "httpsredirect"
      rule: "HostRegexp(` + "`{host:.+}`" + `)"
      service: "noop@internal"`
