package ozonparser

import (
	"Parser/config"
	"Parser/database"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type Parser struct {
	doc       *html.Node
	numOfProd int
	db        database.DBIFace
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Init(settings *config.ParserSettings, db database.DBIFace) error {
	doc, err := htmlquery.LoadURL(settings.Link)
	if err != nil {
		return err
	}
	p.doc = doc
	p.numOfProd = settings.Num
	p.db = db

	return nil
}


