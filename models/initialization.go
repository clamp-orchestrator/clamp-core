package models

import "strings"

// Arguments holds a key value pair that represent an argument with its respective value
// There can be only one argument with that name. The name is case insensitive
type Arguments map[string]string

// CLIArguments is a list if arguments provided as an array of strings. It can process those arguments
// and convert into an models.Arguments type which hold them as key value pairs
type CLIArguments []string

// Parse parses all the arguments within models.CliArguments type and return them as models.Argument type
func (cliArgs CLIArguments) Parse() Arguments {
	args := Arguments{}

	for _, val := range cliArgs {
		splits := strings.Split(val, "=")
		if len(splits) > 1 {
			args[cleanUp(splits[0])] = splits[1]
		}
	}

	return args
}

// Find the argument value using name of the argument as lookupName. defaultValue is the value that is returned when the
// argument with lookupName is not found in the models.Argument type
func (args Arguments) Find(lookupName string, defaultValue string) string {
	if val, ok := args[lookupName]; ok {
		return val
	}

	return defaultValue
}

func cleanUp(value string) string {
	cleanedUpValue := strings.Replace(value, "--", "", -1)
	cleanedUpValue = strings.ToLower(cleanedUpValue)
	return cleanedUpValue
}
