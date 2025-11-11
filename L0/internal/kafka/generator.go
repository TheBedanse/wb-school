package kafka

import (
	"L0/internal/models"

	"github.com/brianvoe/gofakeit/v7"
)

func GenerateTestOrder() *models.Order {
	orderUID := gofakeit.UUID()
	trackNumber := generateTrackNumber()

	return &models.Order{
		OrderUID:    orderUID,
		TrackNumber: trackNumber,
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    gofakeit.Name(),
			Phone:   gofakeit.Phone(),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.Street(),
			Region:  gofakeit.State(),
			Email:   gofakeit.Email(),
		},
		Payment: models.Payment{
			Transaction:  orderUID,
			RequestID:    "",
			Currency:     gofakeit.CurrencyShort(),
			Provider:     gofakeit.Company(),
			Amount:       gofakeit.Number(1000, 10000),
			PaymentDt:    gofakeit.Int64(),
			Bank:         gofakeit.BankName(),
			DeliveryCost: gofakeit.Number(100, 1000),
			GoodsTotal:   gofakeit.Number(10, 1000),
			CustomFee:    gofakeit.Number(0, 100),
		},
		Items:             generateFakeItems(gofakeit.Number(1, 5), trackNumber),
		Locale:            gofakeit.LanguageAbbreviation(),
		InternalSignature: "",
		CustomerID:        gofakeit.UUID(),
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              gofakeit.Number(1, 100),
		DateCreated:       gofakeit.Date(),
		OofShard:          "1",
	}
}

func generateTrackNumber() string {
	return "WB" + gofakeit.DigitN(12)
}

func generateFakeItems(count int, trackNum string) []models.Item {
	items := make([]models.Item, count)

	for i := 0; i < count; i++ {
		items[i] = models.Item{
			ChrtID:      gofakeit.Int64(),
			TrackNumber: trackNum,
			Price:       gofakeit.Number(100, 5000),
			Rid:         gofakeit.UUID(),
			Name:        gofakeit.ProductName(),
			Sale:        gofakeit.Number(0, 50),
			Size:        gofakeit.RandomString([]string{"0", "S", "M", "L", "XL"}),
			TotalPrice:  gofakeit.Number(100, 5000),
			NmID:        gofakeit.Int64(),
			Brand:       gofakeit.Company(),
			Status:      gofakeit.Number(100, 400),
		}
	}

	return items
}
