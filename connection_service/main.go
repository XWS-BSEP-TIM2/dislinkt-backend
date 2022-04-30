package main

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service/startup"
	cfg "github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service/startup/config"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func main() {

	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()

	/*
		fmt.Println("Connection Service")
		helloWorld("bolt://localhost:7687", "neo4j", "connection")
		fmt.Println("END")

	*/

}

func helloWorld(uri, username, password string) (string, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer driver.Close()
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	greeting, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)",
			map[string]interface{}{"message": "SUPER222"})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return "", err
	}

	return greeting.(string), nil
}
