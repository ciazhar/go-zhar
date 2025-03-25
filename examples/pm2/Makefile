install:
	npm install -g pm2

build:
	go build -o myapp main.go

start:
#	pm2 start ./myapp --name myapp
	pm2 start ecosystem.config.js

status:
	pm2 status

logs:
	pm2 logs myapp

monit:
	pm2 monit