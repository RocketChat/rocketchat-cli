package matrix

const registrationDefaults = `id: <id>
hs_token: <hs_token>
as_token: <as_token>
url: http://172.17.0.1:3300
sender_localpart: rocket.cat
de.sorunome.msc2409.push_ephemeral: true
namespaces:
  users:
    - exclusive: false
      regex: .*
  rooms:
    - exclusive: false
      regex: .*
  aliases:
    - exclusive: false
      regex: .*
rocketchat:
  homeserver_url: http://synapse:8008
  homeserver_domain: <domain>
`
