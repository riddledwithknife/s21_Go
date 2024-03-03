package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Database struct {
	knowNodes         []*DatabaseInstance
	replicationFactor int
}

type DatabaseInstance struct {
	host string
	port string
}

func (db *DatabaseInstance) runInstance() {
	// Placeholder for actual logic
}

func main() {
	host := flag.String("host", "127.0.0.1", "Database host")
	port := flag.Int("port", 8080, "Database port")
	flag.Parse()

	db := Database{
		replicationFactor: 2,
	}

	var wg sync.WaitGroup
	for i := 0; i <= db.replicationFactor; i++ {
		instance := DatabaseInstance{
			host: *host,
			port: strconv.Itoa(*port + i),
		}
		db.knowNodes = append(db.knowNodes, &instance)
		wg.Add(1)
		go func() {
			defer wg.Done()
			instance.runInstance()
		}()
	}
	wg.Wait()

	fmt.Printf("Connected to database of Warehouse 13 at %s:%d\n", *host, *port)
	fmt.Println("Known nodes:")
	for _, node := range db.knowNodes {
		fmt.Printf("%s:%s\n", node.host, node.port)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		parts := strings.Split(input, " ")
		command := strings.ToUpper(parts[0])

		switch command {
		case "SET":
			pattern := `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
			matched, _ := regexp.MatchString(pattern, parts[1])
			if !matched {
				fmt.Println("Error: key is not a proper UUID4")
				continue
			}
			// Logic for SET command
		case "GET":
			fmt.Println("GET")
		case "DELETE":
			fmt.Println("DELETE")
		case "EXIT":
			fmt.Println("EXIT")
			return
		default:
			fmt.Println("Unknown command")
		}
	}
}
