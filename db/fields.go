package db

// get in dollars and cents
func (release *Release) GetPriceUSD() float32 {
	return float32(release.Price) / 100
}

// get in dollars and cents
func (purchase *Purchase) GetAmountPaidUSD() float32 {
	return float32(purchase.AmountPaid) / 100
}

// get in dollars and cents
func (purchase *Purchase) GetAuthorsCutUSD() float32 {
	return float32(purchase.AuthorsCut) / 100
}
