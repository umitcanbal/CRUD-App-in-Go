curl -X GET "http://localhost:8080/getusers"

curl -X DELETE "http://localhost:8080/deleteuser/15563"

curl -X POST http://localhost:8080/createuser \
 -H "Content-Type: application/json" \
 -d '{
"Name": "Alice",
"Age": 30
}'

curl -X PUT http://localhost:8080/updateuser \
 -H "Content-Type: application/json" \
 -d '{
"ID": 123,
"Name": "Alice Updated",
"Age": 31
}'
