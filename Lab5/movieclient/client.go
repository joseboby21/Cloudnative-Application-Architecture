// Package main imlements a client for movieinfo service
package main

import (
	"Labs/Lab5/movieapi"
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
)

const (
	address      = "localhost:50051"
	defaultTitle = "Pulp fiction"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := movieapi.NewMovieInfoClient(conn)

	// Contact the server and print out its response.
	title := defaultTitle
	if len(os.Args) > 1 {
		title = os.Args[1]
	}
	// Timeout if server doesn't respond
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	r, err := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: title})
	if err != nil {
		log.Fatalf("could not get movie info: %v", err)
	}
	log.Printf("Movie Info for %s %d %s %v", title, r.GetYear(), r.GetDirector(), r.GetCast())
	cancel()

	// Adding new movie detais to DB
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	r2, err := c.SetMovieInfo(ctx, &movieapi.MovieData{
		Title:    "Batman Begins",
		Year:     2005,
		Director: "Christopher Nolan",
		Cast:     []string{"Christian Bale", "Liam Neeson", "Gary Oldman", "Michael Caine"}})
	if err != nil {
		log.Fatalf("could not add movie info %s", err)
	}
	log.Printf("Addition of Movie to Database was %s \n", r2.Stat)
	cancel()

	// Trying to retrive the details of new movie added
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	r, err = c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: "Batman Begins"})
	if err != nil {
		log.Fatalf("could not get movie info: %v", err)
	}
	log.Printf("Movie Info for  Batman Begins %d %s %v", r.GetYear(), r.GetDirector(), r.GetCast())
	cancel()
}
