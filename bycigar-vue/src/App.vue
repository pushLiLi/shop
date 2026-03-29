<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import TheHeader from './components/TheHeader.vue'
import TheFooter from './components/TheFooter.vue'
import Toast from './components/Toast.vue'
import CartDrawer from './components/CartDrawer.vue'
import ChatWidget from './components/ChatWidget.vue'

const route = useRoute()
const isAdminRoute = computed(() => route.path.startsWith('/admin'))
const showScrollTop = ref(false)

const handleScroll = () => {
  showScrollTop.value = window.scrollY > 300
}

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => window.addEventListener('scroll', handleScroll))
onUnmounted(() => window.removeEventListener('scroll', handleScroll))
</script>

<template>
  <div id="app">
    <TheHeader v-if="!isAdminRoute" />
    <main class="main-content">
      <RouterView v-slot="{ Component }">
        <Transition name="page" mode="out-in">
          <component :is="Component" />
        </Transition>
      </RouterView>
    </main>
    <TheFooter v-if="!isAdminRoute" />
    <Toast />
    <CartDrawer />
    <ChatWidget v-if="!isAdminRoute" />
    <Transition name="scroll-top">
      <button
        v-if="showScrollTop && !isAdminRoute"
        class="scroll-top-btn"
        @click="scrollToTop"
        aria-label="回到顶部"
      >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="18 15 12 9 6 15" />
        </svg>
      </button>
    </Transition>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
  background: #0f0f0f;
  color: #fff;
  line-height: 1.6;
}

#app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.main-content {
  flex: 1;
}

.page-enter-active,
.page-leave-active {
  transition: opacity 0.2s ease;
}

.page-enter-from,
.page-leave-to {
  opacity: 0;
}

.scroll-top-btn {
  position: fixed;
  bottom: 30px;
  right: 90px;
  width: 44px;
  height: 44px;
  border-radius: 50%;
  border: none;
  background: #d4a574;
  color: #0f0f0f;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1100;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.4);
  transition: background 0.2s, transform 0.2s;
}

.scroll-top-btn:hover {
  background: #e0b88a;
  transform: translateY(-2px);
}

.scroll-top-enter-active,
.scroll-top-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.scroll-top-enter-from,
.scroll-top-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>
