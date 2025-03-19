mkdir -m 0755 -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
apt-get update
apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

docker network create caddy

apt-get install -y postgresql-client-16
apt-get install -y curl ca-certificates
install -d /usr/share/postgresql-common/pgdg
curl -o /usr/share/postgresql-common/pgdg/apt.postgresql.org.asc --fail https://www.postgresql.org/media/keys/ACCC4CF8.asc

# Create the repository configuration file:
sh -c 'echo "deb [signed-by=/usr/share/postgresql-common/pgdg/apt.postgresql.org.asc] https://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'

# Postgres SSL setup
mkdir -p /root/dev/caddy/certs
openssl req -x509 -newkey rsa:4096 -keyout /root/dev/caddy/certs/server.key -out /root/dev/caddy/certs/server.crt -days 365 -nodes -subj "/CN=postgres"
chmod 600 /root/dev/caddy/certs/server.key
chmod 644 /root/dev/caddy/certs/server.crt
chown 999 /root/dev/caddy/certs/server.key
chown 999 /root/dev/caddy/certs/server.crt
