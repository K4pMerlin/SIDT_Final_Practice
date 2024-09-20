<script setup lang="ts">

import {computed, PropType, ref} from "vue";
import {Items} from "../../types/Items";


//defineProps(['infoCardColor','building','department'])
// const infoCardColor = ref('#A2D2FF')
const {infoCardColor, teachInfo} =
    defineProps({
      infoCardColor: String,
      teachInfo: Object as PropType<Items.TeachInfo>,
    })


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
    <div class="grid grid-cols-2 gap-44
    justify-between place-items-center whitespace-nowrap
 pt-[1vh]">
      <div class="text-[4.5vw]">{{ teachInfo.room }}</div>
      <div class="border-2 rounded-2xl pl-[1.7vw] pr-[1.7vw] pt-[0.5vw] pb-[0.5vw] text-[3.5vw]"
           :style="{'background-color':bgColors[teachInfo.courseType==='通识课'?1:0],'border-color':borderColors[teachInfo.courseType==='通识课'?1:0],'color':fontColors[teachInfo.courseType==='通识课'?1:0]}">
        {{ teachInfo.courseType }}
      </div>
    </div>
    <div class="grid grid-rows-1
    justify-center place-content-center place-items-center ml-[5vw] mr-[5vw] text-center
text-[5.5vw]">
      {{ teachInfo.courseName }}
    </div>
    <div class="grid grid-cols-3 gap-1
    justify-between whitespace-nowrap
    text-[3.2vw] pb-[1.1vh] pt-[1vh]">
      <div class="text-left pl-[2vw]">「{{ teachInfo.faculty }}」</div>
      <div class="text-center">
        <ul>
          <li>
            {{ teachInfo.teacherName }} {{ teachInfo.teacherTitle }}
          </li>
        </ul>

      </div>
      <div class="text-right pr-[2vw]">起止时间: {{ teachInfo.courseTime }}</div>
    </div>
  </div>
</template>

<style scoped>

</style>