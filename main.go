package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func printUsage() {
	fmt.Printf("Markov Chain text generator.\n\n")
	fmt.Printf("Usage:\n")
	fmt.Printf("  markovchain [-w <N>] [-p <S>] [-l <N>]\n")
	fmt.Printf("  markovchain --help\n\n")
	fmt.Printf("Options:\n")
	fmt.Printf("  --help  Show this screen.\n")
	fmt.Printf("  -w N    Number of maximum words (default=100) (max=10000)\n")
	fmt.Printf("  -p S    Starting prefix (default=start of the text)\n")
	fmt.Printf("  -l N    Prefix length (default=2) (max=5)\n")
}

func printWarning(msg string) {
	fmt.Printf("\x1b[35mWarning\u001b[0m: " + msg + "\n")
}

func printErrorAndExit(msg string) {
	fmt.Printf("\u001b[31mError\u001b[0m: " + msg + "\n")
	fmt.Printf("Use 'markov-chain --help' for more information.\n")
	os.Exit(1)
}

func isNum(s string) bool {
	nonzero := false
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
		if c > '0' {
			nonzero = true
		}
	}
	return nonzero
}

func toNum(s string) int {
	res := 0
	for _, c := range s {
		res = res*10 + int(c-'0')
	}
	return res
}

func contains(words []string, prefix []string) bool {
	if len(prefix) == 0 {
		return true
	}
	for i := 0; i+len(prefix) <= len(words); i++ {
		if words[i] == prefix[0] {
			match := true
			for j := range prefix {
				if prefix[j] != words[i+j] {
					match = false
					break
				}
			}
			if match {
				return true
			}
		}
	}
	return false
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--help" {
		if len(os.Args) > 2 {
			printWarning("args after --help are skipped.\n============================")
		}
		printUsage()
		return
	}

	var incorrectArgs []string
	os.Args = os.Args[1:]
	argCnt := len(os.Args)
	wordNum, prefixStart, prefixLength := -1, make([]string, 0), -1
	ind := 0
	for ind < len(os.Args) {
		arg := os.Args[ind]
		if arg == "-w" {
			if wordNum != -1 {
				printErrorAndExit("flag -w is set more than once.")
			} else if ind == argCnt-1 {
				printErrorAndExit("flag -w is specified but not set.")
			} else if !isNum(os.Args[ind+1]) {
				printErrorAndExit("invalid number is provided for flag -w. The number should be integer, positive and not larger than 10000.")
			} else {
				wordNum = toNum(os.Args[ind+1])
				ind += 2
			}
		} else if arg == "-p" {
			if len(prefixStart) > 0 {
				printErrorAndExit("flag -p is set more than once.")
			} else if ind == argCnt-1 {
				printErrorAndExit("flag -p is specified but not set.")
			} else if len(os.Args[ind+1]) == 0 {
				printErrorAndExit("empty string is set for -p.")
			} else {
				prefixStart = strings.Fields(os.Args[ind+1])
				ind += 2
			}
		} else if arg == "-l" {
			if prefixLength != -1 {
				printErrorAndExit("flag -l is set more than once.")
			} else if ind == argCnt-1 {
				printErrorAndExit("flag -l is specified but not set.")
			} else if !isNum(os.Args[ind+1]) {
				printErrorAndExit("invalid number is provided for flag -l. The number should be integer and positive.")
			} else {
				prefixLength = toNum(os.Args[ind+1])
				ind += 2
			}
		} else {
			incorrectArgs = append(incorrectArgs, arg)
			ind++
		}
	}

	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) != 0 {
		printErrorAndExit("no input text.")
	}

	inputbytes, _ := io.ReadAll(os.Stdin)
	text := strings.Fields(string(inputbytes))
	if len(text) == 0 {
		printErrorAndExit("no words in input.")
	}

	if !contains(text, prefixStart) {
		printErrorAndExit("the text does not contain the specified prefix.")
	}

	if wordNum == -1 {
		wordNum = 100
	} else if wordNum > 10000 {
		printWarning("maximum word number is too large; it is now set to 10000.")
		wordNum = 10000
	}

	if prefixLength == -1 {
		prefixLength = 2
	} else if prefixLength > len(text) {
		printErrorAndExit("prefix length exceeds the number of words in text.")
	} else if prefixLength > 5 {
		printWarning("prefix length is too large; it is now set to 5.")
		prefixLength = 5
	}

	if len(prefixStart) == 0 {
		prefixStart = make([]string, prefixLength)
		copy(prefixStart, text[:prefixLength])
	}
}
