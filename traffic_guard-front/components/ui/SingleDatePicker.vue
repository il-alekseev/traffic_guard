<template>
  <div class="single-date-picker">
    <div 
      class="date-picker__input" 
      :class="{ 
        'date-picker__input--disabled': disabled,
        'date-picker__input--error': error 
      }"
      @click="toggleCalendar"
    >
      <CalendarIcon class="date-picker__icon" />
      <span class="date-picker__text">{{ formattedDate }}</span>
    </div>

    <Transition name="calendar">
      <div v-if="isCalendarOpen" class="date-picker__calendar">
        <div class="calendar">
          <div class="calendar__header">
            <button type="button" class="calendar__nav" @click="previousMonth">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <polyline points="15 18 9 12 15 6"></polyline>
              </svg>
            </button>

            <div class="calendar__title">
              <span>{{ monthNames[currentDate.getMonth()] }}</span>
              <select v-model="selectedYear" class="calendar__year-select" @change="changeYear">
                <option v-for="year in yearsRange" :key="year" :value="year">{{ year }}</option>
              </select>
            </div>

            <button type="button" class="calendar__nav" @click="nextMonth">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <polyline points="9 18 15 12 9 6"></polyline>
              </svg>
            </button>
          </div>

          <div class="calendar__weekdays">
            <div v-for="day in weekdays" :key="day" class="calendar__weekday">
              {{ day }}
            </div>
          </div>

          <div class="calendar__days">
            <div
              v-for="day in calendarDays"
              :key="day.key"
              class="calendar__day"
              :class="{
                'calendar__day--disabled': !day.isCurrentMonth,
                'calendar__day--selected': isSelected(day.date)
              }"
              @click="selectDate(day.date, day.isCurrentMonth)"
            >
              {{ day.day }}
            </div>
          </div>

          <div class="calendar__footer">
            <button type="button" class="calendar__button calendar__button--secondary" @click="clearDate">
              Сбросить
            </button>
            <button type="button" class="calendar__button calendar__button--primary" @click="applyDate">
              Применить
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <Transition name="backdrop">
      <div v-if="isCalendarOpen" class="date-picker__backdrop" @click="closeCalendar"></div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import CalendarIcon from "~/assets/img/calendar.svg"
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { makeUTCDate } from "~/helpers"

interface CalendarDay {
  date: Date
  day: number
  isCurrentMonth: boolean
  key: string
}

interface Props {
  modelValue?: Date | null
  disabled?: boolean
  error?: boolean
  placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: null,
  disabled: false,
  error: false,
  placeholder: 'Выберите дату'
})

const emit = defineEmits<{
  'update:modelValue': [value: Date | null]
}>()

const isCalendarOpen = ref(false)
const currentDate = ref(new Date())
const selectedDate = ref<Date | null>(props.modelValue)
const tempSelectedDate = ref<Date | null>(null)

const weekdays = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']

const monthNames = [
  'Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь',
  'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь'
]

const monthNamesShort = [
  'янв', 'фев', 'мар', 'апр', 'мая', 'июн',
  'июл', 'авг', 'сен', 'окт', 'ноя', 'дек'
]

const formattedDate = computed(() => {
  if (!selectedDate.value) {
    return props.placeholder
  }

  const day = selectedDate.value.getDate()
  const month = monthNamesShort[selectedDate.value.getMonth()]
  const year = selectedDate.value.getFullYear()

  return `${day} ${month}, ${year}`
})

const calendarDays = computed((): CalendarDay[] => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth()
  
  const firstDay = makeUTCDate(year, month, 1)
  const lastDay = makeUTCDate(year, month + 1, 0)
  
  let startDay = firstDay.getDay()
  startDay = startDay === 0 ? 6 : startDay - 1
  
  const days: CalendarDay[] = []
  
  const prevMonthLastDay = makeUTCDate(year, month, 0).getDate()
  for (let i = startDay - 1; i >= 0; i--) {
    const day = prevMonthLastDay - i
    const date = makeUTCDate(year, month - 1, day)
    days.push({
      date,
      day,
      isCurrentMonth: false,
      key: `prev-${day}`
    })
  }
  
  for (let day = 1; day <= lastDay.getDate(); day++) {
    const date = makeUTCDate(year, month, day)
    days.push({
      date,
      day,
      isCurrentMonth: true,
      key: `current-${day}`
    })
  }
  
  const remainingDays = 42 - days.length
  for (let day = 1; day <= remainingDays; day++) {
    const date = makeUTCDate(year, month + 1, day)
    days.push({
      date,
      day,
      isCurrentMonth: false,
      key: `next-${day}`
    })
  }
  
  return days
})

const selectedYear = ref(currentDate.value.getFullYear())

const yearsRange = computed(() => {
  const current = new Date().getFullYear()
  const range: number[] = []
  for (let y = current - 10; y <= current + 10; y++) {
    range.push(y)
  }
  return range
})

const toggleCalendar = () => {
  if (props.disabled) return
  isCalendarOpen.value = !isCalendarOpen.value
  if (isCalendarOpen.value) {
    tempSelectedDate.value = selectedDate.value
  }
}

const closeCalendar = () => {
  isCalendarOpen.value = false
}

const previousMonth = () => {
  currentDate.value = new Date(
    currentDate.value.getFullYear(),
    currentDate.value.getMonth() - 1,
    1
  )
}

const nextMonth = () => {
  currentDate.value = new Date(
    currentDate.value.getFullYear(),
    currentDate.value.getMonth() + 1,
    1
  )
}

const changeYear = () => {
  currentDate.value = new Date(
    selectedYear.value,
    currentDate.value.getMonth(),
    1
  )
}

const selectDate = (date: Date, isCurrentMonth: boolean) => {
  if (!isCurrentMonth) return
  tempSelectedDate.value = date
}

const isSelected = (date: Date): boolean => {
  if (!tempSelectedDate.value) return false
  return date.toDateString() === tempSelectedDate.value.toDateString()
}

const clearDate = () => {
  tempSelectedDate.value = null
  selectedDate.value = null
  emit('update:modelValue', null)
  isCalendarOpen.value = false
}

const applyDate = () => {
  selectedDate.value = tempSelectedDate.value
  emit('update:modelValue', selectedDate.value)
  isCalendarOpen.value = false
}

watch(() => props.modelValue, (newValue) => {
  selectedDate.value = newValue
})

onMounted(() => {
  document.addEventListener('keydown', handleEscape)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscape)
})

watch(
  () => [props.modelValue],
  ([selected]) => {
    selectedDate.value = selected ? new Date(selected) : null;
  },
  { immediate: true }
);


const handleEscape = (e: KeyboardEvent) => {
  if (e.key === 'Escape') {
    closeCalendar()
  }
}
</script>

<style lang="scss" scoped>
.single-date-picker {
  position: relative;
  width: 100%;
}

.date-picker__input {
  display: flex;
  flex-direction: row-reverse;
  align-items: center;
  gap: 0.75rem;
  padding: 0.5rem 0.75rem;
  background: #ffffff;
  border: 1px solid #D4D4D8;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;

  &--disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  &--error {
    border-color: #ef4444;
  }
}

.date-picker__icon {
  width: 1.25rem;
  height: 1.25rem;
  color: #A1A1AA;
  flex-shrink: 0;
}

.date-picker__text {
  font-size: 0.875rem;
  color: #3F3F46;
  flex: 1;
}

.date-picker__calendar {
  position: absolute;
  top: calc(100% + 0.5rem);
  left: 0;
  z-index: 1000;
}

.date-picker__backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 999;
}

.calendar {
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
  padding: 1rem;
  width: 20rem;

  &__header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1.25rem;
  }

  &__nav {
    background: none;
    border: none;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s;

    &:hover {
      background-color: #f3f4f6;
    }

    svg {
      width: 1.25rem;
      height: 1.25rem;
      color: #A1A1AA;
    }
  }

  &__title {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 1rem;
    font-weight: 600;
    color: #3F3F46;
  }

  &__year-select {
    font-size: 0.875rem;
    padding: 0.25rem 0.5rem;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    background: #fff;
    cursor: pointer;
    color: #3F3F46;

    &:hover {
      background-color: #f9fafb;
    }
  }

  &__weekdays {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: 0.25rem;
    margin-bottom: 0.5rem;
  }

  &__weekday {
    text-align: center;
    font-size: 0.75rem;
    font-weight: 600;
    color: #A1A1AA;
    padding: 0.5rem 0;
  }

  &__days {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: 0.25rem;
  }

  &__day {
    aspect-ratio: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.875rem;
    color: #3F3F46;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;

    &:hover:not(&--disabled) {
      background-color: #f3f4f6;
    }

    &--disabled {
      color: #d1d5db;
      cursor: not-allowed;

      &:hover {
        background-color: transparent;
      }
    }

    &--selected {
      background-color: #2563eb;
      color: #ffffff;
      font-weight: 600;

      &:hover {
        background-color: #1d4ed8;
      }
    }
  }

  &__footer {
    display: flex;
    gap: 0.75rem;
    margin-top: 1.25rem;
    padding-top: 1.25rem;
    border-top: 1px solid #e5e7eb;
  }

  &__button {
    flex: 1;
    padding: 0.625rem 1rem;
    border-radius: 6px;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    border: none;

    &--secondary {
      background-color: #f3f4f6;
      color: #374151;

      &:hover {
        background-color: #e5e7eb;
      }
    }

    &--primary {
      background-color: #2563eb;
      color: #ffffff;

      &:hover {
        background-color: #1d4ed8;
      }
    }
  }
}

.calendar-enter-active,
.calendar-leave-active {
  transition: all 0.3s ease;
}

.calendar-enter-from,
.calendar-leave-to {
  opacity: 0;
  transform: translateY(-1rem);
}

.backdrop-enter-active,
.backdrop-leave-active {
  transition: opacity 0.2s ease;
}

.backdrop-enter-from,
.backdrop-leave-to {
  opacity: 0;
}
</style>
