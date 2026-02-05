package account

type AccountDTO struct {
	ID         uint       `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Role       AccountRole `json:"role"`
	MFAEnabled bool       `json:"mfa_enabled"`
	CreatedAt  string	  `json:"created_at"`
	UpdatedAt  string	  `json:"updated_at"`
}

func ToAccountDTO(acc *Account) *AccountDTO {
	return &AccountDTO{
		ID: acc.ID,
		Name: acc.Name,
		Email: acc.Email,
		Role: acc.Role,
		MFAEnabled: acc.MFAEnabled,
		CreatedAt: acc.CreatedAt.String(),
		UpdatedAt: acc.UpdatedAt.String(),
	}
}

func ToAccountDTOs(accounts *[]Account) *[]AccountDTO {
	if accounts == nil {
		return &[]AccountDTO{}
	}

	dtos := make([]AccountDTO, len(*accounts))
	for i, acc := range *accounts {
		dtos[i] = AccountDTO{
			ID:         acc.ID,
			Name:       acc.Name,
			Email:      acc.Email,
			Role:       acc.Role,
			MFAEnabled: acc.MFAEnabled,
			CreatedAt: acc.CreatedAt.String(),
			UpdatedAt: acc.UpdatedAt.String(),
		}
	}
	return &dtos
}