# Compartment Roadmap

## Phase 1 — Bug Fixes

- [ ] 1. Fix `-e` flag not being passed through (`getServiceForCommand` passes `nil` instead of `envs` to `NewService`)
- [ ] 2. Fix `ExistingContainer` always setting `StateRunning` regardless of actual container state
- [ ] 3. Fix `NewContainer` not capturing the image from inspection (inconsistency with `ExistingContainer`)
- [ ] 4. Add error handling for `os.MkdirAll` calls in `preparePostgresDataDir`, `prepareRedisDataDir`, and `config.go`
- [ ] 5. Add "already running" message when `Start()` is called on a container that's already up
- [ ] 6. Document the `Start()` behavior of removing stopped containers and recreating them

## Phase 2 — Tests

- [ ] 7. Extract a Docker client interface from `Container` to enable mocking
- [ ] 8. Add unit tests for `pkg/service` (service creation, env merging, volume preparation)
- [ ] 9. Add unit tests for `pkg/container` (state transitions, create/start/stop/remove logic)
- [ ] 10. Add unit tests for `cmd` (argument parsing, command routing)
- [ ] 11. Add unit tests for `pkg/check` (check flow logic)

## Phase 3 — Container Labeling & List Command

- [ ] 12. Add a `compartment=true` label (and service type/version labels) to containers at creation time
- [ ] 13. Implement a `list` command that queries Docker for containers with the compartment label
- [ ] 14. Update `cmd/compartment.go` to register the `list` command

## Phase 4 — Health Check on Start

- [ ] 15. Accept health check config in `Service` struct (command, interval, retries, timeout)
- [ ] 16. Configure default health checks per service type (e.g., `pg_isready` for postgres, `redis-cli ping` for redis)
- [ ] 17. Make `Start()` poll container health status and wait until healthy (with timeout)

## Phase 5 — Connect Command

- [ ] 18. Add a `ServiceType` field to `Service` struct (and persist it as a container label) so we know what client to use
- [ ] 19. Implement `connect` command that runs `docker exec -it <container> <client>` — `psql -U postgres` for postgres, `redis-cli` for redis, `mysql -u root` for mysql
- [ ] 20. Update `cmd/compartment.go` to register the `connect` command and help text

## Phase 6 — New Services

- [ ] 21. Add MySQL service (image, default env, volumes, health check, connect client)
- [ ] 22. Add MongoDB service (with `mongosh` for connect)
- [ ] 23. Add Elasticsearch service
- [ ] 24. Update the service factory switch and help text for all new services

## Phase 7 — Documentation

- [ ] 25. Update README with `list`, `connect` commands and new services
- [ ] 26. Update CLI help text to reflect all new commands and services
