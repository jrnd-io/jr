package types

type User struct {
	Guid      string  `json:"guid"`
	IsActive  bool    `json:"isActive"`
	Balance   float64 `json:"balance"`
	Age       int     `json:"age"`
	EyeColor  string  `json:"eyeColor"`
	Name      string  `json:"name"`
	Company   string  `json:"company"`
	Email     string  `json:"email"`
	Alt_email string  `json:"alt_email"`
	About     string  `json:"about"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
