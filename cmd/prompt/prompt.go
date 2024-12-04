package prompt

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"strings"
)

func Select(label string, items []string) (int, string) {
	p := promptui.Select{
		Label: label,
		Items: items,
		Size:  10,
		Searcher: func(input string, index int) bool {
			item := items[index]
			name := strings.Replace(strings.ToLower(item), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)
			return strings.Contains(name, input)
		},
	}
	i, result, err := p.Run()
	if err != nil {
		fmt.Printf("%s: %v\n", label, err)
		os.Exit(1)
	}
	return i, result
}

func Confirm(label string) {
	p := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}
	if _, err := p.Run(); err != nil {
		// no error, user just decided not to continue
		os.Exit(0)
	}
}
