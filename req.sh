#! /bin/bash

# curl localhost:3300/users | jq 

# curl -v -X POST \
#     -H "Content-Type: application/json" \
#     -d '{
#       "name": "Womens Clothing",
#       "description": "Apparel designed for female individuals, including dresses, skirts, blouses, and tops."
#     }' \
#     localhost:3300/categories | jq 


# curl -v -X DELETE localhost:3300/users/5 | jq 

curl -v -X POST \
  -H "Content-Type: application/json" \
  -d ' {
    "Name": "Samsung Galaxy S23 Ultra",
    "Summary": "Flagship smartphone with advanced features.",
    "Description": "The Samsung Galaxy S23 Ultra is a premium smartphone...",
    "Price": 1399,
    "Quantity": 80,
    "ProductCode" : "prod002",
    "Brand": {
      "Id": 1
    },
    "Category": {
      "Id": 1
    },
    "ImageUrl": "https://example.com/galaxy_s23_ultra.jpg"
  }' \
  localhost:3300/products | jq 