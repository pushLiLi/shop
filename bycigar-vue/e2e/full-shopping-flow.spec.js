import { test, expect } from '@playwright/test'
import { clearCart, deleteAllAddresses, createAddress, loginViaAPI, injectAuth } from './helpers'

const API_BASE = 'http://localhost:3000/api'

test.describe('Full Shopping Flow', () => {
  let token, user

  test.beforeAll(async ({ request }) => {
    const body = await loginViaAPI(request)
    token = body.token
    user = body.user
  })

  test.beforeEach(async ({ request, page }) => {
    await clearCart(request, token)
    await deleteAllAddresses(request, token)
    await page.goto('/')
    await injectAuth(page, token, user)
  })

  test('complete shopping journey', async ({ page, request }) => {
    await page.goto('/')
    await expect(page.locator('.home-page')).toBeVisible()
    await expect(page.locator('.slide').first()).toBeVisible()
    await expect(page.locator('.products-section .product-card').first()).toBeVisible()

    await page.locator('.search-input').first().fill('高希霸')
    await page.locator('.search-input').first().press('Enter')
    await page.waitForURL(/\/search/)
    await expect(page.locator('.search-page')).toBeVisible()
    await expect(page.locator('.keyword')).toContainText('高希霸')

    await expect(page.locator('.product-card').first()).toBeVisible()
    await page.locator('.product-card .product-name a').first().click()
    await page.waitForURL(/\/products\/\d+/)
    await expect(page.locator('.product-detail-page')).toBeVisible()
    await expect(page.locator('.product-title')).toBeVisible()

    const priceText = await page.locator('.product-price-main').textContent()
    expect(priceText).toMatch(/\$/)

    await page.locator('.qty-btn').last().click()
    await page.locator('.buy-btn').click()
    await expect(page.locator('.toast')).toBeVisible()

    const badge = page.locator('.header-right .icon-badge').last()
    await expect(badge).toBeVisible()
    const badgeText = await badge.textContent()
    expect(parseInt(badgeText)).toBeGreaterThanOrEqual(1)

    await page.locator('.header-right .icon-item').last().click()
    await expect(page.locator('.cart-drawer')).toBeVisible()
    await expect(page.locator('.cart-item').first()).toBeVisible()

    const subtotal = await page.locator('.item-total').first().textContent()
    expect(subtotal).toBeTruthy()

    await page.locator('.checkout-btn').click()
    await page.waitForURL('/checkout')
    await expect(page.locator('.checkout-page')).toBeVisible()

    await createAddress(request, token)

    await page.reload()
    await expect(page.locator('.address-option').first()).toBeVisible()

    await page.locator('.address-option').first().click()
    await page.locator('.submit-btn').click()
    await page.waitForURL(/\/orders/, { timeout: 15000 })
    await expect(page.locator('.orders-page')).toBeVisible()
    await expect(page.locator('.order-card').first()).toBeVisible()

    const orderTotal = await page.locator('.order-total').first().textContent()
    expect(orderTotal).toMatch(/总计/)
  })

  test('browse homepage and view product detail', async ({ page }) => {
    await page.goto('/')
    await expect(page.locator('.home-page')).toBeVisible()
    await expect(page.locator('.slide').first()).toBeVisible()

    await page.locator('.slider-btn.next').click()
    await page.waitForTimeout(500)

    await page.locator('.products-section .product-card .product-name a').first().click()
    await page.waitForURL(/\/products\/\d+/)
    await expect(page.locator('.product-title')).toBeVisible()
    await expect(page.locator('.product-price-main')).toBeVisible()
  })

  test('cart quantity update works', async ({ page, request }) => {
    await page.goto('/')

    await page.locator('.products-section .product-card .add-cart-btn').first().click()
    await expect(page.locator('.toast')).toBeVisible()

    await page.locator('.header-right .icon-item').last().click()
    await expect(page.locator('.cart-drawer')).toBeVisible()

    const plusBtn = page.locator('.quantity-control button').last()
    await plusBtn.click()
    await page.waitForTimeout(500)

    const qtyInput = page.locator('.quantity-control input').first()
    const qty = await qtyInput.inputValue()
    expect(parseInt(qty)).toBeGreaterThanOrEqual(2)
  })

  test('empty cart shows empty state', async ({ page }) => {
    await page.goto('/')

    await page.locator('.header-right .icon-item').last().click()
    await expect(page.locator('.cart-drawer')).toBeVisible()
    await expect(page.locator('.empty-cart')).toBeVisible()
    await expect(page.locator('.empty-cart p')).toContainText('购物车是空的')
  })
})
