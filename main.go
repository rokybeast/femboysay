package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	defaultArtFile = "/usr/share/femboysay/femboy.txt"
	defaultWidth = 40
)

type Message struct {
	lines    []string
	maxWidth int
}

// new message; empty
func newMessage() *Message {
	return &Message{
		lines:    make([]string, 0),
		maxWidth: 0,
	}
}

// adds a line to our message and keeps track of the widest line
func (m *Message) addLine(line string) {
	m.lines = append(m.lines, line)
	if len(line) > m.maxWidth {
		m.maxWidth = len(line)
	}
}

// wrapper
func (m *Message) wrapText(text string, width int) {
	paragraphs := strings.Split(text, "\n")

	for _, para := range paragraphs {
	// skip empty paras
		if len(strings.TrimSpace(para)) == 0 {
			continue
		}

		words := strings.Fields(para)
		if len(words) == 0 {
			continue
		}

		currentLine := ""

		for _, word := range words {
			testLine := currentLine
			if len(currentLine) > 0 {
				testLine += " "
			}
			testLine += word

			if len(testLine) <= width {
				currentLine = testLine
			} else {
				if len(currentLine) > 0 {
					m.addLine(currentLine)
				}
				currentLine = word
			}
		}

		if len(currentLine) > 0 {
			m.addLine(currentLine)
		}
	}
}

func printBorder(width int, left, middle, right rune) {
	fmt.Printf("%c", left)
	for i := 0; i < width+2; i++ {
		fmt.Printf("%c", middle)
	}
	fmt.Printf("%c\n", right)
}

func (m *Message) printBubble() {
	width := m.maxWidth

	// top border
	printBorder(width, ' ', '_', ' ')

	for i, line := range m.lines {
		var leftChar, rightChar rune

		// single line gets angle brackets
		if len(m.lines) == 1 {
			leftChar = '<'
			rightChar = '>'
		} else if i == 0 {
			// first line of multiple gets forward slashes
			leftChar = '/'
			rightChar = '\\'
		} else if i == len(m.lines)-1 {
			// last line gets backslashes
			leftChar = '\\'
			rightChar = '/'
		} else {
			// middle lines get pipes
			leftChar = '|'
			rightChar = '|'
		}

		padding := width - len(line)
		fmt.Printf("%c %s%s %c\n", leftChar, line, strings.Repeat(" ", padding), rightChar)
	}

	// bottom border
	printBorder(width, ' ', '-', ' ')
}

func printArtFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("couldn't open art file '%s': %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading art file: %v", err)
	}

	return nil
}

func readStdin() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}

func main() {
	width := flag.Int("w", defaultWidth, "maximum width of the speech bubble")
	artFile := flag.String("f", defaultArtFile, "path to ASCII art file")
	showHelp := flag.Bool("h", false, "show this help message")

	flag.Parse()

	if *showHelp {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [message]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nIf no message is provided, reads from stdin.\n")
		os.Exit(0)
	}

	if *width < 10 || *width > 200 {
		fmt.Fprintf(os.Stderr, "Error: width must be between 10 and 200\n")
		os.Exit(1)
	}

	var messageText string

	if flag.NArg() > 0 {
		messageText = strings.Join(flag.Args(), " ")
	} else {
		var err error
		messageText, err = readStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
	}

	messageText = strings.TrimSpace(messageText)
	if len(messageText) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no message provided\n")
		fmt.Fprintf(os.Stderr, "Try '%s -h' for help\n", os.Args[0])
		os.Exit(1)
	}

	msg := newMessage()
	msg.wrapText(messageText, *width)

	// print everything
	msg.printBubble()

	if err := printArtFile(*artFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
