package retailDb

import (
	"github.com/pkg/errors"

	"aaronmeza.io/myRetail/pkg/product"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//go:generate counterfeiter ./ PriceDatabase
type PriceDatabase interface {
	GetPrice(productId string) (product.Price, error)
	SetPrice(productId string, price product.Price) (bool, error)
	Close()
}
type priceDatabase struct {
	session *mgo.Session
	db      *mgo.Collection
}

func NewPriceDatabase() (PriceDatabase, error) {
	dbUrl := "127.0.0.1"

	session, err := mgo.Dial(dbUrl)
	if err != nil {
		panic(errors.Wrapf(err, "server url %s", dbUrl))
	}
	db := session.DB("demo").C("products")

	return &priceDatabase{
		session: session,
		db:      db,
	}, err
}

func (m *priceDatabase) GetPrice(productId string) (product.Price, error) {
	result := product.Product{}
	err := m.db.Find(bson.M{"_id": productId}).One(&result)
	if err != nil {
		return product.Price{}, err
	}
	return result.CurrentPrice, err
}

func (m *priceDatabase) SetPrice(productId string, price product.Price) (bool, error) {
	_, err := m.GetPrice(productId)
	if err != nil && err.Error() == "not found" {
		return true, m.db.Insert(&product.Product{Id: productId, Name: "", CurrentPrice: price})
	} else {
		return false, m.db.Update(bson.M{"_id": productId}, bson.M{"$set": bson.M{"current_price": price}})
	}
}

func (m *priceDatabase) Close() {
	m.session.Close()
}
