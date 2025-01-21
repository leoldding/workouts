package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/term"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "setup",
				Usage: "write database name to env file",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					// get path
					exePath, err := os.Executable()
					if err != nil {
						log.Fatal(err)
					}
					dirPath := filepath.Dir(exePath)
					fmt.Println(exePath)

					// check that dynamodb table name exists
					input := cmd.Args().First()
					if input == "" {
						log.Fatal(errors.New("database name cannot be empty"))
					}

					dirPath = filepath.Join(dirPath, "../envs")
					if _, err := os.Stat(dirPath); os.IsNotExist(err) {
						// directory doesn't exist, create it
						err := os.Mkdir(dirPath, 0755)
						if err != nil {
							fmt.Println("Error creating directory:", err)
						}
					}
					filePath := filepath.Join(dirPath, "workouts.env")

					// write to env file
					err = os.WriteFile(filePath, []byte("DYNAMODB_TABLENAME="+input), 0666)
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add new workout to database",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					tableName := getTableName()

					// create dynamodb client
					cfg, err := config.LoadDefaultConfig(context.TODO())
					if err != nil {
						log.Fatal(err)
					}
					client := dynamodb.NewFromConfig(cfg)

					var newWorkout Workout

					newWorkout.fill()

					newWorkout.print()

					fmt.Println("Upload the above workout? (y/n) ")
					if confirm() {
						upload(newWorkout, client, tableName)
					}

					return nil
				},
			},
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func upload(workout Workout, client *dynamodb.Client, tableName string) {
	av, err := attributevalue.MarshalMap(workout)
	if err != nil {
		log.Fatal("failed to marshal Record, %w", err)
	}

	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Workout has been uploaded successfully")
}

func getTableName() string {
	// get path
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dirPath := filepath.Dir(exePath)
	filePath := filepath.Join(dirPath, "../envs/workouts.env")
	// load table name
	err = godotenv.Load(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return os.Getenv("DYNAMODB_TABLENAME")
}

func confirm() bool {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	reader := os.Stdin
	buffer := make([]byte, 1)

	for {
		_, err := reader.Read(buffer)
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}

		char := buffer[0]

		if char == 'y' {
			return true
		} else if char == 'n' {
			return false
		}
	}
	return false
}
