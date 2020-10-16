package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldInitializeArgumentsWithOnlyMigrationArgument(t *testing.T) {
	expectedLen := 1
	expectedValue := "true"

	var cliArguments CLIArguments = []string{"migrate", "--migration-only=true"}

	arguments := cliArguments.Parse()
	migrateOnlyValue := arguments.Find("migration-only", "false")

	assert.Equal(t, expectedLen, len(arguments), fmt.Sprintf("Expected %d argument to be available but there were %d", expectedLen, len(arguments)))
	assert.Equal(t, expectedValue, migrateOnlyValue, fmt.Sprintf("Expected %s to be the value for migration-only argument but was %s", expectedValue, migrateOnlyValue))
}

func TestShouldEnsureThatArgumentNameIsCaseSensitiveAndValueIsCaseSensitive(t *testing.T) {
	expectedValue := "InSaNe"

	var cliArguments CLIArguments = []string{"logDebug", "--joggingLevel=InSaNe"}
	karmaArgsValue := cliArguments.Parse().Find("jogginglevel", "rookie")

	assert.Equal(t, expectedValue, karmaArgsValue, fmt.Sprintf("Expected %s to be the value for joggingLevel argument but was %s", expectedValue, karmaArgsValue))
}

func TestShouldReturnDefaultValueWhenArgumentNotPresent(t *testing.T) {
	expectedValue := "no"

	var cliArguments CLIArguments = []string{"logDebug", "--karma=3"}
	karmaArgsValue := cliArguments.Parse().Find("mana", "no")

	assert.Equal(t, expectedValue, karmaArgsValue, fmt.Sprintf("Expected %s to be the value for mana argument but was %s", expectedValue, karmaArgsValue))
}

func TestShouldEnsureThatArgumentsCanBeAnEmptyList(t *testing.T) {
	expectedLen := 0

	var cliArguments CLIArguments = []string{}
	arguments := cliArguments.Parse()

	assert.Equal(t, expectedLen, len(arguments), fmt.Sprintf("Expected %d to be the size for arguments argument but was %d", expectedLen, len(arguments)))
}
