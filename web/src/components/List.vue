<script setup lang="ts">
import router from '@/router'
import type { Accounting } from '@/types/accounting'
import { onMounted, ref } from 'vue'

const tableData = ref<Accounting[]>([])
const total = ref(0)
const loading = ref(false)
const page = ref(1)
const perPage = ref(10)

onMounted(() => {
  fetchAll()
})

router.afterEach((to, from) => {
  fetchAll()
})

function handlePageUpdate(p: number) {
  page.value = p || 1
  router.push({ query: { page: page.value, perPage: perPage.value } })
}

function handlePerPageUpdate(pp: number) {
  perPage.value = pp || perPage.value
  router.push({ query: { page: page.value, perPage: perPage.value } })
  fetchData(page.value, perPage.value)
}

function fetchAll() {
  page.value = parseInt(router.currentRoute.value.query.page as string) || 1
  perPage.value = parseInt(router.currentRoute.value.query.perPage as string) || 10

  fetchData(page.value, perPage.value)
  fetchCount()
}

function fetchData(page: number, perPage: number) {
  loading.value = true
  tableData.value = []
  fetch(`/api/accounting?page=${page}&perPage=${perPage}`)
    .then((response) => response.json())
    .then((data) => {
      tableData.value.push(...data.map((item: any) => item as Accounting))
      loading.value = false
    })
    .catch((error) => {
      console.error('Error:', error)
      loading.value = false
    })
}

function fetchCount() {
  fetch('/api/accounting/count')
    .then((response) => response.text())
    .then((data) => {
      total.value = parseInt(data)
    })
    .catch((error) => {
      console.error('Error:', error)
    })
}
</script>
<template>
  <el-table v-loading="loading" :data="tableData" stripe style="width: 100%; margin-bottom: 2em">
    <el-table-column prop="account_number" label="Number" width="100" />
    <el-table-column prop="account_name" label="Name" />
    <el-table-column prop="iban" label="IBAN" />
    <el-table-column prop="address" label="Address" />
    <el-table-column prop="amount" label="Amount" width="80" />
    <el-table-column prop="type" label="Type" width="100" />
  </el-table>
  <el-pagination
    background
    layout="prev, pager, next"
    :total="total"
    :current-page="page"
    :page-size="perPage"
    @update:current-page="handlePageUpdate"
    @update:page-size="handlePerPageUpdate"
  />
</template>
