<script setup lang="ts">
import { computed, Ref, ref } from "vue";

import { webGetCurTime } from "../api/req";
import {
  adjustCacheKey,
  baseURL,
  curTimeCacheKey,
  isAdjust,
  teachInfosCacheKey,
  validData,
} from "../api/globalConst";
const curDepartment = ref("");

const changeDepartment = (title) => {
  curDepartment.value = title;
};

const weekday = ref("星期日");
const lessonNum = ref(1);
const weekdayMap = [
  "星期日",
  "星期一",
  "星期二",
  "星期三",
  "星期四",
  "星期五",
  "星期六",
];

const lessonTime = computed(() => {
  switch (lessonNum.value) {
    case -2:
      return "早晨";
    case -3:
      return "中午";
    case -4:
      return "晚饭";
    case -5:
      return "夜晚";
    default:
      return `第 ${lessonNum.value} 节`;
  }
});
console.log(baseURL);
console.log("恭喜你发现彩蛋，交流群: 240728548");

// 使用缓存，而非默认值
const localTime: Ref<Array<number>> = ref(
  JSON.parse(localStorage.getItem(curTimeCacheKey))
);
if (localTime.value != null && localTime.value.length == 3) {
  weekday.value = weekdayMap[localTime.value[1]];
  lessonNum.value = localTime.value[2];
}

// 缓存2，避免加载慢时标题闪动
const adjustValue = localStorage.getItem(adjustCacheKey);
if (adjustValue != null && adjustValue === "true") {
  console.log(adjustValue);
  if (isAdjust.value !== true) {
    isAdjust.value = true;
  }
} else {
  if (isAdjust.value !== false) {
    isAdjust.value = false;
  }
}

webGetCurTime()
  .then((data) => {
    weekday.value = weekdayMap[data.weekday];
    lessonNum.value = data.lessonNum;
    if (isAdjust.value != data.isAdjust) {
      isAdjust.value = data.isAdjust;
    }
    if (!data.valid) {
      setInterval(() => {
        location.reload();
      }, 2000);
    }
    let curTime = localStorage.getItem(curTimeCacheKey);
    let t: Array<number> = JSON.parse(curTime);
    // console.log(t)
    // console.log(data)

    // 有缓存且缓存与当前时间匹配，则不处理
    if (
      curTime != null &&
      curTime != "" &&
      t[0] === data.weekNum &&
      t[1] === data.weekday &&
      t[2] === data.lessonNum
    ) {
      console.log("11");
      // 无需更新缓存
      return;
    }
    // 没有时间缓存或时间缓存和最新时间不匹配，则更新时间缓存

    // 有一个不匹配就更新缓存
    // set直接刷新

    // localStorage.removeItem(curTimeCacheKey)
    // 先把旧的数据缓存删掉，再更新时间缓存
    localStorage.removeItem(teachInfosCacheKey);

    localStorage.setItem(
      curTimeCacheKey,
      JSON.stringify([data.weekNum, data.weekday, data.lessonNum])
    );
    localStorage.setItem(adjustCacheKey, data.isAdjust ? "true" : "");

    // 数据标记为失效，后续请求更新
    validData.value = false;
  })
  .catch((err) => {
    console.log(err);
  });
</script>
<template>
  <!--  <div class="bg-gradient-to-tl from-cyan-200/50 to-blue-400/60 h-[100vh]">-->
  <div class="bg-[#ffffff] text-[#000001]">
    <div class="bg-[#dda15e]/50 h-[100%]">
      <div
        class="flex w-[100vw] h-[11vw] bg-[#dda15e] border-[0.5vw] border-[#bc6c25] place-items-center place-content-center mt">
        <h1 class="text-[5vw]">辅助听课工具{{ isAdjust ? "「调休」" : "" }}</h1>
      </div>

      <div
        class="flex h-[11vw] bg-[#dda15e]/50 place-items-center place-content-center rounded-3xl mt-[3.5vw] ml-[6vw] mr-[6vw]">
        <div class="text-[4.5vw]">{{ weekday }} {{ lessonTime }}</div>
      </div>


    </div>
    <!--    <div class="bg-[#dda15e]/50 h-[70vh]"></div>-->
  </div>

  <div v-if="curDepartment === ''" class="float-right mr-[4vw]">
    交流群: 240728548
  </div>
</template>
<style>
html,
body,
#app {
  width: 100%;
  height: 100%;
}
</style>
