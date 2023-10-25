package tokens

type User struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
	PublicFlags int    `json:"public_flags"`
	MfaEnabled  bool   `json:"mfa_enabled"`
	PremiumType int    `json:"premium_type"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
}

type Billing struct {
	Type int `json:"type"`
}

type Guild struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	Owner                  bool   `json:"owner"`
	Permissions            string `json:"permissions"`
	ApproximateMemberCount int    `json:"approximate_member_count"`
}

type Friend struct {
	ID   string `json:"id"`
	User User   `json:"user,omitempty"`
}

type Invite struct {
	Code string `json:"code"`
}
