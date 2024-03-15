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

// 诺税通saas请求开具发票接口
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

// 诺税通saas发票详情查询接口。
// 调用该接口获取发票开票结果等有关发票信息，部分字段需要配置才返回。
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

type (
	FastInvoiceRedRequest struct {
		OrderNo           string `json:"orderNo,omitempty"`           // 订单号,每个企业唯一
		ExtensionNumber   string `json:"extensionNumber,omitempty"`   // 分机号（只能为空或者数字；不传默认取蓝票的分机，传了则以传入的为准）
		ClerkID           string `json:"clerkId,omitempty"`           // 开票员id（诺诺系统中的id）
		DeptID            string `json:"deptId,omitempty"`            // 部门门店id（诺诺系统中的id）
		OrderTime         string `json:"orderTime,omitempty"`         // 单据时间
		TaxNum            string `json:"taxNum,omitempty"`            // 销方企业税号（需要校验与开放平台头部报文中的税号一致）
		InvoiceCode       string `json:"invoiceCode,omitempty"`       // 对应蓝票发票代码
		InvoiceNumber     string `json:"invoiceNumber,omitempty"`     // 对应蓝票发票号码
		ElecInvoiceNumber string `json:"elecInvoiceNumber,omitempty"` // 对应蓝字数电票号码,蓝票为数电票时，请传入该字段
		InvoiceID         string `json:"invoiceId,omitempty"`         // 对应蓝票发票流水号
		BillNo            string `json:"billNo,omitempty"`            // 红字确认单编号,全电红票必传
		BillUUID          string `json:"billUuid,omitempty"`          // 红字确认单uuid
		InvoiceLine       string `json:"invoiceLine,omitempty"`       // 全电发票票种
		CallBackURL       string `json:"callBackUrl,omitempty"`       // 回调地址
	}

	FastInvoiceRedResponse struct {
		InvoiceSerialNum string `json:"invoiceSerialNum"` // 发票流水号
	}
)

// 诺税通saas发票快捷冲红接口。用于全电发票快捷冲红。
func (c *Client) FastInvoiceRed(
	ctx context.Context, req *FastInvoiceRedRequest,
) (*FastInvoiceRedResponse, error) {
	resp := &FastInvoiceRedResponse{}

	err := c.request(ctx, "nuonuo.OpeMplatform.fastInvoiceRed", req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type (
	SaveInvoiceRedConfirmRequest struct {
		BillID          string `json:"billId,omitempty"` // 红字确认单申请号，需要保持唯一，不传的话系统自动生成一个
		BlueInvoiceLine string `json:"blueInvoiceLine"`  // 对应蓝票发票种类
		ApplySource     string `json:"applySource"`      // 申请方（录入方）身份： 0 销方 1 购方

		// 对应蓝字发票号码(蓝票是增值税发票时必传，长度为8位数字，若传20位数字则 视为是蓝字数电票号码)
		BlueInvoiceNumber string `json:"blueInvoiceNumber,omitempty"`
		BlueInvoiceCode   string `json:"blueInvoiceCode,omitempty"` // 对应蓝字发票代码(蓝票是增值税发票时必传)

		// 对应蓝字数电票号码(数电普票、数电专票、数纸普票、数纸专票都需要传，蓝票是增值税发票时不传)
		BlueElecInvoiceNumber string `json:"blueElecInvoiceNumber,omitempty"`
		BillTime              string `json:"billTime,omitempty"`     // 填开时间（时间戳格式），默认为当前时间
		SellerTaxNo           string `json:"sellerTaxNo"`            // 销方税号
		SellerName            string `json:"sellerName"`             // 销方名称，申请说明为销方申请时可为空
		DepartmentID          string `json:"departmentId,omitempty"` // 部门门店id（诺诺网系统中的id）
		ClerkID               string `json:"clerkId,omitempty"`      // 开票员id（诺诺网系统中的id）
		BuyerTaxNo            string `json:"buyerTaxNo,omitempty"`   // 购方税号
		BuyerName             string `json:"buyerName"`              // 购方名称

		// 蓝字发票增值税用途（预留字段可为空）: 1 勾选抵扣 2 出口退税 3 代办出口退税 4 不抵扣
		VatUsage        string `json:"vatUsage,omitempty"`
		SaleTaxUsage    string `json:"saleTaxUsage,omitempty"`    // 蓝字发票消费税用途（预留字段可为空）
		AccountStatus   string `json:"accountStatus,omitempty"`   // 发票入账状态（预留字段可为空）： 0 未入账 1 已入账
		RedReason       string `json:"redReason"`                 // 冲红原因： 1销货退回 2开票有误 3服务中止 4销售折让
		ExtensionNumber string `json:"extensionNumber,omitempty"` // 分机号
	}

	SaveInvoiceRedConfirmResponse struct {
		BillID string `json:"billId"` // 红字确认单申请号
	}
)

// 诺税通saas红字确认单申请接口
func (c *Client) SaveInvoiceRedConfirm(
	ctx context.Context, req *SaveInvoiceRedConfirmRequest,
) (*SaveInvoiceRedConfirmResponse, error) {
	resp := &SaveInvoiceRedConfirmResponse{}

	err := c.request(ctx, "nuonuo.OpeMplatform.saveInvoiceRedConfirm", req, &resp.BillID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type (
	QueryInvoiceRedConfirmRequest struct {
		Identity      string `json:"identity"`                // 操作方身份：0销方 1购方
		BillStatus    string `json:"billStatus,omitempty"`    // 红字确认单状态（不传则查全部状态）
		BillID        string `json:"billId,omitempty"`        // 红字确认单申请号
		BillNo        string `json:"billNo,omitempty"`        // 红字确认单编号
		BillUUID      string `json:"billUuid,omitempty"`      // 红字确认单 uuid
		BillTimeStart string `json:"billTimeStart,omitempty"` // 填开起始时间
		BillTimeEnd   string `json:"billTimeEnd,omitempty"`   // 填开结束时间
		PageSize      string `json:"pageSize,omitempty"`      // 每页数量（默认10，最大50）
		PageNo        string `json:"pageNo,omitempty"`        // 当前页码（默认1）
	}

	QueryInvoiceRedConfirmResponse struct {
		Total int                      `json:"total"` // 总数
		List  []*InvoiceRedConfirmItem `json:"list"`  // 列表
	}

	InvoiceRedConfirmItem struct {
		BillNo            string `json:"billNo"`            // 红字确认单编号
		BillUUID          string `json:"billUuid"`          // 红字确认单uuid
		BillStatus        string `json:"billStatus"`        // 红字确认单状态
		RequestStatus     string `json:"requestStatus"`     // 操作状态
		OpenStatus        int    `json:"openStatus"`        // 已开具红字发票标记
		ApplySource       int    `json:"applySource"`       // 录入方身份
		BlueInvoiceLine   string `json:"blueInvoiceLine"`   // 蓝字发票票种
		BlueInvoiceNumber string `json:"blueInvoiceNumber"` // 对应蓝票号码
		BlueInvoiceTime   string `json:"blueInvoiceTime"`   // 蓝字发票开票日期
		BillTime          string `json:"billTime"`          // 申请日期
		ConfirmTime       string `json:"confirmTime"`       // 确认日期
		SellerTaxNo       string `json:"sellerTaxNo"`       // 销方税号
		SellerName        string `json:"sellerName"`        // 销方名称
		BuyerTaxNo        string `json:"buyerTaxNo"`        // 购方税号
		BuyerName         string `json:"buyerName"`         // 购方名称
		TaxExcludedAmount string `json:"taxExcludedAmount"` // 冲红合计金额(不含税)
		TaxAmount         string `json:"taxAmount"`         // 冲红合计税额
		RedReason         string `json:"redReason"`         // 冲红原因
		PdfURL            string `json:"pdfUrl"`            // 申请表pdf地址
	}
)

// 诺税通saas红字确认单查询接口
func (c *Client) QueryInvoiceRedConfirm(
	ctx context.Context, req *QueryInvoiceRedConfirmRequest,
) (*QueryInvoiceRedConfirmResponse, error) {
	resp := &QueryInvoiceRedConfirmResponse{}

	err := c.request(ctx, "nuonuo.OpeMplatform.queryInvoiceRedConfirm", req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type (
	ConfirmRedInvoiceRequest struct {
		ExtensionNumber  string `json:"extensionNumber,omitempty"` // 分机号
		ClerkID          string `json:"clerkId,omitempty"`         // 开票员id（诺诺网系统中的id）
		DeptID           string `json:"deptId,omitempty"`          // 部门id（诺诺网系统中的id）
		BillUUID         string `json:"billUuid,omitempty"`        // 红字确认单uuid
		BillID           string `json:"billId,omitempty"`          // 红字确认单申请号
		BillNo           string `json:"billNo,omitempty"`          // 红字确认单编号
		Identity         string `json:"identity"`                  // 操作方（确认方）身份： 0：销方 1：购方
		ConfirmAgreement string `json:"confirmAgreement"`          // 处理意见： 0：拒绝 1：同意
		ConfirmReason    string `json:"confirmReason,omitempty"`   // 处理理由
	}

	ConfirmRedInvoiceResponse struct {
		Result json.RawMessage `json:"result"` // 结果
	}
)

// 诺税通saas红字确认单确认接口
func (c *Client) ConfirmRedInvoice(
	ctx context.Context, req *ConfirmRedInvoiceRequest,
) (*ConfirmRedInvoiceResponse, error) {
	resp := &ConfirmRedInvoiceResponse{}

	err := c.request(ctx, "nuonuo.OpeMplatform.confirm", req, &resp.Result)
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
		List     json.RawMessage `json:"list"`
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
