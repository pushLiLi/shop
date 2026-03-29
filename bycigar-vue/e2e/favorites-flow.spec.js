import { test, expect } from '@playwright/test'
import { loginViaAPI, clearFavorites, injectAuth } from './helpers'

test.describe('Favorites Flow', () => {
  let token, user

  test.beforeAll(async ({ request }) => {
    const body = await loginViaAPI(request)
    token = body.token
    user = body.user
  })

  test.beforeEach(async ({ request, page }) => {
    await clearFavorites(request, token)
    await page.goto('/')
    await injectAuth(page, token, user)
  })

  test('add favorite from product detail and view in favorites page', async ({ page }) => {
    await page.goto('/')
    await page.locator('.products-section .product-card .product-name a').first().click()
    await page.waitForURL(/\/products\/\d+/)

    await page.locator('.purchase-section .favorite-btn').click()
    await page.waitForTimeout(500)

    await page.goto('/favorites')
    await expect(page.locator('.favorites-page')).toBeVisible()
    await expect(page.locator('.favorite-item').first()).toBeVisible()
  })

  test('remove favorite from favorites page', async ({ page, request }) => {
    const productsResp = await request.get('http://localhost:3000/api/products?limit=1')
    const { products } = await productsResp.json()
    const productId = products[0].id

    await request.post('http://localhost:3000/api/favorites', {
      headers: { Authorization: `Bearer ${token}` },
      data: { productId }
    })

    await page.goto('/favorites')
    await expect(page.locator('.favorite-item').first()).toBeVisible()

    await page.locator('.favorite-item .remove-btn').first().click()
    await page.waitForTimeout(500)

    const count = await page.locator('.favorite-item').count()
    expect(count).toBe(0)
  })

  test('empty favorites shows empty state', async ({ page }) => {
    await page.goto('/favorites')
    await expect(page.locator('.favorites-page')).toBeVisible()
    await expect(page.locator('.empty-favorites')).toBeVisible()
    await expect(page.locator('.empty-favorites p')).toContainText('暂无收藏')
  })

  test('batch add to cart from favorites', async ({ page, request }) => {
    const productsResp = await request.get('http://localhost:3000/api/products?limit=2')
    const { products } = await productsResp.json()

    for (const p of products) {
      await request.post('http://localhost:3000/api/favorites', {
        headers: { Authorization: `Bearer ${token}` },
        data: { productId: p.id }
      })
    }

    await page.goto('/favorites')
    await expect(page.locator('.favorite-item').first()).toBeVisible()

    await page.locator('.select-all input[type="checkbox"]').click()

    await page.locator('.batch-cart-btn').click()
    await page.waitForTimeout(500)

    await expect(page.locator('.toast')).toBeVisible()
  })
})
