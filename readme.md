# Description
Planet extractor is a server based idle game. This enables the possiblity of cross platform support and multiplayer.

# Features implemented
## Server
### Resource extraction and extractor building
The game logic is implemented on the server side, meaning that actions and the game loop are all performed by the server.
The server stores a list of players and their resources in memory, and provides methods to harvest resources to add to their storage and build extractors for automatic resource extraction if the player has sufficient resources. The game loop keeps extractors running for passive resource collection.

### User authentication
The server authenticates login requests (currently only through google sign-in) and allows users to act on their respective player accounts once authenticated.

### Cross-platform support
By abstracting the game logic from the client-server communication, different protocols can be implemented for players on any platform to access the game.
For greatest accessibilty, the websocket protocol is currently supported to allow browser play.

## Client
|Feature|Browser|
|:------|:-----:|
|Login  |:heavy_check_mark|