<script setup>
import { ref, onMounted } from 'vue'

const API_BASE = 'http://localhost:3000/api'

const props = defineProps({
  activeSlug: { type: String, default: '' }
})

const categories = ref([])
const loading = ref(true)

onMounted(async () => {
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
</script>

<template>
  <aside class="category-sidebar">
    <div class="sidebar-title">商品分类</div>
    <ul class="category-list" v-if="!loading">
      <li class="category-item">
        <router-link to="/products" class="category-link" :class="{ active: !activeSlug }">
          全部商品
        </router-link>
      </li>
      <li v-for="cat in categories" :key="cat.id" class="category-item">
        <router-link :to="'/category/' + cat.slug" class="category-link" :class="{ active: activeSlug === cat.slug }">
          {{ cat.name }}
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
  display: block;
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

.subcategory-list {
  list-style: none;
  padding: 0;
  margin: 0;
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

@media (max-width: 768px) {
  .category-sidebar {
    background: transparent;
    padding: 0;
    min-width: auto;
    max-width: none;
    position: static;
  }

  .sidebar-title {
    display: none;
  }

  .category-list {
    display: flex;
    gap: 8px;
    overflow-x: auto;
    padding-bottom: 10px;
    flex-wrap: nowrap;
  }

  .category-item {
    margin-bottom: 0;
  }

  .subcategory-list {
    display: none;
  }

  .category-link {
    white-space: nowrap;
    padding: 8px 14px;
    border: 1px solid #333;
    border-radius: 20px;
    font-size: 13px;
  }

  .category-link.active {
    border-color: #d4a574;
    background: rgba(212, 165, 116, 0.15);
  }
}
</style>
