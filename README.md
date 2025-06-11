# Compartment

Compartment is your assistant for spinning up services needed for local development in a Docker-based environment. It allows you to run multiple instances of services in parallel without port conflicts, so you can work on several projects simultaneously.

## Managing Services

### Launching a New Service

To start a service, use:

```
compartment [flags] start <service> [version]
```

Example:

```
compartment start postgresql 17
compartment start redis 8
```

By default, the container is named after the service and its version (e.g., `postgresql17` or `redis8`). You can specify a custom name using the `-n` option (before the command):

```
compartment -n myapp.search start elasticsearch 13
```

To provide environment variables, use the `-e` flag (before the command):

```
compartment -n postgres17 -e POSTGRES_PASSWORD=postgres start postgresql 17
```

### Stopping an Existing Service

To stop a service, use the `stop` command:

```
compartment stop postgresql 17
```

If you provided a custom name for the container, specify it with the `-n` option (before the command):

```
compartment -n myapp.search stop postgresql 17
```

### Getting the Status of a Service

To get information about a service, use the `status` command:

```
compartment status postgresql 17
```

## Accessing a Service

Compartment pairs best with [devdns](https://github.com/ruudud/devdns), and a custom resolver to access containers via their names.

```
docker run -d --rm --name devdns -e DNS_DOMAIN=containers -p 53:53/udp -v /var/run/docker.sock:/var/run/docker.sock:ro ruudud/devdns
```

### Resolver

The custom resolver allows you to route all `.containers` requests to your DNS server running inside the `devdns` container.

#### macOS

Create a file at `/etc/resolver/containers` with the following content:

```
nameserver 127.0.0.1
```

On macOS, you also need [docker-mac-net-connect](https://github.com/chipmk/docker-mac-net-connect) to be able to access containers directly via their IPs.

## Testing

If you have configured everything correctly, you should be able to start the service and access it. For example:

```
compartment start postgresql
psql -U postgres -h postgres.containers
```
