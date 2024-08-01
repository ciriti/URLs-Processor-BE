# CURLS

### Auth

curl -X POST -d "user=admin&pass=password" http://localhost:8080/authenticate

### Test protected ep

curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8080/admin/test-protected

### Test ep

curl http://localhost:8080/test
