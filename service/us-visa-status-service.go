package service

import (
	"context"
	"crawler-visa/models"
	"crawler-visa/utils"
	"encoding/json"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/emersion/go-imap"
	imapID "github.com/emersion/go-imap-id"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	visaAppTypeSelector = `#Visa_Application_Type`
	locationDropdown    = `#Location_Dropdown`
	caseNumberInput     = `#Visa_Case_Number`
	passportNumberInput = `#Passport_Number`
	surnameInput        = `#Surname`
	captchaInput        = `#Captcha`
	statusTranslation   = `#ctl00_ContentPlaceHolder1_ucApplicationStatusView_pTranslation`
	submitDate          = `#ctl00_ContentPlaceHolder1_ucApplicationStatusView_lblSubmitDate`
	statusDate          = `#ctl00_ContentPlaceHolder1_ucApplicationStatusView_lblStatusDate`
	captchaImage        = `#c_status_ctl00_contentplaceholder1_defaultcaptcha_CaptchaImage`
	folderButton        = `#ctl00_ContentPlaceHolder1_imgFolder`
)

func RunVisaStatusCheck(usStatus *models.QueryUsStatus) (models.UsStatus, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false), // 是否启用无头模式
		chromedp.WindowSize(1920, 1080),  // 设置屏幕分辨率
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	statusCheck, err := performVisaStatusCheck(taskCtx, usStatus)
	if err != nil {
		return models.UsStatus{}, err
	}

	return statusCheck, nil
}

func performVisaStatusCheck(taskCtx context.Context, usStatus *models.QueryUsStatus) (models.UsStatus, error) {
	log.Printf("Performing visa status check, Location: %s, Application ID: %s, Passport Number: %s, Surname Initials: %s\n",
		usStatus.Location, usStatus.ApplicationID, usStatus.PassportNumber, usStatus.First5LettersOfSurname)
	var usStatusResult models.UsStatus

	if err := godotenv.Load(".env"); err != nil {
		return usStatusResult, fmt.Errorf("error loading .env file: %w", err)
	}

	client := utils.ChaoJiYing{}
	client.InitWithOptions()

	var imageBuf []byte
	if err := chromedp.Run(taskCtx,
		chromedp.Navigate("https://ceac.state.gov/CEACStatTracker/Status.aspx"),
		chromedp.WaitVisible(visaAppTypeSelector, chromedp.ByID),
		chromedp.SetValue(visaAppTypeSelector, `NIV`, chromedp.ByID),
		chromedp.WaitVisible(locationDropdown, chromedp.ByID),
		chromedp.SetValue(locationDropdown, usStatus.Location, chromedp.ByID),
		chromedp.WaitVisible(caseNumberInput, chromedp.ByID),
		chromedp.SetValue(caseNumberInput, usStatus.ApplicationID, chromedp.ByID),
		chromedp.WaitVisible(passportNumberInput, chromedp.ByID),
		chromedp.SetValue(passportNumberInput, usStatus.PassportNumber, chromedp.ByID),
		chromedp.WaitVisible(surnameInput, chromedp.ByID),
		chromedp.SetValue(surnameInput, usStatus.First5LettersOfSurname, chromedp.ByID),
		chromedp.WaitVisible(captchaImage, chromedp.ByQuery),
		chromedp.Screenshot(captchaImage, &imageBuf, chromedp.NodeVisible),
	); err != nil {
		return usStatusResult, err
	}

	if err := ioutil.WriteFile("captcha.png", imageBuf, 0644); err != nil {
		return usStatusResult, err
	}

	var result models.ChaoJiYing
	response := client.GetPicVal(
		os.Getenv("CJY_USERNAME"),
		os.Getenv("CJY_PASSWORD"),
		os.Getenv("CJY_SOFT_ID"),
		os.Getenv("CJY_CODE_TYPE"),
		os.Getenv("CJY_MIN_LEN"),
		"captcha.png")
	if err := json.Unmarshal(response, &result); err != nil {
		return usStatusResult, fmt.Errorf("failed to get captcha value: %w", err)
	}

	if err := chromedp.Run(taskCtx,
		chromedp.WaitVisible(captchaInput, chromedp.ByID),
		chromedp.SetValue(captchaInput, result.PicStr, chromedp.ByID),
		chromedp.Click(folderButton, chromedp.ByID),
		chromedp.WaitVisible(statusTranslation, chromedp.ByQuery),
		chromedp.Text(statusTranslation, &usStatusResult.Status, chromedp.NodeVisible),
		chromedp.Text(submitDate, &usStatusResult.Created, chromedp.NodeVisible),
		chromedp.Text(statusDate, &usStatusResult.LastUpdated, chromedp.NodeVisible),
	); err != nil {
		return usStatusResult, err
	}
	return usStatusResult, nil
}

func RunVisaEmailTracking(usStatus *models.QueryUsStatus) (models.UsStatus, error) {
	var usStatusResult models.UsStatus

	msg := gomail.NewMessage()
	msg.SetHeader("From", "wuaivisa008@163.com")
	msg.SetHeader("To", "passportstatus@ustraveldocs.com")
	msg.SetHeader("Subject", usStatus.PassportNumber)
	msg.SetBody("text/html", usStatus.PassportNumber)
	n := gomail.NewDialer("smtp.163.com", 465, "wuaivisa008@163.com", "FKKOIOXQCFCRWHFH")
	if err := n.DialAndSend(msg); err != nil {
		log.Printf("发送失败")
	} else {
		log.Println("发送成功")
	}
	time.Sleep(25 * time.Second)

	c, err := client.DialTLS("imap.163.com:993", nil)
	if err != nil {
		return usStatusResult, err
	}

	// 不要忘记退出
	defer c.Logout()

	// 登录
	if err := c.Login("wuaivisa008@163.com", "FKKOIOXQCFCRWHFH"); err != nil {
		return usStatusResult, err

	}
	log.Println("登录")

	// 设置客户端ID
	_, err = imapID.NewClient(c).ID(imapID.ID{
		imapID.FieldName:    "IMAPClient",
		imapID.FieldVersion: "3.1.0",
	})
	if err != nil {
		return usStatusResult, err
	}

	// 选择收件箱
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return usStatusResult, err
	}
	seqset := new(imap.SeqSet)
	seqset.AddNum(mbox.Messages)

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, "BODY[]"}, messages)
	}()

	// 读取并打印邮件标题和内容
	msgEmail := <-messages
	if msg != nil {
		log.Println("标题:", msgEmail.Envelope.Subject)
		if body := msgEmail.GetBody(&imap.BodySectionName{}); body != nil {
			bodyBytes, _ := ioutil.ReadAll(body)
			content := string(bodyBytes)
			r := strings.NewReader(content)
			mr, err := mail.CreateReader(r)
			if err != nil {
				panic(err)
			}
			header := mr.Header
			if received, err := header.Date(); err == nil {
				loc, err := time.LoadLocation("Asia/Shanghai")
				if err != nil {
					fmt.Println("加载时区失败, 使用 UTC 时区:", err)
					loc = time.UTC
				}
				receivedInGMT8 := received.In(loc)
				fmt.Println("收件时间:", receivedInGMT8)
			} else {
				fmt.Println("读取邮件时间失败:", err)
			}
			// 打印邮件正文
			for {
				p, err := mr.NextPart()
				if err != nil {
					break
				}

				switch p.Header.(type) {
				case *mail.InlineHeader:
					// 处理内联部分
					b, _ := ioutil.ReadAll(p.Body)
					fmt.Println("正文:", string(b))
					usStatusResult.Status = string(b)
				case *mail.AttachmentHeader:
					// 忽略附件
					continue
				}
			}
		}
	} else {
		log.Println("没有消息")
	}
	if err := <-done; err != nil {
		return usStatusResult, err
	}
	return usStatusResult, err

}
