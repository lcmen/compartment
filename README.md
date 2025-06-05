# Compartment

Compartment is your assistant for spinning up services needed for local development in a Docker-based environment. It allows you to run multiple instances of services in parallel without port conflicts, so you can work on several projects simultaneously.

## Launching a Service

To start a service, use:

```
compartment start postgresql 17
compartment start redis 8
```

By default, the container is named after the service and its version (e.g., `postgresql17` or `redis8`). You can specify a custom name using the `-n` option:

```
compartment start -n myapp.search elasticsearch 13
```

To provide environment variables, use the same format as Docker:

```
compartment start -n postgres17 -e POSTGRES_PASSWORD=postgres postgresql 17
```

## Stopping a Service

To stop a service, use the `stop` command:

```
compartment stop postgresql 17
```

If you provided a custom name for the container, specify it with the `-n` option:

```
compartment stop -n myapp.search
```

If more than one instance of the service with provided version is running and you do not specify a name, Compartment will prompt you to choose which container to stop.
