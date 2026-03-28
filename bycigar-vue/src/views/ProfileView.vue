<script setup>
import { ref, onMounted, computed, watch } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth";
import AddressForm from "../components/AddressForm.vue";
import { getStateName } from "../utils/states";

const router = useRouter();
const authStore = useAuthStore();

const activeTab = ref("info");
const loading = ref(false);
const orders = ref([]);
const message = ref("");

const user = computed(() => authStore.user);

const editMode = ref(false);
const editForm = ref({
  name: "",
  email: "",
});

const addresses = ref([]);
const showAddressForm = ref(false);
const editingAddress = ref(null);
const addressFormMode = ref("add");

onMounted(async () => {
  if (!authStore.isLoggedIn) {
    router.push("/login");
    return;
  }

  const isValid = await authStore.validateToken();
  if (!isValid) {
    router.push("/login");
    return;
  }

  editForm.value = {
    name: user.value?.name || "",
    email: user.value?.email || "",
  };
  fetchOrders();
  fetchAddresses();
});

async function fetchOrders() {
  loading.value = true;
  try {
    const res = await fetch("http://localhost:3000/api/orders", {
      headers: authStore.getAuthHeaders(),
    });
    if (res.ok) {
      const data = await res.json();
      orders.value = data.orders || [];
    }
  } catch (e) {
    console.error("获取订单失败:", e);
  } finally {
    loading.value = false;
  }
}

async function fetchAddresses() {
  try {
    const res = await fetch("http://localhost:3000/api/addresses", {
      headers: authStore.getAuthHeaders(),
    });
    if (res.ok) {
      const data = await res.json();
      addresses.value = data.addresses || [];
    }
  } catch (e) {
    console.error("获取地址失败:", e);
  }
}

function startEdit() {
  editMode.value = true;
  editForm.value = {
    name: user.value?.name || "",
    email: user.value?.email || "",
  };
}

function cancelEdit() {
  editMode.value = false;
}

async function saveProfile() {
  message.value = "";
  loading.value = true;
  try {
    const res = await fetch("http://localhost:3000/api/auth/profile", {
      method: "PUT",
      headers: authStore.getAuthHeaders(),
      body: JSON.stringify(editForm.value),
    });
    const data = await res.json();
    if (!res.ok) {
      throw new Error(data.error || "更新失败");
    }
    authStore.user = { ...authStore.user, ...editForm.value };
    localStorage.setItem("user", JSON.stringify(authStore.user));
    message.value = "保存成功";
    editMode.value = false;
  } catch (e) {
    message.value = e.message;
  } finally {
    loading.value = false;
  }
}

function handleLogout() {
  authStore.logout();
  router.push("/");
}

function openAddAddress() {
  editingAddress.value = null;
  addressFormMode.value = "add";
  showAddressForm.value = true;
}

function openEditAddress(address) {
  editingAddress.value = address;
  addressFormMode.value = "edit";
  showAddressForm.value = true;
}

function closeAddressForm() {
  showAddressForm.value = false;
  editingAddress.value = null;
}

async function handleSaveAddress(formData) {
  try {
    const url =
      addressFormMode.value === "edit"
        ? `http://localhost:3000/api/addresses/${editingAddress.value.id}`
        : "http://localhost:3000/api/addresses";
    const method = addressFormMode.value === "edit" ? "PUT" : "POST";

    const res = await fetch(url, {
      method,
      headers: authStore.getAuthHeaders(),
      body: JSON.stringify(formData),
    });
    const data = await res.json();
    if (!res.ok) {
      throw new Error(data.error || "保存失败");
    }
    await fetchAddresses();
    closeAddressForm();
  } catch (e) {
    throw e;
  }
}

async function deleteAddress(id) {
  if (!confirm("确定要删除这个地址吗？")) return;

  try {
    const res = await fetch(`http://localhost:3000/api/addresses/${id}`, {
      method: "DELETE",
      headers: authStore.getAuthHeaders(),
    });
    if (res.ok) {
      await fetchAddresses();
    }
  } catch (e) {
    console.error("删除地址失败:", e);
  }
}

async function setDefaultAddress(id) {
  try {
    const res = await fetch(
      `http://localhost:3000/api/addresses/${id}/default`,
      {
        method: "PUT",
        headers: authStore.getAuthHeaders(),
      },
    );
    if (res.ok) {
      await fetchAddresses();
    }
  } catch (e) {
    console.error("设置默认地址失败:", e);
  }
}

function formatAddress(addr) {
  let str = addr.addressLine1;
  if (addr.addressLine2) str += `, ${addr.addressLine2}`;
  str += `, ${addr.city}, ${getStateName(addr.state)} ${addr.zipCode}`;
  return str;
}

const tabs = [
  { id: "info", label: "个人信息" },
  { id: "addresses", label: "地址管理" },
  { id: "orders", label: "我的订单" },
  { id: "security", label: "账户安全" },
];

const addressLimit = 5;

const captchaId = ref("");
const captchaImage = ref("");
const passwordForm = ref({
  oldPassword: "",
  newPassword: "",
  confirmPassword: "",
  captchaCode: "",
});
const passwordMessage = ref("");
const passwordLoading = ref(false);

watch(activeTab, (val) => {
  if (val === "security") {
    refreshCaptcha();
  }
});

async function refreshCaptcha() {
  try {
    const res = await fetch("http://localhost:3000/api/auth/captcha");
    const data = await res.json();
    captchaId.value = data.captchaId;
    captchaImage.value = data.captchaImage;
    passwordForm.value.captchaCode = "";
  } catch (e) {
    console.error("获取验证码失败:", e);
  }
}

async function changePassword() {
  passwordMessage.value = "";
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    passwordMessage.value = "两次输入的新密码不一致";
    return;
  }
  if (passwordForm.value.newPassword.length < 6) {
    passwordMessage.value = "新密码至少需要6个字符";
    return;
  }
  passwordLoading.value = true;
  try {
    const res = await fetch("http://localhost:3000/api/auth/change-password", {
      method: "PUT",
      headers: authStore.getAuthHeaders(),
      body: JSON.stringify({
        oldPassword: passwordForm.value.oldPassword,
        newPassword: passwordForm.value.newPassword,
        captchaId: captchaId.value,
        captchaCode: passwordForm.value.captchaCode,
      }),
    });
    const data = await res.json();
    if (!res.ok) {
      throw new Error(data.error || "密码修改失败");
    }
    passwordMessage.value = "密码修改成功，请重新登录";
    passwordForm.value = {
      oldPassword: "",
      newPassword: "",
      confirmPassword: "",
      captchaCode: "",
    };
    setTimeout(() => {
      authStore.logout();
      router.push("/login");
    }, 2000);
  } catch (e) {
    passwordMessage.value = e.message;
    refreshCaptcha();
  } finally {
    passwordLoading.value = false;
  }
}
</script>

<template>
  <div class="profile-page">
    <div class="profile-container">
      <div class="profile-sidebar">
        <div class="user-avatar">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="48"
            height="48"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.5"
          >
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
            <circle cx="12" cy="7" r="4"></circle>
          </svg>
        </div>
        <div class="user-name">
          {{ user?.name || user?.email?.split("@")[0] || "用户" }}
        </div>
        <div class="user-email">{{ user?.email }}</div>

        <nav class="profile-nav">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            :class="['nav-item', { active: activeTab === tab.id }]"
            @click="activeTab = tab.id"
          >
            {{ tab.label }}
          </button>
          <button class="nav-item logout" @click="handleLogout">
            退出登录
          </button>
        </nav>
      </div>

      <div class="profile-content">
        <div v-if="activeTab === 'info'" class="content-section">
          <h2>个人信息</h2>

          <div
            v-if="message"
            :class="['message', { error: message.includes('失败') }]"
          >
            {{ message }}
          </div>

          <div v-if="editMode" class="edit-form">
            <div class="form-group">
              <label>用户名</label>
              <input
                v-model="editForm.name"
                type="text"
                placeholder="请输入用户名"
              />
            </div>
            <div class="form-group">
              <label>邮箱</label>
              <input
                v-model="editForm.email"
                type="email"
                placeholder="请输入邮箱"
                disabled
              />
            </div>
            <div class="form-actions">
              <button class="btn-cancel" @click="cancelEdit">取消</button>
              <button class="btn-save" @click="saveProfile" :disabled="loading">
                {{ loading ? "保存中..." : "保存" }}
              </button>
            </div>
          </div>

          <div v-else class="info-list">
            <div class="info-item">
              <span class="info-label">用户名</span>
              <span class="info-value">{{ user?.name || "未设置" }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">邮箱</span>
              <span class="info-value">{{ user?.email }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">账户角色</span>
              <span class="info-value">{{
                user?.role === "admin" ? "管理员" : "普通用户"
              }}</span>
            </div>
            <button class="btn-edit" @click="startEdit">编辑信息</button>
          </div>
        </div>

        <div v-if="activeTab === 'orders'" class="content-section">
          <h2>我的订单</h2>

          <div v-if="loading" class="loading">加载中...</div>

          <div v-else-if="orders.length === 0" class="empty-state">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="64"
              height="64"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1"
            >
              <path
                d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
              ></path>
              <polyline points="14 2 14 8 20 8"></polyline>
              <line x1="16" y1="13" x2="8" y2="13"></line>
              <line x1="16" y1="17" x2="8" y2="17"></line>
              <polyline points="10 9 9 9 8 9"></polyline>
            </svg>
            <p>暂无订单记录</p>
            <router-link to="/" class="btn-shop">去购物</router-link>
          </div>

          <div v-else class="orders-list">
            <div v-for="order in orders" :key="order.id" class="order-card">
              <div class="order-header">
                <span class="order-id">订单号: {{ order.orderNo }}</span>
                <span class="order-date">{{
                  new Date(order.createdAt).toLocaleDateString()
                }}</span>
                <span :class="['order-status', order.status]">{{
                  order.status
                }}</span>
              </div>
              <div class="order-items">
                <div
                  v-for="item in order.items"
                  :key="item.id"
                  class="order-item"
                >
                  <img
                    :src="item.product?.imageUrl"
                    :alt="item.product?.name"
                  />
                  <div class="item-info">
                    <span class="item-name">{{ item.product?.name }}</span>
                    <span class="item-qty">x{{ item.quantity }}</span>
                  </div>
                  <span class="item-price"
                    >${{ Number(item.price).toFixed(2) }}</span
                  >
                </div>
              </div>
              <div class="order-footer">
                <span class="order-total"
                  >合计: ${{ Number(order.total).toFixed(2) }}</span
                >
              </div>
            </div>
          </div>
        </div>

        <div v-if="activeTab === 'security'" class="content-section">
          <h2>账户安全</h2>
          <div class="password-form">
            <div
              v-if="passwordMessage"
              :class="[
                'message',
                {
                  error:
                    passwordMessage.includes('失败') ||
                    passwordMessage.includes('错误') ||
                    passwordMessage.includes('不一致'),
                },
              ]"
            >
              {{ passwordMessage }}
            </div>

            <div class="form-group">
              <label>原密码</label>
              <input
                v-model="passwordForm.oldPassword"
                type="password"
                placeholder="请输入原密码"
              />
            </div>

            <div class="form-group">
              <label>新密码</label>
              <input
                v-model="passwordForm.newPassword"
                type="password"
                placeholder="至少6个字符"
              />
            </div>

            <div class="form-group">
              <label>确认新密码</label>
              <input
                v-model="passwordForm.confirmPassword"
                type="password"
                placeholder="再次输入新密码"
              />
            </div>

            <div class="form-group">
              <label>验证码</label>
              <div class="captcha-row">
                <input
                  v-model="passwordForm.captchaCode"
                  type="text"
                  placeholder="请输入验证码"
                  maxlength="4"
                />
                <img
                  v-if="captchaImage"
                  :src="captchaImage"
                  class="captcha-img"
                  @click="refreshCaptcha"
                  title="点击刷新验证码"
                />
                <button
                  type="button"
                  class="btn-refresh-captcha"
                  @click="refreshCaptcha"
                >
                  刷新
                </button>
              </div>
            </div>

            <div class="form-actions">
              <button
                class="btn-save"
                @click="changePassword"
                :disabled="passwordLoading"
              >
                {{ passwordLoading ? "提交中..." : "修改密码" }}
              </button>
            </div>
          </div>
        </div>

        <div v-if="activeTab === 'addresses'" class="content-section">
          <h2>地址管理</h2>

          <div v-if="showAddressForm" class="address-form-overlay">
            <AddressForm
              :address="editingAddress"
              :mode="addressFormMode"
              :save-function="handleSaveAddress"
              @cancel="closeAddressForm"
            />
          </div>

          <div v-else>
            <div class="address-header-row">
              <span class="address-count"
                >{{ addresses.length }} / {{ addressLimit }} 个地址</span
              >
              <button
                v-if="addresses.length < addressLimit"
                class="btn-add-address"
                @click="openAddAddress"
              >
                + 新增地址
              </button>
            </div>

            <div v-if="addresses.length === 0" class="empty-address">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="48"
                height="48"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  d="M21 10c-1.34-.065-.765-.235-1.485-.45A2.07 2 2.07H6a2 2 0 0-2 2v16a2 2 0 0 2 2h12a2 2 0 0 2-2V8z"
                ></path>
                <polyline points="21 10 15 7 15 7 15 7-3 3 3-3 3z"></polyline>
              </svg>
              <p>您还没有保存任何收货地址</p>
              <p class="hint">最多可保存 {{ addressLimit }} 个地址</p>
              <button class="btn-add-address" @click="openAddAddress">
                添加地址
              </button>
            </div>

            <div v-else class="address-list">
              <div
                v-for="addr in addresses"
                :key="addr.id"
                :class="['address-card', { default: addr.isDefault }]"
              >
                <div class="card-header">
                  <span v-if="addr.isDefault" class="default-badge">默认</span>
                  <span class="recipient-name">{{ addr.fullName }}</span>
                </div>
                <div class="card-body">
                  <div class="address-line">
                    <svg
                      class="addr-icon"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                    >
                      <path
                        d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"
                      />
                      <circle cx="12" cy="10" r="3" />
                    </svg>
                    <span class="addr-label">详细地址：</span
                    >{{ addr.addressLine1 }}
                  </div>
                  <div v-if="addr.addressLine2" class="address-line">
                    <svg
                      class="addr-icon"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                    >
                      <path
                        d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"
                      />
                      <circle cx="12" cy="10" r="3" />
                    </svg>
                    <span class="addr-label">公寓单元：</span
                    >{{ addr.addressLine2 }}
                  </div>
                  <div class="address-line">
                    <svg
                      class="addr-icon"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                    >
                      <path d="M3 21h18M5 21V7l8-4v18M19 21V11l-6-4" />
                    </svg>
                    <span class="addr-label">城市邮编：</span>{{ addr.city }},
                    {{ getStateName(addr.state) }} {{ addr.zipCode }}
                  </div>
                  <div class="address-phone">
                    <svg
                      class="addr-icon"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                    >
                      <path
                        d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.127.96.361 1.903.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0 1 22 16.92z"
                      />
                    </svg>
                    <span class="addr-label">联系电话：</span>{{ addr.phone }}
                  </div>
                </div>
                <div class="card-actions">
                  <button
                    v-if="!addr.isDefault"
                    class="btn-set-default"
                    @click="setDefaultAddress(addr.id)"
                  >
                    设为默认
                  </button>
                  <button class="btn-edit-addr" @click="openEditAddress(addr)">
                    编辑
                  </button>
                  <button
                    class="btn-delete-addr"
                    @click="deleteAddress(addr.id)"
                  >
                    删除
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.profile-page {
  min-height: 70vh;
  padding: 40px 20px;
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
}

.profile-container {
  max-width: 1000px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 250px 1fr;
  gap: 30px;
}

.profile-sidebar {
  background: #2d2d2d;
  border-radius: 12px;
  padding: 30px 20px;
  text-align: center;
  height: fit-content;
}

.user-avatar {
  width: 80px;
  height: 80px;
  background: #d4a574;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 15px;
  color: #1a1a1a;
}

.user-name {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 5px;
}

.user-email {
  font-size: 13px;
  color: #888;
  margin-bottom: 30px;
}

.profile-nav {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.nav-item {
  background: transparent;
  border: none;
  padding: 12px 16px;
  color: #aaa;
  font-size: 14px;
  cursor: pointer;
  border-radius: 8px;
  text-align: left;
  transition: all 0.3s;
}

.nav-item:hover {
  background: #3d3d3d;
  color: #fff;
}

.nav-item.active {
  background: #d4a574;
  color: #1a1a1a;
}

.nav-item.logout {
  margin-top: 20px;
  color: #e74c3c;
}

.nav-item.logout:hover {
  background: rgba(231, 76, 60, 0.1);
}

.profile-content {
  background: #2d2d2d;
  border-radius: 12px;
  padding: 30px;
}

.content-section h2 {
  color: #d4a574;
  font-size: 22px;
  margin: 0 0 25px;
  padding-bottom: 15px;
  border-bottom: 1px solid #444;
}

.message {
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 20px;
  background: rgba(212, 165, 116, 0.1);
  color: #d4a574;
}

.message.error {
  background: rgba(231, 76, 60, 0.1);
  color: #e74c3c;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 15px 0;
  border-bottom: 1px solid #3d3d3d;
}

.info-label {
  color: #888;
}

.info-value {
  color: #fff;
}

.btn-edit {
  background: transparent;
  border: 1px solid #d4a574;
  color: #d4a574;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
  margin-top: 10px;
  transition: all 0.3s;
}

.btn-edit:hover {
  background: #d4a574;
  color: #1a1a1a;
}

.edit-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  color: #ccc;
  font-size: 14px;
}

.form-group input {
  padding: 12px 16px;
  background: #1a1a1a;
  border: 1px solid #444;
  border-radius: 8px;
  color: #fff;
  font-size: 15px;
}

.form-group input:focus {
  outline: none;
  border-color: #d4a574;
}

.form-group input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-actions {
  display: flex;
  gap: 15px;
  margin-top: 10px;
}

.btn-cancel,
.btn-save {
  padding: 12px 24px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
}

.btn-cancel {
  background: transparent;
  border: 1px solid #666;
  color: #aaa;
}

.btn-cancel:hover {
  border-color: #888;
  color: #fff;
}

.btn-save {
  background: #d4a574;
  border: none;
  color: #1a1a1a;
  font-weight: 600;
}

.btn-save:hover:not(:disabled) {
  background: #c49564;
}

.btn-save:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #888;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #888;
}

.empty-state svg {
  margin-bottom: 20px;
  opacity: 0.5;
}

.empty-state p {
  margin-bottom: 20px;
}

.btn-shop {
  display: inline-block;
  background: #d4a574;
  color: #1a1a1a;
  padding: 12px 30px;
  border-radius: 6px;
  text-decoration: none;
  font-weight: 600;
  transition: all 0.3s;
}

.btn-shop:hover {
  background: #c49564;
}

.orders-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.order-card {
  background: #1a1a1a;
  border-radius: 10px;
  overflow: hidden;
}

.order-header {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px 20px;
  background: #252525;
  border-bottom: 1px solid #333;
}

.order-id {
  color: #fff;
  font-weight: 500;
}

.order-date {
  color: #888;
  font-size: 13px;
}

.order-status {
  margin-left: auto;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  text-transform: uppercase;
}

.order-status.pending {
  background: rgba(241, 196, 15, 0.2);
  color: #f1c40f;
}

.order-status.paid {
  background: rgba(46, 204, 113, 0.2);
  color: #2ecc71;
}

.order-status.shipped {
  background: rgba(52, 152, 219, 0.2);
  color: #3498db;
}

.order-status.delivered {
  background: rgba(155, 89, 182, 0.2);
  color: #9b59b6;
}

.order-status.cancelled {
  background: rgba(231, 76, 60, 0.2);
  color: #e74c3c;
}

.order-items {
  padding: 15px 20px;
}

.order-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 10px 0;
}

.order-item img {
  width: 50px;
  height: 50px;
  object-fit: cover;
  border-radius: 6px;
  background: #fff;
}

.item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.item-name {
  color: #ccc;
  font-size: 14px;
}

.item-qty {
  color: #888;
  font-size: 12px;
}

.item-price {
  color: #d4a574;
  font-weight: 500;
}

.order-footer {
  padding: 15px 20px;
  border-top: 1px solid #333;
  text-align: right;
}

.order-total {
  color: #fff;
  font-weight: 600;
  font-size: 16px;
}

.password-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 480px;
}

.captcha-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.captcha-row input {
  flex: 1;
}

.captcha-img {
  height: 44px;
  border-radius: 6px;
  cursor: pointer;
  border: 1px solid #444;
  flex-shrink: 0;
  background: #fff;
}

.captcha-img:hover {
  border-color: #d4a574;
}

.btn-refresh-captcha {
  background: transparent;
  border: 1px solid #666;
  color: #aaa;
  padding: 10px 14px;
  border-radius: 6px;
  font-size: 13px;
  white-space: nowrap;
  transition: all 0.3s;
  flex-shrink: 0;
}

.btn-refresh-captcha:hover {
  border-color: #d4a574;
  color: #d4a574;
}

.address-header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.address-count {
  color: #888;
  font-size: 14px;
}

.btn-add-address {
  background: transparent;
  border: 1px solid #d4a574;
  color: #d4a574;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
}

.btn-add-address:hover {
  background: #d4a574;
  color: #1a1a1a;
}

.btn-add-address:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.empty-address {
  text-align: center;
  padding: 60px 20px;
  color: #888;
}

.empty-address svg {
  margin-bottom: 20px;
  opacity: 0.5;
}

.empty-address p {
  margin-bottom: 10px;
}

.empty-address .hint {
  font-size: 13px;
  color: #666;
}

.address-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.address-card {
  background: #1a1a1a;
  border-radius: 10px;
  padding: 20px;
  border: 2px solid transparent;
  transition: all 0.3s;
}

.address-card.default {
  border-color: #d4a574;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.default-badge {
  background: #d4a574;
  color: #1a1a1a;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.recipient-name {
  color: #fff;
  font-size: 16px;
  font-weight: 500;
}

.card-body {
  color: #ccc;
  font-size: 14px;
  line-height: 1.8;
}

.address-line {
  margin-bottom: 4px;
  display: flex;
  align-items: flex-start;
}

.addr-icon {
  width: 16px;
  height: 16px;
  margin-right: 6px;
  flex-shrink: 0;
  margin-top: 2px;
}

.addr-label {
  color: #888;
  margin-right: 4px;
  flex-shrink: 0;
}

.address-phone {
  margin-top: 8px;
  color: #888;
  display: flex;
  align-items: center;
}

.card-actions {
  display: flex;
  gap: 10px;
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #333;
}

.btn-set-default {
  background: transparent;
  border: 1px solid #666;
  color: #aaa;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.btn-set-default:hover {
  border-color: #d4a574;
  color: #d4a574;
}

.btn-edit-addr,
.btn-delete-addr {
  background: transparent;
  border: none;
  padding: 6px 12px;
  cursor: pointer;
  font-size: 12px;
}

.btn-edit-addr {
  color: #d4a574;
}

.btn-edit-addr:hover {
  text-decoration: underline;
}

.btn-delete-addr {
  color: #e74c3c;
}

.btn-delete-addr:hover {
  text-decoration: underline;
}

.address-form-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

@media (max-width: 768px) {
  .profile-container {
    grid-template-columns: 1fr;
  }

  .profile-sidebar {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .profile-nav {
    width: 100%;
    flex-direction: row;
    flex-wrap: wrap;
    justify-content: center;
  }

  .nav-item {
    text-align: center;
  }

  .address-header-row {
    flex-direction: column;
    gap: 10px;
    align-items: flex-start;
  }
}
</style>
