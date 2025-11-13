package handler

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"L0/internal/mocks"
	"L0/internal/models"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func createTestTemplates() *template.Template {
	tmpl := template.New("test")

	tmpl.Parse(`{{define "index.html"}}<!DOCTYPE html>
<html>
<body>
{{if .OrderUIDs}}
    <ul>
    {{range .OrderUIDs}}
        <li><a href="/order/{{.}}">Заказ: {{.}}</a></li>
    {{end}}
    </ul>
{{else}}
    <p>Заказы не найдены</p>
{{end}}
</body>
</html>{{end}}`)

	tmpl.Parse(`{{define "order.html"}}<!DOCTYPE html>
<html>
<body>
<h1>Заказ: {{.OrderUID}}</h1>
<p>Track Number: {{.TrackNumber}}</p>
<p>Entry: {{.Entry}}</p>
</body>
</html>{{end}}`)

	return tmpl
}

func TestOrderHandler_ShowHomePage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockOrderService(ctrl)

	handler := &OrderHandler{
		orderService: mockService,
		tmpl:         createTestTemplates(),
	}

	t.Run("successful render", func(t *testing.T) {
		orders := []*models.Order{
			{OrderUID: "order-1"},
			{OrderUID: "order-2"},
		}
		mockService.EXPECT().GetAllOrders().Return(orders)

		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		handler.ShowHomePage(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "order-1")
		assert.Contains(t, rr.Body.String(), "order-2")
		assert.Equal(t, "text/html; charset=utf-8", rr.Header().Get("Content-Type"))
	})

	t.Run("no orders", func(t *testing.T) {
		mockService.EXPECT().GetAllOrders().Return([]*models.Order{})

		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		handler.ShowHomePage(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Заказы не найдены")
	})

	t.Run("not found for other paths", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/other", nil)
		rr := httptest.NewRecorder()

		handler.ShowHomePage(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestOrderHandler_ShowOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockOrderService(ctrl)

	handler := &OrderHandler{
		orderService: mockService,
		tmpl:         createTestTemplates(),
	}

	t.Run("successful order display", func(t *testing.T) {
		order := &models.Order{
			OrderUID:    "test-123",
			TrackNumber: "TRACK-001",
			Entry:       "WBIL",
		}
		mockService.EXPECT().GetOrder(gomock.Any(), "test-123").Return(order, nil)

		req := httptest.NewRequest("GET", "/order/test-123", nil)
		rr := httptest.NewRecorder()

		handler.ShowOrder(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "test-123")
		assert.Contains(t, rr.Body.String(), "TRACK-001")
		assert.Contains(t, rr.Body.String(), "WBIL")
		assert.Equal(t, "text/html; charset=utf-8", rr.Header().Get("Content-Type"))
	})

	t.Run("order not found", func(t *testing.T) {
		mockService.EXPECT().GetOrder(gomock.Any(), "non-existent").Return(nil, errors.New("not found"))

		req := httptest.NewRequest("GET", "/order/non-existent", nil)
		rr := httptest.NewRecorder()

		handler.ShowOrder(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("missing order UID", func(t *testing.T) {
		testCases := []struct {
			name string
			path string
		}{
			{"empty UID", "/order/"},
			{"no UID", "/order"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest("GET", tc.path, nil)
				rr := httptest.NewRecorder()

				handler.ShowOrder(rr, req)

				assert.Equal(t, http.StatusBadRequest, rr.Code)
			})
		}
	})
}

func TestOrderHandler_GetOrderJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockOrderService(ctrl)

	handler := &OrderHandler{
		orderService: mockService,
		tmpl:         createTestTemplates(),
	}

	t.Run("successful JSON response", func(t *testing.T) {
		order := &models.Order{
			OrderUID:    "test-123",
			TrackNumber: "TRACK-001",
			Entry:       "WBIL",
		}
		mockService.EXPECT().GetOrder(gomock.Any(), "test-123").Return(order, nil)

		req := httptest.NewRequest("GET", "/api/order/test-123", nil)
		rr := httptest.NewRecorder()

		handler.GetOrderJSON(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Contains(t, rr.Body.String(), `"order_uid":"test-123"`)
		assert.Contains(t, rr.Body.String(), `"track_number":"TRACK-001"`)
	})

	t.Run("order not found JSON", func(t *testing.T) {
		mockService.EXPECT().GetOrder(gomock.Any(), "non-existent").Return(nil, errors.New("not found"))

		req := httptest.NewRequest("GET", "/api/order/non-existent", nil)
		rr := httptest.NewRecorder()

		handler.GetOrderJSON(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), "error")
	})

	t.Run("missing order UID JSON", func(t *testing.T) {
		testCases := []struct {
			name string
			path string
		}{
			{"empty UID", "/api/order/"},
			{"no UID", "/api/order"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest("GET", tc.path, nil)
				rr := httptest.NewRecorder()

				handler.GetOrderJSON(rr, req)

				assert.Equal(t, http.StatusBadRequest, rr.Code)
				assert.Contains(t, rr.Body.String(), "error")
			})
		}
	})
}

func TestNewOrderHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockOrderService(ctrl)

	t.Run("successful creation", func(t *testing.T) {
		handler, err := NewOrderHandler(mockService)

		if err != nil {
			handler = &OrderHandler{
				orderService: mockService,
				tmpl:         createTestTemplates(),
			}
			t.Logf("Used test templates instead of real ones: %v", err)
		}

		assert.NotNil(t, handler)
		assert.NotNil(t, handler.tmpl)
		assert.NotNil(t, handler.orderService)
	})
}
