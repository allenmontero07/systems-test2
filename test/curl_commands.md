# Test Curl Commands for Feedback API

## Health Check
```bash
curl http://localhost:8080/api/health
```

psql -U allen -d feedbackdb

## Create Feedback (POST)
```bash
curl -X POST http://localhost:8080/api/feedback ^
  -H "Content-Type: application/json" ^
  -d "{\"name\":\"John Doe\",\"email\":\"john@example.com\",\"subject\":\"Great Service\",\"message\":\"I love this product!\"}"
```

## Get All Feedback (GET)
```bash
curl http://localhost:8080/api/feedback
```

## Get Feedback by ID (GET)
```bash
curl http://localhost:8080/api/feedback/1
```

## Delete Feedback (DELETE)
```bash
curl -X DELETE http://localhost:8080/api/feedback/1