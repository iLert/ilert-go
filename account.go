package ilert

import (
	"encoding/json"
	"fmt"
)

// Account definition https://docs.ilert.com/developer-docs/rest-api/api-reference/account
type Account struct {
	// unique identifier of the account
	ID string `json:"id,omitempty"`

	// name of the organization owning the account
	OrganizationName string `json:"organizationName,omitempty"`

	// default time zone of the account as an IANA time zone name, e.g. "Europe/Berlin"
	Timezone string `json:"timezone,omitempty"`

	// default language of the account as an ISO 639-1 code, either "en" or "de"
	Language string `json:"language,omitempty"`

	// default country of the account as an ISO 3166-1 alpha-2 code, e.g. "DE"
	Region string `json:"region,omitempty"`

	// whether users of the ilert mobile app must unlock it with an additional
	// authentication layer, such as a PIN or biometrics
	EnforceMobileProtection bool `json:"enforceMobileProtection,omitempty"`

	// whether admins may purchase additional seats when creating users, rather than
	// the account owner only. See PurchaseSeat on CreateUserInput.
	AllowAdminSeatPurchase bool `json:"allowAdminSeatPurchase,omitempty"`

	// whether AI features are enabled for the account at all. Returned by the API but
	// absent from the published OpenAPI spec (verified on 16.09.2026).
	AllowAI bool `json:"allowAI,omitempty"`

	// where AI features of the account are allowed to run
	// possible values: "DISABLED", "ALL", "EU"
	AiMode string `json:"aiMode,omitempty"`

	// features unlocked for the account by its plan and add-ons. The set of possible
	// values grows as ilert ships new features and is not part of the stable contract,
	// so it is deliberately not mirrored as a constant list here.
	ApplicationFeatures []string `json:"applicationFeatures,omitempty"`

	// the plan the account is currently subscribed to
	Subscription *AccountSubscription `json:"subscription,omitempty"`
}

// AccountSubscription defines the plan an account is subscribed to
type AccountSubscription struct {
	// display name of the plan
	Name string `json:"name,omitempty"`

	// lifecycle state of the subscription
	// possible values: "TRIAL", "FREE", "ACTIVE", "CANCELED", "UNPAID", "INCOMPLETE", "INCOMPLETE_EXPIRED", "PAST_DUE"
	Status string `json:"status,omitempty"`
}

// AiMode defines where AI features of an account are allowed to run
var AiMode = struct {
	Disabled string
	All      string
	EU       string
}{
	Disabled: "DISABLED",
	All:      "ALL",
	EU:       "EU",
}

// AiModeAll defines the ai mode list
var AiModeAll = []string{
	AiMode.Disabled,
	AiMode.All,
	AiMode.EU,
}

// SubscriptionStatus defines the lifecycle states of an account subscription
var SubscriptionStatus = struct {
	Trial             string
	Free              string
	Active            string
	Canceled          string
	Unpaid            string
	Incomplete        string
	IncompleteExpired string
	PastDue           string
}{
	Trial:             "TRIAL",
	Free:              "FREE",
	Active:            "ACTIVE",
	Canceled:          "CANCELED",
	Unpaid:            "UNPAID",
	Incomplete:        "INCOMPLETE",
	IncompleteExpired: "INCOMPLETE_EXPIRED",
	PastDue:           "PAST_DUE",
}

// SubscriptionStatusAll defines the subscription status list
var SubscriptionStatusAll = []string{
	SubscriptionStatus.Trial,
	SubscriptionStatus.Free,
	SubscriptionStatus.Active,
	SubscriptionStatus.Canceled,
	SubscriptionStatus.Unpaid,
	SubscriptionStatus.Incomplete,
	SubscriptionStatus.IncompleteExpired,
	SubscriptionStatus.PastDue,
}

// GetCurrentAccountInput represents the input of a GetCurrentAccount operation.
type GetCurrentAccountInput struct {
	_ struct{}
}

// GetCurrentAccountOutput represents the output of a GetCurrentAccount operation.
type GetCurrentAccountOutput struct {
	_       struct{}
	Account *Account
}

// GetCurrentAccount gets the account of the currently authenticated caller, including the
// plan it is subscribed to. https://docs.ilert.com/developer-docs/rest-api/api-reference/account
func (c *Client) GetCurrentAccount(input *GetCurrentAccountInput) (*GetCurrentAccountOutput, error) {
	resp, err := c.httpClient.R().Get(fmt.Sprintf("%s/current", apiRoutes.tenants))
	if err != nil {
		return nil, err
	}
	if apiErr := getGenericAPIError(resp, 200); apiErr != nil {
		return nil, apiErr
	}

	account := &Account{}
	err = json.Unmarshal(resp.Body(), account)
	if err != nil {
		return nil, err
	}

	return &GetCurrentAccountOutput{Account: account}, nil
}
