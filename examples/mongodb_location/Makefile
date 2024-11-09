run:
	go run cmd/main.go

insert:
	curl --location 'localhost:3000/location' \
	--header 'Content-Type: application/json' \
	--data-raw '{"name":"Location 1","coordinate":[-73.9712,40.7831]}' && \
	curl --location 'localhost:3000/location' \
	--header 'Content-Type: application/json' \
	--data-raw '{"name":"Location 2","coordinate":[-74.006,40.7128]}' && \
	curl --location 'localhost:3000/location' \
	--header 'Content-Type: application/json' \
	--data-raw '{"name":"Location 3","coordinate":[-73.9776,40.7648]}'

get:
	curl --location 'localhost:3000/location?longitude=-74.006&latitude=40.78&max_distance=10000&limit=2'