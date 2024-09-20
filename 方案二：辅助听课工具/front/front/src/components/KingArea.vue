<script setup lang="ts">

import {computed, onMounted, ref, watch} from "vue";



// const {buildingIconColor} =defineProps(['buildingIconColor']); // ref("#FFAFCC")
const {buildingIconColor, choseBoxColor, department, allBuildings} =
    defineProps({
      buildingIconColor: String,
      choseBoxColor: String,
      department: String,
      allBuildings: Array
    })



// ref("rgb(255,69,137,0.1)")
onMounted(() => {

})

const buildings = ref([])

// console.log(allBuildings)
buildings.value = allBuildings
// emits('getCurBuildingWithDepartment', buildings.value[0], department)

const pageSize = computed(() => {
  return parseInt(String(Math.ceil(buildings.value.length / 5))) // 向上整数
})

const iconLib = ref([
  "bi bi-calendar-heart-fill",
  "bi bi-eraser-fill",
  "bi bi-mortarboard-fill",
  "bi bi-chat-left-text-fill",
  "bi bi-signpost-2-fill",
  "bi bi-slash-square-fill",
])

const curScrollIndex = ref(0)
const handleScroll = (e) => {
  // console.log(e.target);
  // console.log(parseInt(e.target.scrollLeft));

  // console.log(e.target.scrollWidth-e.target.offsetWidth);
  // 进度条过半
  // if (parseInt(e.target.scrollLeft) > e.target.offsetWidth / 2) {
  //   curScrollIndex.value = 1
  // } else {
  //   curScrollIndex.value = 0
  // }

  // console.log(e.target.scrollLeft)
  // const totalScrollLen = e.target.scrollWidth-e.target.offsetWidth
  for (let i = 0; i < pageSize.value; i++) {
    if (parseInt(e.target.scrollLeft) <=
        e.target.scrollWidth / pageSize.value * (i + 1) - e.target.offsetWidth / 2) {
      // e.target.scrollLeft = e.target.scrollWidth * i / pageSize.value
      curScrollIndex.value = i
      return
    }
  }
  // 0-size/3*1  size/3*1-size/3*2 size/3*2-size/3*3


}


const scrollContent = ref()
const onClickProcess = (index: number) => {
  // console.log(num);
  // console.log(pageSize.value);

  // console.log(totalScrollLen);
  scrollContent.value.scrollLeft = scrollContent.value.scrollWidth * index / pageSize.value

}

// const curScrollIndex = computed(()=>{
//     for (let i = 0; i < pageSize.value; i++) {
//     if (parseInt(scrollContent.value.scrollLeft) <=
//         scrollContent.value.scrollWidth / pageSize.value * (i + 1) - scrollContent.value.offsetWidth / 2) {
//       scrollContent.value.scrollLeft = scrollContent.value.scrollWidth * i / pageSize.value
//       console.log(i);
//       return i
//     }
//   }
//     return 0
// })


const partOfBuildings = (n: number): string[] => {
  const len = buildings.value.length
  // console.log(len/2);
  return buildings.value.slice(5 * n, Math.min(len, 5 * (n + 1)))
}

const emits = defineEmits(['getCurBuildingWithDepartment']);

const buildingIndex = ref(0)
const onClickIcon = (index: number) => {
  emits('getCurBuildingWithDepartment', buildings.value[index], department)
  buildingIndex.value = index
}


</script>

<template>
  <div class="relative rounded-xl overflow-auto pb-[2vw] opacity-100">
    <div class="relative w-full flex gap-0 snap-x snap-mandatory overflow-x-auto scroll-smooth" ref="scrollContent"
         @scroll="handleScroll">
      <!--      <div class="snap-center shrink-0">-->
      <!--        <div class="shrink-0 w-4 sm:w-48"></div>-->
      <!--      </div>-->
      <!--      snap-always snap-center-->
      <div class="snap-always snap-center shrink-0 first:pl-0 last:pr-0"
           v-for="i in pageSize">

        <div class="grid grid-cols-5 pt-[5vw] pb-[5.2vw] w-[100vw] border-l-[0.5vw] border-r-[0.5vw]">
          <div v-for="(building,index) in partOfBuildings(i-1)" class="text-center transition-all">
            <div :style="{'border-color': (i-1)*5+index===buildingIndex?buildingIconColor:'rgba(255, 255, 255, 0.0)',
            'border-width':'0.7vw',
            'border-radius':'2vw',
            'background-color':(i-1)*5+index===buildingIndex?choseBoxColor:'rgba(255, 255, 255, 0.0)'}"
                 @click="onClickIcon((i-1)*5+index)">
              <div class="text-[7vw] p-[0.5vw]"><i :class="iconLib[index]" :style="{'color':buildingIconColor}"></i>
              </div>
              <div class="text-[3.3vw] p-[0.5vw] whitespace-nowrap">{{ building }}</div>
            </div>
          </div>
        </div>
      </div>
      <!--      <div class="snap-center shrink-0">-->
      <!--        <div class="shrink-0 w-4 sm:w-48"></div>-->
      <!--      </div>-->
    </div>
  </div>

  <!--   进度条-->
  <div class="flex flex-row place-content-center pb-[4.3vw]" v-if="pageSize!==1">
    <div v-for="i in pageSize">
      <div :class="{
       'w-[20vw] h-[1.4vw] rounded ml-[1vw] mr-[1vw]':true,
       'bg-blue-500':i-1===curScrollIndex,
       'bg-blue-200':i-1!==curScrollIndex
     }" @click="onClickProcess(i-1)"></div>
    </div>
  </div>


</template>

<style scoped>

</style>