
# short url service 

## launch
1. Build the Docker Image
``DOCKER_BUILDKIT=0 docker build -t short-url .``
2. run image
``docker run -p 8080:8080 short-url``
3. Rebuild and Run Docker Compose
docker-compose up --build

## desgin   
### clean architecture
- domain: define the business rules
- usecase: define the business logic
- controller: define the API
- infrastructure: define the data source

### API
- POST: Create a new resource
- GET: Read an existing resource
- PUT: Update an existing resource
- DELETE: Delete an existing resource

### rate limiting(token bucket)
- limit the number of requests per user
- use a token bucket to limit the number of requests per user
- the token bucket is a list of tokens
- each token is a timestamp

### database
- use postgres


