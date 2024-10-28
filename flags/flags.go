package flags

import (
	"errors"
	"strings"
)

type Options struct {
	Long      bool
	Recursive bool
	All       bool
	Reverse   bool
	SortTime  bool
	Paths     []string // Store paths separately from flags
}

func Parse(args []string) (Options, error) {
	options := Options{
		Paths: make([]string, 0),
	}

	// Handle each argument
	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Handle special case of "-" as a directory name
		if arg == "-" {
			options.Paths = append(options.Paths, arg)
			continue
		}

		// Handle flags
		if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			// Handle GNU-style flags (e.g., --help)
			if strings.HasPrefix(arg, "--") {
				return Options{}, errors.New("invalid option -- '" + arg[2:] + "'")
			}

			// Process each character in the flag string
			for _, flag := range arg[1:] {
				switch flag {
				case 'l':
					options.Long = true
				case 'R':
					options.Recursive = true
				case 'a':
					options.All = true
				case 'r':
					options.Reverse = true
				case 't':
					options.SortTime = true
				default:
					return Options{}, errors.New("invalid option -- '" + string(flag) + "'")
				}
			}
		} else {
			// If it's not a flag, it's a path
			options.Paths = append(options.Paths, arg)
		}
	}

	return options, nil
}
