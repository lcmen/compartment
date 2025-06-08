# Compartment

Compartment is your assistant for spinning up services needed for local development in a Docker-based environment. It allows you to run multiple instances of services in parallel without port conflicts, so you can work on several projects simultaneously.

## Launching a Service

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

## Stopping a Service

To stop a service, use the `stop` command:

```
compartment stop postgresql 17
```

If you provided a custom name for the container, specify it with the `-n` option (before the command):

```
compartment -n myapp.search stop postgresql 17
```

If more than one instance of the service with provided version is running and you do not specify a name, Compartment will prompt you to choose which container to stop.
