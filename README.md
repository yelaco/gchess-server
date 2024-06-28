# Go Chess Server

The ```go-chess-server``` project is a backend server application designed to manage chess game and provide a way to create a platform for multiplayer. Built using the Go programming language, the server supports various functionalities essential for hosting online chess matches, including player authentication, game state management, move validation and real-time game updates

## Features

**Player Authentication**
- Login/Registration: Users can create accounts and log in
- Session management: The server handles user sessions to maintain login states and manage active connections

**Game Management**
- Matchmaking: Players can enter matching queue and wait for another player to create a match. If a player leave the match, he/she can come back later by rejoin the match.
- Game state: The server maintains the state of ongoing games, tracking each move and updating the board accordingly.
- Data persistence: After a game ended, its information is saved to database, ensuring that game states are preserved and can be retrieved later for user's analysis purposes.
  
**Move Handling**
- Move validation: The server validates each move to ensure they are legal according to chess rules.
- The server checks for conditions like check, checkmate, and stalemate after each move to determine the game's status.
  
**Real-Time Updates**
- WebSocket integration: The server uses WebSockets to provide real-time updates to connected clients, ensuring players see the latest game state without needing to refresh.

## How to run

### Client

For the client side, I made a CLI for testing purpose. After the server and database are up and running, you can run the ```client.go``` file in the ```test``` folder
```console
$ cd test
$ go run test.go
```

Here is a video that showcase all the functionalities of the CLI

[client](https://github.com/yelaco/go-chess-server/assets/100106895/434673b4-817b-4423-97d0-39a22cecd751)

### Server

**Docker**

Clone the repository and run the command from the root of the project. This will automatically pull the Docker image for this project and the PostgreSQL database to set up the environment using Docker Compose.
```console
$ docker compose up
```

You can also modify the ```docker-compose.yml``` to build from source instead
```yml
services:
  app:
    build: .
    container_name: go-chess-server
    ports:
      - "7202:7202"
      - "7201:7201"
    environment:
      - DATABASE_URL=postgresql://server:chessserver@db:5432/chess
    depends_on:
      - db
```
If you do it this way, you will also need a config file for the server to run
```console
$ mkdir .go-server-chess
$ vim config.json
```
For example
```json
{
    "host": {
        "address": "localhost",
        "game_server_port": "7201",
        "rest_server_port": "7202"
    },
    "database": {
        "name": "postgres",
        "host": "localhost",
        "user": "username",
        "password": "password"
    },
    "game": {
        "matching_timeout": 30
    }
}
```

## API

### REST

 [References](https://documenter.getpostman.com/view/30874401/2sA3duEsiX)
 
- ```POST /api/users```: To register user
- ```POST /api/login```: To log in to the server
- ```GET /api/sessions```: Retrieve match records played by user
- ```GET /api/sessions/{sessionid}```: Retrieve single match record based on ID

### WebSocket

After login, user can now join a match by sending matching request
```json
{
    "action": "matching",
    "data": {
        "playerId": "12345"
    }
}
```

If the ```action``` and ```data``` is valid, server pushes that user into the matching queue. When a match happens, the two connections are forwarded to game management module, where a game instance will be initialized and binded with the player pair. Then, a message is sent back to the user to notify about the match.
```json
{
    "type": "matched",
    "session_id": "1232524",
    "game_state": {
        "status": "ACTIVE",
        "board": [[]]
        "is_white_turn": true,
    },
    "player_state": {
        "is_white_side": true
    }
}
```

On the contrary, if there are any errors in the process or the matching request is timeout, the server replies with
- Error (Note that this error json is universal for all the error response to users)
```json
{
    "type": "error",
    "error": "<err_msg>",
}
```

- Timeout
```json
{
    "type": "timeout",
    "error": "<err_msg>",
}
```

In a match, users can send move request with 
```json
{
    "action": "move",
    "data": {
        "playerId": "12345",
        "sessionId": "1719199808062498696",
        "move": "e2-e4"
    }
}
```

And get resonses as 
```json
{
    "type": "session",
    "game_state": {
        "status": "STALEMATE",
        "board": [[]]
        "is_white_turn": false,
    }
}
```

After the game reaches end state, the server notifies both players and close their connections.
