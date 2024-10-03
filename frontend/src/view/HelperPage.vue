<script setup lang="ts">

import HelperCollapse from "@/components/HelperCollapse.vue";
import {watch} from "vue";
import {GlobalTeachInfosObj} from "@/store/teachInfosObj.ts";
import {GlobalCacheObj} from "../store/cacheObj.ts";
import {validData} from "@/store/globalData.ts";

GlobalCacheObj.loadData()
watch(validData, (newValue) => {
  if (newValue == false) {
    // console.log(validData.value)
    // 第二次请求
    GlobalTeachInfosObj.loadGlobalTeachInfos(false)
    validData.value = true
  }
})
</script>

<template>
  <div class="bg-[#ffffff] text-[#000001]">
    <div class="bg-[#dda15e]/50 h-[100%]">
      <div
          class="flex w-[100vw] h-[11vw] bg-[#dda15e] border-[0.5vw] border-[#bc6c25] place-items-center place-content-center mt"
      >
        <h1 class="text-[5vw]">辅助听课系统</h1>
      </div>

      <div
          class="flex h-[11vw] bg-[#dda15e]/50 place-items-center place-content-center rounded-3xl mt-[3.5vw] ml-[6vw] mr-[6vw] mb-[3.5vw]"
      >
        <div class="text-[4.5vw]">{{ GlobalCacheObj.weekday }} {{ GlobalCacheObj.lessonTime }}</div>
      </div>

      <HelperCollapse/>
    </div>
  </div>

  <div class="float-right mr-[4vw]">
    交流群: 813278961
  </div>

  <div>
    <n-button type="warning">新内容</n-button>
  </div>

</template>

<style scoped>

</style>