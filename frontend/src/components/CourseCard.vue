<script setup lang="ts">

import type { PropType } from 'vue'
import { ref } from 'vue'
import { Items } from "@/types/Items";
const infoCardColor = ref(
    // '#606c38'
    'rgba(96,108,56,0.5)'
)
const {teachInfo} =
    defineProps({
      teachInfo: Object as PropType<Items.TeachInfo>,
    })

    const navigateToRoom = (room: string, department: string, building: string) => {
  if (room) {
    const encodedRoom = encodeURIComponent(room);
    const encodedDepartment = encodeURIComponent(department);
    const encodedBuilding = encodeURIComponent(building);
    const url = `https://ditu.amap.com/search?query=武汉大学${encodedDepartment}${encodedBuilding}${encodedRoom}`;
    window.open(url, '_blank');
  } else {
    console.error('Room information is not available');
  }
};
const bgColors = ref(['rgba(159,75,209,0.5)', 'rgba(60,110,113,0.5)'])
const borderColors = ref(['#AA5ED7', '#8ad96f'])
const fontColors = ref(['#451344', '#052a05'])

// const colorIndex = computed(() => {
//   switch (teachInfo.courseType) {
//     case '专业课':
//       return 0
//     case '通识课':
//       return 1
//     case '英语课':
//       // todo
//       console.log('todo!!!!')
//       break
//   }
//   return 0
// })


</script>

<template>
  <div class="mb-[4vw] rounded-2xl ml-[3vw] mr-[3vw]"
       :style="{'background-color':infoCardColor}">
    <div class="grid grid-cols-3 gap-2 items-center whitespace-nowrap pt-[1vh] px-[2vw]">
      <div class="text-[3.5vw] truncate">{{ teachInfo?.room }}</div>
      <div class="text-[2.5vw] text-center">
        <a href="#" @click="navigateToRoom(teachInfo?.room ?? '', teachInfo?.faculty ?? '','')" class="text-blue-500 hover:text-blue-700">教室在哪</a>
      </div>
      <div class="text-[2.5vw] flex justify-end">
        <span class="border rounded-full px-[1.5vw] py-[0.3vw]"
              :style="{'background-color':bgColors[teachInfo?.courseType==='通识课'?1:0],'border-color':borderColors[teachInfo?.courseType==='通识课'?1:0],'color':fontColors[teachInfo?.courseType==='通识课'?1:0]}">
          {{ teachInfo?.courseType }}
        </span>
      </div>
    </div>
    <div class="grid grid-rows-1
    justify-center place-content-center place-items-center ml-[5vw] mr-[5vw] text-center
    text-[5.5vw]">
      {{ teachInfo?.courseName }}
    </div>
    <div class="grid grid-cols-3 gap-1
    justify-between whitespace-nowrap
    text-[3.2vw] pb-[1.1vh] pt-[1vh]">
      <div class="text-left pl-[2vw]">「{{ teachInfo?.faculty }}」</div>
      <div class="text-center">
        <ul>
          <li>
            {{ teachInfo?.teacherName }} {{ teachInfo?.teacherTitle }}
          </li>
        </ul>

      </div>
      <div class="text-right pr-[2vw]">起止时间: {{ teachInfo?.courseTime }}</div>
    </div>
  </div>
</template>

<style scoped>

</style>