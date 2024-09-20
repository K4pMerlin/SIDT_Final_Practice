<script setup>

import {onMounted, ref, watch} from "vue";

const emits = defineEmits(['changeDepartment']);

const departments = ref([
  "文理学部",
  "工学部",
  "信息学部",
  "医学部",
  "网安基地"
]);


const curDepartment = ref(departments.value[0]);

onMounted(() => {

})

</script>

<template>
  <div class="grid place-content-center">
    <div class="hs-accordion-group">
      <ul v-for="(department,index) in departments">
        <li :class="{'hs-accordion':true,'active':index===0}" :id="department+'-heading'">

          <!--          触发器-->
          <div class="hs-accordion-toggle whitespace-nowrap"
               :aria-controls="department+'-collapse'"
               @click="curDepartment===department?curDepartment='':curDepartment=department">
            <button>{{ department }}</button>
          </div>


          <!--          被折叠的内容-->
          <div :id="department+'-collapse'" :class="{
        'hs-accordion-content overflow-hidden transition-[height] duration-300':true,
        'hidden':index!==0
        }" :aria-labelledby="department+'-heading'">
            好好好
            好好好
            好好好
            好好好
          </div>


        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped>

</style>