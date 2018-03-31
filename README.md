drone-gc is a simple daemon that periodically removes unused docker resources. The garbage collector is optimized for continuous integration systems.

Download the docker image:

```
docker pull drone/gc
```

Start the garbage collector:

```
docker run -d \
  --volume=/var/lib/docker.sock:/var/lib/docker.sock \
  --env=GC_DEBUG=true \
  --env=GC_CACHE=5gb \
  --env=GC_INTERVAL=5m \
  --restart=always \
  --name=gc drone-gc
```

Configuration Parameters:


GC_DEBUG
: Enable debug mode

GC_DEBUG_PRETTY=false
: Pretty print the logs

GC_DEBUG_COLOR=false
: Pretty print the logs with color

GC_IGNORE_IMAGES
: Comma-separated list of images to ignore. Supports globbing.

GC_IGNORE_CONTAINERS
: Comma-separate list of container names to ignore. Support globbing.

GC_INTERVAL=5m
: Interval at which the garbage collector is executed

GC_CACHE=5gb
: Maximum image cache size
