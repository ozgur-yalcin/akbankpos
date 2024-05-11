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
	"net/url"
	"reflect"
	"strconv"
	"strings"
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

var CurrencyISO = map[string]string{
	"949": "TRY",
	"840": "USD",
	"978": "EUR",
	"826": "GBP",
	"392": "JPY",
}

type API struct {
	Mode      string
	SecretKey string
}

type Request struct {
	Version           *string            `json:"version,omitempty"`
	HashItems         *string            `json:"hashItems,omitempty"`
	Lang              *string            `json:"lang,omitempty" form:"lang,omitempty"`
	OkUrl             *string            `json:"okUrl,omitempty" form:"okUrl,omitempty"`
	FailUrl           *string            `json:"failUrl,omitempty" form:"failUrl,omitempty"`
	TxnCode           *string            `json:"txnCode,omitempty" form:"txnCode,omitempty"`
	PaymentModel      *string            `json:"paymentModel,omitempty" form:"paymentModel,omitempty"`
	RequestDateTime   *string            `json:"requestDateTime,omitempty" form:"requestDateTime,omitempty"`
	RandomNumber      *string            `json:"randomNumber,omitempty" form:"randomNumber,omitempty"`
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
	B2B               *B2B               `json:"b2b,omitempty"`
	SGK               *SGK               `json:"sgk,omitempty"`
	Hash              *string            `json:",omitempty" form:"hash,omitempty"`
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
	B2B                      *B2B                  `json:"b2b,omitempty"`
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

type B2B struct {
	IdentityNumber *string `json:"identityNumber,omitempty" form:"b2bIdentityNumber,omitempty"`
}

type Card struct {
	CardHolderName *string `json:"cardHolderName,omitempty"`
	CardNumber     *string `json:"cardNumber,omitempty" form:"creditCard,omitempty"`
	CardCode       *string `json:"cvv2,omitempty" form:"cvv,omitempty"`
	CardExpiry     *string `json:"expireDate,omitempty" form:"expiredDate,omitempty"`
}

type Customer struct {
	EmailAddress *string `json:"emailAddress,omitempty" form:"emailAddress,omitempty"`
	IpAddress    *string `json:"ipAddress,omitempty"`
}

type InsurancePan struct {
	BinNumber         *string `json:"binNumber,omitempty"`
	CardLastFourParam *string `json:"cardLastFourParam,omitempty"`
	IdentityNumber    *string `json:"identityNumber,omitempty"`
}

type Order struct {
	OrderId      *string `json:"orderId,omitempty" form:"orderId,omitempty"`
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

type SGK struct {
	SurchargeAmount *float `json:"surchargeAmount,omitempty"`
}

type SubMerchant struct {
	SubMerchantId *string `json:"subMerchantId,omitempty" form:"subMerchantId,omitempty"`
}

type Terminal struct {
	MerchantSafeId *string `json:"merchantSafeId,omitempty" form:"merchantSafeId,omitempty"`
	TerminalSafeId *string `json:"terminalSafeId,omitempty" form:"terminalSafeId,omitempty"`
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
	CcbRewardAmount        *float  `json:"ccbRewardAmount,omitempty" form:"ccbRewardAmount,omitempty"`
	PcbRewardAmount        *float  `json:"pcbRewardAmount,omitempty" form:"pcbRewardAmount,omitempty"`
	XcbRewardAmount        *float  `json:"xcbRewardAmount,omitempty" form:"xcbRewardAmount,omitempty"`
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
	Amount      *float  `json:"amount,omitempty" form:"amount,omitempty"`
	Currency    *int    `json:"currencyCode,omitempty" form:"currencyCode,omitempty"`
	MotoInd     *int    `json:"motoInd,omitempty"`
	Installment *int    `json:"installCount,omitempty" form:"installCount,omitempty"`
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

func String(v reflect.Value) string {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}
	return fmt.Sprint(v.Interface())
}

func QueryString(v interface{}) (url.Values, error) {
	values := make(url.Values)
	value := reflect.ValueOf(v)
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return values, nil
		}
		value = value.Elem()
	}
	if v == nil {
		return values, nil
	}
	err := reflector(values, value)
	return values, err
}

func reflector(values url.Values, val reflect.Value) error {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		sv := val.Field(i)
		for sv.Kind() == reflect.Ptr {
			if sv.IsNil() {
				break
			}
			sv = sv.Elem()
		}
		if sv.Kind() == reflect.Struct {
			if err := reflector(values, sv); err != nil {
				return err
			}
			continue
		}
		if n, ok := sf.Tag.Lookup("form"); ok {
			ts := strings.Split(n, ",")
			name := ts[0]
			if sv.Kind() == reflect.Float32 {
				value := fmt.Sprintf("%.2f", sv.Interface())
				if len(ts) > 1 && ts[1] == "omitempty" && value != "" {
					values.Add(name, value)
				} else if len(ts) > 1 && ts[1] != "omitempty" {
					values.Add(name, value)
				} else if len(ts) == 1 {
					values.Add(name, value)
				}
			} else {
				value := String(sv)
				if len(ts) > 1 && ts[1] == "omitempty" && value != "" {
					values.Add(name, value)
				} else if len(ts) > 1 && ts[1] != "omitempty" {
					values.Add(name, value)
				} else if len(ts) == 1 {
					values.Add(name, value)
				}
			}
		}
	}
	return nil
}

func Api(merchantid, terminalid, secretkey string) (*API, *Request) {
	api := new(API)
	api.SecretKey = secretkey
	version := "1.00"
	req := new(Request)
	req.Version = &version
	req.Terminal = new(Terminal)
	req.Terminal.MerchantSafeId = &merchantid
	req.Terminal.TerminalSafeId = &terminalid
	return api, req
}

func B64(data string) (hash string) {
	hash = base64.StdEncoding.EncodeToString([]byte(data))
	return hash
}

func D64(data string) []byte {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil
	}
	return b
}

func (api *API) Hash(payload []byte) string {
	hmac := hmac.New(sha512.New, []byte(api.SecretKey))
	hmac.Write(payload)
	return base64.StdEncoding.EncodeToString(hmac.Sum(nil))
}

func (api *API) Hash3D(req url.Values, params []string) string {
	items := []string{}
	for _, param := range params {
		items = append(items, req.Get(param))
	}
	plain := strings.Join(items, "")
	return api.Hash([]byte(plain))
}

func (api *API) Random(n int) string {
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

func (api *API) SetMode(mode string) {
	api.Mode = mode
}

func (req *Request) SetLang(lang string) {
	req.Lang = &lang
}

func (req *Request) SetCardNumber(cardnumber string) {
	if req.Card == nil {
		req.Card = new(Card)
	}
	req.Card.CardNumber = &cardnumber
}

func (req *Request) SetCardExpiry(cardmonth, cardyear string) {
	if req.Card == nil {
		req.Card = new(Card)
	}
	cardexpiry := cardmonth + cardyear
	req.Card.CardExpiry = &cardexpiry
}

func (req *Request) SetCardCode(cardcode string) {
	if req.Card == nil {
		req.Card = new(Card)
	}
	req.Card.CardCode = &cardcode
}

func (req *Request) SetAmount(price, currency string) {
	if req.Transaction == nil {
		req.Transaction = new(Transaction)
	}
	if parse, err := strconv.ParseFloat(price, 32); err == nil {
		amount := float(parse)
		code := CurrencyCode[currency]
		req.Transaction.Amount = &amount
		req.Transaction.Currency = &code
	}
}

func (req *Request) SetInstallment(installment string) {
	if req.Transaction == nil {
		req.Transaction = new(Transaction)
	}
	if parse, err := strconv.Atoi(installment); err == nil {
		req.Transaction.Installment = &parse
	}
}

func (req *Request) SetCustomerIPv4(ipaddress string) {
	if req.Customer == nil {
		req.Customer = new(Customer)
	}
	req.Customer.IpAddress = &ipaddress
}

func (req *Request) SetCustomerEmail(email string) {
	if req.Customer == nil {
		req.Customer = new(Customer)
	}
	req.Customer.EmailAddress = &email
}

func (req *Request) SetOrderId(orderid string) {
	if req.Order == nil {
		req.Order = new(Order)
	}
	req.Order.OrderId = &orderid
}

func (api *API) PreAuth(ctx context.Context, req *Request) (Response, error) {
	date := time.Now().Format("2006-01-02T15:04:05.000")
	rnd := api.Random(128)
	code := "1004"
	motoInd := 0
	req.RequestDateTime = &date
	req.RandomNumber = &rnd
	req.TxnCode = &code
	req.Transaction.MotoInd = &motoInd
	return api.Transaction(ctx, req)
}

func (api *API) Auth(ctx context.Context, req *Request) (Response, error) {
	date := time.Now().Format("2006-01-02T15:04:05.000")
	rnd := api.Random(128)
	code := "1000"
	motoInd := 0
	req.RequestDateTime = &date
	req.RandomNumber = &rnd
	req.TxnCode = &code
	req.Transaction.MotoInd = &motoInd
	return api.Transaction(ctx, req)
}

func (api *API) PreAuth3D(ctx context.Context, req *Request) (Response, error) {
	date := time.Now().Format("2006-01-02T15:04:05.000")
	rnd := api.Random(128)
	code := "1004"
	motoInd := 0
	req.RequestDateTime = &date
	req.RandomNumber = &rnd
	req.TxnCode = &code
	req.Transaction.MotoInd = &motoInd
	return api.Transaction(ctx, req)
}

func (api *API) Auth3D(ctx context.Context, req *Request) (Response, error) {
	date := time.Now().Format("2006-01-02T15:04:05.000")
	rnd := api.Random(128)
	code := "1000"
	motoInd := 0
	req.RequestDateTime = &date
	req.RandomNumber = &rnd
	req.TxnCode = &code
	req.Transaction.MotoInd = &motoInd
	return api.Transaction(ctx, req)
}

func (api *API) PreAuth3Dhtml(ctx context.Context, req *Request) (string, error) {
	date := time.Now().Format("2006-01-02T15:04:05.000")
	rnd := api.Random(128)
	code := "3004"
	model := "3D"
	req.PaymentModel = &model
	req.RequestDateTime = &date
	req.RandomNumber = &rnd
	req.TxnCode = &code
	payload, _ := QueryString(req)
	params := []string{"paymentModel", "txnCode", "merchantSafeId", "terminalSafeId", "orderId", "lang", "amount", "ccbRewardAmount", "pcbRewardAmount", "xcbRewardAmount", "currencyCode", "installCount", "okUrl", "failUrl", "emailAddress", "subMerchantId", "creditCard", "expiredDate", "cvv", "randomNumber", "requestDateTime", "b2bIdentityNumber"}
	hash := api.Hash3D(payload, params)
	req.Hash = &hash
	return api.Transaction3D(ctx, req)
}

func (api *API) Auth3Dhtml(ctx context.Context, req *Request) (string, error) {
	date := time.Now().Format("2006-01-02T15:04:05.000")
	rnd := api.Random(128)
	code := "3000"
	model := "3D"
	req.PaymentModel = &model
	req.RequestDateTime = &date
	req.RandomNumber = &rnd
	req.TxnCode = &code
	payload, _ := QueryString(req)
	params := []string{"paymentModel", "txnCode", "merchantSafeId", "terminalSafeId", "orderId", "lang", "amount", "ccbRewardAmount", "pcbRewardAmount", "xcbRewardAmount", "currencyCode", "installCount", "okUrl", "failUrl", "emailAddress", "subMerchantId", "creditCard", "expiredDate", "cvv", "randomNumber", "requestDateTime", "b2bIdentityNumber"}
	hash := api.Hash3D(payload, params)
	req.Hash = &hash
	return api.Transaction3D(ctx, req)
}

func (api *API) PostAuth(ctx context.Context, req *Request) (Response, error) {
	date := time.Now().Format("2006-01-02T15:04:05.000")
	rnd := api.Random(128)
	code := "1005"
	req.RequestDateTime = &date
	req.RandomNumber = &rnd
	req.TxnCode = &code
	return api.Transaction(ctx, req)
}

func (api *API) Refund(ctx context.Context, req *Request) (Response, error) {
	date := time.Now().Format("2006-01-02T15:04:05.000")
	rnd := api.Random(128)
	code := "1002"
	req.RequestDateTime = &date
	req.RandomNumber = &rnd
	req.TxnCode = &code
	return api.Transaction(ctx, req)
}

func (api *API) Cancel(ctx context.Context, req *Request) (Response, error) {
	date := time.Now().Format("2006-01-02T15:04:05.000")
	rnd := api.Random(128)
	code := "1003"
	req.RequestDateTime = &date
	req.RandomNumber = &rnd
	req.TxnCode = &code
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
	request.Header.Set("auth-hash", api.Hash(payload))
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

func (api *API) Transaction3D(ctx context.Context, req *Request) (res string, err error) {
	payload, err := QueryString(req)
	if err != nil {
		return res, err
	}
	html := []string{}
	html = append(html, `<!DOCTYPE html>`)
	html = append(html, `<html>`)
	html = append(html, `<head>`)
	html = append(html, `<meta http-equiv="Content-Type" content="text/html; charset=utf-8">`)
	html = append(html, `<script type="text/javascript">function submitonload() {document.payment.submit();document.getElementById('button').remove();document.getElementById('body').insertAdjacentHTML("beforeend", "Lütfen bekleyiniz...");}</script>`)
	html = append(html, `</head>`)
	html = append(html, `<body onload="javascript:submitonload();" id="body" style="text-align:center;margin:10px;font-family:Arial;font-weight:bold;">`)
	html = append(html, `<form action="`+EndPoints[api.Mode+"3D"]+`" method="post" name="payment">`)
	for k := range payload {
		html = append(html, `<input type="hidden" name="`+k+`" value="`+payload.Get(k)+`">`)
	}
	html = append(html, `<input type="submit" value="Gönder" id="button">`)
	html = append(html, `</form>`)
	html = append(html, `</body>`)
	html = append(html, `</html>`)
	res = B64(strings.Join(html, "\n"))
	return res, err
}
