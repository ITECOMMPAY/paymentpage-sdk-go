package paymentpage

import (
	"time"
)

// Constants with possible payment params and types of payment
const (
	ParamProjectId            string = "project_id"
	ParamPaymentId            string = "payment_id"
	ParamPaymentAmount        string = "payment_amount"
	ParamPaymentCurrency      string = "payment_currency"
	ParamPaymentDescription   string = "payment_description"
	ParamAccountToken         string = "account_token"
	ParamCardOperationType    string = "card_operation_type"
	ParamBestBefore           string = "best_before"
	ParamCloseOnMissclick     string = "close_on_missclick"
	ParamCssModalWrap         string = "css_modal_wrap"
	ParamCustomerId           string = "customer_id"
	ParamForceAcsNewWindow    string = "force_acs_new_window"
	ParamForcePaymentMethod   string = "force_payment_method"
	ParamLanguageCode         string = "language_code"
	ParamListPaymentBlock     string = "list_payment_block"
	ParamMerchantFailUrl      string = "merchant_fail_url"
	ParamMerchantSuccessUrl   string = "merchant_success_url"
	ParamMode                 string = "mode"
	ParamRecurringRegister    string = "recurring_register"
	ParamCustomerFirstName    string = "customer_first_name"
	ParamCustomerLastName     string = "customer_last_name"
	ParamCustomerPhone        string = "customer_phone"
	ParamCustomerEmail        string = "customer_email"
	ParamCustomerCountry      string = "customer_country"
	ParamCustomerState        string = "customer_state"
	ParamCustomerCity         string = "customer_city"
	ParamCustomerDayOfBirth   string = "customer_day_of_birth"
	ParamCustomerSsn          string = "customer_ssn"
	ParamBillingPostal        string = "billing_postal"
	ParamBillingCountry       string = "billing_country"
	ParamBillingRegion        string = "billing_region"
	ParamBillingCity          string = "billing_city"
	ParamBillingAddress       string = "billing_address"
	ParamRedirect             string = "redirect"
	ParamRedirectFailMode     string = "redirect_fail_mode"
	ParamRedirectFailUrl      string = "redirect_fail_url"
	ParamRedirectOnMobile     string = "redirect_on_mobile"
	ParamRedirectSuccessMode  string = "redirect_success_mode"
	ParamRedirectSuccessUrl   string = "redirect_success_url"
	ParamRedirectTokenizeMode string = "redirect_tokenize_mode"
	ParamRedirectTokenizeUrl  string = "redirect_tokenize_url"
	ParamRegionCode           string = "region_code"
	ParamTargetElement        string = "target_element"
	ParamTerminalId           string = "terminal_id"
	ParamBaseUrl              string = "baseurl"
	ParamPaymentExtraParam    string = "payment_extra_param"

	PaymentTypePurchase  string = "purchase"
	PaymentTypePayout    string = "payout"
	PaymentTypeRecurring string = "recurring"
)

// Structure for preparing payment params
type Payment struct {
	// Map with payment params
	params map[string]interface{}
}

// Setter for payment params
func (p *Payment) SetParam(key string, value interface{}) *Payment {
	if key == ParamBestBefore {
		switch value := value.(type) {
		case time.Time:
			p.params[key] = value.Format(time.RFC3339)
		}
	} else {
		p.params[key] = value
	}

	return p
}

// Method return payment params
func (p *Payment) GetParams() map[string]interface{} {
	return p.params
}

// Constructor for Payment structure
func NewPayment(projectId int, paymentId interface{}) *Payment {
	payment := new(Payment)
	payment.params = map[string]interface{}{}
	payment.SetParam(ParamProjectId, projectId)
	payment.SetParam("interface_type", `{"id": 20}`)

	if paymentId != nil {
		payment.SetParam(ParamPaymentId, paymentId)
	}

	return payment
}
