package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "orderTracking/internal/handler"
    "orderTracking/internal/models"
    "orderTracking/internal/service"
)

type fakeAuth struct{}
func (fakeAuth) CreateUser(user models.User) (int, error) { return 1, nil }
func (fakeAuth) GenerateToken(username, password string) (string, error) { return "tok", nil }
func (fakeAuth) ParseToken(token string) (int, error) { return 1, nil }

type fakeOrder struct{}
func (fakeOrder) CreateOrderItem(orderItem models.OrderItem) (string, error) { return "track-123", nil }
func (fakeOrder) GetOrderDetailsByTrackID(track_id string) (models.OrderItem, error) {
    return models.OrderItem{Track_id: track_id, Status: "created"}, nil
}

func TestCreateOrderAndGetStatus(t *testing.T) {
    svc := &service.Service{Authorization: fakeAuth{}, OrderItem: fakeOrder{}}
    h := handler.NewHandler(svc, nil)
    r := h.InitRoutes()

    body := map[string]any{"title":"t","description":"d","price":10,"from_location":"A","to_location":"B"}
    b, _ := json.Marshal(body)
    req := httptest.NewRequest(http.MethodPost, "/api/lists/add_order", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer test")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK { t.Fatalf("create order status=%d body=%s", w.Code, w.Body.String()) }

    req2 := httptest.NewRequest(http.MethodGet, "/api/lists/status?track_id=track-123", nil)
    req2.Header.Set("Authorization", "Bearer test")
    w2 := httptest.NewRecorder()
    r.ServeHTTP(w2, req2)
    if w2.Code != http.StatusOK { t.Fatalf("get status=%d body=%s", w2.Code, w2.Body.String()) }
}

