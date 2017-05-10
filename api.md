#We will need a API

- Every user can own multiple rooms
- Every user can join multiple rooms
-


##Public endpoints
- [POST] /session -- login
- [DELETE] /session -- logout
- [POST] /sign-up -- create user

##Endpoints behind "auth middleware"
- [GET] /users/:id/rooms -- list user's rooms (joined and owned)
- [GET] /users/:id/rooms/:id -- retrieve user's room details
- [POST] /users/:id/rooms-- create an user's room
- [POST] /users/:id/rooms/:id -- join an user's room
- [DELETE] /users/:id/rooms/:id -- leave an user's room (delete if you're the  owner)
- [PATCH|PUT] /users/:id/rooms/:id -- update the resource;
- [GET]  /users/:id/rooms/:id/message -- list rooms messages
- [POST] /users/:id/rooms/:id/message -- Send message to room
- [PATCH] /users/:id/rooms/:id/message/:id -- Edit a message
- [DELETE] /users/:id/rooms/:id/message/:id -- Delete a message
