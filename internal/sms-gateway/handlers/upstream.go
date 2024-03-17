package handlers

import (
	"time"

	"github.com/capcom6/sms-gateway/internal/sms-gateway/modules/push"
	"github.com/capcom6/sms-gateway/pkg/smsgateway"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type upstreamHandler struct {
	Handler

	config  Config
	pushSvc *push.Service
}

type upstreamHandlerParams struct {
	fx.In

	Config  Config
	PushSvc *push.Service

	Logger    *zap.Logger
	Validator *validator.Validate
}

func newUpstreamHandler(params upstreamHandlerParams) *upstreamHandler {
	return &upstreamHandler{
		Handler: Handler{Logger: params.Logger, Validator: params.Validator},
		config:  params.Config,
		pushSvc: params.PushSvc,
	}
}

//	@Summary		Send push notifications
//	@Description	Enqueues notifications for sending to devices
//	@Tags			Device, Upstream
//	@Accept			json
//	@Produce		json
//	@Param			request	body	smsgateway.UpstreamPushRequest	true	"Push request"
//	@Success		202		"Notification enqueued"
//	@Failure		400		{object}	smsgateway.ErrorResponse	"Invalid request"
//	@Failure		429		{object}	smsgateway.ErrorResponse	"Too many requests"
//	@Failure		500		{object}	smsgateway.ErrorResponse	"Internal server error"
//	@Router			/upstream/v1/push [post]
//
// Send push notifications
func (h *upstreamHandler) postPush(c *fiber.Ctx) error {
	req := smsgateway.UpstreamPushRequest{}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if len(req) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Empty request")
	}

	for _, v := range req {
		if err := h.validateStruct(v); err != nil {
			return err
		}

		if err := h.pushSvc.Enqueue(c.Context(), v.Token, map[string]string{}); err != nil {
			h.Logger.Error("Can't push message", zap.Error(err))
		}
	}

	return c.SendStatus(fiber.StatusAccepted)
}

func (h *upstreamHandler) Register(router fiber.Router) {
	// register only in public mode
	if h.config.GatewayMode != GatewayModePublic {
		return
	}

	router = router.Group("/upstream/v1")

	router.Post("/push", limiter.New(limiter.Config{
		Max:               5,
		Expiration:        60 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}), h.postPush)
}
