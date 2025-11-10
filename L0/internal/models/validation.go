package models

import (
	"fmt"
	"regexp"
	"time"
)

type Validator struct{}

func (v *Validator) ValidateOrder(order *Order) error {
	if order == nil {
		return fmt.Errorf("order is nil")
	}

	if err := v.validateOrderMain(order); err != nil {
		return fmt.Errorf("invalid order main: %w", err)
	}

	if err := v.validateDelivery(&order.Delivery); err != nil {
		return fmt.Errorf("invalid delivery: %w", err)
	}

	if err := v.validatePayment(&order.Payment); err != nil {
		return fmt.Errorf("invalid payment: %w", err)
	}

	if err := v.validateItems(order.Items); err != nil {
		return fmt.Errorf("invalid items: %w", err)
	}

	return nil
}

func (v *Validator) validateOrderMain(order *Order) error {
	if order.OrderUID == "" {
		return fmt.Errorf("order_uid is required")
	}

	if len(order.OrderUID) > 100 {
		return fmt.Errorf("order_uid too long")
	}

	if order.TrackNumber == "" {
		return fmt.Errorf("track_number is required")
	}

	if order.Entry == "" {
		return fmt.Errorf("entry is required")
	}

	if order.Locale == "" {
		return fmt.Errorf("locale is required")
	}

	if order.CustomerID == "" {
		return fmt.Errorf("customer_id is required")
	}

	if order.DeliveryService == "" {
		return fmt.Errorf("delivery_service is required")
	}

	if order.SmID < 0 {
		return fmt.Errorf("sm_id cannot be negative")
	}

	if order.DateCreated.After(time.Now().Add(24 * time.Hour)) {
		return fmt.Errorf("date_created cannot be in the future")
	}

	return nil
}

func (v *Validator) validateDelivery(delivery *Delivery) error {
	if delivery.Name == "" {
		return fmt.Errorf("delivery name is required")
	}

	if delivery.Phone == "" {
		return fmt.Errorf("delivery phone is required")
	}

	if !v.isValidPhone(delivery.Phone) {
		return fmt.Errorf("invalid phone format")
	}

	if delivery.Zip == "" {
		return fmt.Errorf("delivery zip is required")
	}

	if delivery.City == "" {
		return fmt.Errorf("delivery city is required")
	}

	if delivery.Address == "" {
		return fmt.Errorf("delivery address is required")
	}

	if delivery.Region == "" {
		return fmt.Errorf("delivery region is required")
	}

	if delivery.Email == "" {
		return fmt.Errorf("delivery email is required")
	}

	if !v.isValidEmail(delivery.Email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func (v *Validator) validatePayment(payment *Payment) error {
	if payment.Transaction == "" {
		return fmt.Errorf("payment transaction is required")
	}

	if payment.Currency == "" {
		return fmt.Errorf("payment currency is required")
	}

	if len(payment.Currency) != 3 {
		return fmt.Errorf("currency must be 3 characters")
	}

	if payment.Provider == "" {
		return fmt.Errorf("payment provider is required")
	}

	if payment.Amount < 0 {
		return fmt.Errorf("payment amount cannot be negative")
	}

	if payment.PaymentDt <= 0 {
		return fmt.Errorf("payment_dt is required")
	}

	if payment.Bank == "" {
		return fmt.Errorf("payment bank is required")
	}

	if payment.DeliveryCost < 0 {
		return fmt.Errorf("delivery_cost cannot be negative")
	}

	if payment.GoodsTotal < 0 {
		return fmt.Errorf("goods_total cannot be negative")
	}

	if payment.CustomFee < 0 {
		return fmt.Errorf("custom_fee cannot be negative")
	}

	return nil
}

func (v *Validator) validateItems(items []Item) error {
	if len(items) == 0 {
		return fmt.Errorf("at least one item is required")
	}

	for i, item := range items {
		if err := v.validateItem(&item, i); err != nil {
			return err
		}
	}

	return nil
}

func (v *Validator) validateItem(item *Item, index int) error {
	if item.ChrtID <= 0 {
		return fmt.Errorf("item[%d]: chrt_id must be positive", index)
	}

	if item.TrackNumber == "" {
		return fmt.Errorf("item[%d]: track_number is required", index)
	}

	if item.Price < 0 {
		return fmt.Errorf("item[%d]: price cannot be negative", index)
	}

	if item.Rid == "" {
		return fmt.Errorf("item[%d]: rid is required", index)
	}

	if item.Name == "" {
		return fmt.Errorf("item[%d]: name is required", index)
	}

	if item.Sale < 0 {
		return fmt.Errorf("item[%d]: sale must be between 0 and 100", index)
	}

	if item.TotalPrice < 0 {
		return fmt.Errorf("item[%d]: total_price cannot be negative", index)
	}

	if item.NmID <= 0 {
		return fmt.Errorf("item[%d]: nm_id must be positive", index)
	}

	if item.Brand == "" {
		return fmt.Errorf("item[%d]: brand is required", index)
	}

	if item.Status < 0 {
		return fmt.Errorf("item[%d]: status cannot be negative", index)
	}

	return nil
}

func (v *Validator) isValidPhone(phone string) bool {
	// Простая проверка телефона
	phoneRegex := `^\+?[1-9]\d{1,14}$`
	matched, _ := regexp.MatchString(phoneRegex, phone)
	return matched
}

func (v *Validator) isValidEmail(email string) bool {
	// Простая проверка email
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	return matched
}
