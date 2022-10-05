package connector

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type errorResponse struct {
	Err     error  `json:"-"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type Response struct {
	Data   any            `json:"data,omitempty"`
	Error  *errorResponse `json:"error,omitempty"`
	Status int            `json:"status"`
}

func SendError(c echo.Context, err error, status int) error {

	return Response{
		Status: status,
		Error: &errorResponse{
			Message: http.StatusText(status),
			Details: fmt.Sprint(err),
		},
	}.Send(c)

}

func (res Response) Send(c echo.Context) error {
	res.Make()
	return c.JSON(res.Status, res)
}

func (res *Response) Make() {

	if res.Error != nil {
		res.Error.Make()
		if res.Status == 0 {
			res.Status = http.StatusInternalServerError
		}
		res.Error.Message = http.StatusText(res.Status)
	}

	if res.Status == 0 {
		res.Status = http.StatusOK
	}
}

func (err *errorResponse) Make() {
	err.Details = fmt.Sprint(err.Err)
}
