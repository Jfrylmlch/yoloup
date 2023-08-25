package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	flag.Bool("h", false, "Print help information")
	printUpdatedYoloAnnotationFile := flag.Bool("n", false, "Print updated yolo annotation file to stdout")
	requestConfirmation := flag.Bool("i", false, "Request confirmation before attempting to update each file")
	flag.Parse()
	args := flag.Args()
	if len(args) < 3 {
		Usage()
		os.Exit(1)
	}
	var originalClasses map[uint]string = ClassMap{args[0]}.IndexToName()
	var targetClasses map[string]uint = ClassMap{args[1]}.NameToIndex()
	if len(targetClasses) < len(originalClasses) {
		fmt.Printf("Warning: Number of classes in %s is fewer than that in %s. Some classes might not be remapped sucessfully.\n", args[1], args[0])
		forceUpdate := CheckYesNo(Input("Continue [N/y]? "))
		if !forceUpdate {
			fmt.Println("Exit")
			os.Exit(1)
		}
	}
	for _, yoloAnnoationFile := range args[2:] {
		originalYoloAnnotation := ReadFile(yoloAnnoationFile)
		updatedYoloAnnotation, e := UpdateYoloAnnotationFile(originalClasses, targetClasses, originalYoloAnnotation)
		if e != nil {
			log.Fatalf("Error on %s: %v", yoloAnnoationFile, e)
		}
		if *printUpdatedYoloAnnotationFile {
			fmt.Printf("Updated %s: \n", yoloAnnoationFile)
			fmt.Printf("%s\n", updatedYoloAnnotation)
			continue
		}
		if *requestConfirmation {
			fmt.Printf("Original %s: \n", yoloAnnoationFile)
			fmt.Printf("%s", originalYoloAnnotation)
			fmt.Printf("Updated %s: \n", yoloAnnoationFile)
			fmt.Printf("%s\n", updatedYoloAnnotation)
			confirmUpdate := CheckYesNo(Input(fmt.Sprintf("Update '%s' [N/y]? ", yoloAnnoationFile)))
			fmt.Printf("\n")
			if !confirmUpdate {
				continue
			}
		}
		WriteFile(yoloAnnoationFile, updatedYoloAnnotation)
	}

}

func CheckYesNo(input string) bool {
	answer := strings.TrimSpace(input)
	if answer == "Y" || answer == "y" {
		return true
	}
	return false
}

func Input(question string) string {
	fmt.Printf(question)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}
	return input
}

func ReadFile(path string) string {
	f, scanner := FileReader(path)
	builder := strings.Builder{}
	defer f.Close()
	for scanner.Scan() {
		builder.WriteString(fmt.Sprintf("%s\n", scanner.Text()))
	}
	return builder.String()
}

func UpdateYoloAnnotationFile(originalClasses map[uint]string, targetClasses map[string]uint, yoloAnnoation string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(yoloAnnoation))
	lines := make([]string, 0)
	for scanner.Scan() {
		columns := strings.Fields(strings.TrimSpace(scanner.Text()))
		classIdx, e := strconv.Atoi(columns[0])
		if e != nil {
			log.Fatal(e)
		}
		updatedClassIdx, found := targetClasses[originalClasses[uint(classIdx)]]
		if !found {
			return "", fmt.Errorf("Failed to remap class index %d\n", classIdx)
		}

		columns[0] = strconv.FormatUint(uint64(updatedClassIdx), 10)
		lines = append(lines, strings.Join(columns, " "))
	}
	return strings.Join(lines, "\n"), nil
}

func Usage() {
	fmt.Printf("Usage: %s [OPTIONS] [orginal_class_file] [updated_class_file] [yolo_annotation_file ...]\n", os.Args[0])
	flag.PrintDefaults()
}

func WriteFile(path string, content string) {
	w, e := os.Create(path)
	if e != nil {
		log.Fatalln(e)
	}
	defer w.Close()
	_, err := w.WriteString(content)

	if err != nil {
		log.Fatalf("Fail to write to file\n%v", err)
	}
}

type ClassMap struct {
	path string
}

func (x ClassMap) NameToIndex() map[string]uint {
	f, scanner := FileReader(x.path)
	classes := make(map[string]uint, 0)
	var lineNumber uint = 0
	for scanner.Scan() {
		className := strings.TrimSpace(scanner.Text())
		if className != "" {
			classes[className] = lineNumber
			lineNumber++
		}
	}
	f.Close()
	return classes
}

func (x ClassMap) IndexToName() map[uint]string {
	classes := make(map[uint]string, 0)
	f, scanner := FileReader(x.path)
	var lineNumber uint = 0
	for scanner.Scan() {
		className := strings.TrimSpace(scanner.Text())
		if className != "" {
			classes[lineNumber] = className
			lineNumber++
		}
	}
	f.Close()
	return classes
}

func FileReader(path string) (*os.File, *bufio.Scanner) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return f, scanner
}
