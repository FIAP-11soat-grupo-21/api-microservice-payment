package suites

import (
	"context"
	"payment_microservice/internal/bdd/steps"
	"payment_microservice/internal/common/mocks"

	"github.com/cucumber/godog"
	"go.uber.org/mock/gomock"
)

type godogReporter struct{}

func (r *godogReporter) Errorf(_ string, _ ...interface{}) {
}
func (r *godogReporter) Fatalf(format string, _ ...interface{}) { panic("gomock fatal: " + format) }

func InitializeScenario(ctx *godog.ScenarioContext) {
	var helper *steps.PaymentHelper

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		handler := gomock.NewController(&godogReporter{})
		mockPG := new(mocks.MockPaymentGateway)
		mockRepo := new(mocks.MockPaymentRepository)

		helper = &steps.PaymentHelper{
			Ctrl:     handler,
			MockPG:   mockPG,
			MockRepo: mockRepo,
		}
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if helper != nil && helper.Ctrl != nil {
			helper.Ctrl.Finish()
		}
		return ctx, nil
	})

	ctx.Step(`^the payment data is valid with order_id and amount$`, func(amount float64) error {
		return helper.ThePaymentDataIsValid()
	})
	ctx.Step(`^I send a message to create a new PIX billing$`, func() error {
		return helper.SendMessageToCreatePIXBilling()
	})
	ctx.Step(`^the PIX billing should be created successfully$`, func() error {
		return helper.PaymentShouldBeCreated()
	})
}
