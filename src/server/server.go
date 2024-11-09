package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Tsan1024/LLM/generate/streamdemo"
)

type questionAnswerServer struct {
	streamdemo.UnimplementedQuestionAnswerServer
}

// mustEmbedUnimplementedQuestionAnswerServer implements streamdemo.Question

// mustEmbedUnimplementedQuestionAnswerServer implements streamdemo.QuestionAnswerServer.

func (s *questionAnswerServer) Ask(stream streamdemo.QuestionAnswer_AskServer) error {
	for {
		question, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				return nil
			}
			return err
		}

		answer := "Answer to: " + question.Text
		err = stream.Send(&streamdemo.Answer{Text: answer})
		if err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	streamdemo.RegisterQuestionAnswerServer(s, &questionAnswerServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
