package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"github.com/Tsan1024/LLM/generate/streamdemo"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := streamdemo.NewQuestionAnswerClient(conn)

	stream, err := client.Ask(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	questions := []string{"What is your name?", "How are you?", "Where are you from?"}

	go func() {
		for _, q := range questions {
			err := stream.Send(&streamdemo.Question{Text: q})
			if err != nil {
				log.Fatalf("Error sending question: %v", err)
			}
		}
		stream.CloseSend()
	}()

	for {
		answer, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("Error receiving answer: %v", err)
		}
		log.Println(answer.Text)
	}
}
