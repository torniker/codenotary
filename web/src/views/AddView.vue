<script lang="ts" setup>
import { Type, type Accounting } from '@/types/accounting'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { reactive, ref } from 'vue'

const formRef = ref<FormInstance>()
const loading = ref(false)
const formSize = ref('default')
const form = reactive<Accounting>({
  id: '',
  account_number: '',
  account_name: '',
  iban: '',
  address: '',
  amount: 0,
  type: Type.Receiving
})

const rules = reactive<FormRules<Accounting>>({
  account_number: [
    { required: true, message: 'Please input Account number', trigger: 'blur' },
    { min: 3, max: 25, message: 'Length should be 3 to 25', trigger: 'blur' }
  ],
  account_name: [{ required: true, message: 'Please input Account Name', trigger: 'blur' }],
  iban: [
    { required: true, message: 'Please input IBAN', trigger: 'blur' },
    { min: 3, max: 30, message: 'Length should be 3 to 30', trigger: 'blur' }
  ],
  amount: [
    {
      type: 'number',
      required: true,
      trigger: 'blur',
      message: 'Please input amount'
    }
  ],
  type: [
    {
      required: true,
      message: 'Please select type',
      trigger: 'change'
    }
  ]
})

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate(async (valid, fields) => {
    if (valid) {
      loading.value = true
      const response = await fetch('/api/accounting', {
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json'
        },
        method: 'POST',
        body: JSON.stringify(form)
      })
      if (response.status === 200) {
        ElMessage('Accounting added successfully')
        formEl.resetFields()
      } else {
        ElMessage('Error adding accounting')
      }
      loading.value = false
    }
  })
}
</script>
<template>
  <el-form
    v-loading="loading"
    ref="formRef"
    :model="form"
    :rules="rules"
    :size="formSize"
    status-icon
    label-width="150px"
  >
    <el-form-item label="Account Number" prop="account_number">
      <el-input v-model="form.account_number" />
    </el-form-item>
    <el-form-item label="Account Name" prop="account_name">
      <el-input v-model="form.account_name" />
    </el-form-item>
    <el-form-item label="IBAN" prop="iban">
      <el-input v-model="form.iban" />
    </el-form-item>
    <el-form-item label="Address" prop="address">
      <el-input v-model="form.address" />
    </el-form-item>
    <el-form-item label="Amount" prop="amount">
      <el-input-number v-model="form.amount" :min="0" :max="100000000" />
    </el-form-item>
    <el-form-item label="Type" prop="type">
      <el-select v-model="form.type" placeholder="Select type">
        <el-option v-for="t in Type" :key="t" :label="t" :value="t" />
      </el-select>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="submitForm(formRef)">Add</el-button>
    </el-form-item>
  </el-form>
</template>
