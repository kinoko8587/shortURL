
# URL Shortener

## Requirements
Reference: https://bytebytego.com/courses/system-design-interview/design-a-url-shortener

### Problem Understanding and Design Scope

#### Example Usage
- **Original URL**: `https://www.systeminterview.com/q=chatsystem&c=loggedin&v=v3&l=long`
- **Shortened URL**: `https://tinyurl.com/y7keocwj`
- Clicking the shortened URL redirects to the original URL

#### Traffic Requirements
- **URL Generation**: 100 million URLs per day
- **Write Operations**: 100 million / 24 / 3600 = **1,160 per second**
- **Read Operations**: Assuming 10:1 read/write ratio = **11,600 per second**

#### URL Specifications
- **Length**: As short as possible
- **Characters**: Numbers (0-9) and letters (a-z, A-Z)
- **Lifecycle**: URLs cannot be deleted or updated (for simplicity)

### Core Use Cases

1. **URL Shortening**: Given a long URL → return a much shorter URL
2. **URL Redirecting**: Given a shorter URL → redirect to the original URL
3. **High Availability**: System must handle traffic with scalability and fault tolerance

### Capacity Planning (10-year projection)

#### Storage Requirements
- **Total URLs**: 100 million × 365 days × 10 years = **365 billion records**
- **Average URL Length**: 100 bytes
- **Total Storage**: 365 billion × 100 bytes = **36.5 TB**

#### Performance Requirements
- Handle **1,160 writes/second**
- Handle **11,600 reads/second**
- Maintain high availability and low latency for redirects

## API Design

API endpoints facilitate the communication between clients and servers. We will design the APIs REST-style.

### Primary API Endpoints

#### 1. URL Shortening
To create a new short URL, a client sends a POST request with the original long URL.

```
POST api/v1/data/shorten
```

**Request Parameter:**
```json
{
  "longUrl": "longURLString"
}
```

**Response:**
```
shortURL
```

#### 2. URL Redirecting
To redirect a short URL to the corresponding long URL, a client sends a GET request.

```
GET http://short-url.com/{shortUrl}
```

**Example:**
```
GET  http://short-url.com/y7keocwj
```

**Response:**
Returns longURL for HTTP redirection (HTTP 301/302)

**Note on 301 vs 302:**
If the priority is to reduce the server load, using 301 redirect makes sense as only the first request of the same URL is sent to URL shortening servers.

## Launch

### Development with Docker Compose
1. **Start the application**
   ```bash
   docker-compose up --build
   ```

2. **Run in detached mode**
   ```bash
   docker-compose up -d
   ```

3. **Stop the application**
   ```bash
   docker-compose down
   ```

### Alternative: Direct Docker Build
1. **Build the Docker Image**
   ```bash
   DOCKER_BUILDKIT=0 docker build -t short-url .
   ```

2. **Run the image**
   ```bash
   docker run -p 8080:8080 short-url
   ```

## Testing

### API Endpoints Testing

#### 1. Health Check
```bash
curl http://localhost:8080/ping
```
**Expected Response:**
```json
{"message":"pong"}
```

#### 2. Create Short URL
```bash
curl -X POST http://localhost:8080/api/v1/data/shorten \
  -H "Content-Type: application/json" \
  -d '{"longUrl":"https://www.google.com"}'
```
**Expected Response:**
```json
{
  "shortUrl": "http://localhost:8080/abc1234",
  "longUrl": "https://www.google.com"
}
```

#### 3. Test URL Redirection
```bash
# Test redirect (replace abc1234 with actual short code)
curl -I http://localhost:8080/abc1234
```
**Expected Response:**
```
HTTP/1.1 301 Moved Permanently
Location: https://www.google.com
```

#### 4. Test with curl following redirects
```bash
curl -L http://localhost:8080/abc1234
```
This will follow the redirect and show the final page content.

### Rate Limiting Testing

#### Test Rate Limiting with Single User
```bash
# This should show 10 requests returning 200, then 429 for subsequent requests
for i in {1..15}; do
  curl -s -o /dev/null -w "%{http_code}\n" -H "X-User-ID: test-user" http://localhost:8080/ping
done
```
**Expected Output:**
```
200
200
200
200
200
200
200
200
200
200
429
429
429
429
429
```

#### Test with Different Users
```bash
# Different users should have separate rate limits
curl -H "X-User-ID: user1" http://localhost:8080/ping
curl -H "X-User-ID: user2" http://localhost:8080/ping
```

### Database Testing

#### Check Database Content
```bash
# View all stored URLs
docker-compose exec db psql -U postgres -d mydatabase -c "SELECT * FROM urls;"
```

#### Check Database Schema
```bash
# View table structure
docker-compose exec db psql -U postgres -d mydatabase -c "\d urls"
```

### Testing Different URL Types

#### Test Various URLs
```bash
# Long URL
curl -X POST http://localhost:8080/api/v1/data/shorten \
  -H "Content-Type: application/json" \
  -d '{"longUrl":"https://www.example.com/very/long/path/with/parameters?param1=value1&param2=value2"}'

# URL with special characters
curl -X POST http://localhost:8080/api/v1/data/shorten \
  -H "Content-Type: application/json" \
  -d '{"longUrl":"https://www.example.com/search?q=url+shortener&category=tech"}'

# HTTPS URL
curl -X POST http://localhost:8080/api/v1/data/shorten \
  -H "Content-Type: application/json" \
  -d '{"longUrl":"https://docs.google.com/document/d/1234567890/edit"}'
```

### Error Testing

#### Test Invalid Requests
```bash
# Empty longUrl
curl -X POST http://localhost:8080/api/v1/data/shorten \
  -H "Content-Type: application/json" \
  -d '{"longUrl":""}'

# Missing longUrl field
curl -X POST http://localhost:8080/api/v1/data/shorten \
  -H "Content-Type: application/json" \
  -d '{}'

# Invalid JSON
curl -X POST http://localhost:8080/api/v1/data/shorten \
  -H "Content-Type: application/json" \
  -d '{"longUrl":"https://example.com"'
```

#### Test Non-existent Short URLs
```bash
# Should return 404
curl -I http://localhost:8080/nonexistent
```

### Load Testing (Optional)

#### Simple Concurrent Testing
```bash
# Test with multiple concurrent requests
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/data/shorten \
    -H "Content-Type: application/json" \
    -d "{\"longUrl\":\"https://example$i.com\"}" &
done
wait
```

### Environment Testing

#### Test with Different Base URLs
```bash
# Set custom base URL
export BASE_URL=https://myshortener.com
docker-compose up --build

# Test the response includes the custom base URL
curl -X POST http://localhost:8080/api/v1/data/shorten \
  -H "Content-Type: application/json" \
  -d '{"longUrl":"https://www.google.com"}'
```

## desgin   
### clean architecture
- domain: define the business rules
- usecase: define the business logic
- controller: define the API
- infrastructure: define the data source

### Project Structure
```
shortURL/
├── cmd/                       # Application entry points
├── internal/                  # Private application code
│   ├── domain/               # Enterprise business rules
│   │   └── url.go           # URL entity and repository interface
│   ├── usecases/            # Application business rules
│   │   ├── create_short_url.go    # Create short URL use case
│   │   └── redirect_short_url.go  # Redirect URL use case
│   ├── controller/          # Interface adapters (HTTP handlers)
│   │   └── controller.go    # Gin HTTP handlers
│   ├── infrastructure/      # Frameworks and drivers
│   │   └── store.go        # PostgreSQL repository implementation
│   ├── middleware/          # HTTP middleware
│   │   └── middleware.go   # Rate limiting and other middleware
│   └── app.go              # Application initialization
├── main.go                   # Main entry point
├── go.mod                    # Go module file
├── go.sum                    # Go module checksums
├── Dockerfile               # Docker configuration
├── docker-compose.yaml      # Docker Compose configuration
└── readme.md               # Project documentation
```

#### Layer Responsibilities

**Domain Layer (`internal/domain/`)**
- Contains enterprise business rules
- Defines entities and repository interfaces
- Has no dependencies on other layers
- Example: URL entity, URLRepository interface

**Use Cases Layer (`internal/usecases/`)**
- Contains application-specific business rules
- Orchestrates the flow of data between entities
- Depends only on the domain layer
- Example: CreateShortURL, RedirectURL use cases

**Controller Layer (`internal/controller/`)**
- Handles HTTP requests and responses
- Converts data between HTTP format and use case format
- Depends on use cases layer
- Example: Gin handlers for API endpoints

**Infrastructure Layer (`internal/infrastructure/`)**
- Implements interfaces defined in domain layer
- Contains database access, external services
- Depends on domain layer for interfaces
- Example: PostgreSQL repository implementation

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

#### Database Schema
```sql
CREATE TABLE urls (
    id BIGSERIAL PRIMARY KEY,
    shortURL VARCHAR(7) NOT NULL UNIQUE,
    longURL TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Fields:**
- `id`: Auto-incrementing primary key
- `shortURL`: The shortened URL identifier (e.g., "y7keocwj")
- `longURL`: The original long URL

## User Flow

### URL Redirection Flow

1. A user clicks a short URL link: `https://tinyurl.com/zn9edcu`

2. The load balancer forwards the request to web servers.

3. If a shortURL is already in the cache, return the longURL directly.

4. If a shortURL is not in the cache, fetch the longURL from the database. If it is not in the database, it is likely a user entered an invalid shortURL.

5. The longURL is returned to the user.

## URL Shortening Algorithm

### Hash + Collision Resolution Approach

We use a hash-based approach with collision resolution to generate short URLs:

1. **Generate Hash**: Use CRC32 or MD5 to hash the long URL
2. **Extract Short Code**: Take the first 7 characters of the hash (base62 encoded)
3. **Collision Detection**: Check if the short URL already exists in the database
4. **Collision Resolution**: If collision occurs, append a counter or use linear probing
5. **Store Mapping**: Save the short URL to long URL mapping in the database

**Advantages:**
- Fixed short URL length
- No need for a counter/sequence
- Deterministic for the same input (with caching benefits)

**Implementation Steps:**
```
longURL → Hash(MD5/CRC32) → Base62 Encode → Take first 7 chars → Check DB → Resolve collision if needed
```
