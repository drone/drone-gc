![build status](https://beta.drone.io/api/badges/drone/drone-gc/status.svg)

__drone-gc__ is a simple daemon that periodically removes unused docker resources. The garbage collector is optimized for continuous integration systems. It uses an lrfu algorithm to control the size of your image cache, while retaining the most frequently used images.

Installation:

```
docker run -d \
  --volume=/var/run/docker.sock:/var/run/docker.sock \
  --env=GC_DEBUG=true \
  --env=GC_CACHE=5gb \
  --env=GC_INTERVAL=5m \
  --restart=always \
  --name=gc drone/gc
```

Configuration:

<dl>
<dt><code>GC_DEBUG</code></dt>
<dd>Enable debug mode</dd>

<dt><code>GC_DEBUG_PRETTY=false</code></dt>
<dd>Pretty print the logs</dd>

<dt><code>GC_DEBUG_COLOR=false</code></dt>
<dd>Pretty print the logs with color</dd>

<dt><code>GC_IGNORE_IMAGES</code></dt>
<dd>Comma-separated list of images to ignore. Supports globbing.</dd>

<dt><code>GC_IGNORE_CONTAINERS</code></dt>
<dd>Comma-separate list of container names to ignore. Support globbing.</dd>

<dt><code>GC_INTERVAL=5m</code></dt>
<dd>Interval at which the garbage collector is executed</dd>

<dt><code>GC_CACHE=5gb</code></dt>
<dd>Maximum image cache size</dd>
</dl>

__Need help?__ Please post questions or comments to our [community forum](https://discourse.drone.io/).
