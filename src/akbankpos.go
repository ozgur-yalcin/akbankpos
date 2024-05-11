package akbankpos

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var EndPoints = map[string]string{
	"TEST":   "https://apipre.akbank.com",
	"TEST3D": "https://virtualpospaymentgatewaypre.akbank.com/securepay",

	"PROD":   "https://api.akbank.com",
	"PROD3D": "https://virtualpospaymentgateway.akbank.com/securepay",
}

var CurrencyCode = map[string]int{
	"TRY": 949,
	"YTL": 949,
	"TRL": 949,
	"TL":  949,
	"USD": 840,
	"EUR": 978,
	"GBP": 826,
	"JPY": 392,
}

type API struct {
	Mode       string
	MerchantId string
	TerminalId string
	SecretKey  string
}

type Request struct {
	Version           *string            `json:"version,omitempty"`
	HashItems         *string            `json:"hashItems,omitempty"`
	Lang              *string            `json:"lang,omitempty"`
	OkUrl             *string            `json:"okUrl,omitempty"`
	FailUrl           *string            `json:"failUrl,omitempty"`
	TxnCode           *string            `json:"txnCode,omitempty"`
	PaymentModel      *string            `json:"paymentModel,omitempty"`
	RequestDateTime   *string            `json:"requestDateTime,omitempty"`
	RandomNumber      *string            `json:"randomNumber,omitempty"`
	InstitutionCode   *string            `json:"institutionCode,omitempty"`
	Terminal          *Terminal          `json:"terminal,omitempty"`
	Card              *Card              `json:"card,omitempty"`
	InsurancePan      *InsurancePan      `json:"insurancePan,omitempty"`
	Order             *Order             `json:"order,omitempty"`
	Reward            *Reward            `json:"reward,omitempty"`
	Transaction       *Transaction       `json:"transaction,omitempty"`
	Customer          *Customer          `json:"customer,omitempty"`
	Recurring         *Recurring         `json:"recurring,omitempty"`
	PlannedDate       *PlannedDate       `json:"plannedDate,omitempty"`
	PayByLink         *PayByLink         `json:"payByLink,omitempty"`
	SecureTransaction *SecureTransaction `json:"secureTransaction,omitempty"`
	SubMerchant       *SubMerchant       `json:"subMerchant,omitempty"`
	B2b               *B2b               `json:"b2b,omitempty"`
	Sgk               *Sgk               `json:"sgk,omitempty"`
}

type Response struct {
	TxnCode                  *string               `json:"txnCode,omitempty"`
	ResponseCode             *string               `json:"responseCode,omitempty"`
	Hash                     *string               `json:"hash,omitempty"`
	ResponseMessage          *string               `json:"responseMessage,omitempty"`
	HostResponseCode         *string               `json:"hostResponseCode,omitempty"`
	HostMessage              *string               `json:"hostMessage,omitempty"`
	TxnDateTime              *string               `json:"txnDateTime,omitempty"`
	Terminal                 *Terminal             `json:"terminal,omitempty"`
	Card                     *Card                 `json:"card,omitempty"`
	Order                    *Order                `json:"order,omitempty"`
	Transaction              *Transaction          `json:"transaction,omitempty"`
	Campaign                 *Campaign             `json:"campaign,omitempty"`
	Reward                   *Reward               `json:"reward,omitempty"`
	Recurring                *Recurring            `json:"recurring,omitempty"`
	PlannedDate              *PlannedDate          `json:"plannedDate,omitempty"`
	Interest                 *Interest             `json:"interest,omitempty"`
	SubMerchant              *SubMerchant          `json:"subMerchant,omitempty"`
	B2b                      *B2b                  `json:"b2b,omitempty"`
	LinkValidTerm            *float                `json:"linkValidTerm,omitempty"`
	MerchantId               *float                `json:"merchantId,omitempty"`
	LinkExpireDate           *string               `json:"linkExpireDate,omitempty"`
	MerchantOrderId          *string               `json:"merchantOrderId,omitempty"`
	ReferenceId              *string               `json:"referenceId,omitempty"`
	Token                    *string               `json:"token,omitempty"`
	Header                   *Header               `json:"header,omitempty"`
	LinkDetail               *LinkDetail           `json:"linkDetail,omitempty"`
	InstallmentConditionList []*InstallmentCond    `json:"installmentConditionList,omitempty"`
	TxnDetailList            []*TxnDetailListInner `json:"txnDetailList,omitempty"`
	Error                    *Error                `json:"error,omitempty"`
}

type B2b struct {
	IdentityNumber *string `json:"identityNumber,omitempty"`
}

type Card struct {
	CardHolderName *string `json:"cardHolderName,omitempty"`
	CardNumber     *string `json:"cardNumber,omitempty"`
	CardCode       *string `json:"cvv2,omitempty"`
	CardExpiry     *string `json:"expireDate,omitempty"`
}

type Customer struct {
	EmailAddress *string `json:"emailAddress,omitempty"`
	IpAddress    *string `json:"ipAddress,omitempty"`
}

type InsurancePan struct {
	BinNumber         *string `json:"binNumber,omitempty"`
	CardLastFourParam *string `json:"cardLastFourParam,omitempty"`
	IdentityNumber    *string `json:"identityNumber,omitempty"`
}

type Order struct {
	OrderId      *string `json:"orderId,omitempty"`
	OrderTrackId *string `json:"orderTrackId,omitempty"`
}

type PayByLink struct {
	LinkTxnCode       *string `json:"linkTxnCode,omitempty"`
	LinkTransferType  *string `json:"linkTransferType,omitempty"`
	MobilePhoneNumber *string `json:"mobilePhoneNumber,omitempty"`
	Email             *string `json:"email,omitempty"`
}

type PlannedDate struct {
	FirstPlannedDate *string `json:"firstPlannedDate,omitempty"`
}

type Recurring struct {
	NumberOfPayments  *int    `json:"numberOfPayments,omitempty"`
	FrequencyInterval *int    `json:"frequencyInterval,omitempty"`
	FrequencyCycle    *string `json:"frequencyCycle,omitempty"`
	RecurringOrder    *int    `json:"recurringOrder,omitempty"`
}

type SecureTransaction struct {
	SecureId      *string `json:"secureId,omitempty"`
	SecureEcomInd *string `json:"secureEcomInd,omitempty"`
	SecureData    *string `json:"secureData,omitempty"`
	SecureMd      *string `json:"secureMd,omitempty"`
}

type Sgk struct {
	SurchargeAmount *float `json:"surchargeAmount,omitempty"`
}

type SubMerchant struct {
	SubMerchantId *string `json:"subMerchantId,omitempty"`
}

type Terminal struct {
	MerchantSafeId *string `json:"merchantSafeId,omitempty"`
	TerminalSafeId *string `json:"terminalSafeId,omitempty"`
}

type Campaign struct {
	AdditionalInstallment *int    `json:"additionalInstallCount,omitempty"`
	DeferingDate          *string `json:"deferingDate,omitempty"`
	DeferingMonth         *int    `json:"deferingMonth,omitempty"`
}

type Header struct {
	ReturnCode    *string `json:"returnCode,omitempty"`
	ReturnMessage *string `json:"returnMessage,omitempty"`
}

type InstallmentCond struct {
	InstallmentCount *float  `json:"installmentCount,omitempty"`
	InstallmentType  *string `json:"installmentType,omitempty"`
	CardType         *string `json:"cardType,omitempty"`
}

type Interest struct {
	InterestRate   *float `json:"interestRate,omitempty"`
	InterestAmount *float `json:"interestAmount,omitempty"`
}

type LinkDetail struct {
	LinkTransferType  *string `json:"linkTransferType,omitempty"`
	MobilePhoneNumber *string `json:"mobilePhoneNumber,omitempty"`
	Email             *string `json:"email,omitempty"`
	LinkValidTerm     *float  `json:"linkValidTerm,omitempty"`
	Amount            *float  `json:"amount,omitempty"`
	Currency          *int    `json:"currencyCode,omitempty"`
	InstallmentCount  *float  `json:"installmentCount,omitempty"`
	ReferenceId       *string `json:"referenceId,omitempty"`
	ErrorCode         *string `json:"errorCode,omitempty"`
	ErrorMessage      *string `json:"errorMessage,omitempty"`
	LinkExpireDate    *string `json:"linkExpireDate,omitempty"`
	LinkStatus        *string `json:"linkStatus,omitempty"`
	InstallmentType   *float  `json:"installmentType,omitempty"`
}

type Reward struct {
	CcbRewardAmount        *float  `json:"ccbRewardAmount,omitempty"`
	PcbRewardAmount        *float  `json:"pcbRewardAmount,omitempty"`
	XcbRewardAmount        *float  `json:"xcbRewardAmount,omitempty"`
	CcbEarnedRewardAmount  *float  `json:"ccbEarnedRewardAmount,omitempty"`
	CcbBalanceRewardAmount *float  `json:"ccbBalanceRewardAmount,omitempty"`
	CcbRewardDesc          *string `json:"ccbRewardDesc,omitempty"`
	PcbEarnedRewardAmount  *float  `json:"pcbEarnedRewardAmount,omitempty"`
	PcbBalanceRewardAmount *float  `json:"pcbBalanceRewardAmount,omitempty"`
	PcbRewardDesc          *string `json:"pcbRewardDesc,omitempty"`
	XcbEarnedRewardAmount  *float  `json:"xcbEarnedRewardAmount,omitempty"`
	XcbBalanceRewardAmount *float  `json:"xcbBalanceRewardAmount,omitempty"`
	XcbRewardDesc          *string `json:"xcbRewardDesc,omitempty"`
}

type Transaction struct {
	Amount      *float  `json:"amount,omitempty"`
	Currency    *int    `json:"currencyCode,omitempty"`
	MotoInd     *int    `json:"motoInd,omitempty"`
	Installment *int    `json:"installCount,omitempty"`
	AuthCode    *string `json:"authCode,omitempty"`
	Rrn         *string `json:"rrn,omitempty"`
	BatchNumber *int    `json:"batchNumber,omitempty"`
	Stan        *int    `json:"stan,omitempty"`
}

type TxnDetailListInner struct {
	TxnCode                    *string `json:"txnCode,omitempty"`
	ResponseCode               *string `json:"responseCode,omitempty"`
	ResponseMessage            *string `json:"responseMessage,omitempty"`
	HostResponseCode           *string `json:"hostResponseCode,omitempty"`
	HostMessage                *string `json:"hostMessage,omitempty"`
	TxnDateTime                *string `json:"txnDateTime,omitempty"`
	PlannedDateTime            *string `json:"plannedDateTime,omitempty"`
	TerminalSafeId             *string `json:"terminalSafeId,omitempty"`
	MerchantSafeId             *string `json:"merchantSafeId,omitempty"`
	OrderId                    *string `json:"orderId,omitempty"`
	OrderTrackId               *string `json:"orderTrackId,omitempty"`
	AuthCode                   *string `json:"authCode,omitempty"`
	Rrn                        *string `json:"rrn,omitempty"`
	BatchNumber                *int    `json:"batchNumber,omitempty"`
	Stan                       *int    `json:"stan,omitempty"`
	SettlementId               *string `json:"settlementId,omitempty"`
	TxnStatus                  *string `json:"txnStatus,omitempty"`
	Amount                     *float  `json:"amount,omitempty"`
	Currency                   *int    `json:"currencyCode,omitempty"`
	MotoInd                    *int    `json:"motoInd,omitempty"`
	Installment                *int    `json:"installCount,omitempty"`
	CcbRewardAmount            *float  `json:"ccbRewardAmount,omitempty"`
	PcbRewardAmount            *float  `json:"pcbRewardAmount,omitempty"`
	XcbRewardAmount            *float  `json:"xcbRewardAmount,omitempty"`
	PreAuthStatus              *string `json:"preAuthStatus,omitempty"`
	PreAuthCloseAmount         *float  `json:"preAuthCloseAmount,omitempty"`
	PreAuthPartialCancelAmount *float  `json:"preAuthPartialCancelAmount,omitempty"`
	PreAuthCloseDate           *string `json:"preAuthCloseDate,omitempty"`
	MaskedCardNumber           *string `json:"maskedCardNumber,omitempty"`
	RecurringOrder             *int    `json:"recurringOrder,omitempty"`
	RequestType                *string `json:"requestType,omitempty"`
	RequestStatus              *string `json:"requestStatus,omitempty"`
	CancelDate                 *string `json:"cancelDate,omitempty"`
	TryCount                   *int    `json:"tryCount,omitempty"`
	Xid                        *string `json:"xid,omitempty"`
	PaymentModel               *string `json:"paymentModel,omitempty"`
	Eci                        *string `json:"eci,omitempty"`
	SecureData                 *string `json:"secureData,omitempty"`
	OrgOrderId                 *string `json:"orgOrderId,omitempty"`
}

type Error struct {
	Code    string   `json:"code,omitempty"`
	Message string   `json:"message,omitempty"`
	Errors  []Errors `json:"errors,omitempty"`
}

type Errors struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type float float32

func (f float) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.2f", float32(f))), nil
}

func Random(n int) string {
	const alphanum = "0123456789ABCDEF"
	var bytes = make([]byte, n)
	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func Hash(payload []byte, secretkey string) string {
	hmac := hmac.New(sha512.New, []byte(secretkey))
	hmac.Write(payload)
	return base64.StdEncoding.EncodeToString(hmac.Sum(nil))
}

func Api(merchantid, terminalid, secretkey string) (*API, *Request) {
	api := new(API)
	api.MerchantId = merchantid
	api.TerminalId = terminalid
	api.SecretKey = secretkey
	req := new(Request)
	req.Terminal = new(Terminal)
	req.Card = new(Card)
	req.Transaction = new(Transaction)
	req.Customer = new(Customer)
	req.Reward = new(Reward)
	version := "1.00"
	req.Version = &version
	req.Terminal.MerchantSafeId = &merchantid
	req.Terminal.TerminalSafeId = &terminalid
	return api, req
}

func (api *API) SetMode(mode string) {
	api.Mode = mode
}

func (req *Request) SetCardNumber(cardnumber string) {
	req.Card.CardNumber = &cardnumber
}

func (req *Request) SetCardExpiry(cardmonth, cardyear string) {
	cardexpiry := cardmonth + cardyear
	req.Card.CardExpiry = &cardexpiry
}

func (req *Request) SetCardCode(cardcode string) {
	req.Card.CardCode = &cardcode
}

func (req *Request) SetAmount(amount float, currency string) {
	motoInd := 0
	code := CurrencyCode[currency]
	req.Transaction.MotoInd = &motoInd
	req.Transaction.Amount = &amount
	req.Transaction.Currency = &code
}

func (req *Request) SetInstallment(installment int) {
	req.Transaction.Installment = &installment
}

func (req *Request) SetCustomerIPv4(ipaddress string) {
	req.Customer.IpAddress = &ipaddress
}

func (req *Request) SetCustomerEmail(email string) {
	req.Customer.EmailAddress = &email
}

func (req *Request) SetOrderId(orderid string) {
	req.Order.OrderId = &orderid
}

func (api *API) PreAuth(ctx context.Context, req *Request) (Response, error) {
	date := time.Now().Format("2006-01-02T15:04:05.000")
	rnd := Random(128)
	txnCode := "1004"
	req.RequestDateTime = &date
	req.RandomNumber = &rnd
	req.TxnCode = &txnCode
	return api.Transaction(ctx, req)
}

func (api *API) Auth(ctx context.Context, req *Request) (Response, error) {
	date := time.Now().Format("2006-01-02T15:04:05.000")
	rnd := Random(128)
	txnCode := "1000"
	req.RequestDateTime = &date
	req.RandomNumber = &rnd
	req.TxnCode = &txnCode
	req.Reward.CcbRewardAmount = new(float)
	req.Reward.PcbRewardAmount = new(float)
	req.Reward.XcbRewardAmount = new(float)
	return api.Transaction(ctx, req)
}

func (api *API) PostAuth(ctx context.Context, req *Request) (Response, error) {
	date := time.Now().Format("2006-01-02T15:04:05.000")
	rnd := Random(128)
	txnCode := "1005"
	req.RequestDateTime = &date
	req.RandomNumber = &rnd
	req.TxnCode = &txnCode
	return api.Transaction(ctx, req)
}

func (api *API) Refund(ctx context.Context, req *Request) (Response, error) {
	date := time.Now().Format("2006-01-02T15:04:05.000")
	rnd := Random(128)
	txnCode := "1002"
	req.RequestDateTime = &date
	req.RandomNumber = &rnd
	req.TxnCode = &txnCode
	return api.Transaction(ctx, req)
}

func (api *API) Cancel(ctx context.Context, req *Request) (Response, error) {
	date := time.Now().Format("2006-01-02T15:04:05.000")
	rnd := Random(128)
	txnCode := "1003"
	req.RequestDateTime = &date
	req.RandomNumber = &rnd
	req.TxnCode = &txnCode
	return api.Transaction(ctx, req)
}

func (api *API) Transaction(ctx context.Context, req *Request) (res Response, err error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return res, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", EndPoints[api.Mode]+"/api/v1/payment/virtualpos/transaction/process", bytes.NewReader(payload))
	if err != nil {
		return res, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("auth-hash", Hash(payload, api.SecretKey))
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	if response.StatusCode == http.StatusOK {
		if err := decoder.Decode(&res); err == nil {
			return res, nil
		}
	} else {
		if err := decoder.Decode(&res.Error); err == nil {
			return res, errors.New(res.Error.Message)
		}
	}
	return res, errors.New("unknown error")
}
