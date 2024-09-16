

1. Build the Docker Image
``DOCKER_BUILDKIT=0 docker build -t short-url .``


2. run image
``docker run -p 8080:8080 short-url``

3. Rebuild and Run Docker Compose
docker-compose up --build

List
1.docker
2.db
3.rate limiting(token bucket)
