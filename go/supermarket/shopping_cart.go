package supermarket

import (
	"fmt"
	"math"
)

type ProductQuantity struct {
	product  Product
	quantity float64
}

type ShoppingCart struct {
	items             []ProductQuantity
	productQuantities map[Product]float64
}

func NewShoppingCart() *ShoppingCart {
	var s ShoppingCart
	s.items = []ProductQuantity{}
	s.productQuantities = make(map[Product]float64)
	return &s
}

func (c *ShoppingCart) addItemQuantity(product Product, amount float64) {
	c.items = append(c.items, ProductQuantity{product: product, quantity: amount})
	currentAmount, ok := c.productQuantities[product]
	if ok {
		c.productQuantities[product] = currentAmount + amount
	} else {
		c.productQuantities[product] = amount
	}
}

func (c *ShoppingCart) handleOffers(receipt *Receipt, offers map[Product]SpecialOffer, catalog Catalog) {
	for p, _ := range c.productQuantities {
		var quantity = c.productQuantities[p]
		if offer, ok := offers[p]; ok {
			var unitPrice = catalog.unitPrice(p)
			var quantityAsInt = int(math.Round(quantity))
			var discount *Discount = nil
			var x = 1
			if offer.offerType == ThreeForTwo {
				x = 3

			} else if offer.offerType == TwoForAmount {
				discount = getNForAmountDiscount(2, p, quantity, unitPrice, offer.offerType, offer.argument)

			}
			if offer.offerType == FiveForAmount {
				x = 5
			}
			var numberOfXs int = quantityAsInt / x
			if offer.offerType == ThreeForTwo && quantityAsInt > 2 {
				var discountAmount = quantity*unitPrice - (float64(numberOfXs*2)*unitPrice + float64(quantityAsInt%3)*unitPrice)
				discount = &Discount{product: p, description: "3 for 2", discountAmount: -discountAmount}
			}
			if offer.offerType == TenPercentDiscount {
				discount = &Discount{product: p, description: fmt.Sprintf("%.0f %% off", offer.argument), discountAmount: -quantity * unitPrice * offer.argument / 100.0}
			}
			if offer.offerType == FiveForAmount && quantityAsInt >= 5 {
				var discountTotal = unitPrice*quantity - (offer.argument*float64(numberOfXs) + float64(quantityAsInt%5)*unitPrice)
				discount = &Discount{product: p, description: fmt.Sprintf("%d for %.2f", x, offer.argument), discountAmount: -discountTotal}
			}
			if discount != nil {
				receipt.addDiscount(*discount)
			}

		}

	}
}

func getNForAmountDiscount(N int, p Product, quantity float64, unitPrice float64, offerType SpecialOfferType, offerArg float64) *Discount {
	x := 1
	quantityAsInt := int(quantity)
	if quantityAsInt < N {
		return nil
	}
	if offerType == TwoForAmount {
		x = N

		offercount := int(quantityAsInt / x)
		var total = offerArg*float64(offercount) + float64(quantityAsInt-(offercount*2))*unitPrice
		var discountN = unitPrice*quantity - total
		return &Discount{product: p, description: fmt.Sprintf("2 for %.2f", offerArg), discountAmount: -discountN}
	}
	return nil
}
