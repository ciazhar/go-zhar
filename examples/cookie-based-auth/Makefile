login:
	curl -X POST http://localhost:3000/login -d '{"username": "admin", "password": "password"}' -H "Content-Type: application/json"

logout:
	curl -X POST http://localhost:3000/logout --cookie "session_id=54100451-9198-4dd1-a9ab-617f7797ad53"

protected:
	curl -X GET http://localhost:3000/protected --cookie "session_id=54100451-9198-4dd1-a9ab-617f7797ad53"