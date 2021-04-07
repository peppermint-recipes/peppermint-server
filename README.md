# peppermint-server

The backend for [peppermint-clients](https://github.com/peppermint-recipes/peppermint-clients).


## Development

Run `docker-compose up` to start a MongoDB instance.
Then `make run` to start the server.


## Deployment

```sh
version: '3'

services:
  mongo:
    image: mongo:4.4-bionic
    ports:
      - 127.0.0.1:27017:27017
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: example

  mongo-express:
    image: mongo-express:0.54.0
    ports:
        - 127.0.0.1:8081:8081
    environment:
      ME_CONFIG_MONGODB_ADMINUSERNAME: root
      ME_CONFIG_MONGODB_ADMINPASSWORD: example
    depends_on:
      - mongo

  peppermint-server:
    image: ghcr.io/peppermint-recipes/peppermint-server
    ports:
        - 0.0.0.0:8080:8080
    environment:
        DATABASE_USERNAME: root
        DATABASE_PASSWORD: example
        DATABASE_ENDPOINT: mongo:27017
        WEBSERVER_PORT: 8080
        WEBSERVER_ADDRESS: 0.0.0.0
    depends_on:
      - mongo

```
