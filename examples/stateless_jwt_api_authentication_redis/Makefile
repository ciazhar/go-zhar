login:
	curl -X POST http://localhost:3000/login \
      -H "Content-Type: application/json" \
      -d '{ \
        "username": "john_doe", \
        "password": "123456" \
      }'

protected:
	curl -X GET http://localhost:3000/user/profile \
      -H "Authorization: Bearer <ACCESS_TOKEN>"

refresh:
	curl -X POST http://localhost:3000/refresh-token \
      -H "Content-Type: application/json" \
      -d '{ \
        "refresh_token": "<REFRESH_TOKEN>" \
      }'

logout:
	curl -X POST http://localhost:3000/logout \
      -H "Authorization: Bearer <ACCESS_TOKEN>"
