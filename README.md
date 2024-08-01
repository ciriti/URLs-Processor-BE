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

### Get URLs

curl -X GET -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3Nzc1NTcsInVzZXIiOiJhZG1pbkBleGFtcGxlLmNvbSJ9.3lhCE9AHyp0uf-lZasZLTqsWnDMqPBnDH42vs0LHh61" http://localhost:8080/api/urls

### GetStatus

curl -X GET -H "Authorization: Bearer YOUR_JWT_TOKEN" "http://localhost:8080/api/checkStatus?id=TASK_ID"

curl -X GET -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3Nzc1NTcsInVzZXIiOiJhZG1pbkBleGFtcGxlLmNvbSJ9.3lhCE9AHyp0uf-lZasZLTqsWnDMqPBnDH42vs0LHh61" "http://localhost:8080/api/checkStatus?id=1"

### StartComputation

curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3Nzc1NTcsInVzZXIiOiJhZG1pbkBleGFtcGxlLmNvbSJ9.3lhCE9AHyp0uf-lZasZLTqsWnDMqPBnDH42vs0LHh61" -d '{
"url": "http://example.com"
}' http://localhost:8080/api/startComputation

curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3Nzc1NTcsInVzZXIiOiJhZG1pbkBleGFtcGxlLmNvbSJ9.3lhCE9AHyp0uf-lZasZLTqsWnDMqPBnDH42vs0LHh61" -d '{
"id": 1
}' http://localhost:8080/api/startComputation
