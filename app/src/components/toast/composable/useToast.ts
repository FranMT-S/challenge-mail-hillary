import { reactive } from 'vue'

const state = reactive({
  show: false,
  message: '',
  color: 'red-darken-2',
  timeout: 4000
})

function open(message: string, color: string = 'red-darken-2') {
  state.message = message
  state.color = color
  state.show = true
}

function close() {
  state.show = false
}

export function useToast() {
  return { state, open, close }
}
