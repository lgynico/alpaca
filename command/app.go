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
		fmt.Printf("Write data error: %v\n", err)
	}

	if err := p.serverCodeWriter.Write(metas); err != nil {
		fmt.Printf("Write server code error: %v\n", err)
	}

	if err := p.clientCodeWriter.Write(metas); err != nil {
		fmt.Printf("Write client code error: %v\n", err)
	}

	return nil
}

func (p *app) checkFlags(ctx *cli.Context) error {
	p.input = ctx.String(FlagInput)
	p.output = ctx.String(FlagOutput)

	p.dataWriter = writer.NewJsonWriter(p.output)

	c, cTag := spiltTag(ctx.String(FlagClient))
	p.clientCodeWriter = p.getWriter(codeType(c), consts.SideClient, cTag)

	s, sTag := spiltTag(ctx.String(FlagServer))
	p.serverCodeWriter = p.getWriter(codeType(s), consts.SideServer, sTag)

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
	}
}

func (p *app) getWriter(codeType CodeType, side consts.Side, tag ...string) writer.FileWriter {
	var _tag string
	if len(tag) > 0 {
		_tag = tag[0]
	}

	switch codeType {
	case CodeGolang:
		return writer.NewGoWriter(p.output, side)
	case CodeCSharp:
		return writer.NewCSharpWriter(p.output, side, _tag)
	}

	return &writer.NoneWriter{}
}
