package main

import (
	"context"
	"fmt"
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

	doUnary(c)
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
