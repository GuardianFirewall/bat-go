package rewards

import (
	"context"
	"errors"
	"net/http"

	appctx "github.com/brave-intl/bat-go/utils/context"
	"github.com/brave-intl/bat-go/utils/handlers"
	"github.com/brave-intl/bat-go/utils/inputs"
	"github.com/brave-intl/bat-go/utils/logging"
)

// GetParametersHandler - handler to get reward parameters
func GetParametersHandler(service *Service) handlers.AppHandler {
	return handlers.AppHandler(func(w http.ResponseWriter, r *http.Request) *handlers.AppError {
		// get context from request
		ctx := r.Context()

		var (
			currencyInput = r.URL.Query().Get("currency")

			// response structure
			parameters *ParametersV1
		)

		if currencyInput == "" {
			currencyInput = "USD"
		}

		// get logger from context
		logger, err := appctx.GetLogger(ctx)
		if err != nil {
			ctx, logger = logging.SetupLogger(ctx)
		}

		// in here we need to validate our currency
		var currency = new(BaseCurrency)
		if err = inputs.DecodeAndValidate(ctx, currency, []byte(currencyInput)); err != nil {
			if errors.Is(err, ErrBaseCurrencyInvalid) {
				logger.Error().Err(err).Msg("invalid currency input from caller")
				return handlers.ValidationError(
					"Error validating currency url parameter",
					map[string]interface{}{
						"err":      err.Error(),
						"currency": "invalid currency",
					},
				)
			}
			// degraded, unknown error when validating/decoding
			logger.Error().Err(err).Msg("unforseen error in decode and validation")
			return handlers.WrapError(err, "degraded: ", http.StatusInternalServerError)
		}

		parameters, err = service.GetParameters(ctx, currency)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get reward parameters")
			return handlers.WrapError(err, "failed to get parameters", http.StatusInternalServerError)
		}
		return handlers.RenderContent(ctx, parameters, w, http.StatusOK)
	})
}

// GetParametersHandlerV2 - handler to get reward parameters with relevant sku values
func GetParametersHandlerV2(service *Service) handlers.AppHandler {
	return handlers.AppHandler(func(w http.ResponseWriter, r *http.Request) *handlers.AppError {
		// get context from request
		ctx := r.Context()
		currency, appErr := checkCurrency(ctx, r.URL.Query().Get("currency"))
		if appErr != nil {
			return appErr
		}
		parameters, err := service.GetParametersV2(ctx, currency)
		if err != nil {
			return handlers.WrapError(err, "failed to get parameters", http.StatusInternalServerError)
		}
		return handlers.RenderContent(ctx, parameters, w, http.StatusOK)
	})
}

func checkCurrency(ctx context.Context, currencyInput string) (*BaseCurrency, *handlers.AppError) {
	// in here we need to validate our currency
	logger, _ := appctx.GetLogger(ctx)

	var currency = new(BaseCurrency)

	if currencyInput == "" {
		currencyInput = "USD"
	}

	if err := inputs.DecodeAndValidate(ctx, currency, []byte(currencyInput)); err != nil {
		if errors.Is(err, ErrBaseCurrencyInvalid) {
			logger.Error().Err(err).Msg("invalid currency input from caller")
			return currency, handlers.ValidationError(
				"Error validating currency url parameter",
				map[string]interface{}{
					"err":      err.Error(),
					"currency": "invalid currency",
				},
			)
		}
		// degraded, unknown error when validating/decoding
		logger.Error().Err(err).Msg("unforseen error in decode and validation")
		return currency, handlers.WrapError(err, "degraded: ", http.StatusInternalServerError)
	}
	return currency, nil
}
