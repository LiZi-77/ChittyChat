

// Package main implements a simple gRPC client that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It interacts with the route guide service whose definition can be found in routeguide/route_guide.pb.
package main

import (
	"bufio"
	"crypto/sha256"
	"flag"
	"fmt"
	"os"

	"encoding/hex"
	"log"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "chittychat"
)

var participant pb.ChittyChatClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

func connect(user *pb.User) error {
	var streamError error
	fmt.Println(*user)
	stream, err := participant.Join(context.Background(), &pb.Connect{
		User:   user,
		Active: true,
	})

	if err != nil {
		return fmt.Errorf("Connect failed: %v", err)
	}

	wait.Add(1)

	go func(str pb.ChittyChat_JoinClient) {
		defer wait.Done()

		for {
			msg, err := str.Recv()

			if err != nil {
				streamError = fmt.Errorf("Error reading message: %v", err)
				break
			}

			//         fmt.Printf("%v : %s\n", msg.User.DisplayName, msg.Message)
			fmt.Printf("%v \n", msg.Message)
		}
	}(stream)

	return streamError
}

func main() {

	var isConnected bool
	isConnected = false
	err3 := error(nil)
	var dummymsg *pb.Close

	ts := time.Now()
	done := make(chan int)

	name := flag.String("N", "Anonymous", "")
	flag.Parse()

	id := sha256.Sum256([]byte(ts.String() + *name))
	user := &pb.User{
		Id:          hex.EncodeToString(id[:]),
		DisplayName: *name,
	}

	wait.Add(1)

	go func() {
		defer wait.Done()
		scanner := bufio.NewScanner(os.Stdin)
		ts := time.Now()
		msgID := sha256.Sum256([]byte(ts.String() + *name))
		for scanner.Scan() {
			msg := &pb.Message{
				Id:        hex.EncodeToString(msgID[:]),
				User:      user,
				Message:   scanner.Text(),
				Timestamp: ts.String(),
			}

			if msg.Message == "Join" && !isConnected {
				fmt.Printf("Message = Join: %v", msg.Message)
				conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
				if err != nil {
					log.Fatalf("Could not connect to server %v", err)
				}
				defer conn.Close()
				participant = pb.NewChittyChatClient(conn)
				isConnected = true
				connect(user)
				wait.Add(1)
			} else {
				if isConnected {
					if msg.Message == "Leave" {
						isConnected = false
						dummymsg, err3 = participant.Leave(context.Background(), msg)
						if err3 != nil {
							fmt.Printf("Error leaving Chitty-Chat: %v", err3)
							break
						} else {
							//					break
						}
					} else {
						dummymsg, err3 = participant.Publish(context.Background(), msg)
					}
					if err3 != nil {
						fmt.Printf("Error sending this message: %v", err3)
						break
					}
				}
			}

			_, err := error(nil), error(nil)
			if err != nil {
				fmt.Printf("Error sending this message: %v", err)
				break
			}
		}
	}()

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
}
