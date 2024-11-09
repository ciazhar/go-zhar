# Temporary file to store the last inserted ID
ID_FILE := .last_inserted_id

run-mongodb:
	cd ../../deployments/mongodb && docker-compose up

run:
	go run cmd/main.go

get:
	curl --location 'localhost:3000/person'

insert:
	@curl --silent --location 'localhost:3000/person' \
    --header 'Content-Type: application/json' \
    --data-raw '{ \
        "name": "Woods Peppermint", \
        "age": 18, \
        "email": "woods.peppermint@example.com" \
    }' | jq -r '.data.id' | tee $(ID_FILE) | xargs -I {} echo "Inserted ID: {}"

insert-batch:
	@curl --silent -X POST \
      http://localhost:3000/person/batch \
      -H 'Content-Type: application/json' \
      -d '[\
      {\
        "address": {\
          "street": "12",\
          "city": "California",\
          "country": "US"\
        },\
        "age": 35,\
        "email": "john.doe@example.com",\
        "hobbies": ["Cooking", "Gaming"],\
        "name": "John Doe"\
      },\
      {\
        "address": {\
          "street": "32",\
          "city": "Kendal",\
          "country": "Indonesia"\
        },\
        "age": 25,\
        "email": "janesmith@example.com",\
        "hobbies": ["Cooking"],\
        "name": "Jane Smith"\
      }\
    ]' | jq -r '.data[0].id' | tee $(ID_FILE) | xargs -I {} echo "First Inserted ID: {}"

find-one:
	@if [ -f $(ID_FILE) ]; then \
		ID=$$(cat $(ID_FILE)); \
		curl --location "localhost:3000/person/$${ID}"; \
	else \
		echo "No ID found. Run 'make insert' first."; \
	fi

find-all:
	curl --location 'localhost:3000/person?name=John%20Doe&age=35'

find-country:
	curl --location 'localhost:3000/person/country?country=Indonesia'

find-age-range:
	curl --location 'localhost:3000/person/age-range?startAge=27&endAge=40'

find-hobby:
	curl --location 'localhost:3000/person/hobby?hobby=Gaming%2CCooking'

find-minified:
	curl --location 'localhost:3000/person/minified'

update:
	@if [ -f $(ID_FILE) ]; then \
		ID=$$(cat $(ID_FILE)); \
		curl --location --request PUT "localhost:3000/person/$${ID}" \
	    --header 'Content-Type: application/json' \
	    --data-raw '{\
	        "address": {\
	            "street": "12",\
	            "city": "California",\
	            "country": "US"\
	        },\
	        "age": 35,\
	        "email": "john.doe@example.com",\
	        "hobbies": [\
	            "Cooking",\
	            "Gaming"\
	        ],\
	        "name": "John Doel"\
	    }'; \
	else \
		echo "No ID found. Run 'make insert' first."; \
	fi

delete:
	@if [ -f $(ID_FILE) ]; then \
		ID=$$(cat $(ID_FILE)); \
		curl --location --request DELETE "localhost:3000/person/$${ID}"; \
	else \
		echo "No ID found. Run 'make insert' first."; \
	fi
