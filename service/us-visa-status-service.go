package service

import (
	"context"
	"crawler-visa/models"
	"crawler-visa/utils"
	"encoding/json"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
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

func RunVisaApplicationCheck(usStatus *models.QueryUsStatus) (models.UsStatus, error) {
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
		chromedp.SetValue(locationDropdown, `BEJ`, chromedp.ByID),
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
