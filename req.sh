#! /bin/bash

# curl localhost:3300/users | jq 

curl -v -X PUT \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Apple",
        "description": "Renowned for its sophisticated technology products, including iPhones, MacBooks, and iPads."
    }' \
    localhost:3300/brands/4 | jq 


# curl -v -X DELETE localhost:3300/users/5 | jq 