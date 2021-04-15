package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"aaronmeza.io/myRetail/pkg/config"

	"aaronmeza.io/myRetail/pkg/product"

	"aaronmeza.io/myRetail/pkg/retailDb"

	"go.uber.org/zap"
)

//go:generate counterfeiter ./ API
type API interface {
	Products() http.Handler
}

type api struct {
	config *config.MyRetailConfig
	db     retailDb.PriceDatabase
	logger *zap.SugaredLogger
}

func NewAPI(conf *config.MyRetailConfig, db retailDb.PriceDatabase, logger *zap.SugaredLogger) API {
	return &api{
		config: conf,
		db:     db,
		logger: logger,
	}
}

func (api *api) Products() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		switch r.Method {
		case "GET":
			err = api.GetProduct(w, r)
			if err != nil {
				w.WriteHeader(404)
				return
			}
		case "PUT":
			isNew, err := api.SetPrice(r)
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte("Could not set price"))
			} else {
				if isNew {
					w.WriteHeader(201)
				} else {
					w.WriteHeader(200)
				}

			}
		default:
			w.WriteHeader(405)
			w.Header().Set("Allow", "GET, PUT")
		}
	})
}

func (api *api) GetProduct(w http.ResponseWriter, r *http.Request) error {
	product1 := product.Product{Id: parseProductId(r.URL), Name: "", CurrentPrice: product.Price{}}
	var err error
	product1.Name, err = api.getProductName(product1.Id)
	if err != nil {
		api.logger.Errorf("name not found for product " + product1.Id)
		return err
	}

	product1.CurrentPrice, err = api.db.GetPrice(product1.Id)
	if err != nil {
		api.logger.Errorf("price not found for product " + product1.Id)
		return err
	}

	w.WriteHeader(200)
	jsonData, err := json.Marshal(product1)
	if err != nil {
		return err
	}
	w.Write(jsonData)

	return nil
}

func (api *api) getProductName(productId string) (string, error) {
	var client = &http.Client{Timeout: 10 * time.Second}
	productUrl := api.config.ApiAddress + productId + "?" + api.config.ApiQuery

	api.logger.Debug("DEBUG: " + productUrl)

	resp, err := client.Get(productUrl)
	if err != nil {
		api.logger.Errorf("failed to get myProduct information for productId #{productId}, #{err}")
		return "", errors.New("could not get myProduct information")
	}
	defer resp.Body.Close()

	api.logger.Debugf("Debug: body[%+v]", resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		api.logger.Error(err.Error())
	}
	var redSkyProduct1 = product.RedSky{}
	err = json.Unmarshal(body, &redSkyProduct1)
	if err != nil {
		api.logger.Error(err.Error())
	}

	return redSkyProduct1.Product.Item.Description.Title, err
}

func (api *api) SetPrice(r *http.Request) (bool, error) {

	productId := parseProductId(r.URL)
	amount := parsePrice(r)
	price := product.Price{Amount: amount, CurrencyCode: "USD"}

	isNew, err := api.db.SetPrice(productId, price)
	if err != nil {
		api.logger.Errorf("failed to set price for productId #{productId}, #{err}")
		return false, errors.New("could not parse product Id")
	}
	return isNew, err
}

func parseProductId(url *url.URL) string {
	return path.Base(url.Path)
}
func parsePrice(r *http.Request) string {
	r.ParseForm()
	return r.Form.Get("price")
}
