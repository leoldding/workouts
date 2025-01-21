package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Workout struct {
	Id        string     `dynamodbav:"id"`
	Name      string     `dynamodbav:"name"`
	Unix      int64      `dynamodbav:"unix"`
	Notes     string     `dynamodbav:"notes"`
	Exercises []Exercise `dynamodbav:"exercises"`
}

func (workout *Workout) fill() {
	reader := bufio.NewReader(os.Stdin)

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Printf("Failed to generate workout ID: %v\n", err)
		return
	}
	workout.Id = hex.EncodeToString(b)

	for workout.Name == "" {
		fmt.Print("Workout Name: ")
		temp, _ := reader.ReadString('\n')
		workout.Name = strings.TrimSpace(temp)
	}

	workout.Unix = time.Now().Unix()

	fmt.Print("Workout Notes: ")
	temp, _ := reader.ReadString('\n')
	workout.Notes = strings.TrimSpace(temp)

	done := false
	id := 1
	for !done {
		var exercise Exercise

		exercise.Id = id
		id++

		for exercise.Name == "" {
			fmt.Print("Exercise Name: ")
			temp, _ := reader.ReadString('\n')
			exercise.Name = strings.TrimSpace(temp)
		}

		fmt.Print("Exercise Notes: ")
		temp, _ := reader.ReadString('\n')
		exercise.Notes = strings.TrimSpace(temp)

		var setString string
		exercise.Sets = 0
		for exercise.Sets == 0 {
			fmt.Print("Exercise Sets: ")
			temp, _ := reader.ReadString('\n')
			setString = strings.TrimSpace(temp)

			sets, _ := strconv.Atoi(setString)
			exercise.Sets = sets
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
				for i := 0; i < exercise.Sets; i++ {
					exercise.Reps = append(exercise.Reps, reps[0])
				}
			} else if len(reps) == exercise.Sets {
				for i := 0; i < exercise.Sets; i++ {
					exercise.Reps = append(exercise.Reps, reps[i])
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
				for i := 0; i < exercise.Sets; i++ {
					exercise.Values = append(exercise.Values, values[0])
				}
			} else if len(split) == exercise.Sets {
				for i := 0; i < exercise.Sets; i++ {
					exercise.Values = append(exercise.Values, values[i])
				}
			} else {
				valueString = ""
				continue
			}
		}

		for exercise.Unit == "" {
			fmt.Print("Exercise Units: ")
			temp, _ := reader.ReadString('\n')
			exercise.Unit = strings.TrimSpace(temp)
		}

		workout.Exercises = append(workout.Exercises, exercise)

		fmt.Println("Add Another? (y/n) ")
		if !confirm() {
			done = true
		}
	}

}

func (workout *Workout) print() {
	fmt.Println("Workout Name:", workout.Name)

	fmt.Println("\tDate:", time.Unix(workout.Unix, 0))

	if workout.Notes != "" {
		fmt.Println("\tNotes:", workout.Notes)
	}

	for _, exercise := range workout.Exercises {
		fmt.Println("Exercise Name:", exercise.Name)

		if exercise.Notes != "" {
			fmt.Println("\tNotes:", exercise.Notes)
		}

		for i := 0; i < exercise.Sets; i++ {
			fmt.Println("\tSet", strconv.Itoa(i+1)+":", exercise.Values[i], exercise.Unit, "x", exercise.Reps[i])
		}
	}
}
