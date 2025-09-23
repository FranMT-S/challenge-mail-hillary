package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"api/logger"
	"api/models"
	"api/services"

	"github.com/go-chi/render"
)

// MailController handles mail-related operations
type MailController struct {
	MailService services.EmailService
}

// NewMailController creates a new MailController
func NewMailController(mailService services.EmailService) *MailController {
	return &MailController{
		MailService: mailService,
	}
}

// SearchMails searches for emails based on a query
func (c *MailController) SearchMails(w http.ResponseWriter, r *http.Request) {
	cancelationToken, cancel := context.WithCancel(r.Context())
	defer cancel()

	var apiError *models.ApiError
	var query models.QuerySearch

	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		response := models.NewResponse(models.StatusError, models.MailResponse{
			Mails: []models.Email{},
			Total: 0,
		}, "The request is not valid")

		logger.Logger().Error().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Interface("body_parsed", query).
			Msg(err.Error())

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response)
		return
	}

	query.Normalize()

	getEmailsResponse, err := c.MailService.SearchEmails(cancelationToken, query)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.As(err, &apiError) {
			response := models.NewResponse(models.StatusError, models.MailResponse{
				Mails: []models.Email{},
				Total: 0,
			}, apiError.Message)

			logger.Logger().Error().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Interface("body_parsed", query).
				Err(apiError.Caused).
				Msg(apiError.Message)

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response)
			return
		}

		response := models.NewResponse(models.StatusError, models.MailResponse{
			Mails: []models.Email{},
			Total: 0,
		}, "Internal Server Error")

		logger.Logger().Error().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Interface("body_parsed", query).
			Err(err).
			Msg(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response)
		return
	}

	if len(getEmailsResponse.Emails) == 0 {
		response := models.NewResponse(models.StatusNoData, models.MailResponse{
			Mails: []models.Email{},
			Total: 0,
		}, "")

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response)
		return
	}

	response := models.NewResponse(models.StatusSuccess, models.MailResponse{Mails: getEmailsResponse.Emails, Total: getEmailsResponse.Total}, "")
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, response)
}
