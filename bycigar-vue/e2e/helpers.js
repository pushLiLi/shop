const API_BASE = 'http://localhost:3000/api'

export async function loginViaAPI(request, email = 'admin@admin.com', password = '123456') {
  const resp = await request.post(`${API_BASE}/auth/login`, {
    data: { email, password }
  })
  const body = await resp.json()
  return body
}

export async function loginViaUI(page, email = 'admin@admin.com', password = '123456') {
  await page.goto('/login')
  await page.locator('input[type="email"]').fill(email)
  await page.locator('input[type="password"]').fill(password)
  await page.locator('form .submit-btn').click()
  await page.waitForURL(/\/(admin|$)/, { timeout: 10000 })
}

export async function clearCart(request, token) {
  if (!token) return
  const resp = await request.get(`${API_BASE}/cart`, {
    headers: { Authorization: `Bearer ${token}` }
  })
  const data = await resp.json()
  const items = Array.isArray(data?.items) ? data.items : (Array.isArray(data) ? data : [])
  for (const item of items) {
    await request.delete(`${API_BASE}/cart/${item.id}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
  }
}

export async function clearFavorites(request, token) {
  if (!token) return
  const resp = await request.get(`${API_BASE}/favorites`, {
    headers: { Authorization: `Bearer ${token}` }
  })
  const data = await resp.json()
  const items = Array.isArray(data?.items) ? data.items : (Array.isArray(data) ? data : [])
  for (const item of items) {
    await request.delete(`${API_BASE}/favorites/${item.productId}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
  }
}

export async function getAddresses(request, token) {
  if (!token) return []
  const resp = await request.get(`${API_BASE}/addresses`, {
    headers: { Authorization: `Bearer ${token}` }
  })
  const data = await resp.json()
  return Array.isArray(data?.addresses) ? data.addresses : (Array.isArray(data) ? data : [])
}

export async function deleteAllAddresses(request, token) {
  const addresses = await getAddresses(request, token)
  for (const addr of addresses) {
    await request.delete(`${API_BASE}/addresses/${addr.id}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
  }
}

export async function createAddress(request, token, overrides = {}) {
  const data = {
    fullName: 'Test User',
    phone: '13800138000',
    addressLine1: '123 Main St',
    addressLine2: '',
    city: 'New York',
    state: 'NY',
    zipCode: '10001',
    isDefault: true,
    ...overrides
  }
  const resp = await request.post(`${API_BASE}/addresses`, {
    headers: { Authorization: `Bearer ${token}` },
    data
  })
  return resp.json()
}

export async function injectAuth(page, token, user) {
  await page.evaluate(({ token, user }) => {
    localStorage.setItem('token', token)
    localStorage.setItem('user', JSON.stringify(user))
  }, { token, user })
}

export async function cleanupTestData(request, token) {
  const headers = { Authorization: `Bearer ${token}` }

  const productsResp = await request.get(`${API_BASE}/admin/products?search=E2E Test`, { headers })
  const productsData = await productsResp.json()
  const products = Array.isArray(productsData?.products) ? productsData.products : []
  for (const p of products) {
    if (p.name && p.name.startsWith('E2E Test')) {
      await request.delete(`${API_BASE}/admin/products/${p.id}`, { headers })
    }
  }

  const catsResp = await request.get(`${API_BASE}/admin/categories`, { headers })
  const catsData = await catsResp.json()
  const categories = Array.isArray(catsData) ? catsData : []
  for (const c of categories) {
    if (c.name && c.name.startsWith('E2E Test')) {
      await request.delete(`${API_BASE}/admin/categories/${c.id}`, { headers })
    }
  }

  const addresses = await getAddresses(request, token)
  for (const addr of addresses) {
    if (addr.fullName && (addr.fullName === 'E2E Test' || addr.fullName === 'ToDelete' || addr.fullName === 'After Edit' || addr.fullName === 'Before Edit' || addr.fullName === 'Test User')) {
      await request.delete(`${API_BASE}/addresses/${addr.id}`, {
        headers: { Authorization: `Bearer ${token}` }
      }).catch(() => {})
    }
  }
}
