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
	flag.Bool("h", false, "Print help information.")
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
		if *requestConfirmation {
			confirmUpdate := CheckYesNo(Input(fmt.Sprintf("Update '%s' [N/y]? ", yoloAnnoationFile)))
			if !confirmUpdate {
				continue
			}
		}
		UpdateYoloAnnotationFile(originalClasses, targetClasses, yoloAnnoationFile)
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

func UpdateYoloAnnotationFile(originalClasses map[uint]string, targetClasses map[string]uint, yoloAnnoationFile string) {
	f, scanner := FileReader(yoloAnnoationFile)
	defer f.Close()

	lines := make([]string, 0)
	for scanner.Scan() {
		columns := strings.Fields(strings.TrimSpace(scanner.Text()))
		classIdx, e := strconv.Atoi(columns[0])
		if e != nil {
			log.Fatal(e)
		}
		updatedClassIdx := targetClasses[originalClasses[uint(classIdx)]]
		columns[0] = strconv.FormatUint(uint64(updatedClassIdx), 10)
		lines = append(lines, strings.Join(columns, " "))
	}

	WriteFile(yoloAnnoationFile, strings.Join(lines, "\n"))

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
