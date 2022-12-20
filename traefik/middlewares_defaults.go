package traefik

const middlewaresDefaults = `http:
  middlewares:
    httpsredirect:
      redirectScheme:
        scheme: https
        permanent: true`
