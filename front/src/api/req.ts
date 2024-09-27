import {Items} from "@/types/Items";
import {request} from "./myAxois.ts";

export const webGetTeachInfos = ()
    : Promise<Items.BuildingTeachInfos[][]> => {
    return request({
        method: 'post',
        url: '/teach-infos',
    })
};

// export const webGetBuildings = (department: string)
//     : Promise<Items.RespBuildings> => {
//     return request({
//         method: 'get',
//         url: '/buildings',
//         params: {department}
//     })
// };

export const webGetCurTime = ()
    : Promise<Items.ReqCurTime> => {
    return request({
        method: 'get',
        url: '/cur-time',
    })
};