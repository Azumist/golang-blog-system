Basic blogging system made in Go to familiarize myself with the launguage. Uses sqlite and basic cookie-based authentication. 

Start with:
`go run cmd/server/main.go`

TODO: 
- comment deletion

## .env file structure
```env
# Server Configuration
PORT=8080

# Database Configuration
DB_PATH=blog.db

# Authentication
ADMIN_PASSWORD=admin
SESSION_SECRET=secret-change-me
```
## Routes
- GET /api/auth/status
- POST /api/auth/login
- POST /api/auth/logout
- GET /api/articles
- GET /api/articles/{id}
- POST /api/articles/{id}
- PUT /api/articles/{id}
- DELETE /api/articles/{id}
- POST /api/articles/{id}/comments


## Example requests
### Logging in
```
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "password": "admin"
  }' \
  -c cookies.txt
```

### Creating new article
```
curl -X POST http://localhost:8080/api/articles \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "title":"Protected Post",
    "content":"Only admin can create this",
    "author":"Admin",
    "tags":["secure"]
  }'
```