# Qdrant init job

This job creates collection in Qdrant if it does not exist.

## Environment variables

- `QDRANT_ADDR` - Qdrant host (example: `http://qdrant:6333`)
- `QDRANT_INSECURE` - if true request sends without TLS verify
- `COLLECTION_NAME` - collection name
- `VECTORS_SIZE` - vector size
- `VECTORS_DISTANCE` - distance metric (for example `Cosine`, `Dot`, `Euclid`, `Manhattan`)

## Use with swarm-deploy

```yaml
services:
  my-super-service:
    image: myorg/my-super-service
    x-init-deploy-jobs:
      - name: create-qdrant-collection
        image: swarmdeployorg/init-jobs-qdrant:v0.1.0
        env:
          - QDRANT_ADDR=http://qdrant:6333
          - COLLECTION_NAME=my-super-service
          - VECTORS_SIZE=1536
          - VECTORS_DISTANCE=Cosine
```
