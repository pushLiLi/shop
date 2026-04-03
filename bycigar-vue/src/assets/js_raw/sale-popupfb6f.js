document.addEventListener('DOMContentLoaded', function() {
  const salePopup = document.getElementById('salePopup');
  const closeSalePopup = document.getElementById('closeSalePopup');
  const saleProductsContainer = document.getElementById('saleProductsContainer');
  // 获取所有特卖图标（包括PC端和手机端）
  const saleIcons = document.querySelectorAll('.header-sale-icon');
  let popupTimer;

  // 使用全局 lang 对象进行翻译

  // 全局变量
  let currentSlide = 0;
  let slideInterval;
  let shouldShowPopup = true;

  // 加载特卖商品
  function loadSaleProducts() {
    // 从API获取特卖商品数据
    // 使用相对路径，确保语言代码格式正确（小写）
    const currentLang = (document.documentElement.lang || 'zh-cn').toLowerCase();
    const apiUrl = `/${currentLang}/products/sale`;
    
    fetch(apiUrl)
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.json();
      })
      .then(data => {
        const saleProducts = data.data || [];

        // 如果没有特卖商品，不展示窗口
        if (saleProducts.length === 0) {
          console.log('没有特卖商品');
          return;
        }

        // 渲染商品
        renderSaleProducts(saleProducts);
        
        // 检查是否是当天第一次访问，如果是则显示窗口
        if (isFirstVisitToday()) {
          showSalePopup();
        }
      })
      .catch(error => {
        console.error('获取特卖商品失败:', error);
        // 加载失败时不显示窗口
      });
  }
  

  // 渲染特卖商品
  function renderSaleProducts(products) {
    // 清空容器
    saleProductsContainer.innerHTML = '';
    
    // 创建商品包装器
    const productsWrapper = document.createElement('div');
    productsWrapper.className = 'sale-products-wrapper';
    
    // 添加商品
    products.forEach(product => {
      // 直接使用接口返回的数据
      const { id, skuId, brand, brandLogo, name, price, originalPrice, image, detailUrl, hasStock } = product;

      const productItem = document.createElement('div');
      productItem.className = 'sale-product-item';
      productItem.innerHTML = `
        <div class="product-brand">
          ${brandLogo ? `<img src="${brandLogo}" alt="${brand}" class="brand-logo">` : ''}
          ${brand}
        </div>
        <div class="product-name">${name}</div>
        <div class="product-image">
          <img src="${image || 'https://via.placeholder.com/200x200?text=Product'}" alt="${name}">
        </div>
        <div class="product-price">
          <span class="current-price">¥${parseFloat(price).toFixed(2)}</span>
          ${originalPrice > 0 ? `<span class="original-price">¥${parseFloat(originalPrice).toFixed(2)}</span>` : ''}
        </div>
        <div class="product-actions">
          <button class="add-to-cart-btn" data-sku-id="${skuId}" ${!hasStock ? 'disabled' : ''}>${hasStock ? lang.sale.add_to_cart : lang.sale.out_of_stock}</button>
          <button class="view-detail-btn" data-detail-url="${detailUrl}">${lang.sale.details}</button>
        </div>
      `;
      productsWrapper.appendChild(productItem);
    });
    
    // 将包装器添加到容器
    saleProductsContainer.appendChild(productsWrapper);

    // 绑定事件
    bindProductEvents();
    
    // 初始化幻灯片
    initSlideShow(products.length);
  }

  // 初始化幻灯片
  function initSlideShow(totalSlides) {
    // 重置当前幻灯片索引
    currentSlide = 0;
    
    // 清除之前的定时器
    if (slideInterval) {
      clearInterval(slideInterval);
      slideInterval = null;
    }
    
    // 只有一个商品时不需要自动切换
    if (totalSlides > 1) {
      // 设置自动切换定时器
      slideInterval = setInterval(nextSlide, 3000);
    }
  }

  // 切换到下一张幻灯片
  function nextSlide() {
    const wrapper = document.querySelector('.sale-products-wrapper');
    if (!wrapper) return;
    
    const totalSlides = document.querySelectorAll('.sale-product-item').length;
    if (totalSlides <= 1) return;
    
    currentSlide = (currentSlide + 1) % totalSlides;
    updateSlidePosition();
  }

  // 切换到上一张幻灯片
  function prevSlide() {
    const wrapper = document.querySelector('.sale-products-wrapper');
    if (!wrapper) return;
    
    const totalSlides = document.querySelectorAll('.sale-product-item').length;
    if (totalSlides <= 1) return;
    
    currentSlide = (currentSlide - 1 + totalSlides) % totalSlides;
    updateSlidePosition();
  }

  // 更新幻灯片位置
  function updateSlidePosition() {
    const wrapper = document.querySelector('.sale-products-wrapper');
    if (!wrapper) return;
    
    const slide = document.querySelector('.sale-product-item');
    if (!slide) return;
    
    const slideWidth = slide.offsetWidth;
    wrapper.style.transform = `translateX(-${currentSlide * slideWidth}px)`;
  }

  // 停止幻灯片自动切换
  function stopSlideShow() {
    if (slideInterval) {
      clearInterval(slideInterval);
      slideInterval = null;
    }
  }

  // 重启幻灯片自动切换
  function startSlideShow() {
    // 清除之前的定时器
    if (slideInterval) {
      clearInterval(slideInterval);
      slideInterval = null;
    }
    
    const totalSlides = document.querySelectorAll('.sale-product-item').length;
    if (totalSlides > 1) {
      slideInterval = setInterval(nextSlide, 3000);
    }
  }

  // 绑定商品事件
  function bindProductEvents() {
    // 加入购物车按钮
    document.querySelectorAll('.add-to-cart-btn').forEach(btn => {
      btn.addEventListener('click', function() {
        const skuId = this.getAttribute('data-sku-id');
        // 调用加入购物车接口
        addToCart(skuId);
      });
    });

    // 详情按钮
    document.querySelectorAll('.view-detail-btn').forEach(btn => {
      btn.addEventListener('click', function() {
        const detailUrl = this.getAttribute('data-detail-url');
        window.location.href = detailUrl;
      });
    });
  }

  // 加入购物车函数
  function addToCart(productId) {
    // 获取CSRF token
    const csrfToken = document.querySelector('meta[name="csrf-token"]').getAttribute('content');
    
    // 准备请求数据
    const requestData = {
      sku_id: productId,
      quantity: 1,
      buy_now: false
    };
    
    // 发送请求
    fetch(urls.cart_add, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-CSRF-TOKEN': csrfToken,
        'X-Requested-With': 'XMLHttpRequest'
      },
      body: JSON.stringify(requestData)
    })
    .then(response => response.json())
    .then(data => {
      console.log('Cart add response:', data);
      if (data.success) {
          // 显示成功消息
          layer.msg(lang.cart.added_to_cart, {icon: 1});
          // 跳转到购物车页面
          setTimeout(() => {
            window.location.href = urls.cart;
          }, 500);
        } else {
          // 显示错误消息
          layer.msg(lang.cart.add_to_cart_failed, {icon: 2});
        }
    });
  }

  // 检查是否是当天第一次访问
  function isFirstVisitToday() {
    const today = new Date().toISOString().split('T')[0];
    const lastVisit = localStorage.getItem('lastSalePopupVisit');
    
    if (lastVisit !== today) {
      // 不是当天访问过，更新访问记录
      localStorage.setItem('lastSalePopupVisit', today);
      return true;
    }
    
    return false;
  }

  // 显示特卖窗口
  function showSalePopup(forceShow = false, iconElement = null) {
    // 只有当shouldShowPopup为true且（是第一次访问或强制显示）时才显示
    if (!shouldShowPopup && !forceShow) return;
    
    // 计算特卖图标位置，使窗口与图标对齐
    let saleIcon = iconElement || document.querySelector('.header-sale-icon');
    if (saleIcon) {
      const rect = saleIcon.getBoundingClientRect();
      
      // 检查是否为移动设备
      const isMobile = window.innerWidth <= 768;
      
      if (isMobile) {
        // 在移动设备上，居中显示窗口
        const popupWidth = salePopup.offsetWidth || 280;
        const popupHeight = salePopup.offsetHeight || 350;
        
        // 计算居中位置
        const top = (window.innerHeight - popupHeight) / 2 + window.scrollY;
        const left = (window.innerWidth - popupWidth) / 2 + window.scrollX;
        
        // 应用居中位置
        salePopup.style.top = `${top}px`;
        salePopup.style.left = `${left}px`;
        salePopup.style.right = 'auto';
      } else {
        // 在PC端，以特卖图标为相对位置，确保位置不变
        // 由于头部是固定定位，特卖窗口也应该使用固定定位
        // 计算相对于视口的位置
        const top = rect.bottom + 20;
        const right = window.innerWidth - rect.right;
        
        // 应用固定定位
        salePopup.style.position = 'fixed';
        salePopup.style.top = `${top}px`;
        salePopup.style.right = `${right}px`;
        salePopup.style.left = 'auto';
        salePopup.style.bottom = 'auto';
      }
    }
    
    salePopup.style.display = 'block';
    
    // 重启幻灯片自动切换
    startSlideShow();
    
    // 5分钟后自动关闭
    clearTimeout(popupTimer);
    popupTimer = setTimeout(() => {
      hideSalePopup();
    }, 300000);
  }

  // 隐藏特卖窗口
  function hideSalePopup() {
    salePopup.style.display = 'none';
    
    // 停止幻灯片自动切换
    stopSlideShow();
  }

  // 关闭按钮事件
  closeSalePopup.addEventListener('click', hideSalePopup);

  // 左右切换按钮事件
  if (slideBtnLeft) {
    slideBtnLeft.addEventListener('click', function() {
      prevSlide();
      // 点击后重置自动切换定时器
      startSlideShow();
    });
  }

  if (slideBtnRight) {
    slideBtnRight.addEventListener('click', function() {
      nextSlide();
      // 点击后重置自动切换定时器
      startSlideShow();
    });
  }

  // 为所有特卖图标添加点击事件（包括PC端和手机端）
  saleIcons.forEach(function(saleIcon) {
    saleIcon.addEventListener('click', function(e) {
      // 阻止事件冒泡
      e.stopPropagation();
      // 点击图标时重置shouldShowPopup为true并强制显示窗口
      shouldShowPopup = true;
      showSalePopup(true, this);
    });
  });

  // 初始化
  loadSaleProducts();

  // 模拟倒计时
  function updateCountdown() {
    // 这里可以实现真实的倒计时逻辑
    console.log('更新倒计时');
  }

  // 每秒更新倒计时
  setInterval(updateCountdown, 1000);
});
