package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"

	"github.com/urfave/cli/v3"
)

type Workout struct {
	id        string
	name      string
	unix      int64
	notes     string
	exercises []Exercise
}

type Exercise struct {
	id     int
	name   string
	notes  string
	sets   int
	reps   []int
	values []int
	unit   string
}

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
					filePath := filepath.Join(dirPath, "gym.env")

					// write to env file
					err = os.WriteFile(filePath, []byte("DYNAMODB_TABLENAME="+input), 0666)
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:  "add",
				Usage: "add new workout to database",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					var newWorkout Workout

					newWorkout.fill()

					newWorkout.print()

					fmt.Println("Upload the above workout? (y/n) ")
					if confirm() {
						newWorkout.upload()
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

func (workout *Workout) fill() {
	reader := bufio.NewReader(os.Stdin)

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Printf("Failed to generate workout ID: %v\n", err)
		return
	}
	workout.id = hex.EncodeToString(b)

	for workout.name == "" {
		fmt.Print("Workout Name: ")
		temp, _ := reader.ReadString('\n')
		workout.name = strings.TrimSpace(temp)
	}

	workout.unix = time.Now().Unix()

	fmt.Print("Workout Notes: ")
	temp, _ := reader.ReadString('\n')
	workout.notes = strings.TrimSpace(temp)

	done := false
	id := 1
	for !done {
		var exercise Exercise

		exercise.id = id
		id++

		for exercise.name == "" {
			fmt.Print("Exercise Name: ")
			temp, _ := reader.ReadString('\n')
			exercise.name = strings.TrimSpace(temp)
		}

		fmt.Print("Exercise Notes: ")
		temp, _ := reader.ReadString('\n')
		exercise.notes = strings.TrimSpace(temp)

		var setString string
		exercise.sets = 0
		for exercise.sets == 0 {
			fmt.Print("Exercise Sets: ")
			temp, _ := reader.ReadString('\n')
			setString = strings.TrimSpace(temp)

			sets, _ := strconv.Atoi(setString)
			exercise.sets = sets
		}

		var repString string
		for repString == "" {
			fmt.Print("Exercise Reps: ")
			temp, _ := reader.ReadString('\n')
			repString = strings.TrimSpace(temp)

			split := strings.Split(repString, ",")
			reps := []int{}
			for _, str := range split {
				rep, err := strconv.Atoi(str)
				if err != nil {
					repString = ""
					break
				}
				reps = append(reps, rep)
			}
			if repString == "" {
				continue
			}

			if len(reps) == 1 {
				for i := 0; i < exercise.sets; i++ {
					exercise.reps = append(exercise.reps, reps[0])
				}
			} else if len(reps) == exercise.sets {
				for i := 0; i < exercise.sets; i++ {
					exercise.reps = append(exercise.reps, reps[i])
				}
			} else {
				repString = ""
				continue
			}
		}

		var valueString string
		for valueString == "" {
			fmt.Print("Exercise Values: ")
			temp, _ := reader.ReadString('\n')
			valueString = strings.TrimSpace(temp)

			split := strings.Split(valueString, ",")
			values := []int{}
			for _, str := range split {
				value, err := strconv.Atoi(str)
				if err != nil {
					valueString = ""
					break
				}
				values = append(values, value)
			}
			if valueString == "" {
				continue
			}

			if len(values) == 1 {
				for i := 0; i < exercise.sets; i++ {
					exercise.values = append(exercise.values, values[0])
				}
			} else if len(split) == exercise.sets {
				for i := 0; i < exercise.sets; i++ {
					exercise.values = append(exercise.values, values[i])
				}
			} else {
				valueString = ""
				continue
			}
		}

		for exercise.unit == "" {
			fmt.Print("Exercise Units: ")
			temp, _ := reader.ReadString('\n')
			exercise.unit = strings.TrimSpace(temp)
		}

		workout.exercises = append(workout.exercises, exercise)

		fmt.Println("Add Another? (y/n) ")
		if !confirm() {
			done = true
		}
	}

}

func (workout *Workout) print() {
	fmt.Println("Workout Name:", workout.name)

	fmt.Println("\tDate:", time.Unix(workout.unix, 0))

	if workout.notes != "" {
		fmt.Println("\tNotes:", workout.notes)
	}

	for _, exercise := range workout.exercises {
		fmt.Println("Exercise Name:", exercise.name)

		if exercise.notes != "" {
			fmt.Println("\tNotes:", exercise.notes)
		}

		for i := 0; i < exercise.sets; i++ {
			fmt.Println("\tSet", strconv.Itoa(i+1)+":", exercise.values[i], exercise.unit, "x", exercise.reps[i])
		}
	}
}

func (workout *Workout) upload() {
	fmt.Println("uploading workout...")
}
