package main

import (
	"strings"
	"testing"
	"time"
)

func TestNewContract(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		customerName  string
		startDate     time.Time
		endDate       time.Time
		wantErr       bool
		errorContains string
	}{
		{
			name:         "valid contract",
			id:           "C001",
			customerName: "John Doe",
			startDate:    time.Now(),
			endDate:      time.Now().AddDate(1, 0, 0),
			wantErr:      false,
		},
		{
			name:          "empty ID",
			customerName:  "John Doe",
			startDate:     time.Now(),
			endDate:       time.Now().AddDate(1, 0, 0),
			wantErr:       true,
			errorContains: "contract ID cannot be empty",
		},
		{
			name:          "empty customer name",
			id:            "C001",
			startDate:     time.Now(),
			endDate:       time.Now().AddDate(1, 0, 0),
			wantErr:       true,
			errorContains: "customer name cannot be empty",
		},
		{
			name:          "end date before start date",
			id:            "C001",
			customerName:  "John Doe",
			startDate:     time.Now(),
			endDate:       time.Now().AddDate(-1, 0, 0),
			wantErr:       true,
			errorContains: "end date must be after start date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contract, err := NewContract(tt.id, tt.customerName, tt.startDate, tt.endDate)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewContract() error = nil, want error containing %q", tt.errorContains)
					return
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("NewContract() error = %v, want error containing %q", err, tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("NewContract() unexpected error = %v", err)
				return
			}

			if contract.ID != tt.id {
				t.Errorf("contract.ID = %v, want %v", contract.ID, tt.id)
			}
			if contract.CustomerName != tt.customerName {
				t.Errorf("contract.CustomerName = %v, want %v", contract.CustomerName, tt.customerName)
			}
			if !contract.StartDate.Equal(tt.startDate) {
				t.Errorf("contract.StartDate = %v, want %v", contract.StartDate, tt.startDate)
			}
			if !contract.EndDate.Equal(tt.endDate) {
				t.Errorf("contract.EndDate = %v, want %v", contract.EndDate, tt.endDate)
			}
		})
	}
}

func TestContract_IsActive(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		contract Contract
		want     bool
	}{
		{
			name: "active contract",
			contract: Contract{
				StartDate: now.AddDate(0, -1, 0), // 1 month ago
				EndDate:   now.AddDate(0, 1, 0),  // 1 month from now
			},
			want: true,
		},
		{
			name: "future contract",
			contract: Contract{
				StartDate: now.AddDate(0, 1, 0), // 1 month from now
				EndDate:   now.AddDate(0, 2, 0), // 2 months from now
			},
			want: false,
		},
		{
			name: "expired contract",
			contract: Contract{
				StartDate: now.AddDate(0, -2, 0), // 2 months ago
				EndDate:   now.AddDate(0, -1, 0), // 1 month ago
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.contract.IsActive(); got != tt.want {
				t.Errorf("Contract.IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContract_RemainingDays(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		contract Contract
		want     int
	}{
		{
			name: "active contract",
			contract: Contract{
				StartDate: now.AddDate(0, -1, 0), // 1 month ago
				EndDate:   now.AddDate(0, 1, 0),  // 1 month from now (approximately 30 days)
			},
			want: 30,
		},
		{
			name: "expired contract",
			contract: Contract{
				StartDate: now.AddDate(0, -2, 0), // 2 months ago
				EndDate:   now.AddDate(0, -1, 0), // 1 month ago
			},
			want: 0,
		},
		{
			name: "future contract",
			contract: Contract{
				StartDate: now.AddDate(0, 1, 0), // 1 month from now
				EndDate:   now.AddDate(0, 2, 0), // 2 months from now
			},
			want: 60, // approximately
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.contract.RemainingDays()
			// Allow for some flexibility in the day calculation due to months having different lengths
			if got < tt.want-2 || got > tt.want+2 {
				t.Errorf("Contract.RemainingDays() = %v, want approximately %v", got, tt.want)
			}
		})
	}
}
