# postgres init job

This job creates database in Postgres.

Job expects DSN in secret file `/run/secrets/dsn` in formats:
- `postgres://admin:admin@postgres18/my-super-service`
- `host=postgres18 user=admin password=admin dbname=my-super-service sslmode=disable`

## Use with swarm-deploy

```yaml
services:
  my-super-service:
    image: myorg/my-super-service
    x-init-deploy-jobs:
      - name: create-db
        image: swarmdeployorg/init-jobs-postgres:v0.1.0
        secrets:
          - source: my-super-service-pg-dsn
            target: /run/secrets/dsn
```
