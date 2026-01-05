package proposalmodel

// Status represents the lifecycle state of a proposal.
type Status string

const (
	// StatusPending indicates the proposal is awaiting an owner decision.
	StatusPending Status = "pending"
	// StatusAccepted indicates the owner accepted the proposal.
	StatusAccepted Status = "accepted"
	// StatusRejected indicates the owner rejected the proposal.
	StatusRejected Status = "rejected"
	// StatusCancelled indicates the realtor cancelled the proposal before an owner decision.
	StatusCancelled Status = "cancelled"
	// StatusExpired indicates the proposal expired based on its configured expiration date.
	StatusExpired Status = "expired"
)

// String returns the string representation of a Status value.
func (s Status) String() string { return string(s) }

// TransactionType identifies whether the proposal is a sale or rent operation.
type TransactionType string

const (
	// TransactionTypeSale represents a sale proposal.
	TransactionTypeSale TransactionType = "sale"
	// TransactionTypeRent represents a rent proposal.
	TransactionTypeRent TransactionType = "rent"
)

// PaymentMethod captures how the buyer intends to pay.
type PaymentMethod string

const (
	// PaymentMethodCash represents a full upfront payment.
	PaymentMethodCash PaymentMethod = "cash"
	// PaymentMethodFinancing represents bank financing as the primary payment method.
	PaymentMethodFinancing PaymentMethod = "financing"
	// PaymentMethodInstallments represents installment payments directly with the owner.
	PaymentMethodInstallments PaymentMethod = "installments"
	// PaymentMethodCashAndFinancing represents a mix of down payment plus bank financing.
	PaymentMethodCashAndFinancing PaymentMethod = "cashAndFinancing"
	// PaymentMethodExchange represents an exchange (property swap) as payment.
	PaymentMethodExchange PaymentMethod = "exchange"
	// PaymentMethodCashAndExchange represents a mix of cash plus exchange.
	PaymentMethodCashAndExchange PaymentMethod = "cashAndExchange"
)

// GuaranteeType defines guarantees for rent scenarios.
type GuaranteeType string

const (
	// GuaranteeTypeSecurityDeposit represents a security deposit (caucao).
	GuaranteeTypeSecurityDeposit GuaranteeType = "securityDeposit"
	// GuaranteeTypeSuretyBond represents a surety bond (fianca).
	GuaranteeTypeSuretyBond GuaranteeType = "suretyBond"
	// GuaranteeTypeRentalInsurance represents rental insurance.
	GuaranteeTypeRentalInsurance GuaranteeType = "rentalInsurance"
	// GuaranteeTypeGuarantor represents a human guarantor.
	GuaranteeTypeGuarantor GuaranteeType = "guarantor"
)
