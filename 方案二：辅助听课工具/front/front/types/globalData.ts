import {ref} from "vue";
import {Items} from "./Items";


export const teachInfos =
    ref(new Map<string, Map<string, Items.TeachInfo[]>>())