	package main

	import (
		// "bufio"
		"encoding/csv"
		"flag"
		"fmt"
		"os"
		"strings"
		"time"
	)

	func main() {
		csvFile := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
		timeLimit := flag.Int("limit", 10, "total time limit for the quiz in seconds")
		flag.Parse()

		file, err := os.Open(*csvFile)
		if err != nil {
			fmt.Println("Failed to open the CSV file:", err)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			fmt.Println("Failed to read the CSV file:", err)
			return
		}

		score := 0
		timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	problemLoop:
		for i, row := range records {
			answerCh:= make(chan string,1)
			fmt.Printf("Problem #%d: %s = ?\n", i+1, row[0])
			go func() {
				var answer string
				fmt.Scanln(&answer)
				answerCh <- strings.TrimSpace(answer) //closure cz here we are using a varaible from the outer scope
			}()

			select {
			case <-timer.C:
				fmt.Println("\nTime's up!")
				break problemLoop

			case answer := <-answerCh:
				if answer == row[1] {
					fmt.Println("Correct!")
					score++
				} else {
					fmt.Printf("Wrong! The correct answer is %s\n", row[1])
				}
			}
		}

		fmt.Printf("\nYou scored %d out of %d\n", score, len(records))
	}
