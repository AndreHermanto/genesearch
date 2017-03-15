Create an Elastic Search container

### Caveats:
Elastic search may require the following:

```bash
sysctl -w vm.max_map_count=262144
```

Read more here: https://www.elastic.co/guide/en/elasticsearch/reference/5.x/_maximum_map_count_check.html

### Basic Setup:
create a volume only container so that we can share the data
```bash
docker create -v /usr/share/elasticsearch/data --name elasticdnaindex elasticsearch /bin/true
```

run the container with the volume only container
(ensure you point to a config directory, there is an example on in this dir)
```bash
docker run -d --volumes-from elasticdnaindex --name elasticdna -p 32840:9200 -v "$PWD/config":/usr/share/elasticsearch/config elasticsearch
```

### Deployment tips
The Elastic search API does not support any access control. For internal services
it's probably ok to access it directly but for external services or clients
you probably want to put an application layer in front of it. An easy way to
do this is to setup a reverse proxy Apache Server with SSL and access control.
