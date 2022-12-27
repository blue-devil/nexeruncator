package extractjs

import (
	"log"

	"github.com/blue-devil/nexeruncator/nexeruncator/pkg/utils"

	"github.com/urfave/cli/v2"
)

const (
	CmdExtract = "extract"

	flagFile = "file"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        CmdExtract,
		HelpName:    CmdExtract,
		Action:      Action,
		ArgsUsage:   ` `,
		Usage:       `extracts javascript source`,
		Description: `Extracts javascript source file from nexe-compiled binary`,
		Flags:       Flags(),
	}
}

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     flagFile,
			Usage:    "nexe-compiled binary",
			Required: true,
		},
	}
}

func Action(c *cli.Context) error {
	var nexeFile string
	var err error

	if c.IsSet(flagFile) {
		nexeFile = c.String(flagFile)
	}

	if err = utils.ExtractJS(nexeFile); err != nil {
		log.Printf("[-] Failed to extract javasript source file from %s, %v", nexeFile, err)
	}

	return nil
}
