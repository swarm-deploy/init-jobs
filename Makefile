build/qdrant:
	docker build -f ./qdrant/Dockerfile ./qdrant -t swarmdeployorg/init-jobs-qdrant:local