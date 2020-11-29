package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/yara-andrian/go-grpc-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	// doServerStreaming(c)

	// doClientStreaming(c)

	// doBidiStreaming(c)
	doErrorUnary(c)
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

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do client streaming RPC")

	numbers := []int32{3, 4, 9, 34}

	stream, err := c.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("error")
	}

	for _, number := range numbers {
		fmt.Printf("sending req: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("error while receiving response from ComputeAverage: %v", err)
	}

	fmt.Printf("ComputeAverage Response: %v\n", res)
}

func doBidiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Executing Bidi function")

	stream, err := c.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream and find maximum: %v", err)
	}

	waitc := make(chan struct{})

	// send data to server
	go func() {
		numbers := []int32{1, 2, 3, 4, 5}

		for _, number := range numbers {
			fmt.Printf("The number: %v\n", number)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive data from server
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error in receiving server stream: %v", res)
				break
			}

			maximum := res.GetMaximum()
			fmt.Printf("Maximum is: %v\n", maximum)
		}
		close(waitc)
	}()

	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Error unary rpc...")
	// success case
	doErrorCall(c, 10)
	// error case
	doErrorCall(c, -10)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})

	if err != nil {
		respErr, ok := status.FromError(err)

		if ok {
			// actual error from gRPC
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent negative number")
				return
			}
		} else {
			log.Fatalf("Big error calling square root %v\n", respErr)
			return
		}
	}

	fmt.Printf("Result of square root of %v: %v\n", n, res.GetNumberRoot())
}
