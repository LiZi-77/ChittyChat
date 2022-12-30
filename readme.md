### ChittyChat

To run this program,

- Start a server by running 

  ```
  cd Server
  go run server.go
  ```

  This will create a server listening at port 8080.

  The server has a connection list, which maintains all the clients info, including the username and the status of that user

- start some clients

  ```
  cd Client
  go run client.go -UserName user1
  ```

  This will create a client whose name is user1, if you don't want to use -UserName, the name of this user will be anonymous.

- Join/ Leave/ Publish

  after creating clients, you can input several command in these clients' terminal

  ```Go
  Join // this will add this user to the server connection list
  Leave // this will change the status of this client's connection to false, which means it will not receive the published messages from other clients.
  
  some message... // any other input will be a message which will be broadcast to all clients.
  ```

  

here's a simple example

Server side:

```
$  cd Server
$  go run server.go

Server running at port :8080
Broadcasting message: [[user1] joined Chitty-Chat (Lamport time: 1)] (Lamport time: 2) 
Broadcasting message: [[user2] joined Chitty-Chat (Lamport time: 3)] (Lamport time: 4) 
Broadcasting message: [[user2] joined Chitty-Chat (Lamport time: 3)] (Lamport time: 5) 
Broadcasting message: [[user1]: User1 want to say hello! (Lamport time: 6)] (Lamport time: 7) 
Broadcasting message: [[user1]: User1 want to say hello! (Lamport time: 6)] (Lamport time: 8) 
Broadcasting message: [[user2]: user2 want to say hi to user1 too. (Lamport time: 9)] (Lamport time: 10) 
Broadcasting message: [[user2]: user2 want to say hi to user1 too. (Lamport time: 9)] (Lamport time: 11) 
Broadcasting message: [[user2] left Chitty-Chat. (Lamport time: 12)] (Lamport time: 13) 
Broadcasting message: [[user2] left Chitty-Chat. (Lamport time: 12)] (Lamport time: 14) 
```



User1:

```
$  cd Client
$  go run client.go -UserName user1
Join
[user1] joined Chitty-Chat (Lamport time: 1) 
[user2] joined Chitty-Chat (Lamport time: 3) 
User1 want to say hello!
[user1]: User1 want to say hello! (Lamport time: 6) 
[user2]: user2 want to say hi to user1 too. (Lamport time: 9) 
[user2] left Chitty-Chat. (Lamport time: 12) 
```



User2:

```
$  cd Client
$  go run client.go -UserName user2
Join
[user2] joined Chitty-Chat (Lamport time: 3) 
[user1]: User1 want to say hello! (Lamport time: 6) 
user2 want to say hi to user1 too.
[user2]: user2 want to say hi to user1 too. (Lamport time: 9) 
Leave
[user2] left Chitty-Chat. (Lamport time: 12) 
```



and here's the related log info

```
2022/12/30 20:08:50 Server running at port :8080
2022/12/30 20:09:15 [user1] joined Chitty-Chat (Lamport time: 1)
2022/12/30 20:09:15 Broadcasting message:  [user1] joined Chitty-Chat (Lamport time: 1) (Lamport time:  2 )
2022/12/30 20:09:50 [user2] joined Chitty-Chat (Lamport time: 3)
2022/12/30 20:09:50 Broadcasting message:  [user2] joined Chitty-Chat (Lamport time: 3) (Lamport time:  4 )
2022/12/30 20:09:50 Broadcasting message:  [user2] joined Chitty-Chat (Lamport time: 3) (Lamport time:  5 )
2022/12/30 20:10:33 [ user1 ] publish a new message User1 want to say hello! (Lamport time:  6 )
2022/12/30 20:10:33 Broadcasting message:  [user1]: User1 want to say hello! (Lamport time: 6) (Lamport time:  7 )
2022/12/30 20:10:33 Broadcasting message:  [user1]: User1 want to say hello! (Lamport time: 6) (Lamport time:  8 )
2022/12/30 20:10:57 [ user2 ] publish a new message user2 want to say hi to user1 too. (Lamport time:  9 )
2022/12/30 20:10:57 Broadcasting message:  [user2]: user2 want to say hi to user1 too. (Lamport time: 9) (Lamport time:  10 )
2022/12/30 20:10:57 Broadcasting message:  [user2]: user2 want to say hi to user1 too. (Lamport time: 9) (Lamport time:  11 )
2022/12/30 20:11:09 [user2] left Chitty-Chat (Lamport time: 12 )
2022/12/30 20:11:09 Broadcasting message:  [user2] left Chitty-Chat. (Lamport time: 12) (Lamport time:  13 )
2022/12/30 20:11:09 Broadcasting message:  [user2] left Chitty-Chat. (Lamport time: 12) (Lamport time:  14 )

```

