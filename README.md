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
compartment start postgres 17
compartment start redis 8
```

By default, the container is named `<version>.<service>` (e.g., `17.postgres` or `8.redis`). You can specify a custom prefix using the `-n` option — the service name is always appended automatically:

```
compartment -n myapp start redis 8
# → container name: myapp.redis
```

To provide environment variables, use the `-e` flag:

```
compartment -e POSTGRES_PASSWORD=secret start postgres 17
```

### Stopping an Existing Service

To stop a service, use the `stop` command:

```
compartment stop postgres 17
```

If you used a custom name prefix, specify it with `-n`:

```
compartment -n myapp stop redis 8
```

### Getting the Status of a Service

To get information about a service, use the `status` command:

```
compartment status postgres 17
```

### Listing All Managed Containers

To see all compartment-managed containers across all projects:

```
compartment list
```

Example output:

```
NAME                 SERVICE      VERSION    STATE
17.postgres          postgres     17         running
8.redis              redis        8          stopped
devdns               devdns       latest     running
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

##### Docker Desktop for Mac

*This is not needed for OrbStack as it provides access to containers via their IPs by default.*

If you use Docker Desktop on macOS, you also need [docker-mac-net-connect](https://github.com/chipmk/docker-mac-net-connect) to be able to access containers directly via their IPs.

**Important:** You must disable "Resource Saver" in Docker Desktop for Mac for docker-mac-net-connect to work properly. Go to Docker Desktop Settings → Resources → Advanced and uncheck "Resource Saver". See [this issue](https://github.com/chipmk/docker-mac-net-connect/issues/36) for details.

**Quickstart for docker-mac-net-connect:**

```sh
brew install chipmk/tap/docker-mac-net-connect
sudo docker-mac-net-connect
```

Or refer to their [installation guide](https://github.com/chipmk/docker-mac-net-connect#installation).

#### Linux

If you are using Linux, you can usually access containers directly via their IPs without extra setup.

Add the following to `/etc/network/interfaces`:

```
auto p3p1
iface p3p1 inet dhcp
dns-search container
dns-nameservers 127.0.0.1
```

## Testing

If you have configured everything correctly, you should be able to start the service and access it. For example:

```
compartment start postgres
psql -U postgres -h 18.postgres.container
```

### Verifying DNS Resolution

To test that DNS is working correctly:

```bash
# This should return your container's IP address
host 18.postgres.container

# Test with dig as well
dig 18.postgres.container
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
