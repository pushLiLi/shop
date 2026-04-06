export function formatPrice(product) {
  if (!product) return '¥0.00'
  if (product.currency === 'USD' && product.priceUsd > 0) {
    return '$' + Number(product.priceUsd).toFixed(2)
  }
  return '¥' + Number(product.price).toFixed(2)
}

export function formatPriceByCurrency(price, currency) {
  if (currency === 'USD') {
    return '$' + Number(price).toFixed(2)
  }
  return '¥' + Number(price).toFixed(2)
}

export function getCurrencySymbol(currency) {
  return currency === 'USD' ? '$' : '¥'
}
