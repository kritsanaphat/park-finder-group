package utility

import (
	"fmt"
	"net/smtp"
)

func SendMailOTP(receiverEmail string, otp string, username string, password string, smtpServer string, smtpPort string) {
	fmt.Println("send mail OTP..")
	senderEmail := "easy.park.finder@gmail.com"
	subject := "รหัส OTP ในการเข้าสู่ระบบของ PARKFINDER"
	htmlBody := `
	<!DOCTYPE html>
	<html>
        <body>
            <div style="width: 100%; 
                        height: 1000px; 
                        background-color: #d4d4d4">

                <div style="width: 650px; 
							height: 300px;
                            margin: auto; 
                            background-color: #ffffff;
                            padding-bottom: 25px;">

                    <div style="background-color: #6828dc; 
                                margin: auto; 
                                width: 650px; 
                                padding-top: 10px; 
                                padding-bottom: 10px; 
                                border-radius: 0px 0px 20px 20px;
                                text-align: center;">

                        <p style="font-weight: 700; font-size: 24px; color:#ffffff;">รหัส OTP ในการเข้าสู่ระบบของ PARKFINDER</p>
                    </div>
                    <div style="text-align: center; ">
                        <div style="margin-top: 20px">
                            <span style="font-weight: 700; font-size: 32px; color: #6828dc;">PARK</span>
                            <span style="font-weight: 700; font-size: 32px">FINDER</span>
                        </div>
                    
                        <div style="margin-top: 10px; width: 75px; display: inline-block; height: 2px; background-color: #6828dc"></div>
                    
                        <div style="margin-top: 15px">
                            <span style="font-weight: 700; font-size: 32px">` + otp + `</span>
                        </div>
                        <div style="margin-top: 15px">
                            <span style="font-weight: 500;">กรุณากรอกรหัส OTP ภายใน 5 นาที</span>
                        </div>
                    </div>
                </div>
            </div>
        </body>
    </html>`

	message := "From: " + senderEmail + "\n" +
		"To: " + receiverEmail + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		htmlBody

	auth := smtp.PlainAuth("", username, password, smtpServer)
	err := smtp.SendMail(smtpServer+":"+smtpPort,
		auth,
		username, []string{receiverEmail}, []byte(message))

	if err != nil {
		fmt.Println("smtp error:", err)
		return
	}

	fmt.Println("Email sent successfully!")
}

func SendMailOTPPassword(receiverEmail string, otp string, username string, password string, smtpServer string, smtpPort string) {
	fmt.Println("send mail OTP..")
	senderEmail := "easy.park.finder@gmail.com"
	subject := "รหัส OTP ในการเปลี่ยนรหัสผ่าน PARKFINDER"
	htmlBody := `
	<!DOCTYPE html>
	<html>
        <body>
            <div style="width: 100%; 
                        height: 1000px; 
                        background-color: #d4d4d4">

                <div style="width: 650px; 
							height: 300px;
                            margin: auto; 
                            background-color: #ffffff;
                            padding-bottom: 25px;">

                    <div style="background-color: #6828dc; 
                                margin: auto; 
                                width: 650px; 
                                padding-top: 10px; 
                                padding-bottom: 10px; 
                                border-radius: 0px 0px 20px 20px;
                                text-align: center;">

                        <p style="font-weight: 700; font-size: 24px; color:#ffffff;">รหัส OTP ในการเข้าสู่ระบบของ PARKFINDER</p>
                    </div>
                    <div style="text-align: center; ">
                        <div style="margin-top: 20px">
                            <span style="font-weight: 700; font-size: 32px; color: #6828dc;">PARK</span>
                            <span style="font-weight: 700; font-size: 32px">FINDER</span>
                        </div>
                    
                        <div style="margin-top: 10px; width: 75px; display: inline-block; height: 2px; background-color: #6828dc"></div>
                    
                        <div style="margin-top: 15px">
                            <span style="font-weight: 700; font-size: 32px">` + otp + `</span>
                        </div>
                        <div style="margin-top: 15px">
                            <span style="font-weight: 500;">กรุณากรอกรหัส OTP ภายใน 5 นาที</span>
                        </div>
                    </div>
                </div>
            </div>
        </body>
    </html>`

	message := "From: " + senderEmail + "\n" +
		"To: " + receiverEmail + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		htmlBody

	auth := smtp.PlainAuth("", username, password, smtpServer)
	err := smtp.SendMail(smtpServer+":"+smtpPort,
		auth,
		username, []string{receiverEmail}, []byte(message))

	if err != nil {
		fmt.Println("smtp error:", err)
		return
	}

	fmt.Println("Email sent successfully!")
}
func SendMailVerifyCustomer(receiverEmail, userType, username, password, smtpServer, smtpPort, host string) {
	fmt.Println("send mail OTP..")

	senderEmail := "easy.park.finder@gmail.com"
	subject := "ยืนยันเพื่อลงทะเบียนเข้าใช้งาน PARKFINDER"
	htmlBody := `
	<!DOCTYPE html>
	<html>
        <body>
            <div style="width: 100%; 
                        height: 1000px; 
                        background-color:  #ffffff">

                <div style="width: 650px; 
							height: 300px;
                            margin: auto; 
                            background-color: #ffffff;
                            padding-bottom: 25px;">

                    <div style="background-color: #6d00fc; 
                                margin: auto; 
                                width: 650px; 
                                padding-top: 10px; 
                                padding-bottom: 10px; 
                                border-radius: 0px 0px 20px 20px;
                                text-align: center;">

                        <p style="font-weight: 700; font-size: 24px; color:#ffffff;">ยืนยันเพื่อลงทะเบียนเข้าใช้งาน PARKFINDER</p>
                    </div>
                    <div style="text-align: center; ">
                        <div style="margin-top: 20px">
                            <span style="font-weight: 700; font-size: 32px; color: #6d00fc;">PARK</span>
                            <span style="font-weight: 700; font-size: 32px">FINDER</span>
                        </div>
                    
                        <div style="margin-top: 10px; width: 75px; display: inline-block; height: 2px; background-color:  #6d00fc"></div>
                    
                        <div style="margin-top: 30px;  ">
							<a href="http://` + host + `/` + userType + `/verify_email/` + receiverEmail + `" target="_blank" style="
                                display: inline-block;
                                padding: 15px 0;
								width: 250px;
								border-radius: 10px;
								text-decoration: none; 
								font-weight: 700;
								color: #ffffff;
								background-color:  #6d00fc;">
							ยืนยันการลงทะเบียน</a>
                        </div>
                    </div>
                </div>
            </div>
        </body>
    </html>`

	message := "From: " + senderEmail + "\n" +
		"To: " + receiverEmail + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		htmlBody
	auth := smtp.PlainAuth("", username, password, smtpServer)
	err := smtp.SendMail(smtpServer+":"+smtpPort,
		auth,
		username, []string{receiverEmail}, []byte(message))

	if err != nil {
		fmt.Println("smtp error:", err)
		return
	}

	fmt.Println("Email sent successfully!")
}

func SendMailVerifyProvider(receiverEmail, userType, username, password, smtpServer, smtpPort, host string) {
	fmt.Println("send mail OTP..")

	senderEmail := "easy.park.finder@gmail.com"
	subject := "ยืนยันเพื่อลงทะเบียนเข้าใช้งาน PARKFINDER"
	htmlBody := `
	<!DOCTYPE html>
	<html>
        <body>
            <div style="width: 100%; 
                        height: 1000px; 
                        background-color:  #ffffff">

                <div style="width: 650px; 
							height: 300px;
                            margin: auto; 
                            background-color: #ffffff;
                            padding-bottom: 25px;">

                    <div style="background-color: #174dff; 
                                margin: auto; 
                                width: 650px; 
                                padding-top: 10px; 
                                padding-bottom: 10px; 
                                border-radius: 0px 0px 20px 20px;
                                text-align: center;">

                        <p style="font-weight: 700; font-size: 24px; color:#ffffff;">ยืนยันเพื่อลงทะเบียนเข้าใช้งาน PARKFINDER</p>
                    </div>
                    <div style="text-align: center; ">
                        <div style="margin-top: 20px">
                            <span style="font-weight: 700; font-size: 32px; color: #174dff;">PARK</span>
                            <span style="font-weight: 700; font-size: 32px">FINDER</span>
                        </div>
                    
                        <div style="margin-top: 10px; width: 75px; display: inline-block; height: 2px; background-color:  #174dff"></div>
                    
                        <div style="margin-top: 30px;  ">
							<a href="http://` + host + `/` + userType + `/verify_email/` + receiverEmail + `" target="_blank" style="
                                display: inline-block;
                                padding: 15px 0;
								width: 250px;
								border-radius: 10px;
								text-decoration: none; 
								font-weight: 700;
								color: #ffffff;
								background-color:  #174dff;">
							ยืนยันการลงทะเบียน</a>
                        </div>
                    </div>
                </div>
            </div>
        </body>
    </html>`

	message := "From: " + senderEmail + "\n" +
		"To: " + receiverEmail + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		htmlBody
	auth := smtp.PlainAuth("", username, password, smtpServer)
	err := smtp.SendMail(smtpServer+":"+smtpPort,
		auth,
		username, []string{receiverEmail}, []byte(message))

	if err != nil {
		fmt.Println("smtp error:", err)
		return
	}

	fmt.Println("Email sent successfully!")
}
