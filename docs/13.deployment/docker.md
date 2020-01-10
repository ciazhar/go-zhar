# Docker

### Running Docker In Background
```bash
docker run -d container-name
```    

### Check Process Running In Docker Container
```bash
docker exec container-name ps aux
```

### Creating Networks Between Containers using Networks
```bash
docker run --link container-name:
```