import {ref} from "vue";

export const baseURL =
    // window.location.origin;
'http://localhost:80'
// 'https://cengkehelper.top';

export const isApiError = ref(false)
export const hasTeachInfo = ref(new Map<string,boolean>())

export const validData = ref(true)
export const isAdjust = ref(false)

export const teachInfosCacheKey = "teachInfosCache"
export const curTimeCacheKey = "curTimeCache"
export const adjustCacheKey = "adjustCache"

