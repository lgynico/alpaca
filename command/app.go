package command

import (
	"fmt"

	"github.com/lgynico/alpaca/consts"
	"github.com/lgynico/alpaca/writer"

	"github.com/urfave/cli/v2"
)

func Run(args []string) error {
	_app := &app{}
	_app.Init()
	return _app.Run(args)
}

type app struct {
	cli.App
	input, output    string
	templatePath     string
	dataWriter       writer.FileWriter
	clientCodeWriter writer.FileWriter
	serverCodeWriter writer.FileWriter
}

func (p *app) Init() {
	p.Name = "alpaca"
	p.Usage = "A game config parser for parse excel files to json and language codes"
	p.UsageText = "eg. alpaca -i=./config -o=./out -c=c# -s=go"
	p.Flags = p.flags()
	p.Action = p.action
}

func (p *app) action(ctx *cli.Context) error {
	if err := p.checkFlags(ctx); err != nil {
		return err
	}

	metas, err := p.parseFiles()
	if err != nil {
		return err
	}

	if err := p.dataWriter.Write(metas); err != nil {
		fmt.Printf("\033[31m[ERROR] Write data error: %v\033[0m\r\n", err)
	}

	if err := p.serverCodeWriter.Write(metas); err != nil {
		fmt.Printf("\033[31m[ERROR] Write server code error: %v\033[0m\r\n", err)
	}

	if err := p.clientCodeWriter.Write(metas); err != nil {
		fmt.Printf("\033[31m[ERROR] Write client code error: %v\033[0m\r\n", err)
	}

	return nil
}

func (p *app) checkFlags(ctx *cli.Context) error {
	p.input = ctx.String(FlagInput)
	p.output = ctx.String(FlagOutput)
	p.templatePath = ctx.String(FlagTemplate)

	c, cTag := spiltTag(ctx.String(FlagClient))
	p.clientCodeWriter = p.getWriter(codeType(c), consts.SideClient, cTag)

	s, sTag := spiltTag(ctx.String(FlagServer))
	p.serverCodeWriter = p.getWriter(codeType(s), consts.SideServer, sTag)

	p.dataWriter = writer.NewJsonWriter(p.output, len(c) > 0, len(s) > 0)

	return nil
}

func (p *app) flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     FlagInput,
			Aliases:  []string{"i"},
			Usage:    "path to excel config files",
			Required: true,
		},
		&cli.StringFlag{
			Name:        FlagOutput,
			Aliases:     []string{"o"},
			Usage:       "path to output generation files",
			DefaultText: ".",
		},
		&cli.StringFlag{
			Name:    FlagServer,
			Aliases: []string{"s"},
			Usage:   "server side language code to generate [go/c#]",
		},
		&cli.StringFlag{
			Name:    FlagClient,
			Aliases: []string{"c"},
			Usage:   "client side language code to generate [go/c#]",
		},
		&cli.StringFlag{
			Name:    FlagTemplate,
			Aliases: []string{"t"},
			Usage:   "template path to generate code, if not specified, use buildin templates",
		},
	}
}

func (p *app) getWriter(codeType CodeType, side consts.Side, tag string) writer.FileWriter {
	switch codeType {
	case CodeGolang:
		return writer.NewGoWriter(p.output, side, p.templatePath)
	case CodeCSharp:
		return writer.NewCSharpWriter(p.output, side, p.templatePath, tag)
	}

	return &writer.NoneWriter{}
}
