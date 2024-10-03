import {computed, Ref, ref} from "vue";
import {webGetCurTime} from "@/api/req.ts";
import {curTimeCacheKey, teachInfosCacheKey} from "@/api/globalConst.ts";
import {GlobalTeachInfosObj} from "@/store/teachInfosObj.ts";
import {validData} from "@/store/globalData.ts";


export const GlobalCacheObj = (() => {
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

    function loadData(){
        // 使用缓存，而非默认值

        GlobalTeachInfosObj.loadGlobalTeachInfos(true)

        const localTimeStr = localStorage.getItem(curTimeCacheKey)
        const localTime: Ref<Array<number>> = ref(
            localTimeStr ? JSON.parse(localTimeStr) : []
        );
        if (localTime.value && localTime.value.length === 3) {
            weekday.value = weekdayMap[localTime.value[1]];
            lessonNum.value = localTime.value[2];
        }

        webGetCurTime()
            .then((data) => {
                weekday.value = weekdayMap[data.weekday];
                lessonNum.value = data.lessonNum;
                if (!data.valid) {
                    setInterval(() => {
                        location.reload();
                    }, 2000);
                }
                let curTime = localStorage.getItem(curTimeCacheKey);
                let t: Array<number> = [];
                if (curTime !== null) {
                    t = JSON.parse(curTime);
                }
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

                // 数据标记为失效，后续请求更新
                validData.value = false;
            })
            .catch((err) => {
                console.log(err);
            });

    }

    return {
        weekday,
        lessonTime,
        loadData
    }
})();
