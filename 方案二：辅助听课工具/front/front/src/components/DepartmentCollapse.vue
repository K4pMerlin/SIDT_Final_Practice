<script setup lang="ts">

import {computed, onMounted, ref, watch} from "vue";
import KingArea from "./KingArea.vue";
import CourseCard from "./CourseCard.vue";
import {hasTeachInfo, isApiError, teachInfosCacheKey, validData} from "../../api/globalConst";
import {teachInfos} from "../../types/globalData";
import {webGetTeachInfos} from "../../api/req";
import {Items} from "../../types/Items";

const emits = defineEmits(['changeDepartment']);

const departments = ref([
  "文理学部",
  "工学部",
  "信息学部",
  "医学部",
  "网安基地"
]);
const departCopy: any = ref([])
onMounted(() => {
  departCopy.value = departments.value
})

'#d5d8dc'

'#c48833'

// 配色方案
// const collapseHeaderBgColor = ref(['bg-[#CDB4DB] rounded-l-[2vw]', 'bg-[#d6ccc2] rounded-l-[2vw]', 'bg-[#77A1B8] rounded-l-[2vw]',
//   'bg-[#344e41]/70 rounded-l-[2vw]', 'bg-[#d5d8dc] rounded-l-[2vw]']) // #dee2e6

const collapseHeaderBgColor = ref(
    'bg-[#dda15e]'
) // #dee2e6


// const collapseKingBgColor = ref(['bg-[#FFC8DD]/50', 'bg-[#f5ebe0]/70',
//   'bg-[#dbfaff]/25', 'bg-[#a3b18a]/70', 'bg-[#f8f9fa]'])
const collapseKingBgColor = ref(
    'bg-[#fefae0]'
)


const collapseBuildingIconColor = ref(
    '#606c38'
)
const collapseChoseBoxColor = ref(
    'rgb(255,69,137,0.1)'
)
const collapseContentBgColor = ref(
    'bg-[#dda15e]/50'
) // adb5bd
const infoCardColor = ref(
    // '#606c38'
    'rgba(96,108,56,0.5)'
)
'#dbfaff'

const searchText = ref('')
const copyInfo = ref(new Map())

const curDepartment = ref(departments.value[0]);
emits('changeDepartment', curDepartment.value)

const curBuildingDepartmentMap = ref(new Map<string, string>())

const getCurBuildingWithDepartment = (building, department) => {
  curBuildingDepartmentMap.value.set(department, building)
  // console.log(curBuildingDepartmentMap.value);
}


const buildingsMap = ref(new Map<string, string[]>())

const curBuildings = computed(() => {
  return (department) => {
    return buildingsMap.value.get(department)
  }
})


const getBuildings = () => {
  for (const department of departments.value) {
    for (let building of teachInfos.value.get(department).keys()) {
      buildingsMap.value.get(department).push(building)
    }

    // 后端已排序
    // 原地排序
    // buildingsMap.value.get(department).sort();
  }
}


// 场景：父组件和子组件在页面展示内容，
// 按照正常的生命周期顺序是 =》父组件创建=》子组件创建=》子组件挂载=》父组件挂载

// hasTeachInfo.value.set("文理学部",true)

// 这里是创建父组件时执行的逻辑，先初始化buildingsMap，使子组件得以有效
for (const department of departments.value) {
  buildingsMap.value.set(department, [])
}

const initTeachInfos = (data: Items.BuildingTeachInfos[][]) => {
  for (let i = 0; i < 5; i++) {
    // console.log(data[i])

    let tempMap = new Map<string, Items.TeachInfo[]>()

    // 是否无课
    hasTeachInfo.value.set(departments.value[i], data[i].length !== 0)

    // 初始化当前选中的教学楼
    curBuildingDepartmentMap.value.set(departments.value[i],
        data[i].length !== 0 ? data[i][0].building : '')

    // console.log(hasTeachInfo.value.get(departments.value[i]));
    data[i].forEach(t => {
      tempMap.set(t.building, t.infos)
    })
    teachInfos.value.set(departments.value[i], tempMap)
    copyInfo.value.set(departments.value[i], tempMap)
  }
  getBuildings()
}

const reqTeachInfos = () => {
  // 重置buildingsMap
  for (const department of departments.value) {
    buildingsMap.value.set(department, [])
  }
  webGetTeachInfos()
      .then((data) => {
        isApiError.value = false
        if (data.length != 5) {
          console.log("断言失败，数据长度应当为5！")
          return
        }

        localStorage.setItem(teachInfosCacheKey, JSON.stringify(data))
        console.log("refresh")

        initTeachInfos(data)

        // console.log(buildingsMap.value);
      })
      .catch((err) => {
        isApiError.value = true
        console.log(err)
      })
}

let data = localStorage.getItem(teachInfosCacheKey)
if (data != null && data != '') {

  initTeachInfos(JSON.parse(data))
} else {
  // 没有缓存，网络获取信息
  reqTeachInfos()
}

const searchFun = () => {
  // 
  if(!searchText.value) {
    departments.value = departCopy.value
    teachInfos.value = copyInfo.value
  }
  let departBool = departCopy.value.some(item => item.indexOf(searchText.value) !== -1)
  if(departBool) {
    // 查询到部门
    departments.value = departCopy.value.filter(item => item.indexOf(searchText.value) > -1)
  } else {
     // 未查到 查询时间
    let list = []
    let mapReset = new Map()
    copyInfo.value.forEach((item, keys) => {
      let floors: any = copyInfo.value.get(keys)
      let child = new Map()
      for (let [key ,value] of floors) {
        list = list.concat(value)
        let filterList = value.filter((item: any) => item.courseTime.indexOf(searchText.value) !== -1)
        child.set(key, filterList)
      }
      mapReset.set(keys, child)
    })
    let timeBool = list.some((item: any) => item.courseTime.indexOf(searchText.value) !== -1)
    if(timeBool) {
      // 查询到时间
      teachInfos.value = mapReset
      departments.value = departCopy.value
    } else {
      // 时间与部门都未查询到
      departments.value = []    
    }
  }
}
const isSubmit = (e) => {
  e.preventDefault(); 
}

watch(validData, (newValue) => {
  if (newValue == false) {
    // console.log(validData.value)
    // 第二次请求
    reqTeachInfos()
    validData.value = true
  }
})


</script>

<template>
  <div class="searchs">
    <form action="" @submit="isSubmit($event)">
      <input placeholder="请搜索部门、时间(第x-x节)" type="search" v-model="searchText" @change="searchFun" class="search-row"/>
    </form>
  </div>
  <div class="w-[100vw] mt-[4.5vw]">
    <div class="hs-accordion-group">
      <ul v-for="(department,index) in departments">
        <li :class="{'hs-accordion':true,'active':index===0}" :id="department+'-heading'">

          <div class="border-2 border-[#bc6c25] mb-[1vw]">

            <!--          触发器样式容器-->
            <div :class='collapseHeaderBgColor'>
              <!--            触发器-->
              <div class="hs-accordion-toggle whitespace-nowrap p-[3vw]"
                   :aria-controls="department+'-collapse'"
                   @click="curDepartment===department?curDepartment='':curDepartment=department;
             emits('changeDepartment',curDepartment)">
                <div class=""> <!-- 容器 -->
                  <button v-if="department==='文理学部'" class="mr-[1vw]">
                    <svg t="1709564852350" class="icon" viewBox="0 0 1152 1024" version="1.1"
                         xmlns="http://www.w3.org/2000/svg" p-id="14572" width="5vw" height="5vw">
                      <path
                          d="M576 800.548571a34.742857 34.742857 0 0 1-16.457143-3.291428L217.234286 655.140571a43.008 43.008 0 0 1-27.940572-40.740571V43.885714a43.885714 43.885714 0 0 1 19.748572-36.571428 44.617143 44.617143 0 0 1 41.910857-3.657143L438.857143 81.627429a41.910857 41.910857 0 0 1 22.454857 23.332571 43.885714 43.885714 0 0 1-56.173714 57.270857L277.942857 109.714286V585.142857l253.659429 105.033143v-504.685714a43.081143 43.081143 0 0 1 28.013714-40.301715L901.046857 3.657143a47.177143 47.177143 0 0 1 41.984 3.657143 43.885714 43.885714 0 0 1 19.748572 36.571428v570.514286a42.569143 42.569143 0 0 1-28.013715 40.301714l-186.806857 76.8a43.885714 43.885714 0 0 1-56.246857-57.197714 41.910857 41.910857 0 0 1 22.528-23.332571l159.670857-66.56v-475.428572L620.470857 214.381714v541.257143a45.641143 45.641143 0 0 1-44.470857 44.690286z m0 0V1024a49.956571 49.956571 0 0 1-15.652571-2.56l-533.211429-203.995429A43.885714 43.885714 0 0 1 0 776.996571V265.581714a43.885714 43.885714 0 0 1 42.788571-44.397714 43.885714 43.885714 0 0 1 42.788572 44.397714v481.206857l490.057143 187.757715 490.057143-187.757715V299.885714a42.788571 42.788571 0 0 1 85.577142 0v476.598857a43.885714 43.885714 0 0 1-27.136 40.521143l-532.48 204.068572a43.885714 43.885714 0 0 1-15.652571 2.925714z m0 0"
                          fill="#333333" p-id="14573"></path>
                    </svg>
                  </button>


                  <button v-else-if="department==='工学部'" class="mr-[1vw]">
                    <svg t="1709564738216" class="icon" viewBox="0 0 1024 1024" version="1.1"
                         xmlns="http://www.w3.org/2000/svg" p-id="10570" width="5vw" height="5vw">
                      <path
                          d="M377.856 1005.056c-3.072 0-6.144-0.512-8.704-1.024-78.336-23.04-151.552-65.024-210.944-121.856-9.728-9.216-12.8-23.552-7.168-35.84 17.408-39.936 14.848-83.968-6.656-121.344-21.504-37.376-58.88-61.44-101.888-66.56-13.312-1.536-24.576-11.264-27.648-24.064C5.12 593.92 0 552.96 0 512c0-40.96 5.12-81.92 14.336-121.344 3.072-12.8 14.336-22.528 27.648-24.064 43.008-5.12 80.384-29.184 101.888-66.56 21.504-37.376 24.064-81.92 6.656-121.344-5.12-12.288-2.56-26.624 7.168-35.84C217.088 85.504 289.792 43.008 368.64 20.48c12.8-3.584 26.624 1.024 34.816 11.776C429.056 67.072 468.992 87.04 512 87.04s82.944-19.968 108.544-54.784c8.192-10.752 22.016-15.36 34.816-11.776 78.336 23.04 151.552 65.024 210.944 121.856 9.728 9.216 12.8 23.552 7.168 35.84-17.408 39.936-14.848 83.968 6.656 121.344 21.504 37.376 58.88 61.952 101.888 66.56 13.312 1.536 24.576 11.264 27.648 24.064C1018.88 430.08 1024 471.04 1024 512s-5.12 81.92-14.336 121.344c-3.072 12.8-14.336 22.528-27.648 24.064-43.008 5.12-80.384 29.184-101.888 66.56-21.504 37.376-24.064 81.92-6.656 121.344 5.12 12.288 2.56 26.624-7.168 35.84-59.392 56.832-132.096 98.816-210.944 121.856-12.8 3.584-26.624-1.024-34.816-11.776-25.6-34.816-65.536-54.784-108.544-54.784-43.008 0-82.944 19.968-108.544 54.784-6.144 9.216-15.36 13.824-25.6 13.824zM217.6 849.92c43.52 37.888 94.72 67.584 149.504 86.528 37.376-39.936 90.112-62.976 144.896-62.976s107.52 23.04 144.896 62.976c54.784-18.944 105.472-48.128 149.504-86.528-15.872-52.224-9.216-109.568 17.92-157.184 27.648-47.616 73.728-81.92 126.976-94.208 5.632-28.16 8.192-57.344 8.192-86.016 0-29.184-2.56-57.856-8.192-86.528-53.248-12.288-99.328-46.592-126.976-94.208-27.648-47.616-33.792-104.96-17.92-157.184-43.52-37.888-94.72-67.584-149.504-86.528C619.52 128 566.784 151.04 512 151.04S404.48 128 367.104 88.064C312.32 107.008 261.12 136.192 217.088 174.592c15.872 52.224 9.728 109.568-17.92 157.184s-73.728 81.92-126.976 94.208c-5.632 28.16-8.192 57.344-8.192 86.016 0 29.184 3.072 58.368 8.192 86.528 53.248 12.288 99.328 46.592 126.976 94.208 27.648 47.616 34.304 104.96 18.432 157.184z m634.88-141.312z"
                          fill="#333333" p-id="10571"></path>
                      <path
                          d="M512 151.04c-62.976 0-122.88-30.208-160.256-80.896-10.24-14.336-7.68-34.304 6.656-44.544 14.336-10.24 34.304-7.68 44.544 6.656C429.056 67.072 468.992 87.04 512 87.04s82.944-19.968 108.544-54.784c10.24-14.336 30.72-17.408 44.544-6.656 14.336 10.24 17.408 30.72 6.656 44.544-36.864 50.688-96.768 80.896-159.744 80.896zM978.432 429.568h-3.584c-62.464-7.168-118.784-43.52-150.016-98.304-31.232-54.272-35.328-121.344-10.24-179.2 7.168-16.384 26.112-23.552 41.984-16.384 16.384 7.168 23.552 26.112 16.384 41.984-17.408 39.936-14.848 83.968 6.656 121.344 21.504 37.376 58.88 61.952 101.888 66.56 17.408 2.048 30.208 17.92 28.16 35.328-1.536 16.896-15.36 28.672-31.232 28.672zM844.288 890.368c-12.288 0-24.064-7.168-29.184-19.456-25.088-57.856-21.504-124.416 9.728-179.2 31.232-54.272 87.552-91.136 150.016-98.304 17.408-2.048 33.28 10.752 35.328 28.16 2.048 17.408-10.752 33.28-28.16 35.328-43.008 5.12-80.384 29.184-101.888 66.56-21.504 37.376-24.064 81.92-6.656 121.344 7.168 16.384-0.512 34.816-16.384 41.984-4.608 3.072-8.704 3.584-12.8 3.584zM377.856 1005.056c-6.656 0-13.312-2.048-18.944-6.144-14.336-10.24-17.408-30.72-6.656-44.544 37.376-50.688 97.28-80.896 160.256-80.896s122.88 30.208 160.256 80.896c10.752 14.336 7.68 34.304-6.656 44.544-14.336 10.752-34.304 7.68-44.544-6.656-25.6-34.816-65.536-54.784-108.544-54.784-43.008 0-82.944 19.968-108.544 54.784-7.168 8.192-16.896 12.8-26.624 12.8zM180.224 890.88c-4.096 0-8.704-1.024-12.8-2.56-16.384-7.168-23.552-26.112-16.384-41.984 17.408-39.936 14.848-83.968-6.656-121.344-21.504-37.376-58.88-61.44-101.888-66.56-17.408-2.048-30.208-17.92-28.16-35.328 2.048-17.408 17.92-30.208 35.328-28.16 62.464 7.168 118.784 43.52 150.016 98.304 31.232 54.272 35.328 121.344 10.24 179.2-5.632 11.264-17.408 18.432-29.696 18.432zM45.568 430.08c-15.872 0-29.696-12.288-31.744-28.16-2.048-17.92 10.752-33.28 28.16-35.328 43.008-5.12 80.384-29.184 101.888-66.56 21.504-37.376 24.064-81.92 6.656-121.344-7.168-16.384 0.512-34.816 16.384-41.984 16.384-7.168 34.816 0.512 41.984 16.384 25.088 57.856 21.504 124.416-9.728 179.2C167.936 386.048 111.616 422.912 49.152 430.08h-3.584zM512 680.96c-93.184 0-168.96-75.776-168.96-168.96s75.776-168.96 168.96-168.96 168.96 75.776 168.96 168.96-75.776 168.96-168.96 168.96z m0-273.92c-57.856 0-104.96 47.104-104.96 104.96s47.104 104.96 104.96 104.96 104.96-47.104 104.96-104.96-47.104-104.96-104.96-104.96z"
                          fill="#333333" p-id="10572"></path>
                    </svg>
                  </button>


                  <button v-else-if="department==='信息学部'" class="mr-[1vw]">
                    <svg t="1709563353091" class="icon" viewBox="0 0 1024 1024" version="1.1"
                         xmlns="http://www.w3.org/2000/svg" p-id="1967" width="5vw" height="5vw">
                      <path
                          d="M85.333333 213.333333a85.333333 85.333333 0 0 1 85.333334-85.333333h682.666666a85.333333 85.333333 0 0 1 85.333334 85.333333v469.333334a85.333333 85.333333 0 0 1-85.333334 85.333333H170.666667a85.333333 85.333333 0 0 1-85.333334-85.333333V213.333333z m768 0H170.666667v469.333334h682.666666V213.333333zM938.666667 853.333333a42.666667 42.666667 0 0 1-42.666667 42.666667H128a42.666667 42.666667 0 1 1 0-85.333333h768a42.666667 42.666667 0 0 1 42.666667 42.666666z"
                          fill="#333333" p-id="1968"></path>
                    </svg>
                  </button>


                  <button v-else-if="department==='医学部'" class="mr-[1vw]">
                    <svg t="1709565117364" class="icon" viewBox="0 0 1024 1024" version="1.1"
                         xmlns="http://www.w3.org/2000/svg" p-id="3339" width="5vw" height="5vw">
                      <path
                          d="M802.133 426.667H233.244c-17.066 0-31.857 12.515-31.857 28.444s14.79 28.445 31.857 28.445h568.89c17.066 0 31.857-12.516 31.857-28.445s-14.791-28.444-31.858-28.444zM604.16 667.876h-58.027V608.71c0-15.929-12.515-28.444-28.444-28.444s-28.445 13.653-28.445 28.444v59.165H430.08c-15.929 0-28.444 12.515-28.444 28.444s12.515 28.444 28.444 28.444h59.164v59.165c0 15.929 12.516 28.444 28.445 28.444s28.444-12.515 28.444-28.444v-59.165h59.165c15.929 0 28.444-12.515 28.444-28.444s-13.653-28.444-29.582-28.444z"
                          fill="#2C2C2C" p-id="3340"></path>
                      <path
                          d="M863.573 227.556h-42.097l-22.756-95.574c-13.653-58.026-64.853-98.986-124.018-98.986H360.676c-59.165 0-111.503 40.96-124.018 98.986l-22.756 95.574h-42.098c-63.715 0-114.915 51.2-114.915 114.915V839.68c0 63.716 51.2 114.916 114.915 114.916h691.77c63.715 0 114.915-51.2 114.915-114.916V342.471c0-63.715-51.2-114.915-114.916-114.915zM312.89 149.049c5.689-22.756 25.031-38.685 47.787-38.685h314.026c22.756 0 43.236 15.93 47.787 38.685l18.204 77.369H294.684l18.205-77.37z m593.92 690.631c0 23.893-19.342 43.236-43.236 43.236H171.804c-23.893 0-43.235-19.343-43.235-43.236V342.471c0-23.893 19.342-43.235 43.235-43.235h691.77c23.893 0 43.235 19.342 43.235 43.235V839.68z"
                          fill="#333333" p-id="3341"></path>
                    </svg>
                  </button>


                  <button v-else-if="department==='网安基地'" class="mr-[1vw]">
                    <svg t="1709564674791" class="icon" viewBox="0 0 1024 1024" version="1.1"
                         xmlns="http://www.w3.org/2000/svg" p-id="4215" width="5vw" height="5vw">
                      <path
                          d="M512 825.376c0 31.84 13.632 64.32 34.816 95.264-26.88 4.8-54.528 7.36-82.816 7.36C223.968 928 26.496 745.728 2.464 512H0v-64C8.704 199.136 213.12 0 464 0c247.168 0 449.216 193.28 463.232 436.96a271.808 271.808 0 0 1-8.512-3.232 698.336 698.336 0 0 1-59.328-25.056 395.84 395.84 0 0 0-60.8-162.592l-6.656 4.32a728 728 0 0 1-121.824 62.656c7.904 34.176 13.184 69.76 15.872 106.656-10.56 4.608-22.016 9.344-34.176 14.016-8.256 2.848-16.096 5.376-23.488 7.68a630.368 630.368 0 0 0-14.528-109.312 518.88 518.88 0 0 1-101.76 18.432L512 825.376z m-229.184-108l-4.224 1.6c-27.2 10.56-55.2 23.584-84.256 39.36a397.504 397.504 0 0 0 170.08 92.384 475.104 475.104 0 0 1-81.6-133.344z m166.08-35.072a517.12 517.12 0 0 0-109.696 16.864 408.544 408.544 0 0 0 109.088 151.008l0.64-167.872zM67.648 512a397.216 397.216 0 0 0 78.528 193.6 895.36 895.36 0 0 1 114.72-54.624C249.6 608 242.88 561.696 240.736 512.032h-173.12z m381.92 0l-149.952 0.032c2.016 42.88 7.68 82.816 17.12 119.936a566.784 566.784 0 0 1 132.48-23.36L449.536 512z m-134.24-184.32l-1.376 5.76A637.28 637.28 0 0 0 299.36 448h150.4l0.384-96.576a569.504 569.504 0 0 1-134.848-23.744zM130.048 245.184l-2.784 4.224A397.12 397.12 0 0 0 65.056 448h175.456c1.728-48.832 7.904-95.424 18.496-139.712a918.176 918.176 0 0 1-128.96-63.104z m321.088-166.624c-50.432 52.416-88.256 111.232-113.6 176.768l-0.928 2.432c39.072 11.264 77.024 17.856 113.824 19.84l0.704-199.04z m60.928 40.192L512 276.576a463.808 463.808 0 0 0 80.96-14.56 519.296 519.296 0 0 0-80.896-143.264z m44.352-43.264l1.44 1.856A587.296 587.296 0 0 1 649.6 243.04c34.56-13.664 69.44-31.36 104.544-53.152a398.048 398.048 0 0 0-197.728-114.4z m-184.832 0.032l-3.072 0.704a398.848 398.848 0 0 0-193.376 112.16c35.84 20.16 70.784 36.832 104.96 50.176a586.144 586.144 0 0 1 91.52-163.04z"
                          fill="#444444" p-id="4216"></path>
                      <path
                          d="M822.208 465.44L789.12 448l-30.784 18.56-4.896 2.976a455.456 455.456 0 0 1-53.248 26.464l-9.92 3.968c-13.12 4.608-25.6 8.224-36.608 10.816l-8 1.76a151.808 151.808 0 0 1-11.776 2.016L576 519.264v272.48c0 73.504 63.296 141.696 184.192 209.088l30.208 16.416 19.04-10.304c128.192-69.376 195.36-139.488 195.36-215.2v-272.224l-54.592-4.8-3.52-0.16c-7.488 0-23.84-2.88-47.392-11.808l-17.696-6.912c-12.576-5.12-24.064-10.496-34.4-15.904l-7.52-4.064a233.056 233.056 0 0 1-12.48-7.36l-4.992-3.072z m-35.52 70.176l4.352-2.56 3.68 2.24c13.984 8.192 33.088 17.856 56.576 27.392l20.32 7.968a284.896 284.896 0 0 0 48.64 13.952l7.744 1.216v205.92l-0.16 4.096c-2.816 37.216-43.84 82.08-127.136 131.36l-10.304 5.952-0.864-0.48c-92.864-53.248-136.736-101.6-136.736-140.928v-206.272l3.744-0.64a416.64 416.64 0 0 0 61.664-17.088l9.024-3.584 3.904-1.6a518.784 518.784 0 0 0 55.552-26.944z"
                          fill="#333333" p-id="4217"></path>
                      <path
                          d="M857.856 646.784l54.944 50.848-133.088 129.152-95.744-86.176 52.864-52.768 40.832 36.736z"
                          fill="#444444" p-id="4218"></path>
                    </svg>
                  </button>

                  <button class="ml-[1vw] text-[4.5vw]">
                    {{ department }}
                  </button>
                  <span class="overflow-auto float-right mr-[2vw] text-[4.9vw] text-black">
              <button :class="{
            'transition duration-400 pl-[1vw] pr-[1vw]':true,
            'rotate-x-180':curDepartment===department}">
              <i class="bi bi-chevron-down"></i>
            </button>
            </span>

                </div>


              </div>

            </div>

          </div>


          <!--          被折叠的内容-->
          <div :id="department+'-collapse'" :class="{
        'hs-accordion-content overflow-hidden transition-[height] duration-300':true,
        'hidden':index!==0
        }" :aria-labelledby="department+'-heading'">
            <div> <!-- 被折叠的内容容器-->

              <div :class="collapseKingBgColor" v-if="curBuildings(department).length>0">
                <KingArea :building-icon-color="collapseBuildingIconColor"
                          :chose-box-color="collapseChoseBoxColor"
                          :department="department"
                          @get-cur-building-with-department="getCurBuildingWithDepartment"
                          :all-buildings="curBuildings(department)"
                ></KingArea>
              </div>

              <!--              课程卡片-->
              <div :class="collapseContentBgColor" style="padding-top: 2vw">
                <div v-if="isApiError">
                  网络异常，请通过正确的域名访问~
                </div>
                <div v-else-if="!hasTeachInfo.get(department)">
                  <!--                  {{ hasTeachInfo.get(department) === true }}-->
                  <div class="pb-[10vw]">
                    <div class="flex place-content-center">
                      <img class="w-[50vw]" src="/src/assets/desk3.png" alt="">
                    </div>
                    <div class="text-center text-[6vw]">该学部这个时间没有课~</div>
                  </div>
                  <!---->
                  <!---->
                  <!--                  text-[#c48833]-->
                </div>
                <div v-else>
                  <div>
                    <div class="mt-[3vw] pb-[1vw]"
                         v-for="teachInfo in teachInfos.get(department)?.get(curBuildingDepartmentMap?.get(department))">
                      <CourseCard :info-card-color="infoCardColor"
                                  :teach-info="teachInfo">
                      </CourseCard>
                    </div>
                  </div>
                  <!---->
                </div>
              </div>
            </div>
          </div>


        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
.search-row {
  width: 100%;
  height: 10vw;
  background-color: #fff;
  box-sizing: border-box;
  padding: 0 3vw;
  font-size: 3.6vw;
  border-radius: 1.5rem;
  outline: none;
}
.searchs {
  width: 100vw;
  box-sizing: border-box;
  padding: 3vw 6vw 0 6vw;
}
</style>