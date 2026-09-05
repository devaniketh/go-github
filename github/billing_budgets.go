// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// OrganizationBudgetAlerting represents alerting settings for a GitHub organization budget.
type OrganizationBudgetAlerting struct {
	WillAlert       *bool    `json:"will_alert,omitempty"`
	AlertRecipients []string `json:"alert_recipients,omitempty"`
}

// OrganizationBudget represents a GitHub organization budget.
type OrganizationBudget struct {
	ID                  *string                     `json:"id,omitempty"`
	BudgetType          *string                     `json:"budget_type,omitempty"`
	BudgetProductSKU    *string                     `json:"budget_product_sku,omitempty"`
	BudgetProductSkus   []string                    `json:"budget_product_skus,omitempty"`
	BudgetScope         *string                     `json:"budget_scope,omitempty"`
	BudgetEntityName    *string                     `json:"budget_entity_name,omitempty"`
	BudgetAmount        *int                        `json:"budget_amount,omitempty"`
	PreventFurtherUsage *bool                       `json:"prevent_further_usage,omitempty"`
	BudgetAlerting      *OrganizationBudgetAlerting `json:"budget_alerting,omitempty"`
	User                *string                     `json:"user,omitempty"`
	ExpiresAt           *string                     `json:"expires_at,omitempty"`
}

func (b OrganizationBudget) String() string {
	return Stringify(b)
}

// OrganizationListBudgets represents a collection of GitHub organization budgets.
type OrganizationListBudgets struct {
	Budgets     []*OrganizationBudget `json:"budgets"`
	HasNextPage *bool                 `json:"has_next_page,omitempty"`
	TotalCount  *int                  `json:"total_count,omitempty"`
}

func (b OrganizationListBudgets) String() string {
	return Stringify(b)
}

// OrganizationListBudgetsOptions specifies the optional parameters to the
// BillingService.ListOrganizationBudgets method.
type OrganizationListBudgetsOptions struct {
	Scope string `url:"scope,omitempty"`
	User  string `url:"user,omitempty"`

	ListOptions
}

// OrganizationCreateBudget represents the payload to create a GitHub organization budget.
type OrganizationCreateBudget struct {
	BudgetAmount        int                         `json:"budget_amount"`
	PreventFurtherUsage bool                        `json:"prevent_further_usage"`
	BudgetAlerting      *OrganizationBudgetAlerting `json:"budget_alerting,omitempty"`
	BudgetScope         string                      `json:"budget_scope"`
	BudgetEntityName    *string                     `json:"budget_entity_name,omitempty"`
	BudgetType          string                      `json:"budget_type"`
	BudgetProductSKU    *string                     `json:"budget_product_sku,omitempty"`
	User                *string                     `json:"user,omitempty"`
	ExpiresAt           *string                     `json:"expires_at,omitempty"`
}

// OrganizationUpdateBudget represents the payload to update a GitHub organization budget.
type OrganizationUpdateBudget struct {
	BudgetAmount        *int                        `json:"budget_amount,omitempty"`
	PreventFurtherUsage *bool                       `json:"prevent_further_usage,omitempty"`
	BudgetAlerting      *OrganizationBudgetAlerting `json:"budget_alerting,omitempty"`
	BudgetScope         *string                     `json:"budget_scope,omitempty"`
	BudgetEntityName    *string                     `json:"budget_entity_name,omitempty"`
	BudgetType          *string                     `json:"budget_type,omitempty"`
	BudgetProductSKU    *string                     `json:"budget_product_sku,omitempty"`
	User                *string                     `json:"user,omitempty"`
	ExpiresAt           *string                     `json:"expires_at,omitempty"`
}

// OrganizationCreateOrUpdateBudgetResponse represents the response when creating or updating an organization budget.
type OrganizationCreateOrUpdateBudgetResponse struct {
	Message string              `json:"message"`
	Budget  *OrganizationBudget `json:"budget"`
}

// OrganizationDeleteBudgetResponse represents the response when deleting an organization budget.
type OrganizationDeleteBudgetResponse struct {
	Message  string  `json:"message"`
	BudgetID *string `json:"budget_id,omitempty"`
	ID       *string `json:"id,omitempty"`
}

// ListOrganizationBudgets gets all budgets for an organization.
//
// GitHub API docs: https://docs.github.com/rest/billing/budgets?apiVersion=2022-11-28#get-all-budgets-for-an-organization
//
//meta:operation GET /organizations/{org}/settings/billing/budgets
func (s *BillingService) ListOrganizationBudgets(ctx context.Context, org string, opts *OrganizationListBudgetsOptions) (*OrganizationListBudgets, *Response, error) {
	u := fmt.Sprintf("organizations/%v/settings/billing/budgets", org)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var budgets *OrganizationListBudgets
	resp, err := s.client.Do(req, &budgets)
	if err != nil {
		return nil, resp, err
	}

	return budgets, resp, nil
}

// GetOrganizationBudget gets a budget by ID for an organization.
//
// GitHub API docs: https://docs.github.com/rest/billing/budgets?apiVersion=2022-11-28#get-a-budget-by-id-for-an-organization
//
//meta:operation GET /organizations/{org}/settings/billing/budgets/{budget_id}
func (s *BillingService) GetOrganizationBudget(ctx context.Context, org, budgetID string) (*OrganizationBudget, *Response, error) {
	u := fmt.Sprintf("organizations/%v/settings/billing/budgets/%v", org, budgetID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var budget *OrganizationBudget
	resp, err := s.client.Do(req, &budget)
	if err != nil {
		return nil, resp, err
	}

	return budget, resp, nil
}

// CreateOrganizationBudget creates a new budget for an organization.
//
// GitHub API docs: https://docs.github.com/rest/billing/budgets?apiVersion=2022-11-28#create-a-budget-for-an-organization
//
//meta:operation POST /organizations/{org}/settings/billing/budgets
func (s *BillingService) CreateOrganizationBudget(ctx context.Context, org string, body OrganizationCreateBudget) (*OrganizationCreateOrUpdateBudgetResponse, *Response, error) {
	u := fmt.Sprintf("organizations/%v/settings/billing/budgets", org)

	req, err := s.client.NewRequest(ctx, "POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	var createBudgetResponse *OrganizationCreateOrUpdateBudgetResponse
	resp, err := s.client.Do(req, &createBudgetResponse)
	if err != nil {
		return nil, resp, err
	}

	return createBudgetResponse, resp, nil
}

// UpdateOrganizationBudget updates an existing budget for an organization.
//
// GitHub API docs: https://docs.github.com/rest/billing/budgets?apiVersion=2022-11-28#update-a-budget-for-an-organization
//
//meta:operation PATCH /organizations/{org}/settings/billing/budgets/{budget_id}
func (s *BillingService) UpdateOrganizationBudget(ctx context.Context, org, budgetID string, body OrganizationUpdateBudget) (*OrganizationCreateOrUpdateBudgetResponse, *Response, error) {
	u := fmt.Sprintf("organizations/%v/settings/billing/budgets/%v", org, budgetID)

	req, err := s.client.NewRequest(ctx, "PATCH", u, body)
	if err != nil {
		return nil, nil, err
	}

	var updateBudgetResponse *OrganizationCreateOrUpdateBudgetResponse
	resp, err := s.client.Do(req, &updateBudgetResponse)
	if err != nil {
		return nil, resp, err
	}

	return updateBudgetResponse, resp, nil
}

// DeleteOrganizationBudget deletes a budget by ID for an organization.
//
// GitHub API docs: https://docs.github.com/rest/billing/budgets?apiVersion=2022-11-28#delete-a-budget-for-an-organization
//
//meta:operation DELETE /organizations/{org}/settings/billing/budgets/{budget_id}
func (s *BillingService) DeleteOrganizationBudget(ctx context.Context, org, budgetID string) (*OrganizationDeleteBudgetResponse, *Response, error) {
	u := fmt.Sprintf("organizations/%v/settings/billing/budgets/%v", org, budgetID)

	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var deleteBudgetResponse *OrganizationDeleteBudgetResponse
	resp, err := s.client.Do(req, &deleteBudgetResponse)
	if err != nil {
		return nil, resp, err
	}

	return deleteBudgetResponse, resp, nil
}
