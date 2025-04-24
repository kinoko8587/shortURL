
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
- use gin
- POST: Create a new resource
- GET: Read an existing resource
- PUT: Update an existing resource
- DELETE: Delete an existing resource

### rate limiting(token bucket)
- limit the number of requests per user
- use a token bucket to limit the number of requests per user
- the token bucket is a list of tokens
- each token is a timestamp

#### testing
- use curl to test the API
- use docker compose to run the service and test the API
- With a for loop to simulate spam, should return 429
``
for i in {1..20}; do
  curl -s -o /dev/null -w "%{http_code}\n" -H "X-User-ID: test-user" http://localhost:8080/ping
done
``
- with Different IPs / Users
``
curl -H "X-User-ID: user1" http://localhost:8080/ping
curl -H "X-User-ID: user2" http://localhost:8080/ping
``

### database
- use postgres
