#! /bin/bash

# curl localhost:3300/users | jq 

curl -v -X POST \
    -H "Content-Type: application/json" \
    -d '{"name":"Patrizia Lamport","age":43,"sex":"Female","username":"plamport3","email":"plamport3@globo.com","phone":"8415083252","password":"mdqi","city":"Tianmen","country":"China"}' \
    localhost:3300/users | jq


# curl -X PUT \
#    -H "Content-Type: application/json" \
#    -d '{"name" : "Leyenda Troth", "age" : 34, "sex" : "Female"}' \
#     localhost:3300/users/5 | jq 


# curl -v -X DELETE localhost:3300/users/5 | jq 