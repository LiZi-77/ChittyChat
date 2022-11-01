package main

import (
	pb "chittychat"
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"sync"

	"google.golang.org/grpc"
)

var grpcServer pb.ChittyChatServer

type Connection struct {
	stream pb.ChittyChat_JoinServer
	id     string
	username string
	active bool
	err    chan error
}

type Server struct {
	Connection []*Connection
	pb.UnimplementedChittyChatServer
	local_timestamp int64
}

func GetTimestamp(s *Server, i int64) int64 {
	l := float64(s.local_timestamp)
	i_ := float64(i)
	f := math.Max(l, i_) + 1
	return int64(f)
}

func logToFile() {
	// create if not exist
	file, err := os.OpenFile(fmt.Sprint("log.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	} 
	log.SetOutput(file)

}

func IsContain(items []Connection, item Connection) bool {
	for _, eachItem := range items {
		if eachItem.id == item.id {
			return true
		}
	}
	return false
}

func (s *Server) Join(pconn *pb.Connect, stream pb.ChittyChat_JoinServer) error {

	var msg pb.Message
	var ctx context.Context
	
	next_timestamp, _ := strconv.Atoi(msg.GetTimestamp())

	conn := &Connection{
		stream: stream,
		id:     pconn.User.Id,
		username: pconn.User.DisplayName,
		active: true,
		err:    make(chan error),
	}
	s.Connection = append(s.Connection, conn)
	str := strconv.FormatInt(s.local_timestamp, 10)
	
	msg.Message ="[" + conn.username + "] joined Chitty-Chat (Lamport time: " + str + ")"
	log.Println("[" + conn.username + "] joined Chitty-Chat (Lamport time: " + str + ")")
	s.local_timestamp = GetTimestamp(s, int64(next_timestamp))
	
	s.Broadcast(ctx, &msg)

	 return <-conn.err
}

func (s *Server) Broadcast(ctx context.Context, msg *pb.Message) (*pb.Close, error) {

	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range s.Connection {
		wait.Add(1)

		func(msg *pb.Message, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(msg)
				 
				fmt.Printf("Broadcasting message: [%v] (Lamport time: %v) \n", msg.Message, s.local_timestamp)
				
				log.Println("Broadcasting message: ", msg.Message, "(Lamport time: ", s.local_timestamp, ")")
				s.local_timestamp = GetTimestamp(s, s.local_timestamp)
				

				if err != nil {
					log.Fatalf("Error with stream %v. Error: %v", conn.stream, err)
					conn.active = false
					conn.err <- err
				}
			}
		}(msg, conn)
	}

	go func() {
		wait.Wait()
		close(done)
	}()

	//maybe timestamp++ here?

	<-done

	return &pb.Close{}, nil
}

func (s *Server) Publish(ctx context.Context, msg *pb.Message) (*pb.Close, error) {
	incoming_timestamp, _ := strconv.Atoi(msg.GetTimestamp())

	log.Println("[", msg.User.DisplayName, "] publish a new message", msg.Message, "(Lamport time: ", s.local_timestamp, ")")

	str := strconv.FormatInt(s.local_timestamp, 10)
	msg.Message = "[" + msg.User.DisplayName + "]: " + msg.Message + " (Lamport time: " + str + ")"

	s.local_timestamp = GetTimestamp(s, int64(incoming_timestamp)) 

	s.Broadcast(ctx, msg)
	return &pb.Close{}, nil
}

func (s *Server) Leave(ctx context.Context, msg *pb.Message) (*pb.Close, error) {
	incoming_timestamp, _ := strconv.Atoi(msg.GetTimestamp())

	str := strconv.FormatInt(s.local_timestamp, 10)
	msg.Message = "[" + msg.User.DisplayName + "] left Chitty-Chat. (Lamport time: " + str + ")"

	log.Println("[" + msg.User.DisplayName+"] left Chitty-Chat"+" (Lamport time:", s.local_timestamp, ")")
	s.local_timestamp = GetTimestamp(s, int64(incoming_timestamp))
	s.Broadcast(ctx, msg)

	for _, conn := range s.Connection {
		if conn.id == msg.User.DisplayName {
			conn.active = false
		}
	}
	return &pb.Close{}, nil
}

func main() {
	logToFile()
	var connections []*Connection
	var ThisBroadcastServer pb.UnimplementedChittyChatServer

	server := &Server{connections, ThisBroadcastServer, 0}
	server.local_timestamp++

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error creating the server %v", err)
	}

	log.Println("Server running at port :8080")
	fmt.Println("Server running at port :8080")

	pb.RegisterChittyChatServer(grpcServer, server)
	grpcServer.Serve(listener)
}
