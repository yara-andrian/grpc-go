package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/yara-andrian/go-grpc-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Calculator client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// fmt.Printf("Created clienc: %f", c)

	// doUnary(c)

	doServerStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SUM unary rpc...")

	req := &calculatorpb.SumRequest{
		FirstNumber:  32,
		SecondNumber: 64,
	}

	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("err while calling Greeting RPC: %v", err)
	}

	log.Printf("Response from Sum: %v", res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a PrimeDecomposition Server Streaming RPC ...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		PrimeNumber: 1200,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling PrimeDecomposition RPC: %v", err)
	}

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}

		fmt.Println(res.GetPrimeResult())
	}
}
