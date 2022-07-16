package ozonparser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

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

func (p *Parser) Parse() error {
	products, err := htmlquery.QueryAll(p.doc, "//*[@class=\"ju3 u3j\"]")
	if err != nil {
		return err
	}

	for i := 0; i < p.numOfProd; i++ {
		err = getProductInfo(products[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func getProductInfo(node *html.Node) error {
	urlImg, err := getUrlImg(node)
	if err != nil {
		return err
	}

	name, err := getName(node)
	if err != nil {
		return nil
	}

	price, err := getPrice(node)
	if err != nil {
		return err
	}

	urlProduct, err := getUrl(node)
	if err != nil {
		return err
	}

	URL, err := url.Parse(urlProduct)
	if err != nil {
		return nil
	}
	urlProduct = "https://www.ozon.ru" + URL.Path

	id, err := getId(URL)
	if err != nil {
		return err
	}

	fmt.Printf("id: %d\nname: %s\nurl: %s\nurl_img: %s\nprice: %s\n", id, name, urlProduct, urlImg, price)

	return nil
}

func getUrlImg(product *html.Node) (string, error) {
	img, err := htmlquery.Query(product, "//img/@src")
	if err != nil {
		return "", err
	}

	return img.FirstChild.Data, nil
}

func getName(product *html.Node) (string, error) {
	name, err := htmlquery.Query(product, "//*[@class=\"j4u\"]/a/span/span")
	if err != nil {
		return "", err
	}

	return name.FirstChild.Data, nil
}

func getPrice(product *html.Node) (string, error) {
	price, err := htmlquery.QueryAll(product, "//*[@class=\"j4u\"]/div/span")
	if err != nil {
		return "", err
	}

	return price[0].FirstChild.Data, nil
}

func getUrl(product *html.Node) (string, error) {
	urlNode, err := htmlquery.Query(product, "/a/@href")
	if err != nil {
		return "", err
	}

	return urlNode.FirstChild.Data, nil
}

func getId(url *url.URL) (int, error) {
	parsed := strings.Split(url.Path, "-")
	idString := parsed[len(parsed)-1]
	id, err := strconv.Atoi(idString[:len(idString)-1])
	if err != nil {
		return 0, err
	}

	return id, nil
}
