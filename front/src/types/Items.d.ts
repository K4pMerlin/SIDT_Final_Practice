import {ref} from "vue";

export declare module Items {
    export interface TeachInfo {
        room: string;
        faculty: string;
        courseName: string;
        teacherName: string;
        teacherTitle: string;
        courseTime: string;
        courseType: string;
    }

    export interface ReqCurTime {
        weekNum: number;
        weekday: number;
        lessonNum: number;
        valid: boolean;
        isAdjust: boolean;
    }

    export interface BuildingTeachInfos {
        building: string;
        infos: Items.TeachInfo[];
    }

}