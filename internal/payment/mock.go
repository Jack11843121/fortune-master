package payment

// Mock 模拟支付，用于本地开发与演示：直接返回一个“确认支付”地址。
type Mock struct{}

func NewMock() *Mock { return &Mock{} }

func (m *Mock) Name() string { return "mock" }

func (m *Mock) CreateOrder(o OrderInfo) (*PayResult, error) {
	return &PayResult{
		Provider: "mock",
		PayParams: map[string]any{
			"amount":      o.Amount,
			"description": o.Desc,
		},
		MockPayURL: "/api/v1/pay/mock/confirm?order_no=" + o.OrderNo,
	}, nil
}