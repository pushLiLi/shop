<script setup>
const currentYear = new Date().getFullYear()

const footerLinks = [
  { name: '关于我们', path: '/about' },
  { name: '服务条款', path: '/services' },
  { name: '隐私政策', path: '/privacy-policy' },
  { name: '严正声明', path: '/returns-policy' }
]
</script>

<template>
  <footer class="site-footer">
    <div class="container">
      <div class="footer-content">
        <div class="footer-links-section">
          <div class="row">
            <div class="col-md-7">
              <div class="row">
                <div v-for="link in footerLinks" :key="link.path" class="col-6 col-md-3">
                  <router-link :to="link.path" class="footer-link">
                    {{ link.name }}
                  </router-link>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="footer-description">
          <p>BYCIGAR是中国领先的雪茄文化与在线购物平台。我们提供最新、最专业的雪茄测评、品牌新闻与养护知识，并为您甄选全球优质雪茄及配件，支持便捷在线购买。加入我们的雪茄社区，探索醇香世界。</p>
          <p class="service-time"><strong>客服在线时间每周一至周六 9:00到18:00</strong></p>
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

.footer-links-section {
  border-bottom: 1px solid #333;
  padding-bottom: 30px;
}

.footer-link {
  display: block;
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
  .footer-links-section .col-6 {
    margin-bottom: 15px;
  }

  .footer-brand {
    margin-top: 20px;
  }


}
</style>
