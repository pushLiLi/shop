import { useToastStore } from '../stores/toast'

export function useShare() {
  const toast = useToastStore()

  async function copyToClipboard(text) {
    try {
      await navigator.clipboard.writeText(text)
      return true
    } catch {
      const textarea = document.createElement('textarea')
      textarea.value = text
      textarea.style.position = 'fixed'
      textarea.style.opacity = '0'
      document.body.appendChild(textarea)
      textarea.select()
      const success = document.execCommand('copy')
      document.body.removeChild(textarea)
      return success
    }
  }

  async function share(options) {
    const { title, text, url } = options

    if (navigator.share) {
      try {
        await navigator.share({ title, text, url })
        return true
      } catch (e) {
        if (e.name === 'AbortError') return false
      }
    }

    const shareUrl = url || window.location.href
    const copied = await copyToClipboard(shareUrl)
    if (copied) {
      toast.success('链接已复制到剪贴板')
    } else {
      toast.error('复制失败，请手动复制')
    }
    return copied
  }

  async function shareProduct(product) {
    const url = `${window.location.origin}/products/${product.id}`
    return share({
      title: product.name,
      text: product.description ? product.description.substring(0, 100) : `${product.name} - ¥${product.price}`,
      url
    })
  }

  return { share, shareProduct, copyToClipboard }
}
