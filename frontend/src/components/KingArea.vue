<script setup lang="ts">

import {computed, onMounted, ref} from "vue";

// 样式
const buildingIconColor = '#606c38'

const choseBoxColor = 'rgb(255,69,137,0.1)'

const iconLib = ref([
  "bi bi-calendar-heart-fill",
  "bi bi-eraser-fill",
  "bi bi-mortarboard-fill",
  "bi bi-chat-left-text-fill",
  "bi bi-signpost-2-fill",
  "bi bi-slash-square-fill",
])

// 数据
const {allBuildings} =
    defineProps({
      allBuildings: Array
    })

const curBuilding = defineModel()

const buildings = ref<any[]>([])
onMounted(() => {
  if (allBuildings) {
    buildings.value = allBuildings
  }
})


const pageSize = computed(() => {
  return parseInt(String(Math.ceil(buildings.value.length / 5))) // 向上整数
})

// 动作
const curScrollIndex = ref(0)
const handleScroll = (e: Event) => {
  const target = e.target as HTMLElement;
  for (let i = 0; i < pageSize.value; i++) {
    if (parseInt(target.scrollLeft.toString()) <=
        target.scrollWidth / pageSize.value * (i + 1) - target.offsetWidth / 2) {
      curScrollIndex.value = i
      return
    }
  }
  // 0-size/3*1  size/3*1-size/3*2 size/3*2-size/3*3
}


const scrollContent = ref()
const onClickProcess = (index: number) => {
  scrollContent.value.scrollLeft = scrollContent.value.scrollWidth * index / pageSize.value
}

/// 部分教学楼，返回分页后的数据，n为页号，从0开始
const partOfBuildings = (n: number): string[] => {
  const len = buildings.value.length
  return buildings.value.slice(5 * n, Math.min(len, 5 * (n + 1)))
}

// 返回当前是哪个教学楼，便于展示数据（每次点击图标触发）
// const emits = defineEmits(['clickIcon']);

const buildingIndex = ref(0)
const onClickIcon = (index: number) => {
  // emits('clickIcon', buildings.value[index])
  console.log(curBuilding.value)
  buildingIndex.value = index
  curBuilding.value  = buildings.value[index]
}


</script>

<template>
  <div class="relative rounded-xl overflow-auto pb-[2vw] opacity-100">
    <div class="relative w-full flex gap-0 snap-x snap-mandatory overflow-x-auto scroll-smooth"
         ref="scrollContent"
         @scroll="handleScroll">
      <div class="snap-always snap-center shrink-0 first:pl-0 last:pr-0"
           v-for="i in pageSize">

        <div class="grid grid-cols-5 pt-[5vw] pb-[5.2vw] w-[100vw] border-l-[0.5vw] border-r-[0.5vw]">
          <div v-for="(building,index) in partOfBuildings(i-1)" class="text-center transition-all">
            <div :style="{
              'border-color': (i-1)*5+index===buildingIndex?buildingIconColor:'rgba(255, 255, 255, 0.0)',
            'border-width':'0.7vw',
            'border-radius':'2vw',
            'background-color':(i-1)*5+index===buildingIndex?choseBoxColor:'rgba(255, 255, 255, 0.0)'}"
                 @click="onClickIcon((i-1)*5+index)">
              <div class="text-[7vw] p-[0.5vw]">
                <i :class="iconLib[index]" :style="{'color':buildingIconColor}"></i>
              </div>
              <div class="text-[3.3vw] p-[0.5vw] whitespace-nowrap">{{ building }}</div>
            </div>
          </div>
        </div>
      </div>
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