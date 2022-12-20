package matrix

const homeserverSynapseDefaults = `server_name: "<server_name>"
pid_file: /data/homeserver.pid
listeners:
  - port: 8008
    tls: false
    type: http
    x_forwarded: true
    resources:
      - names: [client, federation]
        compress: false
log_config: "/data/<server_name>.log.config"
media_store_path: /data/media_store
registration_shared_secret: "pU+G+Q@qHtlPf;hG3.W@Rz;G=K6d+UpA&KjnVIjhOAM^w2oZXo"
report_stats: true
macaroon_secret_key: "sv;qbjyug*z;4V4dZ7o#57d_i9i@XQ&b3=FGXWCv-lc8Ya~9e*"
form_secret: "DlSe=VIx0UuG=btT_Qyt6Yv-H@UX=ktUNvb3;:KKz_SNApp9m#"
signing_key_path: "/data/<server_name>.signing.key"
trusted_key_servers:
  - server_name: "matrix.org"
app_service_config_files:
  - /data/registration.yaml
retention:
  enabled: true
enable_registration: true
enable_registration_without_verification: true
suppress_key_server_warning: true
federation_ip_range_blacklist:
  - '127.0.0.0/8'
  - '10.0.0.0/8'
  - '172.16.0.0/12'
  - '192.168.0.0/16'
  - '100.64.0.0/10'
  - '169.254.0.0/16'
  - '::1/128'
  - 'fe80::/64'
  - 'fc00::/7'
database:
  name: psycopg2
  args:
    user: synapse
    password: itsasecret
    database: synapse
    host: postgres
    cp_min: 5
    cp_max: 10
redis:
  enabled: true
  host: redis
  port: 6379
rc_message:
  per_second: 25
  burst_count: 500
rc_registration:
  per_second: 25
  burst_count: 50
rc_registration_token_validity:
  per_second: 10
  burst_count: 50
rc_login:
  address:
    per_second: 0.17
    burst_count: 3
  account:
    per_second: 0.17
    burst_count: 3
  failed_attempts:
    per_second: 0.17
    burst_count: 3
rc_admin_redaction:
  per_second: 10
  burst_count: 100
rc_joins:
  local:
    per_second: 10
    burst_count: 100
  remote:
    per_second: 10
    burst_count: 100
rc_3pid_validation:
  per_second: 0.003
  burst_count: 5
rc_invites:
  per_room:
    per_second: 10
    burst_count: 100
  per_user:
    per_second: 5
    burst_count: 50
rc_federation:
  window_size: 10000
  sleep_limit: 10
  sleep_delay: 500
  reject_limit: 500
  concurrent: 30
federation_rr_transactions_per_room_per_second: 100
`
