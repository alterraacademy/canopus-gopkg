package canopusgo

// CartPayload payload to create cart
type CartPayload struct {
	CartDetails struct {
		ID      string `json:"id" validate:"required"`
		Payment struct {
			Key  string `json:"key,omitempty"`
			Type string `json:"type,omitempty"`
		} `json:"payment"`
		Amount    float64 `json:"amount" validate:"required"`
		Title     string  `json:"title" validate:"required"`
		Currency  string  `json:"currency" validate:"required"`
		ExpiredAt string  `json:"expiredAt" validate:"required"`
	} `json:"cartDetails" validate:"required"`
	ItemDetails     []CartPayloadItemDetail `json:"itemDetails" validate:"required"`
	CustomerDetails struct {
		FirstName      string `json:"firstName" validate:"required"`
		LastName       string `json:"lastName,omitempty"`
		Email          string `json:"email" validate:"required"`
		Phone          string `json:"phone" validate:"required"`
		BillingAddress struct {
			FirstName  string `json:"firstName,omitempty"`
			LastName   string `json:"lastName,omitempty"`
			Phone      string `json:"phone,omitempty"`
			Address    string `json:"address,omitempty"`
			City       string `json:"city,omitempty"`
			PostalCode string `json:"postalCode,omitempty"`
		} `json:"billingAddress,omitempty"`
		ShippingAddress struct {
			FirstName  string `json:"firstName,omitempty"`
			LastName   string `json:"lastName,omitempty"`
			Phone      string `json:"phone,omitempty"`
			Address    string `json:"address,omitempty"`
			City       string `json:"city,omitempty"`
			PostalCode string `json:"postalCode,omitempty"`
		} `json:"shippingAddress,omitempty"`
	} `json:"customerDetails" validate:"required"`
	URL struct {
		ReturnURL       string `json:"returnURL" validate:"required"`
		CancelURL       string `json:"cancelURL" validate:"required"`
		NotificationURL string `json:"notificationURL" validate:"required"`
	} `json:"url" validate:"required"`
	ExtendInfo struct {
		AdditionalPrefix string `json:"additionalPrefix,omitempty"`
	} `json:"extendInfo"`
}

// CartPayloadItemDetail item cart detail
type CartPayloadItemDetail struct {
	Name           string  `json:"name" validate:"required"`
	Desc           string  `json:"desc"`
	Price          float64 `json:"price" validate:"required"`
	Quantity       int     `json:"quantity" validate:"required"`
	SKU            string  `json:"SKU" validate:"required"`
	AdditionalInfo struct {
		NoHandphone string `json:"No Handphone"`
	} `json:"additionalInfo"`
}
