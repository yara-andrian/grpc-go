package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/yara-andrian/go-grpc-course/blog/blogpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// Create Blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Andrian",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}

	blogResponse, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error during creation: %v \n", err)
	}
	fmt.Printf("Blog has been created: %v \n", blogResponse)
	blogID := blogResponse.GetBlog().GetId()

	// Read Blog
	fmt.Println("Read the blog")

	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "5fc8e480b715938c2ef19298"})

	if err2 != nil {
		fmt.Printf("Error happened while reading2: %v \n", err2)
	}

	readBlogReq := &blogpb.ReadBlogRequest{
		BlogId: blogID,
	}

	readBlogRes, readBlogErr := c.ReadBlog(context.Background(),
		readBlogReq,
	)

	if readBlogErr != nil {
		fmt.Printf("Error happened while reading3: %v", readBlogErr)
	}

	fmt.Printf("Blog was read: %v", readBlogRes)

	// Update Blog
	fmt.Println("Update Blog")

	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "New Author",
		Title:    "My First Blog (edited)",
		Content:  "Content of the frist blog, with some awesome additions",
	}

	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})

	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v \n", updateErr)
	}

	fmt.Printf("Blog was updated: %v\n", updateRes)

	// delete Blog
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})

	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v \n", updateErr)
	}
	fmt.Printf("Blog was deleted: %v \n", deleteRes)

	// list blog
	fmt.Println("List Blog started")
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})

	if err != nil{
		log.Fatalf("Error while calling ListBlog RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF{
			fmt.Println("Nothing")
			break
		}

		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}
}
