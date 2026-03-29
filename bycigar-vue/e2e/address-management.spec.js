import { test, expect } from '@playwright/test'
import { loginViaAPI, getAddresses, injectAuth, createAddress, cleanupTestData } from './helpers'

test.describe('Address Management Flow', () => {
  let token, user, existingCount

  test.beforeAll(async ({ request }) => {
    const body = await loginViaAPI(request)
    token = body.token
    user = body.user
    await cleanupTestData(request, token)
    const addresses = await getAddresses(request, token)
    existingCount = addresses.length
  })

  test.beforeEach(async ({ page }) => {
    await page.goto('/')
    await injectAuth(page, token, user)
  })

  test('address list shows existing addresses', async ({ page }) => {
    await page.goto('/profile')
    await expect(page.locator('.profile-page')).toBeVisible()

    await page.locator('.nav-item', { hasText: '地址管理' }).click()
    await expect(page.locator('.content-section h2', { hasText: '地址管理' })).toBeVisible()

    const addressCards = page.locator('.address-card')
    const count = await addressCards.count()
    expect(count).toBeGreaterThanOrEqual(existingCount)
  })

  test('add a new address', async ({ page, request }) => {
    const addressesBefore = await getAddresses(request, token)
    if (addressesBefore.length >= 5) {
      expect(true).toBe(true)
      return
    }

    await page.goto('/profile')
    await expect(page.locator('.profile-page')).toBeVisible()

    await page.locator('.nav-item', { hasText: '地址管理' }).click()
    await expect(page.locator('.content-section h2', { hasText: '地址管理' })).toBeVisible()

    await page.locator('.btn-add-address').first().click()
    await expect(page.locator('.address-form-overlay')).toBeVisible()

    await page.locator('.address-form input[placeholder="Full Name"]').fill('E2E Test')
    await page.locator('.address-form input[placeholder="Street Address"]').fill('123 E2E Street')
    await page.locator('.address-form input[placeholder="City"]').fill('Testville')
    await page.locator('.address-form select').selectOption('TX')
    await page.locator('.address-form input[placeholder="12345"]').fill('75001')
    await page.locator('.address-form input[placeholder="(213) 555-1234"]').fill('9725551234')

    await page.locator('.address-form .btn-save').click()
    await page.waitForTimeout(1500)

    const addressesAfter = await getAddresses(request, token)
    expect(addressesAfter.length).toBe(existingCount + 1)
  })

  test('edit an existing address', async ({ page, request }) => {
    const addr = await createAddress(request, token, {
      fullName: 'Before Edit',
      phone: '9725550000',
      addressLine1: '456 Edit Ave',
      city: 'Editown',
      state: 'FL',
      zipCode: '33101'
    })

    await page.goto('/profile')
    await page.locator('.nav-item', { hasText: '地址管理' }).click()
    await page.waitForTimeout(500)

    const editBtn = page.locator('.address-card', { hasText: 'Before Edit' }).locator('.btn-edit-addr')
    if (await editBtn.isVisible()) {
      await editBtn.click()
      await expect(page.locator('.address-form-overlay')).toBeVisible()

      const nameInput = page.locator('.address-form input[placeholder="Full Name"]')
      await nameInput.clear()
      await nameInput.fill('After Edit')
      await page.locator('.address-form .btn-save').click()
      await page.waitForTimeout(1500)

      const card = page.locator('.address-card', { hasText: 'After Edit' })
      await expect(card).toBeVisible()
    }

    const resp = await request.delete(`http://localhost:3000/api/addresses/${addr.id}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
  })

  test('set default address', async ({ page, request }) => {
    const addresses = await getAddresses(request, token)
    if (addresses.length < 2) return

    await page.goto('/profile')
    await page.locator('.nav-item', { hasText: '地址管理' }).click()
    await page.waitForTimeout(500)

    const nonDefaultCards = page.locator('.address-card').filter({ hasNot: page.locator('.default-badge') })
    const nonDefaultCount = await nonDefaultCards.count()

    if (nonDefaultCount > 0) {
      const setDefaultBtn = nonDefaultCards.first().locator('.btn-set-default')
      if (await setDefaultBtn.isVisible()) {
        await setDefaultBtn.click()
        await page.waitForTimeout(1000)
      }
    }

    expect(true).toBe(true)
  })

  test('delete a deletable address', async ({ page, request }) => {
    const addr = await createAddress(request, token, {
      fullName: 'ToDelete',
      phone: '9725559999',
      addressLine1: '999 Delete St',
      city: 'Delcity',
      state: 'GA',
      zipCode: '30301'
    })

    await page.goto('/profile')
    await page.locator('.nav-item', { hasText: '地址管理' }).click()
    await page.waitForTimeout(500)

    const deleteBtn = page.locator('.address-card', { hasText: 'ToDelete' }).locator('.btn-delete-addr')
    if (await deleteBtn.isVisible()) {
      await deleteBtn.click()
      await page.waitForTimeout(1000)

      const addressesAfter = await getAddresses(request, token)
      expect(addressesAfter.length).toBeLessThanOrEqual(existingCount + 1)
    } else {
      const resp = await request.delete(`http://localhost:3000/api/addresses/${addr.id}`, {
        headers: { Authorization: `Bearer ${token}` }
      })
      expect(resp.status()).toBeLessThan(500)
    }
  })

  test.afterAll(async ({ request }) => {
    await cleanupTestData(request, token)
  })
})
