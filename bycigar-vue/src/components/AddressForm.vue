<script setup>
import { ref, computed, watch } from 'vue'
import { US_STATES } from '../utils/states'

const props = defineProps({
  address: {
    type: Object,
    default: null
  },
  mode: {
    type: String,
    default: 'add'
  },
  saveFunction: {
    type: Function,
    required: true
  }
})

const emit = defineEmits(['save', 'cancel'])

const form = ref({
  fullName: '',
  addressLine1: '',
  addressLine2: '',
  city: '',
  state: '',
  zipCode: '',
  phone: '',
  isDefault: false
})

const loading = ref(false)
const error = ref('')

watch(() => props.address, (newAddress) => {
  if (newAddress) {
    form.value = { ...newAddress }
  }
}, { immediate: true })

function formatPhone(value) {
  const numbers = value.replace(/\D/g, '')
  if (numbers.length <= 3) return numbers
  if (numbers.length <= 6) return `(${numbers.slice(0, 3)}) ${numbers.slice(3)}`
  return `(${numbers.slice(0, 3)}) ${numbers.slice(3, 6)}-${numbers.slice(6, 10)}`
}

function handlePhoneInput(e) {
  form.value.phone = formatPhone(e.target.value)
}

function formatZip(value) {
  return value.replace(/\D/g, '').slice(0, 5)
}

function handleZipInput(e) {
  form.value.zipCode = formatZip(e.target.value)
}

function validate() {
  if (!form.value.fullName.trim()) return '请输入收件人姓名'
  if (!form.value.addressLine1.trim()) return '请输入街道地址'
  if (!form.value.city.trim()) return '请输入城市'
  if (!form.value.state) return '请选择州'
  if (!form.value.zipCode.trim()) return '请输入邮编'
  if (form.value.zipCode.length < 5) return '邮编格式不正确'
  if (!form.value.phone.trim()) return '请输入电话号码'
  if (form.value.phone.replace(/\D/g, '').length < 10) return '电话号码格式不正确'
  return null
}

async function handleSubmit() {
  error.value = ''
  
  const validationError = validate()
  if (validationError) {
    error.value = validationError
    return
  }
  
  loading.value = true
  try {
    await props.saveFunction({ ...form.value })
  } catch (e) {
    error.value = e.message || '保存失败'
    loading.value = false
  }
}

function handleCancel() {
  emit('cancel')
}
</script>

<template>
  <div class="address-form">
    <div v-if="error" class="error-message">{{ error }}</div>
    
    <div class="form-row">
      <div class="form-group full">
        <label>收件人姓名 *</label>
        <input 
          v-model="form.fullName" 
          type="text" 
          placeholder="Full Name"
          maxlength="50"
        >
      </div>
    </div>
    
    <div class="form-row">
      <div class="form-group full">
        <label>街道地址 *</label>
        <input 
          v-model="form.addressLine1" 
          type="text" 
          placeholder="Street Address"
          maxlength="100"
        >
      </div>
    </div>
    
    <div class="form-row">
      <div class="form-group full">
        <label>公寓/单元号 (可选)</label>
        <input 
          v-model="form.addressLine2" 
          type="text" 
          placeholder="Apt, Suite, Unit, etc."
          maxlength="50"
        >
      </div>
    </div>
    
    <div class="form-row">
      <div class="form-group">
        <label>城市 *</label>
        <input 
          v-model="form.city" 
          type="text" 
          placeholder="City"
          maxlength="50"
        >
      </div>
      
      <div class="form-group">
        <label>州 *</label>
        <select v-model="form.state">
          <option value="">选择州</option>
          <option v-for="s in US_STATES" :key="s.code" :value="s.code">
            {{ s.code }} - {{ s.name }}
          </option>
        </select>
      </div>
      
      <div class="form-group zip">
        <label>邮编 *</label>
        <input 
          :value="form.zipCode"
          @input="handleZipInput"
          type="text" 
          placeholder="12345"
          maxlength="5"
        >
      </div>
    </div>
    
    <div class="form-row">
      <div class="form-group full">
        <label>电话号码 *</label>
        <input 
          :value="form.phone"
          @input="handlePhoneInput"
          type="text" 
          placeholder="(213) 555-1234"
          maxlength="14"
        >
      </div>
    </div>
    
    <div class="form-row checkbox-row">
      <label class="checkbox-label">
        <input type="checkbox" v-model="form.isDefault">
        <span>设为默认地址</span>
      </label>
    </div>
    
    <div class="form-actions">
      <button class="btn-cancel" @click="handleCancel">取消</button>
      <button class="btn-save" @click="handleSubmit" :disabled="loading">
        {{ loading ? '保存中...' : (mode === 'add' ? '添加地址' : '保存修改') }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.address-form {
  max-width: 500px;
}

.error-message {
  background: rgba(231, 76, 60, 0.1);
  color: #e74c3c;
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.form-row {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
}

.form-group {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group.full {
  flex: 1;
}

.form-group.zip {
  flex: 0 0 100px;
}

.form-group label {
  color: #888;
  font-size: 13px;
}

.form-group input,
.form-group select {
  padding: 10px 14px;
  background: #1a1a1a;
  border: 1px solid #444;
  border-radius: 6px;
  color: #fff;
  font-size: 14px;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #d4a574;
}

.form-group select {
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%23888' stroke-width='2'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
  padding-right: 36px;
}

.form-group select option {
  background: #1a1a1a;
  color: #fff;
}

.checkbox-row {
  margin: 20px 0;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  color: #ccc;
  font-size: 14px;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  accent-color: #d4a574;
  cursor: pointer;
}

.form-actions {
  display: flex;
  gap: 15px;
  margin-top: 25px;
}

.btn-cancel,
.btn-save {
  padding: 12px 24px;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
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

@media (max-width: 500px) {
  .form-row {
    flex-direction: column;
    gap: 15px;
  }
  
  .form-group.zip {
    flex: 1;
  }
}
</style>
