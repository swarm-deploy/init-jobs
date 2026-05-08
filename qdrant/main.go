package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/swarm-deploy/init-jobs/qdrant/internal"
)

var (
	Version   = "v0.1.0"
	BuildTime = "2026-05-08 00:00:00"
)

type Config struct {
	Qdrant internal.Config `envPrefix:"QDRANT_"`

	// CollectionName is the collection to ensure in Qdrant.
	CollectionName string `env:"COLLECTION_NAME,notEmpty,required"`
	// VectorsSize is the vector dimensionality.
	VectorsSize int `env:"VECTORS_SIZE,notEmpty,required"`
	// VectorsDistance is the Qdrant distance metric name.
	VectorsDistance string `env:"VECTORS_DISTANCE,notEmpty,required"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	slog.InfoContext(ctx, "running job", slog.String("version", Version), slog.String("build_time", BuildTime))

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		slog.ErrorContext(ctx, "failed to parse env", slog.Any("err", err))
		os.Exit(1)
	}

	if cfg.VectorsSize <= 0 {
		slog.ErrorContext(ctx, "vectors size must be positive", slog.Int("vectors_size", cfg.VectorsSize))
		os.Exit(1)
	}

	client, err := internal.NewClient(cfg.Qdrant)
	if err != nil {
		slog.ErrorContext(ctx, "invalid qdrant host", slog.Any("err", err))
		os.Exit(1)
	}

	exists, err := client.CollectionExists(ctx, cfg.CollectionName)
	if err != nil {
		slog.ErrorContext(ctx, "failed to check collection exists", slog.Any("err", err))
		os.Exit(1)
	}

	if exists {
		slog.InfoContext(ctx, "collection already exists", slog.String("collection_name", cfg.CollectionName))
		return
	}

	if err := client.CreateCollection(ctx, cfg.CollectionName, cfg.VectorsSize, cfg.VectorsDistance); err != nil {
		slog.ErrorContext(ctx, "failed to create collection", slog.Any("err", err))
		os.Exit(1)
	}

	slog.InfoContext(ctx, "collection created", slog.String("collection_name", cfg.CollectionName))
}
