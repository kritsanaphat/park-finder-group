package routers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
)

func (cr CustomerRouters) customerRegisterHandler(c echo.Context) error {
	context := c.Request().Context()

	request := new(models.RegisterAccountRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	account := cr.CustomerServices.CheckExistEmail(context, request.Email)
	if account != nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "This email is already exists"})
	}

	if err := cr.CustomerServices.CustomerRegister(context, request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	go utility.SendMailVerifyCustomer(request.Email, "customer", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_SERVER"), os.Getenv("SMTP_PORT"), os.Getenv("HOST"))

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) customerLoginHandler(c echo.Context) error {
	context := c.Request().Context()

	request := new(models.LoginRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	account := cr.CustomerServices.CheckExistEmail(context, request.Email)
	if account == nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "This email is not registered"})
	}

	if !cr.CustomerServices.CheckPassword(*account, request.Password) {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Incorrect password"})
	}

	if !account.Verify {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "This email is not verify"})
	}

	accessToken := utility.GenerateToken(request.Email, "customer")
	if accessToken == "" {
		return c.JSON(http.StatusInternalServerError, models.MessageResponse{Message: "Cannot generate new token"})
	}

	err := cr.CustomerServices.AddToken(context, account, accessToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.MessageResponse{Message: "Error while add user token"})
	}

	otp := utility.GenerateOTP()
	fmt.Println(otp)
	err = cr.CustomerServices.SaveOTP(request.Email, otp)
	if err != nil {
		fmt.Println("Error with cache OTP")
	}

	go utility.SendMailOTP(request.Email, otp, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_SERVER"), os.Getenv("SMTP_PORT"))

	return c.JSON(http.StatusOK, models.LoginResponse{AccessToken: accessToken, Data: account.ToMap()})
}

func (cr CustomerRouters) customerLogoutHandler(c echo.Context) error {
	ctx := c.Request().Context()

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}

	email := utility.GetEmailFromToken(token.Raw)
	user := cr.CustomerServices.CheckExistEmail(ctx, email)
	if user == nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "User not found"})
	}

	err := cr.CustomerServices.RevokeToken(ctx, user, token.Raw)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.MessageResponse{Message: "Error while revoke user token"})
	}

	return c.JSON(http.StatusOK, models.MessageResponse{Message: "Logout succuess"})
}

func (cr CustomerRouters) updateCustomerPassword(c echo.Context) error {
	ctx := c.Request().Context()

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Invalid jwt token"})
	}
	email := utility.GetEmailFromToken(token.Raw)

	request := new(models.ResetPasswordRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	account := cr.CustomerServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusNoContent, models.MessageResponse{Message: "This email is not registered"})
	}
	if !cr.CustomerServices.CheckPassword(*account, request.OldPassword) {
		return c.JSON(http.StatusUnauthorized, models.MessageResponse{Message: "Incorrect password"})
	}

	if err := cr.CustomerServices.UpdateCustomerPassword(ctx, request, email); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) sendOTPForgotHandler(c echo.Context) error {
	ctx := c.Request().Context()
	request := new(models.ForgotPasswordEmailRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	account := cr.CustomerServices.CheckExistEmail(ctx, request.Email)
	if account == nil {
		return c.JSON(http.StatusNoContent, models.MessageResponse{Message: "This email is not registered"})
	}

	otp := utility.GenerateOTP()
	fmt.Println(otp)
	err := cr.CustomerServices.SaveOTPForgot(request.Email, otp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	go utility.SendMailOTPPassword(request.Email, otp, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_SERVER"), os.Getenv("SMTP_PORT"))

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) resendVerifsyHandler(c echo.Context) error {
	ctx := c.Request().Context()
	request := new(models.ForgotPasswordEmailRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	account := cr.CustomerServices.CheckExistEmail(ctx, request.Email)
	if account == nil {
		return c.JSON(http.StatusNoContent, models.MessageResponse{Message: "This email is not registered"})
	}

	go utility.SendMailVerifyCustomer(request.Email, "customer", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_SERVER"), os.Getenv("SMTP_PORT"), os.Getenv("HOST"))

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) resendOTPLoginHandler(c echo.Context) error {
	ctx := c.Request().Context()
	request := new(models.ForgotPasswordEmailRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	account := cr.CustomerServices.CheckExistEmail(ctx, request.Email)
	if account == nil {
		return c.JSON(http.StatusNoContent, models.MessageResponse{Message: "This email is not registered"})
	}

	otp := utility.GenerateOTP()
	fmt.Println(otp)
	err := cr.CustomerServices.RemoveOTP(request.Email, otp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	err = cr.CustomerServices.SaveOTP(request.Email, otp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	go utility.SendMailOTP(request.Email, otp, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_SERVER"), os.Getenv("SMTP_PORT"))

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) resendOTPForgotHandler(c echo.Context) error {
	ctx := c.Request().Context()
	request := new(models.ForgotPasswordEmailRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	account := cr.CustomerServices.CheckExistEmail(ctx, request.Email)
	if account == nil {
		return c.JSON(http.StatusNoContent, models.MessageResponse{Message: "This email is not registered"})
	}

	otp := utility.GenerateOTP()
	fmt.Println(otp)
	err := cr.CustomerServices.RemoveOTPForgot(request.Email, otp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	err = cr.CustomerServices.SaveOTPForgot(request.Email, otp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	go utility.SendMailOTPPassword(request.Email, otp, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_SERVER"), os.Getenv("SMTP_PORT"))

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) updateCustomerPasswordByForgotHandler(c echo.Context) error {
	ctx := c.Request().Context()

	request := new(models.LoginRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}

	account := cr.CustomerServices.CheckExistEmail(ctx, request.Email)
	if account == nil {
		return c.JSON(http.StatusNoContent, models.MessageResponse{Message: "This email is not registered"})
	}
	if err := cr.CustomerServices.UpdateCustomerPasswordByForgot(ctx, request.Password, request.Email); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) verifyCustomerOTPForgotHandler(c echo.Context) error {

	request := new(models.OTPRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}
	if request.OTP != os.Getenv("AUTOMATE_OTP_TEST") {
		err := cr.CustomerServices.CheckOTPForgot(request.Email, request.OTP)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "invalid or expired otp"})
		}
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) verifyCustomerOTP(c echo.Context) error {

	request := new(models.OTPRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
	}
	if request.OTP != os.Getenv("AUTOMATE_OTP_TEST") {
		err := cr.CustomerServices.CheckOTP(request.Email, request.OTP)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "invalid or expired otp"})
		}
	}

	return c.JSON(http.StatusOK, models.StatusResponse{Message: "ok", Status: http.StatusOK})
}

func (cr CustomerRouters) verifyCustomerEmail(c echo.Context) error {
	ctx := c.Request().Context()
	email := c.Param("email")

	account := cr.CustomerServices.CheckExistEmail(ctx, email)
	if account == nil {
		return c.JSON(http.StatusConflict, models.MessageResponse{Message: "This email is not registered"})
	}

	if err := cr.CustomerServices.UpdateCustomerVerify(ctx, email); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.StatusResponse{Message: err.Error(), Status: http.StatusUnprocessableEntity})
	}

	template, err := template.New("verifyEmailTemplate").Parse(htmlTemplate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.StatusResponse{Message: "Failed to parse verify page template: " + err.Error(), Status: http.StatusInternalServerError})
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.StatusResponse{Message: "Failed to load verify page :" + err.Error(), Status: http.StatusInternalServerError})
	}

	return template.Execute(c.Response().Writer, nil)
}

// errorHandler handle expired jwt token by update token status and stamp revoke date
func (cr CustomerRouters) errorHandler(c echo.Context, err error) error {
	message := "invalid or expired jwt"
	switch err.(type) {
	case *echojwt.TokenExtractionError:
		message = "missing or malformed jwt"
	default:

		ctx := c.Request().Context()
		auth := c.Request().Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			message = "missing or malformed jwt"
			return c.JSON(http.StatusUnauthorized, models.MessageResponse{
				Message: message,
			})
		}
		tokenStr := strings.Split(auth, " ")[1]
		token := cr.CustomerServices.CheckExistToken(ctx, tokenStr)
		if token != nil && token.Valid {
			expire := utility.GetExpireDateFromToken(tokenStr)
			expireDate := time.Unix(expire, 0)
			cr.CustomerServices.RevokeExpireToken(ctx, tokenStr, expireDate)
		}
	}
	return c.JSON(http.StatusUnauthorized, models.MessageResponse{
		Message: message,
	})
}

// parseToken validate the given auth token, Return token if recieved token is valid otherwise return nil
func (cr CustomerRouters) parseToken(c echo.Context, auth string) (interface{}, error) {

	token, err := jwt.Parse(auth, utility.KeyFunc)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	ctx := c.Request().Context()
	dbToken := cr.CustomerServices.CheckExistToken(ctx, token.Raw)
	if !dbToken.Valid {
		return nil, fmt.Errorf("token not valid")
	}

	return token, nil
}

var htmlTemplate = `<html>
    <body style="width: 100%; height: 1000px; background-color: #d4d4d4; font-family: Arial, Helvetica, sans-serif;">
        <div
            style="
                width: 650px;
                height: 300px;
                margin: auto;
                background-color: #ffffff;
            "
        >
            <!-- Header -->
            <div
                style="
                    height: 100px;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    border-radius: 0px 0px 20px 20px;
                    background-color: #6828dc;
                "
            >
                <span
                    style="
                        font-weight: 700;
                        font-size: 24px;
                        color: #ffffff;
                    "
                    >ยืนยันเพื่อลงทะเบียนเข้าใช้งาน PARKFINDER</span
                >
            </div>

            <!-- Logo + OTP + Text -->
            <div
                style="
                    display: flex;
                    flex-direction: column;
                    align-items: center;
                "
            >
                <!-- Logo Text -->
                <div style="margin-top: 30px">
                    <span
                        style="
                            font-weight: 700;
                            font-size: 32px;
                            color: #6828dc;
                        "
                        >PARK</span
                    >
                    <span style="font-weight: 700; font-size: 32px"
                        >FINDER</span
                    >
                </div>
                <!-- Line -->
                <div
                    style="width: 75px; height: 2px; background-color: #6828dc"
                ></div>
                <!-- Button -->
                <div style="margin-top: 30px">
                    <span style="font-weight: 500; font-size: 18px;">ยืนยันการลงทะเบียนเสร็จเรียบร้อยแล้ว</span>
                </div>
                <div>
                    <span style="font-weight: 500; font-size: 18px;">คุณสามารถเข้าสู่ระบบกับ PARKFINDER ได้</span>
                </div>
            </div>
        </div>
    </body>
</html>
`
