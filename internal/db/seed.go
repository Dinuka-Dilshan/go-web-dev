package db

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Dinuka-Dilshan/go-web-dev/internal/store"
)

var userNames = []string{
	"alice", "bob", "charlie", "david", "emma",
	"frank", "grace", "harry", "isabella", "jack",
	"karen", "liam", "mia", "noah", "olivia",
	"peter", "quinn", "rachel", "sam", "tina",
	"ursula", "victor", "wendy", "xavier", "yara",
	"zach", "aryan", "bianca", "carl", "diana",
	"elijah", "fiona", "george", "hannah", "ivan",
	"jasmine", "kevin", "lara", "mason", "nina",
	"oscar", "paula", "ron", "sophia", "tom",
	"uma", "violet", "william", "yasmin", "zeke",
}

type Post struct {
	Title   string
	Content string
	Tags    []string
}

var postContents = []Post{
	{Title: "Introduction to Go", Content: "A beginner-friendly overview of the Go programming language and its core features.", Tags: []string{"go", "basics", "programming"}},
	{Title: "Understanding Goroutines", Content: "Exploring lightweight concurrency in Go and how goroutines work behind the scenes.", Tags: []string{"go", "concurrency", "goroutines"}},
	{Title: "Mastering Channels", Content: "A detailed guide on Go channels and their role in concurrent communication.", Tags: []string{"go", "channels", "concurrency"}},
	{Title: "REST API with Go", Content: "How to build a scalable and clean REST API using Go’s standard library.", Tags: []string{"api", "web", "go"}},
	{Title: "Error Handling Patterns", Content: "Best practices for error handling in Go and avoiding common pitfalls.", Tags: []string{"go", "errors", "clean-code"}},

	{Title: "Working with Interfaces", Content: "Understanding interfaces and how Go uses them for polymorphism.", Tags: []string{"go", "interfaces", "oop"}},
	{Title: "Structs and Methods", Content: "How structs model data and how methods add behavior to them.", Tags: []string{"go", "structs", "methods"}},
	{Title: "Using the Go Module System", Content: "Managing dependencies with Go modules step-by-step.", Tags: []string{"go", "modules", "dependency-management"}},
	{Title: "Building CLI Tools", Content: "How to create powerful command-line tools in Go.", Tags: []string{"cli", "go", "tools"}},
	{Title: "Working with JSON", Content: "Encoding and decoding JSON using Go’s encoding/json package.", Tags: []string{"json", "go", "serialization"}},

	{Title: "Understanding Slices", Content: "Deep dive into slices, underlying arrays, and memory behavior.", Tags: []string{"go", "slices", "memory"}},
	{Title: "Maps in Go", Content: "How Go maps work and how to use them effectively.", Tags: []string{"go", "maps", "data-structures"}},
	{Title: "Testing in Go", Content: "Writing unit tests in Go using the built-in testing package.", Tags: []string{"testing", "go", "tdd"}},
	{Title: "Benchmarking Code", Content: "Measuring performance with Go benchmarks.", Tags: []string{"performance", "benchmarking", "go"}},
	{Title: "Building Microservices", Content: "Using Go to build distributed, scalable microservices.", Tags: []string{"microservices", "go", "architecture"}},

	{Title: "Working with Databases", Content: "Connecting Go applications with PostgreSQL and MySQL.", Tags: []string{"database", "go", "sql"}},
	{Title: "Implementing Clean Architecture", Content: "Applying clean architecture principles in Go applications.", Tags: []string{"clean-architecture", "go", "design"}},
	{Title: "JWT Authentication in Go", Content: "Implementing secure authentication using JWT.", Tags: []string{"security", "jwt", "go"}},
	{Title: "Building Middleware", Content: "Creating custom middleware for logging, authentication, and validation.", Tags: []string{"middleware", "go", "web"}},
	{Title: "File Handling Basics", Content: "Reading and writing files using the os and io packages.", Tags: []string{"files", "go", "io"}},

	{Title: "Working with Context", Content: "Using context for deadlines, cancellation, and request-scoped values.", Tags: []string{"context", "go", "concurrency"}},
	{Title: "Optimizing Go Performance", Content: "Insights into profiling and optimizing Go code.", Tags: []string{"performance", "profiling", "go"}},
	{Title: "Go Memory Management", Content: "Explaining garbage collection and memory allocation in Go.", Tags: []string{"memory", "go", "internals"}},
	{Title: "Reflection in Go", Content: "When and how to use reflection responsibly.", Tags: []string{"reflection", "go"}},
	{Title: "Building WebSockets in Go", Content: "Real-time communication with Go’s WebSocket libraries.", Tags: []string{"websocket", "go", "real-time"}},

	{Title: "Creating a Task Scheduler", Content: "Writing a cron-like task scheduler in Go.", Tags: []string{"scheduler", "go", "tasks"}},
	{Title: "Introduction to Generics", Content: "How Go 1.18+ generics work and practical use cases.", Tags: []string{"generics", "go", "advanced"}},
	{Title: "Logging Best Practices", Content: "How to implement structured and contextual logging.", Tags: []string{"logging", "go", "observability"}},
	{Title: "Graceful Shutdowns", Content: "Ensuring services stop safely using context and signals.", Tags: []string{"shutdown", "go", "services"}},
	{Title: "Using Go Routines Safely", Content: "Avoiding race conditions and ensuring safe concurrency.", Tags: []string{"concurrency", "goroutines", "safety"}},

	{Title: "API Rate Limiting", Content: "Building a simple rate limiter using middleware and Redis.", Tags: []string{"rate-limit", "api", "go"}},
	{Title: "Dockerizing Go Applications", Content: "Packaging Go services in lightweight containers.", Tags: []string{"docker", "go", "deployment"}},
	{Title: "Deploying Go Apps to AWS Lambda", Content: "How to build and deploy Go functions to AWS.", Tags: []string{"aws", "lambda", "go"}},
	{Title: "Using Message Queues", Content: "Working with Kafka and RabbitMQ in Go.", Tags: []string{"queue", "go", "messaging"}},
	{Title: "Building GraphQL APIs", Content: "Creating a GraphQL API using Go libraries.", Tags: []string{"graphql", "api", "go"}},

	{Title: "Event-Driven Architecture", Content: "Implementing event-driven patterns in Go applications.", Tags: []string{"architecture", "events", "go"}},
	{Title: "Using Redis in Go", Content: "Caching and pub/sub essentials with Redis.", Tags: []string{"redis", "go", "cache"}},
	{Title: "Implementing Webhooks", Content: "Building and securing webhook endpoints.", Tags: []string{"webhooks", "api", "go"}},
	{Title: "Building a URL Shortener", Content: "A practical project using Go, Redis, and routing.", Tags: []string{"project", "web", "go"}},
	{Title: "Working with Time in Go", Content: "Parsing, formatting, and manipulating time values.", Tags: []string{"time", "go"}},

	{Title: "Go Project Structure", Content: "Recommended folder structure for scalable Go applications.", Tags: []string{"architecture", "go", "structure"}},
	{Title: "Creating a Custom Router", Content: "Implementing a minimal HTTP router from scratch.", Tags: []string{"router", "http", "go"}},
	{Title: "Validating Input Data", Content: "Techniques for robust request validation.", Tags: []string{"validation", "go"}},
	{Title: "Secure Password Storage", Content: "Using bcrypt and best practices for password hashing.", Tags: []string{"security", "passwords", "go"}},
	{Title: "Using Third-Party Libraries", Content: "How to evaluate and integrate external Go packages.", Tags: []string{"libraries", "go"}},

	{Title: "CI/CD for Go Applications", Content: "Setting up automated pipelines using GitHub Actions.", Tags: []string{"ci-cd", "automation", "go"}},
	{Title: "Understanding Go Races", Content: "Detecting and fixing race conditions using the race detector.", Tags: []string{"concurrency", "races", "go"}},
	{Title: "Building a Chat App", Content: "A real-time chat server using WebSockets.", Tags: []string{"chat", "real-time", "go"}},
	{Title: "Handling Large Files", Content: "Efficient streaming and processing techniques.", Tags: []string{"files", "streaming", "go"}},
	{Title: "Building a Job Queue", Content: "Creating a worker pool and background processing system.", Tags: []string{"workers", "queue", "go"}},
}

var generatedComments = []string{
	"Great explanation, very helpful!",
	"This was really insightful.",
	"I learned something new today.",
	"Please write more on this topic.",
	"Your examples are easy to follow.",
	"This clarified a lot of confusion I had.",
	"Amazing content as always.",
	"Simple and well explained.",
	"This helped me fix a bug in my code.",
	"More content like this, please.",
	"Clear and concise writing.",
	"One of the best explanations I’ve found.",
	"This should be more widely shared.",
	"Loved the breakdown of the concept.",
	"Thanks for sharing your knowledge.",
	"Helped me understand Go better.",
	"Very useful for beginners.",
	"Great job, keep it up!",
	"Well structured and informative.",
	"Exactly what I was looking for.",
}

func Seed(store store.Storage) error {

	ctx := context.Background()
	users := generateUsers(100)

	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			fmt.Println(err)
		}
	}

	posts := generatePosts(250, users)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			fmt.Println(err)
		}
	}

	comments := generateComments(250, users, posts)

	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("seed done")
	return nil

}

func generateUsers(num int) []*store.User {

	users := make([]*store.User, num)

	for i := range num {
		userName := userNames[rand.Intn(len(userNames))] + fmt.Sprintf("%v", i)
		users[i] = &store.User{
			UserName: userName,
			Email:    fmt.Sprintf("%v@example.com", userName),
			Password: "123456",
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {

	posts := make([]*store.Post, num)

	for i := range num {
		postContent := postContents[rand.Intn(len(postContents))]
		posts[i] = &store.Post{
			Content: postContent.Content,
			Title:   postContent.Title,
			Tags:    postContent.Tags,
			UserId:  users[rand.Intn(len(users))].ID,
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {

	comments := make([]*store.Comment, num)

	for i := range num {

		comments[i] = &store.Comment{
			Content: generatedComments[rand.Intn(len(generatedComments))],
			PostId:  posts[rand.Intn(len(posts))].ID,
			UserId:  users[rand.Intn(len(users))].ID,
		}
	}

	return comments
}
