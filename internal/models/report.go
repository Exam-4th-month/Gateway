package models

type GetSpendingReportRequest struct {
    StartDate string `json:"start_date"`
    EndDate   string `json:"end_date"`
}

type GetIncomeReportRequest struct {
    StartDate string `json:"start_date"`
    EndDate   string `json:"end_date"`
}

type GetBudgetPerformanceReportRequest struct {
    StartDate string `json:"start_date"`
    EndDate   string `json:"end_date"`
}

type GetGoalProgressReportRequest struct{}
