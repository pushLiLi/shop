import { test, expect } from '@playwright/test'

test.describe('Search and Filter Flow', () => {
  test('search for a product keyword', async ({ page }) => {
    await page.goto('/')
    await page.locator('.search-input').first().fill('高希霸')
    await page.locator('.search-input').first().press('Enter')
    await page.waitForURL(/\/search/)
    await expect(page.locator('.search-page')).toBeVisible()
    await expect(page.locator('.keyword')).toContainText('高希霸')
    await expect(page.locator('.product-card').first()).toBeVisible()
  })

  test('search with no results', async ({ page }) => {
    await page.goto('/')
    await page.locator('.search-input').first().fill('zzznonexistentproduct')
    await page.locator('.search-input').first().press('Enter')
    await page.waitForURL(/\/search/)
    await expect(page.locator('.search-page')).toBeVisible()
    await expect(page.locator('.no-results')).toBeVisible()
    await expect(page.locator('.no-results p').first()).toContainText('未找到')
  })

  test('sort by price on search results', async ({ page }) => {
    await page.goto('/search?q=高希霸')
    await expect(page.locator('.product-card').first()).toBeVisible()

    await page.locator('.sort-btn', { hasText: '价格' }).click()
    await page.waitForTimeout(500)

    await expect(page.locator('.product-card').first()).toBeVisible()
  })

  test('category page shows products', async ({ page }) => {
    await page.goto('/category/古巴雪茄')
    await expect(page.locator('.category-page')).toBeVisible()
    await expect(page.locator('.page-title')).toBeVisible()
    await expect(page.locator('.product-card').first()).toBeVisible()
  })

  test('sort products on category page', async ({ page }) => {
    await page.goto('/category/古巴雪茄')
    await expect(page.locator('.product-card').first()).toBeVisible()

    await page.locator('.sort-btn', { hasText: '价格' }).click()
    await page.waitForTimeout(500)

    await expect(page.locator('.product-card').first()).toBeVisible()

    await page.locator('.sort-btn', { hasText: '名称' }).click()
    await page.waitForTimeout(500)

    await expect(page.locator('.product-card').first()).toBeVisible()
  })

  test('pagination on category page', async ({ page }) => {
    await page.goto('/')
    const viewMoreLinks = page.locator('.view-more')
    if (await viewMoreLinks.first().isVisible()) {
      await viewMoreLinks.first().click()
      await page.waitForURL(/\/category\//)
      await expect(page.locator('.category-page')).toBeVisible()

      const nextBtn = page.locator('.page-btn', { hasText: '下一页' })
      if (await nextBtn.isVisible() && await nextBtn.isEnabled()) {
        await nextBtn.click()
        await page.waitForTimeout(500)
        await expect(page.locator('.product-card').first()).toBeVisible()
      }
    }
  })

  test('combined search and sort', async ({ page }) => {
    await page.goto('/search?q=雪茄')
    await expect(page.locator('.product-card').first()).toBeVisible()

    await page.locator('.sort-btn', { hasText: '价格' }).click()
    await page.waitForTimeout(500)

    await expect(page.locator('.search-info')).toBeVisible()
    await expect(page.locator('.product-card').first()).toBeVisible()
  })
})
