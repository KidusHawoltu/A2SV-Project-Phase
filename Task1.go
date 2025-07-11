package main

import (
	"fmt"
)

type SubjectGrade struct {
	Name  string
	Grade float64
}

func calculateAvg(studentGrade map[string]SubjectGrade) float64 {
	total, count := 0.0, 0.0
	for _, subjectGrade := range studentGrade {
		total += subjectGrade.Grade
		count += 1
	}
	if count == 0 {
		return total
	}
	return total / count
}

func main() {
	studentName, subjectCount := "", 0
	fmt.Print("Enter You Name and number of subjects you took (e.g. Kidus 3): ")
	fmt.Scanln(&studentName, &subjectCount)

	subjectWidth := 12
	grades := make(map[string]SubjectGrade)
	for i := 0; i < subjectCount; i++ {
		n, g := "", 0.0
		fmt.Printf("Enter Name Subject %d (e.g. Math): ", i+1)
		fmt.Scanln(&n)
		for _, err := grades[n]; err; {
			fmt.Printf("You have already registered %q, Enter antoher Subject: ", n)
			fmt.Scanln(&n)
			_, err = grades[n]
		}
		if len(n) > subjectWidth {
			subjectWidth = len(n)
		}

		fmt.Printf("Enter your score for %q (e.g. 78): ", n)
		fmt.Scanln(&g)
		for g > 100 || g < 0 {
			fmt.Printf("You entered wrong grade (%v)\n", g)
			fmt.Printf("Enter your real score for %q (e.g. 78): ", n)
			fmt.Scanln(&g)
		}
		grades[n] = SubjectGrade{n, g}
	}

	fmt.Printf("\n\n %s's Grade Report:\n", studentName)
	fmt.Print("|")
	for i := 0; i < subjectWidth+15; i++ {
		fmt.Print("-")
	}
	fmt.Println("|")
	fmt.Printf("| %-*s| %-*s|\n", subjectWidth+5, "Subject Name", 7, "Score")
	fmt.Print("|")
	for i := 0; i < subjectWidth+6; i++ {
		fmt.Print("-")
	}
	fmt.Print("|")
	for range 8 {
		fmt.Print("-")
	}
	fmt.Println("|")
	for _, grade := range grades {
		fmt.Printf("| %-*s| %-*.2f|\n", subjectWidth+5, grade.Name, 7, grade.Grade)
	}
	fmt.Print("|")
	for i := 0; i < subjectWidth+15; i++ {
		fmt.Print("-")
	}
	fmt.Println("|")
	fmt.Printf(" Average Score : %.2f\n", calculateAvg(grades))
}
