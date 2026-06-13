package echo

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/himbo22/xoxz/common-service/xoxz/logger"
	"github.com/himbo22/xoxz/common-service/xoxz/model"
	"github.com/himbo22/xoxz/common-service/xoxz/util"
	"github.com/labstack/echo/v5"
)

func ErrorHandler(c *echo.Context, xoxzLogger logger.XoxzLogger, err error) {
	httpCode := http.StatusInternalServerError
	res := model.ResponseModel{
		StatusCode: 9999, // General system error code
		Message:    "System error, please try again later",
	}

	// 1. Check if this is an AppError that we intentionally threw
	if appErr, ok := err.(*util.AppError); ok {
		httpCode = appErr.HTTPCode
		res.StatusCode = appErr.CustomCode
		res.Message = appErr.Message
		res.Error = appErr.Detail
	} else if errors.Is(err, echo.ErrNotFound) {
		// 2. Check for echo.ErrNotFound (default 404)
		httpCode = 404
		res.Message = "Not Found"
	} else if he, ok := err.(*echo.HTTPError); ok {
		// 2. Check if Echo internally threw this error (e.g., 404 wrong URL, 413 body too large)
		httpCode = he.Code
		res.Message = fmt.Sprintf("%v", he.Message)
	} else {
		// 3. If we land here, there's a code bug, panic, or raw DB error.
		// THIS IS WHERE YOU USE ZAP LOGGER TO WRITE LOG FILES AND FIND THE REAL ERROR
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			xoxzLogger.WithEcho().Warn("Client disconnected or timeout", logger.Error(err))
			httpCode = 499 // Nginx standard code for Client Closed Request
			res.Message = "Client disconnected"
		} else {
			// --- THIS IS A REAL BUG ---
			xoxzLogger.WithEcho().Error("Unhandled system error",
				logger.Error(err),
				logger.String("path", c.Request().URL.Path), // Log the URL being called for easier debugging
			)
		}
	}

	// Send JSON response to Frontend
	if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
		// Fixed logic: Only write JSON when HTTP headers have NOT been sent yet (!Committed)
		if !resp.Committed {
			err := c.JSON(httpCode, res)
			if err != nil {
				xoxzLogger.WithEcho().Errorf("Error writing error response: %v", err)
			}
		} else {
			// Response was committed mid-flight (e.g., streaming file then DB crashed)
			xoxzLogger.WithEcho().Errorf("Error %v occurred but HTTP was already committed. Skipping JSON write.", err)
		}
	}
}
