package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type Problem struct {
	question string
	answer   string
}

func prepareQuiz(fileName string) []Problem {
	fmt.Println(fileName)

	// Open the csvFile
	csvFile, err := os.Open(fileName)
	// If an error occurs print it and exit
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(csvFile)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	quiz := make([]Problem, len(records))

	for i := range records {
		quiz[i].question = records[i][0]
		quiz[i].answer = records[i][1]
	}
	return quiz

}

func doQuiz(quiz []Problem, c chan []string) {
	userAns := make([]string, len(quiz))
	for i := range quiz {
		fmt.Printf("Problem #%d: %s = ", i+1, quiz[i].question)
		fmt.Scanln(&userAns[i])
	}
	c <- userAns
}

func ask(p Problem, i int, c chan string) {
	var answer string
	fmt.Printf("Problem #%d: %s = ", i, p.question)
	fmt.Scanf("%s\n", &answer)
	c <- answer
}

func main() {
	// A flag to define the csv filename
	csvFilename := flag.String("f", "problems.csv", "Define the filename for the input quiz")
	timeLimit := flag.Int("l", 30, "time limit to answer the quiz")
	flag.Parse()

	quiz := prepareQuiz(*csvFilename)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	totalGoodQ := 0

	for i, p := range quiz {
		c := make(chan string)
		go ask(p, i+1, c)
		select {
		case <-timer.C:
			fmt.Println("\nTotal Questions: ", len(quiz))
			fmt.Println("Total Good Questions: ", totalGoodQ)
			return
		case answer := <-c:
			if answer == p.answer {
				totalGoodQ++
			}
		}

	}

	fmt.Println("Total Questions: ", len(quiz))
	fmt.Println("Total Good Questions: ", totalGoodQ)

}
