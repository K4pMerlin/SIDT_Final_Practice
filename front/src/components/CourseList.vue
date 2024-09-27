<script setup lang="ts">

import CourseCard from "@/components/CourseCard.vue";
import {computed, onMounted} from "vue";
import {GlobalTeachInfosObj} from "@/store/teachInfosObj.ts";

const emits = defineProps(['curDepartment','curBuilding'])

const hasInfo = computed(()=>{
  return GlobalTeachInfosObj.getBuildings(emits.curDepartment).length!==0
})

const teachInfos = computed(() => {
  return GlobalTeachInfosObj.getTeachInfosByDepartmentAndBuilding(emits.curDepartment, emits.curBuilding)
})

onMounted(() => {
  console.log(GlobalTeachInfosObj.getBuildings(emits.curDepartment));
  console.log(emits.curDepartment)
  console.log(hasInfo.value)
})

const isApiSomeErr = computed(()=>{
  return GlobalTeachInfosObj.apiErrorMsg.value!==''
})


</script>

<template>

  <div class="bg-[#dda15e]/50 pt-[2vw]">
    <div v-if="isApiSomeErr">
      网络异常，请通过正确的域名访问~
      {{GlobalTeachInfosObj.apiErrorMsg}}
    </div>
    <div v-else-if="!hasInfo">

<!--      无课图片-->
      <div class="pb-[10vw]">
        <div class="flex place-content-center">
          <img class="w-[50vw]" src="/src/assets/desk3.png" alt="">
        </div>
        <div class="text-center text-[6vw]">该学部这个时间没有课~</div>
      </div>

    </div>
    <div v-else>
      <div>
        <div class="mt-[3vw] pb-[1vw]"
             v-for="teachInfo in teachInfos">
          <CourseCard :teach-info="teachInfo"/>
        </div>
      </div>
      <!---->
    </div>
  </div>

</template>

<style scoped>

</style>