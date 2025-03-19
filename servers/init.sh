mkdir -m 0755 -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
apt-get update
apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

mkdir -p /root/dev/caddy/certs
openssl req -x509 -newkey rsa:4096 -keyout /root/dev/caddy/certs/server.key -out /root/dev/caddy/certs/server.crt -days 365 -nodes -subj "/CN=postgres"
chmod 600 /root/dev/caddy/certs/server.key
chmod 644 /root/dev/caddy/certs/server.crt
