// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBillingService_ListOrganizationBudgets(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/organizations/o/settings/billing/budgets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page":     "2",
			"per_page": "10",
			"scope":    "organization",
			"user":     "octocat",
		})
		fmt.Fprint(w, `{
			"budgets": [
				{
					"id": "2066deda-923f-43f9-88d2-62395a28c0cdd",
					"budget_type": "ProductPricing",
					"budget_product_sku": "actions",
					"budget_product_skus": ["actions"],
					"budget_scope": "organization",
					"budget_entity_name": "example-org",
					"budget_amount": 1000,
					"prevent_further_usage": true,
					"budget_alerting": {
						"will_alert": true,
						"alert_recipients": ["org-admin"]
					},
					"user": "octocat",
					"expires_at": "2026-12-31"
				}
			],
			"has_next_page": true,
			"total_count": 1
		}`)
	})

	opts := &OrganizationListBudgetsOptions{
		Scope:       "organization",
		User:        "octocat",
		ListOptions: ListOptions{Page: 2, PerPage: 10},
	}
	ctx := t.Context()
	budgets, _, err := client.Billing.ListOrganizationBudgets(ctx, "o", opts)
	if err != nil {
		t.Errorf("Billing.ListOrganizationBudgets returned error: %v", err)
	}

	want := &OrganizationListBudgets{
		Budgets: []*OrganizationBudget{
			{
				ID:                  new("2066deda-923f-43f9-88d2-62395a28c0cdd"),
				BudgetType:          new(BudgetTypeProductPricing),
				BudgetProductSKU:    new("actions"),
				BudgetProductSkus:   []string{"actions"},
				BudgetScope:         new(BudgetScopeOrganization),
				BudgetEntityName:    new("example-org"),
				BudgetAmount:        new(1000),
				PreventFurtherUsage: new(true),
				BudgetAlerting: &OrganizationBudgetAlerting{
					WillAlert:       new(true),
					AlertRecipients: []string{"org-admin"},
				},
				User:      new("octocat"),
				ExpiresAt: new("2026-12-31"),
			},
		},
		HasNextPage: new(true),
		TotalCount:  new(1),
	}
	if !cmp.Equal(budgets, want) {
		t.Errorf("Billing.ListOrganizationBudgets returned %+v, want %+v", budgets, want)
	}

	const methodName = "ListOrganizationBudgets"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.ListOrganizationBudgets(ctx, "o", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.ListOrganizationBudgets(ctx, "\n", opts)
		return err
	})
}

func TestBillingService_ListOrganizationBudgets_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Billing.ListOrganizationBudgets(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestBillingService_GetOrganizationBudget(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/organizations/o/settings/billing/budgets/b-123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": "b-123",
			"budget_type": "ProductPricing",
			"budget_product_sku": "actions_linux",
			"budget_scope": "repository",
			"budget_entity_name": "example-repo",
			"budget_amount": 500,
			"prevent_further_usage": true,
			"budget_alerting": {
				"will_alert": true,
				"alert_recipients": ["mona", "lisa"]
			}
		}`)
	})

	ctx := t.Context()
	budget, _, err := client.Billing.GetOrganizationBudget(ctx, "o", "b-123")
	if err != nil {
		t.Errorf("Billing.GetOrganizationBudget returned error: %v", err)
	}

	want := &OrganizationBudget{
		ID:                  new("b-123"),
		BudgetType:          new(BudgetTypeProductPricing),
		BudgetProductSKU:    new("actions_linux"),
		BudgetScope:         new(BudgetScopeRepository),
		BudgetEntityName:    new("example-repo"),
		BudgetAmount:        new(500),
		PreventFurtherUsage: new(true),
		BudgetAlerting: &OrganizationBudgetAlerting{
			WillAlert:       new(true),
			AlertRecipients: []string{"mona", "lisa"},
		},
	}
	if !cmp.Equal(budget, want) {
		t.Errorf("Billing.GetOrganizationBudget returned %+v, want %+v", budget, want)
	}

	const methodName = "GetOrganizationBudget"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.GetOrganizationBudget(ctx, "o", "b-123")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.GetOrganizationBudget(ctx, "\n", "\n")
		return err
	})
}

func TestBillingService_GetOrganizationBudget_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Billing.GetOrganizationBudget(ctx, "%", "b-123")
	testURLParseError(t, err)
}

func TestBillingService_CreateOrganizationBudget(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	req := OrganizationCreateBudget{
		BudgetAmount:        500,
		PreventFurtherUsage: true,
		BudgetAlerting: &OrganizationBudgetAlerting{
			WillAlert:       new(true),
			AlertRecipients: []string{"org-admin"},
		},
		BudgetScope:      BudgetScopeOrganization,
		BudgetEntityName: new(""),
		BudgetType:       BudgetTypeProductPricing,
		BudgetProductSKU: new("actions"),
	}

	mux.HandleFunc("/organizations/o/settings/billing/budgets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testJSONBody(t, r, req)
		fmt.Fprint(w, `{
			"message": "Budget successfully created.",
			"budget": {
				"id": "b-123",
				"budget_type": "ProductPricing",
				"budget_product_sku": "actions",
				"budget_scope": "organization",
				"budget_entity_name": "example-organization",
				"budget_amount": 500,
				"prevent_further_usage": true,
				"budget_alerting": {
					"will_alert": true,
					"alert_recipients": ["org-admin"]
				}
			}
		}`)
	})

	ctx := t.Context()
	resp, _, err := client.Billing.CreateOrganizationBudget(ctx, "o", req)
	if err != nil {
		t.Errorf("Billing.CreateOrganizationBudget returned error: %v", err)
	}

	want := &OrganizationCreateOrUpdateBudgetResponse{
		Message: "Budget successfully created.",
		Budget: &OrganizationBudget{
			ID:                  new("b-123"),
			BudgetType:          new(BudgetTypeProductPricing),
			BudgetProductSKU:    new("actions"),
			BudgetScope:         new(BudgetScopeOrganization),
			BudgetEntityName:    new("example-organization"),
			BudgetAmount:        new(500),
			PreventFurtherUsage: new(true),
			BudgetAlerting: &OrganizationBudgetAlerting{
				WillAlert:       new(true),
				AlertRecipients: []string{"org-admin"},
			},
		},
	}
	if !cmp.Equal(resp, want) {
		t.Errorf("Billing.CreateOrganizationBudget returned %+v, want %+v", resp, want)
	}

	const methodName = "CreateOrganizationBudget"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.CreateOrganizationBudget(ctx, "o", req)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.CreateOrganizationBudget(ctx, "\n", req)
		return err
	})
}

func TestBillingService_CreateOrganizationBudget_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Billing.CreateOrganizationBudget(ctx, "%", OrganizationCreateBudget{})
	testURLParseError(t, err)
}

func TestBillingService_UpdateOrganizationBudget(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	req := OrganizationUpdateBudget{
		BudgetAmount:        new(10),
		PreventFurtherUsage: new(false),
		BudgetAlerting: &OrganizationBudgetAlerting{
			WillAlert:       new(true),
			AlertRecipients: []string{"org-admin"},
		},
	}

	mux.HandleFunc("/organizations/o/settings/billing/budgets/b-123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testJSONBody(t, r, req)
		fmt.Fprint(w, `{
			"message": "Budget successfully updated.",
			"budget": {
				"id": "b-123",
				"budget_type": "ProductPricing",
				"budget_product_sku": "actions_linux",
				"budget_scope": "repository",
				"budget_entity_name": "example-repo",
				"budget_amount": 10,
				"prevent_further_usage": false,
				"budget_alerting": {
					"will_alert": true,
					"alert_recipients": ["org-admin"]
				}
			}
		}`)
	})

	ctx := t.Context()
	resp, _, err := client.Billing.UpdateOrganizationBudget(ctx, "o", "b-123", req)
	if err != nil {
		t.Errorf("Billing.UpdateOrganizationBudget returned error: %v", err)
	}

	want := &OrganizationCreateOrUpdateBudgetResponse{
		Message: "Budget successfully updated.",
		Budget: &OrganizationBudget{
			ID:                  new("b-123"),
			BudgetType:          new(BudgetTypeProductPricing),
			BudgetProductSKU:    new("actions_linux"),
			BudgetScope:         new(BudgetScopeRepository),
			BudgetEntityName:    new("example-repo"),
			BudgetAmount:        new(10),
			PreventFurtherUsage: new(false),
			BudgetAlerting: &OrganizationBudgetAlerting{
				WillAlert:       new(true),
				AlertRecipients: []string{"org-admin"},
			},
		},
	}
	if !cmp.Equal(resp, want) {
		t.Errorf("Billing.UpdateOrganizationBudget returned %+v, want %+v", resp, want)
	}

	const methodName = "UpdateOrganizationBudget"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.UpdateOrganizationBudget(ctx, "o", "b-123", req)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.UpdateOrganizationBudget(ctx, "\n", "\n", req)
		return err
	})
}

func TestBillingService_UpdateOrganizationBudget_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Billing.UpdateOrganizationBudget(ctx, "%", "b-123", OrganizationUpdateBudget{})
	testURLParseError(t, err)
}

func TestBillingService_DeleteOrganizationBudget(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/organizations/o/settings/billing/budgets/b-123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{
			"message": "Budget successfully deleted.",
			"budget_id": "b-123"
		}`)
	})

	ctx := t.Context()
	resp, _, err := client.Billing.DeleteOrganizationBudget(ctx, "o", "b-123")
	if err != nil {
		t.Errorf("Billing.DeleteOrganizationBudget returned error: %v", err)
	}

	want := &OrganizationDeleteBudgetResponse{
		Message:  "Budget successfully deleted.",
		BudgetID: new("b-123"),
	}
	if !cmp.Equal(resp, want) {
		t.Errorf("Billing.DeleteOrganizationBudget returned %+v, want %+v", resp, want)
	}

	const methodName = "DeleteOrganizationBudget"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Billing.DeleteOrganizationBudget(ctx, "o", "b-123")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Billing.DeleteOrganizationBudget(ctx, "\n", "\n")
		return err
	})
}

func TestBillingService_DeleteOrganizationBudget_invalidOrg(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()
	_, _, err := client.Billing.DeleteOrganizationBudget(ctx, "%", "b-123")
	testURLParseError(t, err)
}
