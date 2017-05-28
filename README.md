Example of tournament service.

⚠️ By default system contains a several players for testing with zero points of amount
Players: P1, P2, P3, P4 and P5.

See [game/maing.go](https://github.com/AlianceBloom/example_api_service/blob/master/game/main.go#L20)

# API
See [api_service/server.go](https://github.com/AlianceBloom/example_api_service/blob/master/api_service/server.go#L43)
Create new user

```
POST /newUser
params: 
PlayerId - string
Points - int
example: 
curl -X PUT "localhost:8080/newUser?PlayerId=P2&Points=21"
```

Get ballance of player
```
GET /balance
params: 
PlayerId - string
Points - int
example: 
curl -X GET "localhost:8080/balance?PlayerId=P1"
```
Take points from player
```
PUT /take
params: 
PlayerId - string
Amount - int
example: 
curl -X PUT "localhost:8080/take?PlayerId=P1&Amount=21"
```

Fund points to player
```
PUT /fund
params: 
PlayerId - string
Amount - int
example: 
curl -X PUT "localhost:8080/fund?PlayerId=P1&Amount=21"
```

Create new tournament 
```
POST /announceTournament
params: 
TournamentId - string
Deposit - int
example: 
curl -X POST "localhost:8080/fund?TournamentId=T1&Deposit=21"
```

Join tournament with player and backers
```
PUT /joinTournament
params: 
TournamentId - string
PlayerId - string
BackerId - list of string
example: 
curl -X PUT "localhost:8080/joinTournament?TournamentId=T1&PlayerId=P1" 
curl -X PUT "localhost:8080/joinTournament?TournamentId=T1&PlayerId=P2&BackerId=P3&BackerId=P4" 
```

Finish tournament and get result
```
POST /resultTournament
params: 
TournamentId - string
PlayerId - string
example: 
curl -X POST "localhost:8080/resultTournament?TournamentId=T1&PlayerId=P1" 
```

# Docker
```
git clone git@github.com:AlianceBloom/example_api_service.git
cd example_api_service
docker build -t example_api_service . 
docker run --publish 8080:8080  -it example_api_service
```
