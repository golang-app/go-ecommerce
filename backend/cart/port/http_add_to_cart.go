package port

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/bkielbasa/go-ecommerce/backend/cart/domain"
	"github.com/bkielbasa/go-ecommerce/backend/internal/https"
)

type AddToCartRequest struct {
	ProductID string `json:"product_id"`
	Qty       int    `json:"quantity"`
}

func (h HTTP) AddToCart(w http.ResponseWriter, r *http.Request) {
	cartID := cartIDFromCookies(w, r)

	req := AddToCartRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		https.BadRequest(w, "serialization-error", err.Error())
		return
	}

	err := h.cart.AddToCart(r.Context(), cartID, req.ProductID, req.Qty)

	if errors.Is(err, domain.ErrProductNotFound) {
		https.NotFound(w, "cart-not-found", err.Error())
		return
	}

	if err != nil {
		https.InternalError(w, "internal-error", err.Error())
		log.Print(err)
		return
	}

	https.NoContent(w)
}
