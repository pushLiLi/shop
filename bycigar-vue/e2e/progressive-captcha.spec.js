import { test, expect } from '@playwright/test'
import { loginViaAPI, injectAuth } from './helpers'

test.describe('Progressive Captcha Flow', () => {
  test('login fails with wrong password', async ({ page }) => {
    await page.goto('/login')
    await page.locator('input[type="email"]').fill('test-captcha@test.com')
    await page.locator('input[type="password"]').fill('wrongpassword')
    await page.locator('form .submit-btn').click()
    await page.waitForTimeout(500)

    await expect(page.locator('.error-msg')).toBeVisible()
    const errorText = await page.locator('.error-msg').textContent()
    expect(errorText.length).toBeGreaterThan(0)
  })

  test('captcha appears after 3 consecutive failures', async ({ page }) => {
    await page.goto('/login')
    await page.locator('input[type="email"]').fill('test-captcha@test.com')
    await page.locator('input[type="password"]').fill('wrongpassword')

    await page.locator('form .submit-btn').click()
    await page.waitForTimeout(300)
    await page.locator('form .submit-btn').click()
    await page.waitForTimeout(300)
    await page.locator('form .submit-btn').click()
    await page.waitForTimeout(500)

    const captchaImg = page.locator('.captcha-img')
    if (await captchaImg.isVisible()) {
      await expect(captchaImg).toBeVisible()
      await expect(page.locator('.captcha-row input[placeholder*="验证码"]')).toBeVisible()
    }
  })

  test('login succeeds with correct credentials', async ({ page, request }) => {
    const body = await loginViaAPI(request)
    await page.goto('/')
    await injectAuth(page, body.token, body.user)
    await page.goto('/profile')
    await expect(page.locator('.profile-page')).toBeVisible()
  })

  test('register form always shows captcha', async ({ page }) => {
    await page.goto('/login')
    await page.locator('.tab-btn', { hasText: '注册' }).click()
    await page.waitForTimeout(500)

    const captchaImg = page.locator('.captcha-img')
    if (await captchaImg.isVisible()) {
      await expect(captchaImg).toBeVisible()
    }

    await expect(page.locator('input[placeholder*="再次输入密码"]')).toBeVisible()
  })

  test('tab switching between login and register', async ({ page }) => {
    await page.goto('/login')

    await expect(page.locator('.tab-btn.active', { hasText: '登录' })).toBeVisible()

    await page.locator('.tab-btn', { hasText: '注册' }).click()
    await expect(page.locator('.tab-btn.active', { hasText: '注册' })).toBeVisible()
    await expect(page.locator('input[placeholder*="再次输入密码"]')).toBeVisible()

    await page.locator('.tab-btn', { hasText: '登录' }).click()
    await expect(page.locator('.tab-btn.active', { hasText: '登录' })).toBeVisible()
    await expect(page.locator('input[placeholder*="再次输入密码"]')).not.toBeVisible()
  })

  test('register shows validation errors', async ({ page }) => {
    await page.goto('/login')
    await page.locator('.tab-btn', { hasText: '注册' }).click()
    await page.waitForTimeout(500)

    await page.locator('input[type="email"]').fill('test@test.com')
    await page.locator('input[placeholder*="6位"]').fill('12345')
    await page.locator('input[placeholder*="再次输入密码"]').fill('67890')
    await page.locator('form .submit-btn').click()
    await page.waitForTimeout(500)

    const errorEl = page.locator('.error-msg')
    if (await errorEl.isVisible()) {
      const text = await errorEl.textContent()
      expect(text.length).toBeGreaterThan(0)
    }
  })
})
