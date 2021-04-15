package api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	myRetailConfig "aaronmeza.io/myRetail/pkg/config"

	"aaronmeza.io/myRetail/pkg/product"
	"aaronmeza.io/myRetail/pkg/retailDb/retailDbfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	. "aaronmeza.io/myRetail/pkg/api"
)

var _ = Describe("Api", func() {
	var fakeDB *retailDbfakes.FakePriceDatabase
	var logger *zap.SugaredLogger
	var testAPIServer *httptest.Server
	var server API
	var myProduct product.Product
	var c *myRetailConfig.MyRetailConfig

	BeforeEach(func() {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		testAPIServer = httptest.NewServer(handler)
		fakeDB = &retailDbfakes.FakePriceDatabase{}
		zapLog, _ := zap.NewProduction()
		defer zapLog.Sync()
		logger = zapLog.Sugar()

		myProduct = product.Product{
			Id:   "13860428",
			Name: "The Big Lebowski (Blu-ray)",
			CurrentPrice: product.Price{
				Amount:       "13.49",
				CurrencyCode: "USD",
			},
		}

	})

	AfterEach(func() {
		fakeDB.Close()
	})

	Context("GetProduct", func() {
		It("product Id invalid", func() {
			testAPIServer.Close()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(404)
			})
			testAPIServer = httptest.NewServer(handler)
			c = &myRetailConfig.MyRetailConfig{
				ApiAddress: testAPIServer.URL,
				ApiQuery:   "",
			}
			server = NewAPI(c, fakeDB, logger)

			req, err := http.NewRequest("GET", "/products/"+myProduct.Id, nil)
			Expect(err).To(BeNil())

			recorder := httptest.NewRecorder()
			apiHandler := server.Products()
			apiHandler.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(404))

		})

		It("Name matches", func() {
			req, err := http.NewRequest("GET", "/products/"+myProduct.Id, nil)
			Expect(err).To(BeNil())
			c := &myRetailConfig.MyRetailConfig{
				ApiAddress: "https://redsky.target.com/v3/pdp/tcin/",
				ApiQuery:   "excludes=taxonomy,price,promotion,bulk_ship,rating_and_review_reviews,rating_and_review_statistics,question_answer_statistics&key=candidate",
			}

			fakeDB.GetPriceReturnsOnCall(0,
				product.Price{
					Amount:       "13.49",
					CurrencyCode: "USD",
				}, nil)
			server = NewAPI(c, fakeDB, logger)
			recorder := httptest.NewRecorder()
			apiHandler := server.Products()
			apiHandler.ServeHTTP(recorder, req)

			body, err := ioutil.ReadAll(recorder.Result().Body)
			Expect(err).To(BeNil())

			var product1 product.Product
			err = json.Unmarshal(body, &product1)
			Expect(err).To(BeNil())

			Expect(recorder.Code).To(Equal(200))
			Expect(myProduct.Name).To(Equal(product1.Name))
			Expect(myProduct.CurrentPrice).To(Equal(product1.CurrentPrice))
			getPrice := fakeDB.GetPriceArgsForCall(0)
			Expect(getPrice).To(Equal("13860428"))
			Expect(fakeDB.GetPriceCallCount()).To(Equal(1))

		})

	})
})
