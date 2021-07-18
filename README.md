# Installation

## Requirement
- Docker Desktop 3.5.2
- Docker 20.10.7
- Docker Compose 1.29.2
- Marvel API credential from https://developer.marvel.com/

## Running
1. Edit `env/common.env` to add your Marvel API credential. Change `<your_private_key>` and `<your_public_key>` accordingly.

2. Run server
```
docker-compose up
```

3. Open http://localhost
You can access containers via clicking the link. Please wait a moment if you can't see any proxy-able containers. Access container via <container-name>.localhost.
