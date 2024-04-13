run:
	go run cmd/main.go

get:
	curl --location 'localhost:3000/person'

insert:
	curl --location 'localhost:3000/person' \
    --header 'Content-Type: application/json' \
    --data-raw '{ \
        "name": "Woods Peppermint", \
        "age": 18, \
        "email": "woods.peppermint@example.com" \
    }'

insert-batch:
	curl -X POST \
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
        "name": "John Doel"\
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
    ]'

find-one:
	curl --location 'localhost:3000/person/661392cdce2ee988af47e4b3'
find-all:
	curl --location 'localhost:3000/person?name=John%20Doe&age=30'
find-country:
	curl --location 'localhost:3000/person/country?country=Indonesia'
find-age-range:
	curl --location 'localhost:3000/person/age-range?startAge=27&endAge=40'
find-hobby:
	curl --location 'localhost:3000/person/hobby?hobby=Gaming%2CCooking'
find-minified:
	curl --location 'localhost:3000/person/minified'
update:
	curl --location --request PUT 'localhost:3000/person/661392cdce2ee988af47e4b3' \
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
    }'
delete:
	curl --location --request DELETE 'localhost:3000/person/661392cdce2ee988af47e4b4'