#! /bin/bash

# curl localhost:3300/users | jq 

curl -v -X POST \
    -H "Content-Type: application/json" \
    -d '{
      "name": "Womens Clothing",
      "description": "Apparel designed for female individuals, including dresses, skirts, blouses, and tops."
    }' \
    localhost:3300/categories | jq 


# curl -v -X DELETE localhost:3300/users/5 | jq 