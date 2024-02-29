<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>
<template>
  <div class="common-layout">
    <el-container>
      <el-header>
        <el-menu class="el-menu-demo" mode="horizontal" :ellipsis="false" :router="true">
          <el-menu-item route="/"> Codenotary </el-menu-item>
          <div style="flex-grow: 1" />
          <el-menu-item index="home" route="/"> Accounting </el-menu-item>
          <el-menu-item index="add" route="/add"> Add </el-menu-item>
        </el-menu>
      </el-header>
      <el-main>
        <RouterView />
      </el-main>
    </el-container>
  </div>
</template>
