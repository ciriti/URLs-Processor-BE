# CURLS

### Auth

curl -X POST -H "Content-Type: application/json" -d '{ "user": "admin@example.com", "pass": "password" }' http://localhost:8080/authenticate

### Post urls

curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3Nzc1NTcsInVzZXIiOiJhZG1pbkBleGFtcGxlLmNvbSJ9.3lhCE9AHyp0uf-lZasZLTqsWnDMqPBnDH42vs0LHh60" -d '{
  "urls": [
    "http://example1.com",
    "http://example2.com",
    "http://example3.com"
  ]
}' http://localhost:8080/api/urls

### Test protected ep

curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8080/admin/test-protected

### Test ep

curl http://localhost:8080/test
