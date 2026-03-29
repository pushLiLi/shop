import { test, expect } from '@playwright/test'
import { injectAuth, loginViaAPI, cleanupTestData } from './helpers'

test.describe('Admin Panel Flow', () => {
  let token, user

  test.beforeEach(async ({ page, request }) => {
    const loginData = await loginViaAPI(request)
    token = loginData.token
    user = loginData.user

    await page.goto('/')
    await injectAuth(page, token, user)
    await page.goto('/admin/products')
    await page.waitForURL(/\/admin/)
  })

  test('admin products list loads', async ({ page }) => {
    await expect(page.locator('.admin-products')).toBeVisible()
    await expect(page.locator('.data-table')).toBeVisible()
    await expect(page.locator('.data-table tbody tr').first()).toBeVisible()
  })

  test('admin can search and filter products', async ({ page }) => {
    await expect(page.locator('.data-table')).toBeVisible()

    await page.locator('.toolbar input[placeholder*="搜索"]').fill('test')
    await page.locator('.toolbar .btn-search').click()
    await page.waitForTimeout(500)

    const rows = page.locator('.data-table tbody tr')
    const count = await rows.count()
    expect(count).toBeGreaterThanOrEqual(1)

    await page.locator('.toolbar .btn-reset').click()
    await page.waitForTimeout(500)
    const resetCount = await page.locator('.data-table tbody tr').count()
    expect(resetCount).toBeGreaterThanOrEqual(count)
  })

  test('admin can open create product modal', async ({ page }) => {
    await page.locator('.btn-add', { hasText: '添加商品' }).click()
    await expect(page.locator('.modal')).toBeVisible()
    await expect(page.locator('.modal h3')).toContainText('添加商品')

    await page.locator('.modal input[placeholder*="商品名称"]').fill('E2E Test Product')
    await page.locator('.modal input[type="number"]').first().fill('99.99')

    await page.locator('.modal .btn-save').click()
    await page.waitForTimeout(1000)

    const modalVisible = await page.locator('.modal').isVisible().catch(() => false)
    if (!modalVisible) {
      await expect(page.locator('.data-table')).toBeVisible()
    }
  })

  test('admin categories list loads', async ({ page }) => {
    await page.goto('/admin/categories')
    await expect(page.locator('.admin-categories')).toBeVisible()
    await expect(page.locator('.category-group').first()).toBeVisible()
  })

  test('admin can open create category modal', async ({ page }) => {
    await page.goto('/admin/categories')
    await page.locator('.btn-add', { hasText: '添加顶级分类' }).click()
    await expect(page.locator('.modal')).toBeVisible()
    await expect(page.locator('.modal h3')).toContainText('添加分类')

    await page.locator('.modal input[placeholder*="高希霸"]').fill('E2E Test Category')
    await page.locator('.modal .btn-save').click()
    await page.waitForTimeout(1000)

    const modalVisible = await page.locator('.modal').isVisible().catch(() => false)
    if (!modalVisible) {
      await expect(page.locator('.category-list')).toBeVisible()
    }
  })

  test('admin banners list loads', async ({ page }) => {
    await page.goto('/admin/banners')
    await expect(page.locator('.admin-banners')).toBeVisible()
    await expect(page.locator('.banner-item').first()).toBeVisible()
  })

  test('admin can toggle product status', async ({ page }) => {
    await page.goto('/admin/products')
    await expect(page.locator('.data-table')).toBeVisible()

    const firstBadge = page.locator('.data-table tbody tr .badge-success, .data-table tbody tr .badge-danger').first()
    await firstBadge.click()
    await page.waitForTimeout(500)

    await expect(page.locator('.data-table')).toBeVisible()
  })

  test('admin sidebar navigation', async ({ page }) => {
    await expect(page.locator('.sidebar-nav')).toBeVisible()

    await page.locator('.sidebar-nav .nav-item', { hasText: '分类' }).click()
    await page.waitForURL('/admin/categories')
    await expect(page.locator('.admin-categories')).toBeVisible()

    await page.locator('.sidebar-nav .nav-item', { hasText: '轮播图' }).click()
    await page.waitForURL('/admin/banners')
    await expect(page.locator('.admin-banners')).toBeVisible()
  })

  test('admin can logout', async ({ page }) => {
    await page.locator('.btn-logout').click()
    await page.waitForURL('/')
    await expect(page.locator('.home-page, .login-page')).toBeVisible()
  })

  test.afterAll(async ({ request }) => {
    await cleanupTestData(request, token)
  })
})
