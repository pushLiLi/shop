<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const API_BASE = '/api'

const props = defineProps({
  activeSlug: { type: String, default: '' }
})

const categories = ref([])
const loading = ref(true)
const expandedCategories = ref(new Set())
const isMobile = ref(window.innerWidth <= 768)
const drawerOpen = ref(false)

function onResize() {
  isMobile.value = window.innerWidth <= 768
}

onMounted(async () => {
  window.addEventListener('resize', onResize)
  try {
    const res = await fetch(`${API_BASE}/categories`)
    const data = await res.json()
    categories.value = Array.isArray(data) ? data : (data.data || [])
  } catch (e) {
    console.error('获取分类失败:', e)
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', onResize)
})

const activeParentId = computed(() => {
  if (!props.activeSlug || !categories.value.length) return null
  for (const cat of categories.value) {
    if (cat.slug === props.activeSlug) return cat.id
    if (cat.children) {
      for (const child of cat.children) {
        if (child.slug === props.activeSlug) return cat.id
      }
    }
  }
  return null
})

const currentCategoryName = computed(() => {
  if (!props.activeSlug) return '全部分类'
  for (const cat of categories.value) {
    if (cat.slug === props.activeSlug) return cat.name
    if (cat.children) {
      for (const child of cat.children) {
        if (child.slug === props.activeSlug) return child.name
      }
    }
  }
  return '全部分类'
})

function openDrawer() {
  drawerOpen.value = true
  document.body.style.overflow = 'hidden'
}

function closeDrawer() {
  drawerOpen.value = false
  document.body.style.overflow = ''
}

function toggleCategory(cat) {
  if (!cat.children || !cat.children.length) return
  if (expandedCategories.value.has(cat.id)) {
    expandedCategories.value.delete(cat.id)
  } else {
    expandedCategories.value.add(cat.id)
  }
}

function handleCategoryClick(cat, event) {
  if (isMobile.value && cat.children && cat.children.length > 0) {
    event.preventDefault()
    toggleCategory(cat)
  } else if (isMobile.value) {
    closeDrawer()
  }
}

function handleSubClick() {
  if (isMobile.value) {
    closeDrawer()
  }
}

function isExpanded(catId) {
  return expandedCategories.value.has(catId)
}
</script>

<template>
  <aside class="category-sidebar">
    <button v-if="isMobile && !loading" class="mobile-category-btn" @click="openDrawer">
      <span class="mobile-btn-text">{{ currentCategoryName }}</span>
      <svg class="mobile-btn-icon" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="6 9 12 15 18 9"></polyline>
      </svg>
    </button>

    <Teleport to="body">
      <Transition name="drawer">
        <div v-if="isMobile && drawerOpen" class="drawer-overlay" @click="closeDrawer">
          <div class="drawer-content" @click.stop>
            <div class="drawer-header">
              <span class="drawer-title">选择分类</span>
              <button class="drawer-close" @click="closeDrawer">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>
            <ul class="drawer-list">
              <li v-for="cat in categories" :key="cat.id" class="drawer-item">
                <router-link
                  :to="'/category/' + cat.slug"
                  class="drawer-link"
                  :class="{ active: activeSlug === cat.slug }"
                  @click="handleCategoryClick(cat, $event)"
                >
                  <span>{{ cat.name }}</span>
                  <svg
                    v-if="cat.children && cat.children.length"
                    class="drawer-expand-icon"
                    :class="{ expanded: isExpanded(cat.id) }"
                    xmlns="http://www.w3.org/2000/svg"
                    width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                  >
                    <polyline points="6 9 12 15 18 9"></polyline>
                  </svg>
                </router-link>
                <Transition name="expand">
                  <ul
                    v-if="cat.children && cat.children.length && isExpanded(cat.id)"
                    class="drawer-sublist"
                  >
                    <li v-for="child in cat.children" :key="child.id">
                      <router-link
                        :to="'/category/' + child.slug"
                        class="drawer-link sub"
                        :class="{ active: activeSlug === child.slug }"
                        @click="handleSubClick"
                      >
                        {{ child.name }}
                      </router-link>
                    </li>
                  </ul>
                </Transition>
              </li>
            </ul>
          </div>
        </div>
      </Transition>
    </Teleport>

    <template v-if="!isMobile">
      <div class="sidebar-title">商品分类</div>
      <ul class="category-list" v-if="!loading">
        <li
          v-for="cat in categories"
          :key="cat.id"
          class="category-item"
          :class="{ 'has-active-child': activeParentId === cat.id }"
        >
          <router-link
            :to="'/category/' + cat.slug"
            class="category-link"
            :class="{ active: activeSlug === cat.slug }"
          >
            <span>{{ cat.name }}</span>
            <svg
              v-if="cat.children && cat.children.length"
              class="category-chevron"
              xmlns="http://www.w3.org/2000/svg"
              width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
            >
              <polyline points="6 9 12 15 18 9"></polyline>
            </svg>
          </router-link>
          <ul v-if="cat.children && cat.children.length" class="subcategory-list">
            <li v-for="child in cat.children" :key="child.id">
              <router-link :to="'/category/' + child.slug" class="category-link sub" :class="{ active: activeSlug === child.slug }">
                {{ child.name }}
              </router-link>
            </li>
          </ul>
        </li>
      </ul>
    </template>
  </aside>
</template>

<style scoped>
.category-sidebar {
  background: #1a1a1a;
  border-radius: 8px;
  padding: 20px;
  min-width: 220px;
  max-width: 240px;
  position: sticky;
  top: 100px;
  align-self: flex-start;
  max-height: calc(100vh - 120px);
  overflow-y: auto;
}

.category-sidebar::-webkit-scrollbar {
  width: 4px;
}

.category-sidebar::-webkit-scrollbar-track {
  background: transparent;
}

.category-sidebar::-webkit-scrollbar-thumb {
  background: #333;
  border-radius: 2px;
}

.category-sidebar::-webkit-scrollbar-thumb:hover {
  background: #555;
}

.sidebar-title {
  color: #d4a574;
  font-size: 16px;
  font-weight: 600;
  border-bottom: 1px solid #333;
  padding-bottom: 12px;
  margin-bottom: 12px;
}

.category-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.category-item {
  margin-bottom: 2px;
}

.category-link {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  color: #ccc;
  text-decoration: none;
  border-radius: 4px;
  transition: all 0.2s;
  font-size: 14px;
}

.category-link:hover {
  color: #d4a574;
  background: rgba(212, 165, 116, 0.08);
}

.category-link.active {
  color: #d4a574;
  background: rgba(212, 165, 116, 0.15);
  font-weight: 500;
}

.category-chevron {
  color: #666;
  transition: transform 0.3s ease;
  flex-shrink: 0;
  margin-left: 4px;
}

.category-item:hover .category-chevron,
.category-item.has-active-child .category-chevron {
  transform: rotate(180deg);
  color: #d4a574;
}

.subcategory-list {
  list-style: none;
  padding: 0;
  margin: 0;
  max-height: 0;
  opacity: 0;
  overflow: hidden;
  transition: max-height 0.3s ease, opacity 0.25s ease;
}

.category-item:hover > .subcategory-list,
.category-item.has-active-child > .subcategory-list {
  max-height: 500px;
  opacity: 1;
}

.category-link.sub {
  padding-left: 28px;
  font-size: 13px;
  color: #999;
}

.category-link.sub:hover {
  color: #d4a574;
}

.category-link.sub.active {
  color: #d4a574;
}

.mobile-category-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
  color: #d4a574;
  font-size: 15px;
  cursor: pointer;
  transition: all 0.2s;
}

.mobile-category-btn:active {
  background: #252525;
}

.mobile-btn-text {
  font-weight: 500;
}

.mobile-btn-icon {
  transition: transform 0.2s;
}

.drawer-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 1000;
  display: flex;
  align-items: flex-end;
}

.drawer-content {
  width: 100%;
  max-height: 70vh;
  background: #1a1a1a;
  border-radius: 16px 16px 0 0;
  overflow-y: auto;
  padding-bottom: env(safe-area-inset-bottom, 20px);
}

.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #333;
  position: sticky;
  top: 0;
  background: #1a1a1a;
}

.drawer-title {
  color: #d4a574;
  font-size: 16px;
  font-weight: 600;
}

.drawer-close {
  background: transparent;
  border: none;
  color: #888;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.2s;
}

.drawer-close:hover {
  color: #fff;
}

.drawer-list {
  list-style: none;
  padding: 8px 0;
  margin: 0;
}

.drawer-item {
  border-bottom: 1px solid #2a2a2a;
}

.drawer-item:last-child {
  border-bottom: none;
}

.drawer-link {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  color: #ccc;
  text-decoration: none;
  font-size: 15px;
  transition: all 0.2s;
  position: relative;
}

.drawer-link:active {
  background: rgba(255, 255, 255, 0.05);
}

.drawer-link.active {
  color: #d4a574;
  font-weight: 500;
}

.drawer-link.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 20px;
  background: #d4a574;
  border-radius: 0 2px 2px 0;
}

.drawer-expand-icon {
  color: #666;
  transition: transform 0.2s;
}

.drawer-expand-icon.expanded {
  transform: rotate(180deg);
}

.drawer-sublist {
  list-style: none;
  padding: 0;
  margin: 0;
  background: rgba(0, 0, 0, 0.2);
}

.drawer-link.sub {
  padding: 14px 20px 14px 36px;
  font-size: 14px;
  color: #999;
}

.drawer-link.sub:active {
  background: rgba(255, 255, 255, 0.03);
}

.drawer-link.sub.active {
  color: #d4a574;
}

.drawer-link.sub.active::before {
  left: 16px;
  height: 16px;
}

.expand-enter-active,
.expand-leave-active {
  transition: all 0.25s ease;
  overflow: hidden;
}

.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  max-height: 0;
}

.expand-enter-to,
.expand-leave-from {
  opacity: 1;
  max-height: 500px;
}

.drawer-enter-active,
.drawer-leave-active {
  transition: opacity 0.3s ease;
}

.drawer-enter-active .drawer-content,
.drawer-leave-active .drawer-content {
  transition: transform 0.3s cubic-bezier(0.32, 0.72, 0, 1);
}

.drawer-enter-from,
.drawer-leave-to {
  opacity: 1;
}

.drawer-enter-from .drawer-content,
.drawer-leave-to .drawer-content {
  transform: translateY(100%);
}

@media (max-width: 768px) {
  .category-sidebar {
    background: transparent;
    padding: 0;
    min-width: auto;
    max-width: none;
    position: static;
    max-height: none;
    overflow: visible;
  }
}
</style>
