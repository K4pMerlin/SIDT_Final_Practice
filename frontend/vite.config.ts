import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import {NaiveUiResolver} from 'unplugin-vue-components/resolvers'
import {ElementPlusResolver} from 'unplugin-vue-components/resolvers'

import ElementPlus from "unplugin-element-plus/vite";

import path from 'path';

const pathSrc = path.resolve(__dirname, "src");
// https://vitejs.dev/config/
export default defineConfig({
    css: {
        preprocessorOptions: {
            scss: {
                api: "modern-compiler" // or 'modern'
            }
        }
    },
    plugins: [
        vue(),
        AutoImport({
            resolvers: [ElementPlusResolver()],
            imports: [
                'vue',
                {
                    'naive-ui': [
                        'useDialog',
                        'useMessage',
                        'useNotification',
                        'useLoadingBar'
                    ]
                }
            ]
        }),
        Components({
            resolvers: [ElementPlusResolver(), NaiveUiResolver()],
        }),
        ElementPlus({}),

    ],
    resolve: {
        alias: {
            "@": pathSrc,
        },
    }
})
