package main

import (
	"context"
	"fmt"
	"syscall"
	"testing"
	"time"

	pb "github.com/vladovsiychuk/demo-grpc/protob/discovery/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestDiscovery(t *testing.T) {
	go func() {
		main()
	}()

	time.Sleep(time.Second)

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()
	client := pb.NewDiscoveryServiceClient(conn)

	// Add posts 5 to 11
	for i := 5; i <= 11; i++ {
		_, err := client.AddPost(context.Background(), &pb.AddPostRequest{
			Post: &pb.AddPost{
				Owner:       "user",
				FrontPicUrl: fmt.Sprintf("front%d.jpg", i),
				BackPicUrl:  fmt.Sprintf("back%d.jpg", i),
			},
		})
		if err != nil {
			t.Fatalf("AddPost %d failed: %v", i, err)
		}
	}

	pageSize := 3
	var cursor string

	// First call: no cursor, expect 11,10,9
	resp, err := client.GetPosts(context.Background(), &pb.GetPostsRequest{})
	if err != nil || !resp.Success {
		t.Fatalf("GetPosts 1 failed: %v", err)
	}
	ids := []string{}
	for _, post := range resp.Data {
		ids = append(ids, post.Id)
	}
	cursor = resp.Cursor
	if len(ids) != pageSize {
		t.Errorf("Expected %d posts, got %d", pageSize, len(ids))
	}
	// Check filenames for first page
	for i, post := range resp.Data {
		expectedFront := fmt.Sprintf("front%d.jpg", 11-i)
		expectedBack := fmt.Sprintf("back%d.jpg", 11-i)
		if post.FrontPicUrl != expectedFront || post.BackPicUrl != expectedBack {
			t.Errorf("Page 1: Post %d: got FrontPicUrl=%s, BackPicUrl=%s; want %s, %s", i, post.FrontPicUrl, post.BackPicUrl, expectedFront, expectedBack)
		}
	}

	// Second call: use cursor, expect 8,7,6
	resp, err = client.GetPosts(context.Background(), &pb.GetPostsRequest{Cursor: cursor})
	if err != nil || !resp.Success {
		t.Fatalf("GetPosts 2 failed: %v", err)
	}
	ids = []string{}
	for _, post := range resp.Data {
		ids = append(ids, post.Id)
	}
	cursor = resp.Cursor
	if len(ids) != pageSize {
		t.Errorf("Expected %d posts, got %d", pageSize, len(ids))
	}
	for i, post := range resp.Data {
		expectedFront := fmt.Sprintf("front%d.jpg", 8-i)
		expectedBack := fmt.Sprintf("back%d.jpg", 8-i)
		if post.FrontPicUrl != expectedFront || post.BackPicUrl != expectedBack {
			t.Errorf("Page 2: Post %d: got FrontPicUrl=%s, BackPicUrl=%s; want %s, %s", i, post.FrontPicUrl, post.BackPicUrl, expectedFront, expectedBack)
		}
	}

	// Third call: use cursor, expect 5
	resp, err = client.GetPosts(context.Background(), &pb.GetPostsRequest{Cursor: cursor})
	if err != nil || !resp.Success {
		t.Fatalf("GetPosts 3 failed: %v", err)
	}
	ids = []string{}
	for _, post := range resp.Data {
		ids = append(ids, post.Id)
	}
	cursor = resp.Cursor
	if len(ids) != 1 {
		t.Errorf("Expected 1 post, got %d", len(ids))
	}
	if len(resp.Data) == 1 {
		post := resp.Data[0]
		if post.FrontPicUrl != "front5.jpg" || post.BackPicUrl != "back5.jpg" {
			t.Errorf("Page 3: got FrontPicUrl=%s, BackPicUrl=%s; want front5.jpg, back5.jpg", post.FrontPicUrl, post.BackPicUrl)
		}
	}

	// Fourth call: use cursor, expect no posts
	resp, err = client.GetPosts(context.Background(), &pb.GetPostsRequest{Cursor: cursor})
	if err != nil || !resp.Success {
		t.Fatalf("GetPosts 4 failed: %v", err)
	}
	if len(resp.Data) != 0 {
		t.Errorf("Expected 0 posts, got %d", len(resp.Data))
	}
	if resp.Cursor != "" {
		t.Errorf("Expected empty cursor, got %q", resp.Cursor)
	}

	_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
}
