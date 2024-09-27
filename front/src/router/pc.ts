import {RouteRecordRaw} from "vue-router";


export const routes:Array<RouteRecordRaw> = [
    {
        path: '/',
        component: () => import('@/view/HelperPage.vue'),
        meta: {hidden: true}
    }
]

export default routes;
