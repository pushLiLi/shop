package handlers

func ResetLoginFailures() {
	loginFailures = make(map[string]*loginAttempt)
}

func GetCaptchaCode(captchaID string) string {
	return captchaInstance.Store.Get(captchaID, false)
}
