package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/parviz-yu/digital-wallet/internal/config"
	"github.com/parviz-yu/digital-wallet/internal/models"
	"github.com/parviz-yu/digital-wallet/internal/service"
	customerrors "github.com/parviz-yu/digital-wallet/pkg/custom-errors"
	"github.com/parviz-yu/digital-wallet/pkg/logger"
	"github.com/parviz-yu/digital-wallet/pkg/security"
)

var (
	ErrInvalidReqBody = errors.New("invalid request body")
	ErrInvalidAmount  = errors.New("invalid  amount")
)

type Handler struct {
	cfg config.Config
	log logger.LoggerI
	svc service.ServiceI
}

func NewHandler(cfg config.Config, log logger.LoggerI, svc service.ServiceI) *Handler {
	return &Handler{
		cfg: cfg,
		log: log,
		svc: svc,
	}
}

func (h *Handler) DoesWalletExists(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.DoesWalletExists"

	log := logger.With(
		h.log,
		logger.String("fn", fn),
		logger.String("request_id", middleware.GetReqID(r.Context())),
	)

	userID := r.Context().Value(ctxKeyUserID).(string)
	_, err := h.svc.DoesWalletExists(r.Context(), userID)
	if errors.Is(err, customerrors.ErrWalletNotFound) {
		log.Warn(err.Error(), logger.String("X-UserID", userID))

		Error(w, r, http.StatusNotFound, nil)
		return
	}
	if err != nil {
		log.Error(err.Error(), logger.String("X-UserID", userID))

		Error(w, r, http.StatusInternalServerError, nil)
		return
	}

	Respond(w, r, http.StatusOK, nil)
}

func (h *Handler) PutFunds() http.HandlerFunc {
	type request struct {
		Amount float64 `json:"amount"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.PutFunds"

		log := logger.With(
			h.log,
			logger.String("fn", fn),
			logger.String("request_id", middleware.GetReqID(r.Context())),
		)

		userID := r.Context().Value(ctxKeyUserID).(string)
		digest := r.Header.Get(digestHeader)
		if digest == "" {
			log.Warn(ErrNoXDigestHeader.Error(), logger.String("X-UserID", userID))

			Error(w, r, http.StatusUnauthorized, ErrNoXDigestHeader)
			return
		}

		req := request{}
		jsonDecoder := json.NewDecoder(r.Body)
		jsonDecoder.DisallowUnknownFields()
		if err := jsonDecoder.Decode(&req); err != nil {
			log.Error(err.Error(), logger.String("X-UserID", userID))

			Error(w, r, http.StatusBadRequest, ErrInvalidReqBody)
			return
		}
		defer r.Body.Close()

		if req.Amount < 1 {
			log.Warn("negative amount", logger.String("X-UserID", userID))

			Error(w, r, http.StatusBadRequest, ErrInvalidAmount)
			return
		}

		reqBody, err := json.Marshal(req)
		if err != nil {
			log.Error(err.Error(), logger.String("X-UserID", userID), logger.Any("reqBody", req))

			Error(w, r, http.StatusInternalServerError, nil)
			return
		}

		if !security.VerifyBody(h.cfg.SecretToket, reqBody, digest) {
			log.Warn(ErrInvalidXDigestHeader.Error(), logger.String("X-UserID", userID))

			Error(w, r, http.StatusUnauthorized, ErrInvalidXDigestHeader)
			return
		}

		paymentReq := models.PaymentReq{
			UserID: userID,
			Amount: req.Amount,
		}

		err = h.svc.PutFunds(r.Context(), &paymentReq)
		var customErr customerrors.ErrLimitExceeded
		if errors.As(err, &customErr) {
			log.Warn(err.Error(), logger.String("X-UserID", userID))
			Error(
				w,
				r,
				http.StatusOK,
				fmt.Errorf("limit exceeded, for %s is %d TJS", customErr.WalletType, customErr.MaxAmount/100))
			return
		}

		if errors.Is(err, customerrors.ErrWalletNotFound) {
			log.Warn(err.Error(), logger.String("X-UserID", userID))

			Error(w, r, http.StatusNotFound, nil)
			return
		}
		if err != nil {
			log.Error(err.Error(), logger.String("X-UserID", userID))

			Error(w, r, http.StatusInternalServerError, nil)
			return
		}

		Respond(w, r, http.StatusOK, nil)
	}
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.GetBalance"

	log := logger.With(
		h.log,
		logger.String("fn", fn),
		logger.String("request_id", middleware.GetReqID(r.Context())),
	)

	userID := r.Context().Value(ctxKeyUserID).(string)
	resp, err := h.svc.GetWalletBalance(r.Context(), userID)
	if errors.Is(err, customerrors.ErrWalletNotFound) {
		log.Warn(err.Error(), logger.String("X-UserID", userID))

		Error(w, r, http.StatusNotFound, customerrors.ErrWalletNotFound)
		return
	}
	if err != nil {
		log.Error(err.Error(), logger.String("X-UserID", userID))

		Error(w, r, http.StatusInternalServerError, nil)
		return
	}

	Respond(w, r, http.StatusOK, resp)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.GetStats"

	log := logger.With(
		h.log,
		logger.String("fn", fn),
		logger.String("request_id", middleware.GetReqID(r.Context())),
	)

	userID := r.Context().Value(ctxKeyUserID).(string)
	resp, err := h.svc.GetWalletStats(r.Context(), userID)
	if errors.Is(err, customerrors.ErrWalletNotFound) {
		log.Warn(err.Error(), logger.String("X-UserID", userID))

		Error(w, r, http.StatusNotFound, customerrors.ErrWalletNotFound)
		return
	}
	if err != nil {
		log.Error(err.Error(), logger.String("X-UserID", userID))

		Error(w, r, http.StatusInternalServerError, nil)
		return
	}

	Respond(w, r, http.StatusOK, resp)
}
