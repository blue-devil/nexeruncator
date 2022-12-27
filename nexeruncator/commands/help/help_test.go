package help_test

import (
	"os"
	"testing"

	"github.com/blue-devil/nexeruncator/nexeruncator/commands/extractjs"
	"github.com/blue-devil/nexeruncator/nexeruncator/commands/help"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestHelp(t *testing.T) {
	exeName, err := os.Executable()
	require.NoError(t, err)

	app := &cli.App{
		Commands: []*cli.Command{
			help.Command(),
			extractjs.Command(),
		},
	}
	testArgs := []string{exeName, "help"}
	require.NoError(t, app.Run(testArgs))

	testArgs = []string{exeName, "help", "key"}
	require.NoError(t, app.Run(testArgs))

}
