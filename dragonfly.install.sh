#Mac
docker run -p 6379:6379 --ulimit memlock=-1 \
    -v <your-path-to-dragonfly-data>:/data \
    docker.dragonflydb.io/dragonflydb/dragonfly \
    --maxmemory 4g \
    --dir /data \
    --snapshot_cron "*/30 * * * *"

#Linux
docker run -d --network=host -p 6379:6379 --ulimit memlock=-1 \
    --name=dragonfly \
    -v <your-path-to-dragonfly-data>:/data \
    docker.dragonflydb.io/dragonflydb/dragonfly \
    --maxmemory 4g \
    --dir /data \
    --snapshot_cron "0 1 * * *"