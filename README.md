# Installation

## Requirement
- Docker Desktop 3.5.2
- Docker 20.10.7
- Docker Compose 1.29.2
- Marvel API credential from https://developer.marvel.com/

## Running
1. Edit `config/common.json` to add your Marvel API credential. Change `<your_private_key>` and `<your_public_key>` accordingly.

2. Run server
```
docker-compose up
```

3. Open API at http://localhost:8080
4. Open SwaggerUI at http://localhost:3000
