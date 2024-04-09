#! /bin/bash

# curl localhost:3300/users | jq 

# curl -X POST \
#     -H "Content-Type: application/json" \
#     -d '{"name":"Pru Wapol","age":32,"sex":"Female"}' \
#     localhost:3300/users | jq


curl -X PUT \
   -H "Content-Type: application/json" \
   -d '{"name" : "Leyenda Troth", "age" : 34, "sex" : "Female"}' \
    localhost:3300/users/5 | jq 