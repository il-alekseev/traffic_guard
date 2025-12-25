<template>
  <div class="date-picker">
    <div class="date-picker__section" @click="toggleCalendar">
      <CalendarIcon />
      <span class="date-picker__text">{{ formattedDateRange }}</span>
    </div>

    <div class="date-picker__section date-picker__section--dropdown" @click="togglePeriodMenu">
      <span class="date-picker__text">{{ selectedPeriod }}</span>
      <ArrowIcon class="date-picker__chevron" :class="{ 'date-picker__chevron--open': isPeriodMenuOpen }" />
    </div>

    <Transition name="dropdown">
      <div v-if="isPeriodMenuOpen" class="date-picker__menu">
        <div
          v-for="period in periods"
          :key="period"
          class="date-picker__menu-item"
          :class="{ 'date-picker__menu-item--active': selectedPeriod === period }"
          @click="selectPeriod(period)"
        >
          {{ period }}
        </div>
      </div>
    </Transition>
 
    <Transition name="calendar">
      <div v-if="isCalendarOpen" class="date-picker__calendar">
        <div class="calendar">
          <div class="calendar__header">
            <button class="calendar__nav" @click="previousMonth">
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

            <button class="calendar__nav" @click="nextMonth">
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
                'calendar__day--selected': isSelected(day.date),
                'calendar__day--in-range': isInRange(day.date),
                'calendar__day--start': isStartDate(day.date),
                'calendar__day--end': isEndDate(day.date)
              }"
              @click="selectDate(day.date, day.isCurrentMonth)"
            >
              {{ day.day }}
            </div>
          </div>

          <div class="calendar__footer">
            <button class="calendar__button calendar__button--secondary" @click="clearDates">
              Очистить
            </button>
            <button class="calendar__button calendar__button--primary" @click="applyDates">
              Применить
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <Transition name="backdrop">
      <div v-if="isCalendarOpen || isPeriodMenuOpen" class="date-picker__backdrop" @click="closeAll"></div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import CalendarIcon from "~/assets/img/calendar.svg"
import ArrowIcon from "~/assets/img/arrow-down.svg";
import { makeUTCDate } from "~/helpers/index";

interface CalendarDay {
  date: Date
  day: number
  isCurrentMonth: boolean
  key: string
}

const props = defineProps<{
  modelValue?: { from: Date | null; to: Date | null }
}>()

const emit = defineEmits<{
  'update:modelValue': [value: { from: Date | null; to: Date | null }]
}>()

const isCalendarOpen = ref(false)
const isPeriodMenuOpen = ref(false)
const currentDate = ref(new Date())
const startDate = ref<Date | null>(props.modelValue?.from || null)
const endDate = ref<Date | null>(props.modelValue?.to || null)
const tempStartDate = ref<Date | null>(null)
const tempEndDate = ref<Date | null>(null)
const selectedPeriod = ref('День')
const periods = ['День', 'Неделя', 'Месяц', 'Год']
const weekdays = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']

const monthNames = [
  'Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь',
  'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь'
]

const monthNamesShort = [
  'янв', 'фев', 'мар', 'апр', 'мая', 'июн',
  'июл', 'авг', 'сен', 'окт', 'ноя', 'дек'
]

const currentMonthYear = computed(() => {
  const month = monthNames[currentDate.value.getMonth()]
  const year = currentDate.value.getFullYear()
  return `${month} ${year}`
})

const formattedDateRange = computed(() => {
  if (!startDate.value || !endDate.value) {
    return 'Выберите даты'
  }

  const formatDate = (date: Date) => {
    const day = date.getDate()
    const month = monthNamesShort[date.getMonth()]
    return `${day} ${month}`
  }

  const start = formatDate(startDate.value)
  const end = formatDate(endDate.value)
  const year = endDate.value.getFullYear()

  return `${start} - ${end}, ${year}`
})

const calendarDays = computed((): CalendarDay[] => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth()
  
  const firstDay = makeUTCDate(year, month, 1);
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

const toggleCalendar = () => {
  isCalendarOpen.value = !isCalendarOpen.value
  isPeriodMenuOpen.value = false
  if (isCalendarOpen.value) {
    tempStartDate.value = startDate.value
    tempEndDate.value = endDate.value
  }
}

const togglePeriodMenu = () => {
  isPeriodMenuOpen.value = !isPeriodMenuOpen.value
  isCalendarOpen.value = false
}

const closeAll = () => {
  isCalendarOpen.value = false
  isPeriodMenuOpen.value = false
}

const selectPeriod = (period: string) => {
  selectedPeriod.value = period
  isPeriodMenuOpen.value = false
  
  const today = new Date()
  const start = new Date(today)
  let end = new Date(today)
  
  switch (period) {
    case 'День':
      break
    case 'Неделя':
      end = new Date(today.setDate(today.getDate() + 7))
      break
    case 'Месяц':
      end = new Date(today.setMonth(today.getMonth() + 1))
      break
    case 'Год':
      end = new Date(today.setFullYear(today.getFullYear() + 1))
      break
  }
  
  startDate.value = start
  endDate.value = end
  emit('update:modelValue', { from: startDate.value, to: endDate.value })
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

const selectDate = (date: Date, isCurrentMonth: boolean) => {
  if (!isCurrentMonth) return
  
  if (!tempStartDate.value || (tempStartDate.value && tempEndDate.value)) {
    tempStartDate.value = date
    tempEndDate.value = null
  } else {
    if (date < tempStartDate.value) {
      tempEndDate.value = tempStartDate.value
      tempStartDate.value = date
    } else {
      tempEndDate.value = date
    }
  }
}

const selectedYear = ref(currentDate.value.getFullYear())

const yearsRange = computed(() => {
  const current = new Date().getFullYear()
  const range: number[] = []
  for (let y = current - 10; y <= current + 10; y++) {
    range.push(y)
  }
  return range
})

const changeYear = () => {
  currentDate.value = new Date(
    selectedYear.value,
    currentDate.value.getMonth(),
    1
  )
}

const isSelected = (date: Date): boolean => {
  if (!tempStartDate.value) return false
  const dateStr = date.toDateString()
  return dateStr === tempStartDate.value.toDateString() || 
         (tempEndDate.value && dateStr === tempEndDate.value.toDateString()) as boolean
}

const isInRange = (date: Date): boolean => {
  if (!tempStartDate.value || !tempEndDate.value) return false
  return date > tempStartDate.value && date < tempEndDate.value
}

const isStartDate = (date: Date): boolean => {
  if (!tempStartDate.value) return false
  return date.toDateString() === tempStartDate.value.toDateString()
}

const isEndDate = (date: Date): boolean => {
  if (!tempEndDate.value) return false
  return date.toDateString() === tempEndDate.value.toDateString()
}

const clearDates = () => {
  tempStartDate.value = null
  tempEndDate.value = null
}

const applyDates = () => {
  startDate.value = tempStartDate.value
  endDate.value = tempEndDate.value
  emit('update:modelValue', { from: startDate.value, to: endDate.value })
  isCalendarOpen.value = false
}

onMounted(() => {
  document.addEventListener('keydown', handleEscape);
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscape)
})

watch(
  () => [props.modelValue?.from, props.modelValue?.to],
  ([from, to]) => {
    startDate.value = from ? new Date(from) : null;
    endDate.value = to ? new Date(to) : null;
  },
  { immediate: true }
);

const handleEscape = (e: KeyboardEvent) => {
  if (e.key === 'Escape') {
    closeAll()
  }
}
</script>

<style lang="scss" scoped>
.date-picker {
  position: relative;
  display: inline-flex;
  gap: 0;
  background: #ffffff;
  border-radius: 6px;
  box-shadow: 0px 1px 2px 0px #0000000D;

  &__section {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    padding: 0.5rem 0.75rem;
    cursor: pointer;
    transition: background-color 0.2s;
    border-right: none;

    &:first-child {
      border-right: 1px solid #e5e7eb;
    }

    &:hover {
      background-color: #f9fafb;
    }

    &:first-child {
      border-radius: 6px 0 0 6px;
    }

    &:last-child {
      border-radius: 0 6px 6px 0;
    }

    &--dropdown {
      gap: 0.375rem;
    }
  }

  &__icon {
    width: 1.25rem;
    height: 1.25rem;
    color: #A1A1AA;
    flex-shrink: 0;
  }

  &__text {
    font-size: 0.875rem;
    line-height: 1.25rem;
    color: #3F3F46;
    font-weight: 500;
    white-space: nowrap;
  }

  &__chevron {
    width: 1.25rem;
    height: 1.25rem;
    color: #A1A1AA;
    transition: transform 0.2s;
    flex-shrink: 0;

    &--open {
      transform: rotate(180deg);
    }
  }

  &__menu {
    position: absolute;
    top: calc(100% + 0.5rem);
    right: 0;
    background: #ffffff;
    border-radius: 6px;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
    overflow: hidden;
    z-index: 1000;
    min-width: 8.75rem;
  }

  &__menu-item {
    padding: 0.75rem 1.25rem;
    font-size: 0.875rem;
    line-height: 1rem;
    color: #3F3F46;
    cursor: pointer;
    transition: background-color 0.2s;

    &:hover {
      background-color: #f3f4f6;
    }

    &--active {
      background-color: #eff6ff;
      color: #1447E6;
      font-weight: 500;
    }
  }

  &__calendar {
    position: absolute;
    top: calc(100% + 0.5rem);
    left: 0;
    z-index: 1000;
  }

  &__backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 999;
  }
}

.calendar {
  background: #ffffff;
  border-radius: 6px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
  padding: 0.75rem;
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
    line-height: 1rem;
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
    line-height: 1rem;
    color: #3F3F46;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
    position: relative;

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

    &--in-range {
      background-color: #dbeafe;
      color: #1e40af;
    }

    &--start {
      background-color: #2563eb;
      color: #ffffff;
      border-radius: 6px 0 0 6px;
    }

    &--end {
      background-color: #2563eb;
      color: #ffffff;
      border-radius: 0 6px 6px 0;
    }

    &--start#{&}--end {
      border-radius: 6px;
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
    line-height: 1rem;
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

.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-0.5rem);
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