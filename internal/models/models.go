package models

type CreateAccountRequest struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Balance  float32 `json:"balance"`
	Currency string  `json:"currency"`
}

type CreateBudgetRequest struct {
    CategoryID string  `json:"category_id"`
    Amount     float32 `json:"amount"`
    Period     string  `json:"period"`
    StartDate  string  `json:"start_date"`
    EndDate    string  `json:"end_date"`
}

type CreateCategoryRequest struct {
    Name string `json:"name"`
    Type string `json:"type"`
}

type CreateGoalRequest struct {
    Name          string  `json:"name"`
    TargetAmount  float32 `json:"target_amount"`
    CurrentAmount float32 `json:"current_amount"`
    Deadline      string  `json:"deadline"`
    Status        string  `json:"status"`
}

type CreateTransactionRequest struct {
    AccountID   string  `json:"account_id"`
    CategoryID  string  `json:"category_id"`
    Amount      float32 `json:"amount"`
    Type        string  `json:"type"`
    Description string  `json:"description"`
    Date        string  `json:"date"`
}
