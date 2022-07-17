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
	var doc *html.Node
	var err error

	if settings.Link != "" {
		doc, err = htmlquery.LoadURL(settings.Link) // Doesn't work because of bots protection
	} else {
		doc, err = htmlquery.LoadDoc(settings.Path)
	}

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

	if len(products) == 0 {
		return fmt.Errorf("there is no products")
	}

	if len(products) < p.numOfProd {
		// If number of products less than number of requested products then save to database all products
		p.numOfProd = len(products)
	}

	var productsArr []*database.Goods
	for i := 0; i < p.numOfProd; i++ {
		product, err := getProductInfo(products[i])
		if err != nil {
			return err
		}
		productsArr = append(productsArr, product)
	}

	err = p.db.AddNewProducts(productsArr)
	if err != nil {
		return err
	}

	return nil
}

func getProductInfo(node *html.Node) (*database.Goods, error) {
	urlImg, err := getUrlImg(node)
	if err != nil {
		return nil, err
	}

	name, err := getName(node)
	if err != nil {
		return nil, err
	}

	price, err := getPrice(node)
	if err != nil {
		return nil, err
	}

	urlProduct, err := getUrl(node)
	if err != nil {
		return nil, err
	}

	URL, err := url.Parse(urlProduct)
	if err != nil {
		return nil, err
	}
	urlProduct = "https://www.ozon.ru" + URL.Path

	id, err := getId(URL)
	if err != nil {
		return nil, err
	}

	return &database.Goods{
		ID:     id,
		Name:   name,
		URL:    urlProduct,
		URLImg: urlImg,
		Price:  price,
	}, nil
}

func getUrlImg(product *html.Node) (string, error) {
	img, err := htmlquery.Query(product, "//img/@src")
	if err != nil {
		return "", err
	}
	if img == nil {
		return "", fmt.Errorf("can't find img")
	}

	return img.FirstChild.Data, nil
}

func getName(product *html.Node) (string, error) {
	name, err := htmlquery.Query(product, "//*[@class=\"j4u\"]/a/span/span")
	if err != nil {
		return "", err
	}
	if name == nil {
		return "", fmt.Errorf("can't find name")
	}

	return name.FirstChild.Data, nil
}

func getPrice(product *html.Node) (string, error) {
	price, err := htmlquery.QueryAll(product, "//*[@class=\"j4u\"]/div/span")
	if err != nil {
		return "", err
	}
	if len(price) == 0 {
		return "", fmt.Errorf("can't find price")
	}

	return price[0].FirstChild.Data, nil
}

func getUrl(product *html.Node) (string, error) {
	urlNode, err := htmlquery.Query(product, "/a/@href")
	if err != nil {
		return "", err
	}
	if urlNode == nil {
		return "", fmt.Errorf("can't find url")
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
