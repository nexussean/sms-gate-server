package handlers

import (
	"github.com/capcom6/sms-gateway/internal/sms-gateway/modules/health"
	"github.com/capcom6/sms-gateway/internal/version"
	"github.com/capcom6/sms-gateway/pkg/maps"
	"github.com/capcom6/sms-gateway/pkg/smsgateway"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type healthHanlderParams struct {
	fx.In

	HealthSvc *health.Service

	Logger *zap.Logger
}

type healthHandler struct {
	Handler

	healthSvc *health.Service

	logger *zap.Logger
}

//	@Summary		Health check
//	@Description	Checks if service is healthy
//	@Tags			System
//	@Produce		json
//	@Success		200	{object}	smsgateway.HealthResponse	"Health check result"
//	@Failure		500	{object}	smsgateway.HealthResponse	"Service is unhealthy"
//	@Router			/health [get]
//	@Router			/api/3rdparty/v1/health [get]
//
// Health check
func (h *healthHandler) getHealth(c *fiber.Ctx) error {
	check, err := h.healthSvc.HealthCheck(c.Context())
	if err != nil {
		return err
	}

	res := smsgateway.HealthResponse{
		Status:    smsgateway.HealthStatus(check.Status),
		Version:   version.AppVersion,
		ReleaseID: version.AppReleaseID(),
		Checks: maps.MapValues(
			check.Checks,
			func(c health.CheckDetail) smsgateway.HealthCheck {
				return smsgateway.HealthCheck{
					Description:   c.Description,
					ObservedUnit:  c.ObservedUnit,
					ObservedValue: c.ObservedValue,
					Status:        smsgateway.HealthStatus(c.Status),
				}
			},
		),
	}

	if check.Status == health.StatusFail {
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *healthHandler) Register(router fiber.Router) {
	router.Get("/health", h.getHealth)
}

func newHealthHandler(params healthHanlderParams) *healthHandler {
	return &healthHandler{
		Handler:   Handler{Logger: params.Logger.Named("HealthHandler"), Validator: nil},
		healthSvc: params.HealthSvc,
		logger:    params.Logger,
	}
}
