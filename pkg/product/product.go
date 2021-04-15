package product

type Product struct {
	Id           string `bson:"_id" json:"id"`
	Name         string `bson:"name" json:"name"`
	CurrentPrice Price  `bson:"current_price" json:"current_price"`
}

type Price struct {
	Amount       string `bson:"value" json:"value"`
	CurrencyCode string `bson:"currency_code" json:"currency_code"`
}

type RedSky struct {
	Product RedSkyProduct `json:"product"`
}
type RedSkyProduct struct {
	Item RedSkyItem `json:"item"`
}

type RedSkyItem struct {
	Id          string             `json:"tcin"`
	Description ProductDescription `json:"product_description"`
}

type ProductDescription struct {
	Title string `json:"title"`
}
