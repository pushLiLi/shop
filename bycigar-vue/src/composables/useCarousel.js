import { ref, computed, onMounted, onUnmounted } from 'vue'

export function useCarousel(options = {}) {
  const {
    autoplay = true,
    interval = 4000,
    pauseOnHover = true
  } = options

  const currentIndex = ref(0)
  const slides = ref([])
  const isPaused = ref(false)
  const isTransitioning = ref(false)
  let timer = null
  let touchStartX = 0
  let touchEndX = 0

  const totalSlides = computed(() => slides.value.length)

  const next = () => {
    if (isTransitioning.value || totalSlides.value === 0) return
    isTransitioning.value = true
    currentIndex.value = (currentIndex.value + 1) % totalSlides.value
    setTimeout(() => { isTransitioning.value = false }, 500)
  }

  const prev = () => {
    if (isTransitioning.value || totalSlides.value === 0) return
    isTransitioning.value = true
    currentIndex.value = (currentIndex.value - 1 + totalSlides.value) % totalSlides.value
    setTimeout(() => { isTransitioning.value = false }, 500)
  }

  const goTo = (index) => {
    if (isTransitioning.value || index === currentIndex.value) return
    isTransitioning.value = true
    currentIndex.value = index
    setTimeout(() => { isTransitioning.value = false }, 500)
  }

  const startAutoplay = () => {
    if (!autoplay) return
    stopAutoplay()
    timer = setInterval(() => {
      if (!isPaused.value) next()
    }, interval)
  }

  const stopAutoplay = () => {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
  }

  const pause = () => { isPaused.value = true }
  const resume = () => { isPaused.value = false }

  const onMouseEnter = () => { if (pauseOnHover) pause() }
  const onMouseLeave = () => { if (pauseOnHover) resume() }

  const onTouchStart = (e) => { touchStartX = e.changedTouches[0].screenX }
  const onTouchEnd = (e) => {
    touchEndX = e.changedTouches[0].screenX
    const diff = touchStartX - touchEndX
    if (Math.abs(diff) > 50) {
      diff > 0 ? next() : prev()
    }
  }

  onMounted(() => { startAutoplay() })
  onUnmounted(() => { stopAutoplay() })

  return {
    currentIndex,
    slides,
    isPaused,
    isTransitioning,
    totalSlides,
    next,
    prev,
    goTo,
    pause,
    resume,
    onMouseEnter,
    onMouseLeave,
    onTouchStart,
    onTouchEnd
  }
}
