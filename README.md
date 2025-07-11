# Compartment

Compartment is your assistant for spinning up local services needed for development in a Docker-based environment. It allows you to run multiple instances of services in parallel without port conflicts, so you can work on several projects simultaneously.

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
compartment -n myapp.cache start redis 8
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
compartment -n myapp.cache stop redis 8
```

### Getting the Status of a Service

To get information about a service, use the `status` command:

```
compartment status postgresql 17
```

## Accessing a Service

Compartment pairs best with [devdns](https://github.com/ruudud/devdns), and a custom resolver to access containers via their names.

**What is devdns?**

> devdns automatically creates DNS records for your running Docker containers, so you can access them as `<container-name>.container` from your host machine.

To start it, run:

```
compartment start devdns
```

### Resolver

The custom resolver allows you to route all `.container` requests to your DNS server running inside the `devdns` container.

#### macOS

Create a file at `/etc/resolver/container` with the following content:

```
nameserver 127.0.0.1
```

On macOS, you also need [docker-mac-net-connect](https://github.com/chipmk/docker-mac-net-connect) to be able to access containers directly via their IPs.

**Important:** You must disable "Resource Saver" in Docker Desktop for Mac for docker-mac-net-connect to work properly. Go to Docker Desktop Settings → Resources → Advanced and uncheck "Resource Saver". See [this issue](https://github.com/chipmk/docker-mac-net-connect/issues/36) for details.

**Quickstart for docker-mac-net-connect:**

```sh
brew install chipmk/tap/docker-mac-net-connect
sudo docker-mac-net-connect
```

Or refer to their [installation guide](https://github.com/chipmk/docker-mac-net-connect#installation).

#### Linux

If you are using Linux, you can usually access containers directly via their IPs without extra setup.

## Testing

If you have configured everything correctly, you should be able to start the service and access it. For example:

```
compartment start postgresql
psql -U postgres -h postgres.container
```

### Verifying DNS Resolution

To test that DNS is working correctly:

```bash
# This should return your container's IP address
host postgres.container

# Test with dig as well
dig postgres.container
```

## Checking Configuration

To check if your development environment is properly configured, use the `check` command:

```
compartment check
```

This command verifies your development environment setup by checking:

- Docker connectivity
- Whether the devdns container is running
- Direct access to containers via their IPs
