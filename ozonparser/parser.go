package ozonparser

import (
	"Parser/config"
	"Parser/database"
	"net/url"
	"strconv"
	"strings"

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
	//doc, err := htmlquery.LoadDoc("C:\\Users\\iljad\\Desktop\\test.html")
	//req, _ := http.NewRequest("GET", settings.Link, nil)
	//client := &http.Client{Timeout: time.Second * 5}
	//resp, err := client.Do(req)
	//if err != nil {
	//	return err
	//}
	//defer resp.Body.Close()
	//
	//var r []byte
	//_, err = resp.Body.Read(r)
	//if err != nil {
	//	return err
	//}
	//
	//hmtlDoc := string(r)
	//splited := strings.Split(hmtlDoc, "<body>")
	//splited = strings.Split(splited[1], "</body>")
	//body := splited[0]
	//doc, err := htmlquery.Parse(strings.NewReader(body))
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
		product, err := getProductInfo(products[i])
		if err != nil {
			return err
		}
		err = p.db.AddNewProduct(product)
		if err != nil {
			return err
		}
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
