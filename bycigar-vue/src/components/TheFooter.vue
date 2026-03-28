<script setup>
import { onMounted } from 'vue'
import { useSettingsStore } from '../stores/useSettingsStore'

const currentYear = new Date().getFullYear()
const settingsStore = useSettingsStore()

const footerLinks = [
  { name: '关于我们', path: '/about' },
  { name: '服务条款', path: '/services' },
  { name: '隐私政策', path: '/privacy-policy' },
  { name: '严正声明', path: '/statement' }
]

onMounted(() => {
  settingsStore.fetchSettings()
})
</script>

<template>
  <footer class="site-footer">
    <div class="container">
      <div class="footer-content">
        <div class="footer-links">
          <router-link
            v-for="link in footerLinks"
            :key="link.path"
            :to="link.path"
            class="footer-link"
          >
            {{ link.name }}
          </router-link>
        </div>

        <div class="footer-description">
          <p>{{ settingsStore.footerDescription }}</p>
          <p class="service-time"><strong>{{ settingsStore.footerServiceTime }}</strong></p>
        </div>


      </div>
    </div>
  </footer>
</template>

<style scoped>
.site-footer {
  background: #1a1a1a;
  color: #ccc;
  padding: 40px 0 20px;
  margin-top: 60px;
}

.footer-content {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.footer-links {
  display: flex;
  justify-content: space-between;
  border-bottom: 1px solid #333;
  padding-bottom: 30px;
}

.footer-link {
  color: #ccc;
  text-decoration: none;
  padding: 5px 0;
  transition: color 0.3s;
  font-size: 14px;
}

.footer-link:hover {
  color: #d4a574;
}

.footer-brand h4 {
  color: #fff;
  font-size: 16px;
  margin: 0;
  font-weight: 500;
}

.footer-description {
  font-size: 14px;
  line-height: 1.8;
  color: #999;
}

.footer-description p {
  margin-bottom: 10px;
}

.service-time {
  color: #28a745;
}

@media (max-width: 768px) {
  .footer-links {
    flex-wrap: wrap;
    justify-content: center;
    gap: 20px 40px;
  }
}
</style>
