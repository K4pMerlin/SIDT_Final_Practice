/** @type {import('tailwindcss').Config} */
import plugin from "tailwindcss/plugin";

// Rotate X utilities
const rotateX = plugin(function ({addUtilities}) {
  addUtilities({
    '.rotate-x-20': {
      transform: 'rotateX(20deg)',
    },
    '.rotate-x-40': {
      transform: 'rotateX(40deg)',
    },
    '.rotate-x-60': {
      transform: 'rotateX(60deg)',
    },
    '.rotate-x-180': {
      transform: 'rotateX(180deg)',
    },
  })
})

export default {
  content: [
    "node_modules/preline/dist/*.js",
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [rotateX, require('preline/plugin')],
}

