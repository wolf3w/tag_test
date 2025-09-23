package main_test

import (
	"testing"

	"google.golang.org/grpc"
)

func TestClient(t *testing.T) {
	addr := "localhost:8087"
	client, err := grpc.NewClient(addr)
}
