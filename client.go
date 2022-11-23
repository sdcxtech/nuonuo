package nuonuo

import (
	"context"
	"crypto/hmac"
	"crypto/sha1" // nolint: gosec
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type Client struct {
	url       string
	appKey    string
	appSecret string
	userTax   string

	tc          TokenController
	restyClient *resty.Client
	rand        *rand.Rand
}

func New(url, appKey, appSecret, userTax string, tc TokenController) *Client {
	return &Client{
		url:         url,
		appKey:      appKey,
		appSecret:   appSecret,
		userTax:     userTax,
		tc:          tc,
		restyClient: resty.New().SetTimeout(5 * time.Second),
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())), // nolint: gosec
	}
}

type (
	OpenInvoiceRequest struct {
		Order *InvoiceOrder `json:"order"`
	}

	InvoiceOrder struct {
		BuyerName        string `json:"buyerName,omitempty"`
		BuyerTaxNum      string `json:"buyerTaxNum,omitempty"`
		BuyerTel         string `json:"buyerTel,omitempty"`
		BuyerAddress     string `json:"buyerAddress,omitempty"`
		BuyerAccount     string `json:"buyerAccount,omitempty"`
		SalerTaxNum      string `json:"salerTaxNum,omitempty"`
		SalerTel         string `json:"salerTel,omitempty"`
		SalerAddress     string `json:"salerAddress,omitempty"`
		SalerAccount     string `json:"salerAccount,omitempty"`
		OrderNo          string `json:"orderNo,omitempty"`
		InvoiceDate      string `json:"invoiceDate,omitempty"`
		InvoiceCode      string `json:"invoiceCode,omitempty"`
		InvoiceNum       string `json:"invoiceNum,omitempty"`
		RedReason        string `json:"redReason,omitempty"`
		BillInfoNo       string `json:"billInfoNo,omitempty"`
		DepartmentID     string `json:"departmentId,omitempty"`
		ClerkID          string `json:"clerkId,omitempty"`
		Remark           string `json:"remark,omitempty"`
		Checker          string `json:"checker,omitempty"`
		Payee            string `json:"payee,omitempty"`
		Clerk            string `json:"clerk,omitempty"`
		ListFlag         string `json:"listFlag,omitempty"`
		ListName         string `json:"listName,omitempty"`
		PushMode         string `json:"pushMode,omitempty"`
		BuyerPhone       string `json:"buyerPhone,omitempty"`
		Email            string `json:"email,omitempty"`
		InvoiceType      string `json:"invoiceType,omitempty"`
		InvoiceLine      string `json:"invoiceLine,omitempty"`
		PaperInvoiceType string `json:"paperInvoiceType,omitempty"`
		SpecificFactor   string `json:"specificFactor,omitempty"`
		ProxyInvoiceFlag string `json:"proxyInvoiceFlag,omitempty"`
		CallBackURL      string `json:"callBackUrl,omitempty"`
		ExtensionNumber  string `json:"extensionNumber,omitempty"`
		TerminalNumber   string `json:"terminalNumber,omitempty"`
		MachineCode      string `json:"machineCode,omitempty"`
		VehicleFlag      string `json:"vehicleFlag,omitempty"`
		HiddenBmbbbh     string `json:"hiddenBmbbbh,omitempty"`
		NextInvoiceCode  string `json:"nextInvoiceCode,omitempty"`
		NextInvoiceNum   string `json:"nextInvoiceNum,omitempty"`
		InvoiceNumEnd    string `json:"invoiceNumEnd,omitempty"`
		SurveyAnswerType string `json:"surveyAnswerType,omitempty"`
		BuyerManagerName string `json:"buyerManagerName,omitempty"`
		ManagerCardType  string `json:"managerCardType,omitempty"`
		ManagerCardNo    string `json:"managerCardNo,omitempty"`

		InvoiceDetail []*GoodsItem `json:"invoiceDetail,omitempty"`

		AdditionalElementName string               `json:"additionalElementName,omitempty"`
		AdditionalElementList []*AdditionalElement `json:"additionalElementList,omitempty"`
	}

	GoodsItem struct {
		GoodsName           string `json:"goodsName,omitempty"`
		GoodsCode           string `json:"goodsCode,omitempty"`
		SelfCode            string `json:"selfCode,omitempty"`
		WithTaxFlag         string `json:"withTaxFlag,omitempty"`
		Price               string `json:"price,omitempty"`
		Num                 string `json:"num,omitempty"`
		Unit                string `json:"unit,omitempty"`
		SpecType            string `json:"specType,omitempty"`
		Tax                 string `json:"tax,omitempty"`
		TaxRate             string `json:"taxRate,omitempty"`
		TaxExcludedAmount   string `json:"taxExcludedAmount,omitempty"`
		TaxIncludedAmount   string `json:"taxIncludedAmount,omitempty"`
		InvoiceLineProperty string `json:"invoiceLineProperty,omitempty"`
		FavouredPolicyFlag  string `json:"favouredPolicyFlag,omitempty"`
		FavouredPolicyName  string `json:"favouredPolicyName,omitempty"`
		Deduction           string `json:"deduction,omitempty"`
		ZeroRateFlag        string `json:"zeroRateFlag,omitempty"`
	}

	AdditionalElement struct {
		ElementName  string `json:"elementName,omitempty"`
		ElementType  string `json:"elementType,omitempty"`
		ElementValue string `json:"elementValue,omitempty"`
	}

	OpenInvoiceResponse struct {
		InvoiceSerialNum string `json:"invoiceSerialNum"`
	}
)

func (c *Client) OpenInvoice(
	ctx context.Context, req *OpenInvoiceRequest,
) (*OpenInvoiceResponse, error) {
	resp := &OpenInvoiceResponse{}

	err := c.request(ctx, "nuonuo.OpeMplatform.requestBillingNew", req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type (
	QueryInvoiceRequest struct {
		SerialNos            []string `json:"serialNos,omitempty"`
		OrderNos             []string `json:"orderNos,omitempty"`
		IsOfferInvoiceDetail string   `json:"isOfferInvoiceDetail,omitempty"`
	}

	InvoiceResultItem struct {
		SerialNo                  string `json:"serialNo"`
		OrderNo                   string `json:"orderNo"`
		Status                    string `json:"status"`
		StatusMsg                 string `json:"statusMsg"`
		FailCause                 string `json:"failCause"`
		PdfURL                    string `json:"pdfUrl"`
		PictureURL                string `json:"pictureUrl"`
		InvoiceTime               int64  `json:"invoiceTime"`
		InvoiceCode               string `json:"invoiceCode"`
		InvoiceNo                 string `json:"invoiceNo"`
		AllElectronicInvoiceNumbe string `json:"allElectronicInvoiceNumbe"`
		ExTaxAmount               string `json:"exTaxAmount"`
		TaxAmount                 string `json:"taxAmount"`
		OrderAmount               string `json:"orderAmount"`
		PayerName                 string `json:"payerName"`
		PayerTaxNo                string `json:"payerTaxNo"`
		Address                   string `json:"address"`
		Telephone                 string `json:"telephone"`
		BankAccount               string `json:"bankAccount"`
		InvoiceKind               string `json:"invoiceKind"`
		CheckCode                 string `json:"checkCode"`
		QrCode                    string `json:"qrCode"`
		MachineCode               string `json:"machineCode"`
		CipherText                string `json:"cipherText"`
		PaperPdfURL               string `json:"paperPdfUrl"`
		OfdURL                    string `json:"ofdUrl"`
		Clerk                     string `json:"clerk"`
		Payee                     string `json:"payee"`
		Checker                   string `json:"checker"`
		SalerAccount              string `json:"salerAccount"`
		SalerTel                  string `json:"salerTel"`
		SalerAddress              string `json:"salerAddress"`
		SalerTaxNum               string `json:"salerTaxNum"`
		SaleName                  string `json:"saleName"`
		Remark                    string `json:"remark"`
		ProductOilFlag            int    `json:"productOilFlag"`
		ImgURLs                   string `json:"imgUrls"`
		ExtensionNumber           string `json:"extensionNumber"`
		TerminalNumber            string `json:"terminalNumber"`
		DeptID                    string `json:"deptId"`
		ClerkID                   string `json:"clerkId"`
		OldInvoiceCode            string `json:"oldInvoiceCode"`
		OldInvoiceNo              string `json:"oldInvoiceNo"`
		OldEleInvoiceNumber       string `json:"oldEleInvoiceNumber"`
		ListFlag                  string `json:"listFlag"`
		ListName                  string `json:"listName"`
		Phone                     string `json:"phone"`
		NotifyEmail               string `json:"notifyEmail"`
		VehicleFlag               string `json:"vehicleFlag"`
		CreateTime                int64  `json:"createTime"`
		UpdateTime                int64  `json:"updateTime"`
		ProxyInvoiceFlag          string `json:"proxyInvoiceFlag"`
		InvoiceDate               int64  `json:"invoiceDate"`
		InvoiceType               string `json:"invoiceType"`
		RedReason                 string `json:"redReason"`
		InvalidTime               string `json:"invalidTime"`
		InvalidSource             string `json:"invalidSource"`
		InvalidReason             string `json:"invalidReason"`
		SpecificReason            string `json:"specificReason"`
		SpecificFactor            int    `json:"specificFactor"`
		BuyerManagerName          string `json:"buyerManagerName"`
		ManagerCardType           string `json:"managerCardType"`
		ManagerCardNo             string `json:"managerCardNo"`
	}
)

func (c *Client) QueryInvoice(
	ctx context.Context, req *QueryInvoiceRequest,
) ([]*InvoiceResultItem, error) {
	resp := []*InvoiceResultItem{}

	err := c.request(ctx, "nuonuo.OpeMplatform.queryInvoiceResult", req, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) sign(senid, nonce, timestamp, content string) (string, error) {
	pairs := [][2]string{
		{"a", "services"},
		{"l", "v1"},
		{"p", "open"},
		{"k", c.appKey},
		{"i", senid},
		{"n", nonce},
		{"t", timestamp},
		{"f", content},
	}

	parts := make([]string, 0, len(pairs))
	for i := range pairs {
		parts = append(parts, strings.Join(pairs[i][:], "="))
	}

	payload := strings.Join(parts, "&")

	h := hmac.New(sha1.New, []byte(c.appSecret))
	_, err := h.Write([]byte(payload))
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature, nil
}

type requestCommon struct {
	senID     string
	nonce     string
	timestamp string
	appKey    string
}

func (c *Client) newRequestCommon() *requestCommon {
	nonce := fmt.Sprintf("%08d", c.rand.Intn(100_000_000)) // nolint: gosec
	if nonce[0] == '0' {
		nonce = strconv.Itoa(c.rand.Intn(9)+1) + nonce[1:]
	}

	return &requestCommon{
		senID:     strings.ReplaceAll(uuid.New().String(), "-", ""),
		nonce:     nonce,
		timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		appKey:    c.appKey,
	}
}

func (c *Client) request(
	ctx context.Context,
	method string,
	reqBody any,
	respPtr any,
) error {
	data, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	content := string(data)

	rc := c.newRequestCommon()
	signature, err := c.sign(rc.senID, rc.nonce, rc.timestamp, content)
	if err != nil {
		return err
	}

	err = c.post(ctx, rc, method, signature, content, respPtr)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) post(
	ctx context.Context,
	rc *requestCommon,
	method string,
	signature string,
	body any,
	resultPtr any,
) error {
	token, err := c.tc.GetToken(ctx)
	if err != nil {
		return fmt.Errorf("get token: %w", err)
	}

	var result struct {
		Code     string          `json:"code"`
		Describe string          `json:"describe"`
		Result   json.RawMessage `json:"result"`
	}

	req := c.restyClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Nuonuo-Sign", signature).
		SetHeader("accessToken", token).
		SetHeader("method", method).
		SetQueryParams(map[string]string{
			"senid":     rc.senID,
			"nonce":     rc.nonce,
			"timestamp": rc.timestamp,
			"appkey":    rc.appKey,
		}).
		ForceContentType("application/json").
		SetBody(body).
		SetResult(&result)

	if c.userTax != "" {
		req.SetHeader("userTax", c.userTax)
	}

	resp, err := req.Post(c.url)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("http status: %s, body: %s", resp.Status(), resp.Body())
	}

	if result.Code != "E0000" {
		return &Error{Code: result.Code, Msg: result.Describe}
	}

	if resultPtr != nil {
		err = json.Unmarshal(result.Result, resultPtr)
		if err != nil {
			return err
		}
	}

	return nil
}
