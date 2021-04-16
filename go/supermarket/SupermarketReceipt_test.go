package supermarket

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestProducts struct {
	product Product
	price   float64
	amount  float64

	specialOfferType SpecialOfferType
	specialOfferArg  float64
	addToCart        bool
}

type ReceiptTest struct {
	name            string
	testProducts    []TestProducts
	totalPrice      float64
	discountsLen    int
	itemsLen        int
	item0Price      float64
	item0TotalPrice float64
	item0Quantity   float64
}

var receiptests = []ReceiptTest{
	{
		name: "Three for two",
		testProducts: []TestProducts{
			{
				product: Product{
					name: "hairbrush",
					unit: Each,
				},
				price:            1.50,
				specialOfferType: ThreeForTwo,
				specialOfferArg:  0.0,
				amount:           3.0,
				addToCart:        true,
			},
		},
		totalPrice:      3.00,
		itemsLen:        1,
		discountsLen:    1,
		item0Price:      1.5,
		item0TotalPrice: 4.5,
		item0Quantity:   3.0,
	},
	{
		name: "Two for amount",
		testProducts: []TestProducts{
			{
				product: Product{
					name: "hand grenade",
					unit: Each,
				},
				price:            1.75,
				specialOfferType: TwoForAmount,
				specialOfferArg:  2.0,
				amount:           3.0,
				addToCart:        true,
			},
		},
		totalPrice:      3.75,
		itemsLen:        1,
		discountsLen:    1,
		item0Price:      1.75,
		item0TotalPrice: 5.25,
		item0Quantity:   3.0,
	},
	{
		name: "Five for amount",
		testProducts: []TestProducts{
			{
				product: Product{
					name: "hens teeth",
					unit: Each,
				},
				price:            3.10,
				specialOfferType: FiveForAmount,
				specialOfferArg:  10.0,
				amount:           6.0,
				addToCart:        true,
			},
		},
		totalPrice:      13.10,
		itemsLen:        1,
		discountsLen:    1,
		item0Price:      3.10,
		item0TotalPrice: 18.60,
		item0Quantity:   6.0,
	},
	{
		name: "No discounts",
		testProducts: []TestProducts{
			{
				product: Product{
					name: "chapati flour",
					unit: Kilo,
				},
				price:            4.95,
				specialOfferType: -1,
				amount:           0.79,
				addToCart:        true,
			},
		},
		totalPrice:      3.9105000000000003,
		itemsLen:        1,
		discountsLen:    0,
		item0Price:      4.95,
		item0TotalPrice: 3.9105000000000003,
		item0Quantity:   0.79,
	},
	{
		name: "Test ten percent discount",
		testProducts: []TestProducts{
			{
				product: Product{
					name: "apples",
					unit: Kilo,
				},
				price:            1.99,
				amount:           2.5,
				specialOfferType: -1,
				addToCart:        true,
			},
			{
				product: Product{
					name: "toothbrush",
					unit: Each,
				},
				price:            1.99,
				amount:           2.5,
				specialOfferType: TenPercentDiscount,
				specialOfferArg:  10.0,
			},
		},
		totalPrice:      4.975,
		itemsLen:        1,
		discountsLen:    0,
		item0Price:      1.99,
		item0TotalPrice: 2.5 * 1.99,
		item0Quantity:   2.5,
	},
}

type FakeCatalog struct {
	_products map[string]Product
	_prices   map[string]float64
}

func (c FakeCatalog) unitPrice(product Product) float64 {
	return c._prices[product.name]
}

func (c FakeCatalog) addProduct(product Product, price float64) {
	c._products[product.name] = product
	c._prices[product.name] = price
}

func NewFakeCatalog() *FakeCatalog {
	var c FakeCatalog
	c._products = make(map[string]Product)
	c._prices = make(map[string]float64)
	return &c
}

func TestReceipts(t *testing.T) {
	for _, tt := range receiptests {
		t.Run(tt.name, func(t *testing.T) {
			var catalog = NewFakeCatalog()
			var teller = NewTeller(catalog)
			var cart = NewShoppingCart()

			for _, p := range tt.testProducts {
				catalog.addProduct(p.product, p.price)
				if p.specialOfferType != -1 {
					teller.addSpecialOffer(p.specialOfferType, p.product, p.specialOfferArg)
				}
				if p.addToCart {
					cart.addItemQuantity(p.product, p.amount)
				}
			}

			var receipt = teller.checksOutArticlesFrom(cart)
			assert.Equal(t, tt.totalPrice, receipt.totalPrice())
			assert.Equal(t, tt.discountsLen, len(receipt.discounts))
			require.Equal(t, tt.itemsLen, len(receipt.items))
			var receiptItem = receipt.items[0]
			assert.Equal(t, tt.item0Price, receiptItem.price)
			assert.Equal(t, tt.item0TotalPrice, receiptItem.totalPrice)
			assert.Equal(t, tt.item0Quantity, receiptItem.quantity)

		})
	}
}
