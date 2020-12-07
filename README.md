# Roulette API


## Technical Choices

* Language - Go
* Postgresql Database
* Docker-compose for running the database
* Prometheus for monitoring


## Setup

### Requirements

* Ensure `Go` and `docker-compose` are installed

### Running the application

In the root directory,
* Run `go install ./...` which will pull and install all the dependencies
* Run `docker-compose up -d ` to start postgres. This should create the database and the tables.
* Run `go run main.go` to start the application. Starts on port 8080 by default. To start on a different port run, `go run main.go -port=portnumber `

### Running the tests

In the root directory,
* Run `go test ./...`
* All the test files are located in the controllers package.
* The tests do not depend on any external dependencies. The db interface has been mocked out in the mocks package.

### Monitoring

* Accessible at /metrics endpoint
* Prometheus metrics have been used for monitoring.
* Uses middleware to monitor the number of hits and latency of each http request.

## Oveview

* A roulette game takes place in a room. 
* A room can be created by specifying the roulette variant type. Currently supports sky_roulette and superboost_roulette. More variants can be added when required (configurable)
* A room has a limit on the number of rounds and the players (depends on the roulette variant)
* A player needs to join a room in order to play. Multiple players can join a room.
* Bets can be placed during each round.
* Currently supports straight up, colour and odd/even bets. More bet types can be added when required (configurable).
* When all the players in a room are ready, spin occurs.
* After a spin, the result of the round is posted to the API server and players are moved to the next round.
* Players can retrieve their results for the round indicating win or loss and the total return.
* Also provides endpoint for retrieving the rooms, room status, player status. 


## API Methods

* JSON is used for all requests and responses
* Status code 200 is used for all successfull responses
* Status codes 404 and 500 used for unsuccessful responses
### GET /room/variants

* Get all the supported variant types


Response:
```javascript
{
    "Data": [
        {
            "variant_type": 1,
            "variant_name": "sky_roulette",
            "max_players": 7,
            "max_rounds": 10
        },
        {
            "variant_type": 2,
            "variant_name": "superboost_roulette",
            "max_players": 10,
            "max_rounds": 15
        }
    ]
}
```



### POST /room

* Creates a new room. A unique roomid is created. 
* Different variants of roulettes can be created by specifying the `variant_type` in the body. Defaults to variant type 1 if it is not provided

Request Body:


```javascript
{
	"variant_type":1
}
```
Response:
```javascript
{
    "Data": {
        "id": "DA279D8EA1",
        "variant": {
            "variant_type": 1,
            "variant_name": "sky_roulette",
            "max_players": 7,
            "max_rounds": 10
        },
        "current_round": 1,
        "players": null,
        "created_at": "2020-12-06T21:40:05.585966557Z"
    }
}
```
### POST /player/join

* Join a room. The room id and the player id must be provided in the request body.

Request Body:


```javascript
{
    "room_id":"DA279D8EA1",
    "id":     "5EA8E64F56"
}
```
Response:
```javascript
{
    "Message": "Successfully joined the room",
    "Data": {
        "id": "5EA8E64F56",
        "room_id": "DA279D8EA1",
        "display_name": "Roulette-681D0F917D-Master",
        "ready_status": false,
        "bets_placed": null,
        "in_room": true,
        "created_at": "2020-12-06T21:46:53.881680105Z"
    }
}
```
Errors : 
1. Room does not exist
2. Player is already in another room. Please exit the room first
3. Player is already in the room

### POST /player/bet

* Place a bet. A player can place multiple bets.
* Supported bets are 


| BetName        | BetType     | Selections               |
| ------------- |:------------:| ------------------------:|
| Straightup    | 1            |  1 - 36                  |
| colour        | 2            |  1 (red) and 2 (black)   |
| oddeven       | 3            |  1 (even) and 2 (odd)    |


* New bettypes can be added by adding them into the `oddsconfig.json` file


Request Body:


```javascript
{
    "id": "5EA8E64F56",
    "room_id": "DA279D8EA1",
    "bets_placed": [
        {
            "bettype": 1,
            "stake": 1.5,
            "selection": 24
        },
         {
            "bettype": 2,
            "stake": 1.5,
            "selection": 2
        }
    ]
}
```
Response:
```javascript
{
    "Message": "Succefully placed the bet"
}
```
Errors :
1. Player is not in the room
2. Bet not accepted since player is already ready for the current spin


### POST /player/ready

* All players must be ready before a spin can occur.

Request Body:


```javascript
{
    "room_id":"DA279D8EA1",
    "id":     "5EA8E64F56"
}
```
Response:
```javascript
{
    "Message": "All the players are ready. Ready to spin"
}
```


### POST /bet/result

* Post the result of the round to the server after the spin


Request Body:


```javascript
{
    "room_id": "DA279D8EA1",
    "bet_result": 
        {
            "number": 24,
            "colour": 1
        }
    
}
```
Response:
```javascript
{
    "Message": "Successfully inserted the result"
}
```


### GET /player/result/{playerid}/{roomid}{roundno}

* Get the result of the player for a specific room and round

Request URL:
```javascript
http://localhost:8080/player/result/5EA8E64F56/DA279D8EA1/1
```
Response:
```javascript
{
    "Message": "Successfully retrieved the result",
    "Data": [
        {
            "room_id": "DA279D8EA1",
            "round_no": 1,
            "bettype": 1,
            "stake": 1.5,
            "odds": 35,
            "liability": 52.5,
            "selection": 24,
            "result": 1,
            "resultext": "WIN",
            "total_return": 54,
            "bet_result": {
                "number": 24,
                "colour": 1,
                "oddeven": 1
            }
        },
        {
            "room_id": "DA279D8EA1",
            "round_no": 1,
            "bettype": 2,
            "stake": 1.5,
            "odds": 0,
            "liability": 1.5,
            "selection": 2,
            "result": 0,
            "resultext": "LOST",
            "total_return": 0,
            "bet_result": {
                "number": 24,
                "colour": 1,
                "oddeven": 1
            }
        }
    ]
}
```


### POST /player/exit

* Exit the player from the room

Request Body:


```javascript
{
    "room_id":"DA279D8EA1",
    "id":     "5EA8E64F56"
}
```
Response:
```javascript
{
    "Message": "Successfully exited from the room"
}
```

### GET /room

* Get the status of all active rooms

Response:
```javascript
{
    "Data": [
        {
            "id": "DA279D8EA1",
            "variant": {
                "variant_type": 1,
                "variant_name": "sky_roulette",
                "max_players": 7,
                "max_rounds": 10
            },
            "current_round": 2,
            "players": null,
            "created_at": "0001-01-01T00:00:00Z"
        },
        {
            "id": "BBDCD6CFCC",
            "variant": {
                "variant_type": 2,
                "variant_name": "superboost_roulette",
                "max_players": 10,
                "max_rounds": 15
            },
            "current_round": 1,
            "players": null,
            "created_at": "0001-01-01T00:00:00Z"
        }
    ]
}
```

Errors : 

1. roomid must be provided

### GET /room/{roomid}

* Get the status of a specific room

Response:
```javascript
{
    "Data": {
        "id": "DA279D8EA1",
        "variant": {
            "variant_type": 1,
            "variant_name": "sky_roulette",
            "max_players": 7,
            "max_rounds": 10
        },
        "current_round": 2,
        "players": null,
        "created_at": "0001-01-01T00:00:00Z"
    }
}
```




### GET /bet/{roomid}

* Get the bet details from specific room 

Response:
```javascript
{
    "Message": "Successfully retrieved the bets for the room",
    "Data": [
        {
            "room_id": "DA279D8EA1",
            "round_no": 1,
            "bettype": 1,
            "stake": 1.5,
            "odds": 0,
            "liability": 52.5,
            "selection": 24,
            "result": 0,
            "resultext": "",
            "total_return": 0,
            "bet_result": {
                "number": 24,
                "colour": 1,
                "oddeven": 1
            }
        },
        {
            "room_id": "DA279D8EA1",
            "round_no": 1,
            "bettype": 2,
            "stake": 1.5,
            "odds": 0,
            "liability": 1.5,
            "selection": 2,
            "result": 0,
            "resultext": "",
            "total_return": 0,
            "bet_result": {
                "number": 24,
                "colour": 1,
                "oddeven": 1
            }
        }
    ]
}
```



### GET /player

* Get all the player details including their bets and results

Response:
```javascript
{
    "Message": "Succefully retrieved the player details",
    "Data": [
        {
            "id": "5EA8E64F56",
            "room_id": "DA279D8EA1",
            "display_name": "Roulette-681D0F917D-Master",
            "ready_status": false,
            "bets_placed": [
                {
                    "room_id": "DA279D8EA1",
                    "round_no": 1,
                    "bettype": 1,
                    "stake": 1.5,
                    "odds": 35,
                    "liability": 52.5,
                    "selection": 24,
                    "result": 1,
                    "resultext": "WIN",
                    "total_return": 54,
                    "bet_result": {
                        "number": 24,
                        "colour": 1,
                        "oddeven": 1
                    }
                },
                {
                    "room_id": "DA279D8EA1",
                    "round_no": 1,
                    "bettype": 2,
                    "stake": 1.5,
                    "odds": 0,
                    "liability": 1.5,
                    "selection": 2,
                    "result": 0,
                    "resultext": "LOST",
                    "total_return": 0,
                    "bet_result": {
                        "number": 24,
                        "colour": 1,
                        "oddeven": 1
                    }
                }
            ],
            "in_room": true,
            "created_at": "2020-12-06T21:46:53.88168Z"
        },
        {
            "id": "1FB9F75F57",
            "room_id": "BBDCD6CFCC",
            "display_name": "Roulette-BC977DEC7E-Master",
            "ready_status": false,
            "bets_placed": [
                {
                    "room_id": "BBDCD6CFCC",
                    "round_no": 1,
                    "bettype": 1,
                    "stake": 1.5,
                    "odds": 0,
                    "liability": 52.5,
                    "selection": 24,
                    "result": 0,
                    "resultext": "LOST",
                    "total_return": 0,
                    "bet_result": {
                        "number": 29,
                        "colour": 2,
                        "oddeven": 2
                    }
                },
                {
                    "room_id": "BBDCD6CFCC",
                    "round_no": 1,
                    "bettype": 2,
                    "stake": 1.5,
                    "odds": 1,
                    "liability": 1.5,
                    "selection": 2,
                    "result": 1,
                    "resultext": "WIN",
                    "total_return": 3,
                    "bet_result": {
                        "number": 29,
                        "colour": 2,
                        "oddeven": 2
                    }
                }
            ],
            "in_room": true,
            "created_at": "2020-12-06T22:22:37.082925Z"
        }
    ]
}
```
