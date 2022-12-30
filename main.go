// nexeruncator - exports/inserts javascript from nexe-compiled binaries
// Copyright (C) 2022  BlueDeviL // SCT <bluedevil.SCT@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"io"
	"log"
	"os"

	"github.com/blue-devil/nexeruncator/nexeruncator/commands/extractjs"
	"github.com/blue-devil/nexeruncator/nexeruncator/commands/help"
	"github.com/blue-devil/nexeruncator/nexeruncator/commands/insertjs"
	"github.com/blue-devil/nexeruncator/nexeruncator/pkg/utils"

	"github.com/urfave/cli/v2"
)

var Version = "0.3.0"
var AppName = "nexeruncator"

func main() {
	utils.PrintBanner()

	app := &cli.App{
		Name:     AppName,
		Usage:    "aberuncator of nexe-compiled binaries",
		Commands: Commands(os.Stdin),
		Version:  Version,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Commands(reader io.Reader) []*cli.Command {
	return []*cli.Command{
		help.Command(),
		extractjs.Command(),
		insertjs.Command(),
	}
}
