package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"L0/internal/interfaces"
)

type OrderHandler struct {
	orderService interfaces.OrderService
	tmpl         *template.Template
}

func NewOrderHandler(orderService interfaces.OrderService) (*OrderHandler, error) {
	tmpl, err := template.ParseFiles(
		"html/index.html",
		"html/order.html",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &OrderHandler{
		orderService: orderService,
		tmpl:         tmpl,
	}, nil
}

func (h *OrderHandler) ShowHomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	orders := h.orderService.GetAllOrders()
	orderUIDs := make([]string, 0, len(orders))
	for _, order := range orders {
		orderUIDs = append(orderUIDs, order.OrderUID)
	}

	data := struct {
		OrderUIDs []string
	}{
		OrderUIDs: orderUIDs,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

func (h *OrderHandler) ShowOrder(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 || pathParts[2] == "" {
		http.Error(w, "Order UID is required", http.StatusBadRequest)
		return
	}
	orderUID := pathParts[2]

	ctx := r.Context()
	order, err := h.orderService.GetOrder(ctx, orderUID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		log.Printf("Error getting order %s: %v", orderUID, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "order.html", order); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

func (h *OrderHandler) GetOrderJSON(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 || pathParts[3] == "" {
		http.Error(w, `{"error": "Order UID is required"}`, http.StatusBadRequest)
		return
	}
	orderUID := pathParts[3]

	ctx := r.Context()
	order, err := h.orderService.GetOrder(ctx, orderUID)
	if err != nil {
		http.Error(w, `{"error": "Order not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, `{"error": "Failed to encode order"}`, http.StatusInternalServerError)
		log.Printf("Error encoding order to JSON: %v", err)
	}
}
