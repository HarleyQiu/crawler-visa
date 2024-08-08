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

func RunVisaApplicationCheck(usStatus *models.QueryUsStatus) (string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false), // 是否启用无头模式
		chromedp.WindowSize(1920, 1080),  // 设置屏幕分辨率
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	if err := chromedp.Run(taskCtx); err != nil {
		return "", err
	}

	statusCheck, err := performVisaStatusCheck(taskCtx, usStatus)
	if err != nil {
		log.Println("Error performing Visa Status Check:", err)
		return "", err
	}
	return statusCheck, nil
}

func performVisaStatusCheck(taskCtx context.Context, usStatus *models.QueryUsStatus) (string, error) {
	if err := godotenv.Load(".env"); err != nil {
		return "", fmt.Errorf("error loading .env file: %w", err)
	}

	client := utils.ChaoJiYing{}
	client.InitWithOptions()

	var imageBuf []byte
	if err := chromedp.Run(taskCtx,
		chromedp.Navigate("https://ceac.state.gov/CEACStatTracker/Status.aspx"),
		chromedp.WaitVisible(`#Visa_Application_Type`, chromedp.ByID),
		chromedp.SetValue(`#Visa_Application_Type`, `NIV`, chromedp.ByID),
		chromedp.WaitVisible(`#Location_Dropdown`, chromedp.ByID),
		chromedp.SetValue(`#Location_Dropdown`, `BEJ`, chromedp.ByID),
		chromedp.WaitVisible(`#Visa_Case_Number`, chromedp.ByID),
		chromedp.SetValue(`#Visa_Case_Number`, usStatus.ApplicationID, chromedp.ByID),
		chromedp.WaitVisible(`#Passport_Number`, chromedp.ByID),
		chromedp.SetValue(`#Passport_Number`, usStatus.PassportNumber, chromedp.ByID),
		chromedp.WaitVisible(`#Surname`, chromedp.ByID),
		chromedp.SetValue(`#Surname`, usStatus.First5LettersOfSurname, chromedp.ByID),
		chromedp.WaitVisible(`#c_status_ctl00_contentplaceholder1_defaultcaptcha_CaptchaImage`, chromedp.ByQuery),
		chromedp.Screenshot(`#c_status_ctl00_contentplaceholder1_defaultcaptcha_CaptchaImage`, &imageBuf, chromedp.NodeVisible),
	); err != nil {
		return "", err
	}

	if err := ioutil.WriteFile("captcha.png", imageBuf, 0644); err != nil {
		return "", err
	}
	log.Println("CAPTCHA saved as captcha.png")

	var result models.ChaoJiYing
	response := client.GetPicVal(
		os.Getenv("CJY_USERNAME"),
		os.Getenv("CJY_PASSWORD"),
		os.Getenv("CJY_SOFT_ID"),
		os.Getenv("CJY_CODE_TYPE"),
		os.Getenv("CJY_MIN_LEN"),
		"captcha.png")
	if err := json.Unmarshal(response, &result); err != nil {
		return "", fmt.Errorf("failed to get captcha value: %w", err)
	}

	if err := chromedp.Run(taskCtx,
		chromedp.WaitVisible(`#Captcha`, chromedp.ByID),
		chromedp.SetValue(`#Captcha`, result.PicStr, chromedp.ByID),
		chromedp.Click(`#ctl00_ContentPlaceHolder1_imgFolder`, chromedp.ByID),
		chromedp.WaitVisible(`#ctl00_ContentPlaceHolder1_ucApplicationStatusView_pTranslation`, chromedp.ByQuery),
		chromedp.Text(`#ctl00_ContentPlaceHolder1_ucApplicationStatusView_pTranslation`, &result.PicStr, chromedp.NodeVisible),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("Result:", result.PicStr)
			return nil
		}),
	); err != nil {
		return "", err
	}
	return result.PicStr, nil
}
