package insertjs

import (
	"log"

	"github.com/blue-devil/nexeruncator/nexeruncator/pkg/utils"

	"github.com/urfave/cli/v2"
)

const (
	CmdInsert = "insert"
	flagDest  = "dest"
	flagSrc   = "src"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        CmdInsert,
		HelpName:    CmdInsert,
		Action:      Action,
		ArgsUsage:   ` `,
		Usage:       `inserts javascript source`,
		Description: `Inserts javascript source file into nexe-compiled binary`,
		Flags:       Flags(),
	}
}

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     flagDest,
			Usage:    "nexe-compiled binary",
			Required: true,
		},
		&cli.StringFlag{
			Name:     flagSrc,
			Usage:    "javascript source",
			Required: true,
		},
	}
}

func Action(c *cli.Context) error {
	var nexeFile string
	var jsFile string
	var err error

	if c.IsSet(flagDest) && c.IsSet(flagSrc) {
		nexeFile = c.String(flagDest)
		jsFile = c.String(flagSrc)
	}
	if err = utils.InsertJS(nexeFile, jsFile); err != nil {
		log.Printf("[-] Failed to insert javascript source file %s into %s, error: %v", jsFile, nexeFile, err)
	}
	return nil
}
