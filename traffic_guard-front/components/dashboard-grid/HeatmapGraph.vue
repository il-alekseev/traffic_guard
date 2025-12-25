<template>
  <div class="heatmap">
    <div class="heatmap__graph" ref="chartContainer"></div>
    <div class="heatmap__legend">
      <div class="heatmap__legend-extremum-values">
        >{{ extremumValues.min }}
      </div>
      <div class="heatmap__legend-color-block heatmap__legend-color-block_1-level"></div>
      <div class="heatmap__legend-color-block heatmap__legend-color-block_2-level"></div>
      <div class="heatmap__legend-color-block heatmap__legend-color-block_3-level"></div>
      <div class="heatmap__legend-color-block heatmap__legend-color-block_4-level"></div>
      <div class="heatmap__legend-extremum-values">
        <{{ extremumValues.max }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import type { ApexOptions } from 'apexcharts'

type Props = {
  height: string
  values: number[]
}

const props = defineProps<Props>()

const chartContainer = ref<HTMLElement | null>(null);
const chartInstance = ref<any>(null);
const monthNames = ['Янв','Фев','Мар','Апр','Май','Июн','Июл','Авг','Сен','Окт','Ноя','Дек'];
const weekDays = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'];

const extremumValues = computed(() => {
  if (!props.values || props.values.length === 0) {
    return { min: 0, max: 0 }
  }
  return { 
    min: Math.min(...props.values), 
    max: Math.max(...props.values) 
  }
})

const colorRanges = computed(() => {
  const min = extremumValues.value.min;
  const max = extremumValues.value.max;
  const step = (max - min) / 4;
  
  return [
    { from: min, to: min + step, color: '#EFF6FF', name: 'Низкий' },
    { from: min + step + 1, to: min + step * 2, color: '#93C5FD', name: 'Средний' },
    { from: min + step * 2 + 1, to: min + step * 3, color: '#3B82F6', name: 'Высокий' },
    { from: min + step * 3 + 1, to: max, color: '#1E40AF', name: 'Очень высокий' }
  ]
})

onMounted(async () => {
  const ApexCharts = (await import('apexcharts')).default;
  
  if (!chartContainer.value) {
    console.error('Chart container not found');
    return;
  }

  if (!props.values || props.values.length === 0) {
    console.error('No values provided');
    return;
  }

  const year = new Date().getFullYear();
  const startDate = new Date(year, 0, 1);
  
  let startDayOfWeek = startDate.getDay();
  startDayOfWeek = startDayOfWeek === 0 ? 6 : startDayOfWeek - 1;
  
  const datesWithValues: { date: Date; value: number; dayOfWeek: number; monthName: string; weekInYear: number }[] = [];
  
  for (let i = 0; i < props.values.length; i++) {
    const date = new Date(year, 0, 1 + i);
    let dayOfWeek = date.getDay();
    dayOfWeek = dayOfWeek === 0 ? 6 : dayOfWeek - 1;
    
    const monthName = monthNames[date.getMonth()];

    const startOfYear = new Date(year, 0, 1);
    const daysSinceStart = Math.floor((date.getTime() - startOfYear.getTime()) / (1000 * 60 * 60 * 24));
    const weekInYear = Math.floor((daysSinceStart + startDayOfWeek) / 7);
    
    datesWithValues.push({
      date,
      value: props.values[i],
      dayOfWeek,
      monthName,
      weekInYear
    })
  }
  
  const totalWeeks = Math.max(...datesWithValues.map(d => d.weekInYear)) + 1;

  const xAxisLabels: string[] = [];
  const processedMonths = new Set<string>();
  
  for (let week = 0; week < totalWeeks; week++) {
    const daysInWeek = datesWithValues.filter(d => d.weekInYear === week);
    
    if (daysInWeek.length > 0) {
      const monthName = daysInWeek[0].monthName;
      
      if (!processedMonths.has(monthName)) {
        xAxisLabels.push(monthName);
        processedMonths.add(monthName);
      } else {
        xAxisLabels.push('');
      }
    } else {
      xAxisLabels.push('');
    }
  }

  const series = weekDays.map((dayName, dayOfWeekIndex) => {
    const weeklyData: { x: string; y: number | null }[] = [];
    
    for (let week = 0; week < totalWeeks; week++) {
      const dayData = datesWithValues.find(d => 
        d.dayOfWeek === dayOfWeekIndex && d.weekInYear === week
      );
      
      weeklyData.push({
        x: `W${week}`,
        y: dayData ? dayData.value : null
      });
    }
    
    return {
      name: dayName,
      data: weeklyData
    }
  })

  const options: ApexOptions = {
    chart: {
      type: 'heatmap',
      height: parseInt(props.height) || 300,
      toolbar: { show: false },
      selection: {
        enabled: false
      },
      events: {
        dataPointSelection: undefined
      }
    },
    plotOptions: {
      heatmap: {
        shadeIntensity: 0.5,
        radius: 4,
        useFillColorAsStroke: false,
        colorScale: {
          ranges: colorRanges.value
        },
        distributed: false,
        enableShades: true
      }
    },
    dataLabels: {
      enabled: false
    },
    stroke: {
      show: true,
      width: 4,
      colors: ['#ffffff']
    },
    xaxis: {
      type: 'category',
      categories: xAxisLabels,
      position: 'top',
      labels: {
        show: true,
        rotate: 0,
        hideOverlappingLabels: true,
        style: {
          colors: '#A1A1AA',
          fontSize: '13px'
        }
      },
      tooltip: {
        enabled: false
      }
    },
    yaxis: {
      labels: {
        show: true,
        style: {
          colors: '#A1A1AA',
          fontSize: '13px'
        },
        offsetX: -10
      },
      reversed: true
    },
    grid: {
      padding: {
        top: 0,
        right: 0,
        bottom: 0,
        left: 10
      }
    },
    states: {
      hover: {
        filter: {
          type: 'none'
        }
      },
      active: {
        filter: {
          type: 'none'
        }
      }
    },
    legend: {
      show: false
    },
    tooltip: {
      enabled: true,
      shared: false,
      intersect: true,
      x: {
        show: false
      },
      custom: function({ series, seriesIndex, dataPointIndex, w }) {
        const value = series[seriesIndex][dataPointIndex];
        
        const dayOfWeek = weekDays[seriesIndex];
        
        const dateInfo = datesWithValues.find(d => 
          d.dayOfWeek === seriesIndex && 
          d.weekInYear === dataPointIndex
        );
        
        if (!value || value === null) {
          return '<div class="apexcharts-tooltip-custom" style="padding: 8px 12px; background: #fff; border: 1px solid #e3e3e3; border-radius: 4px;">Нет данных</div>'
        };
        
        const dateStr = dateInfo ? 
          `${dateInfo.date.getDate()} ${dateInfo.monthName}` : 
          dayOfWeek
        
        return `
          <div class="apexcharts-tooltip-custom" style="padding: 8px 12px; background: #fff; border: 1px solid #e3e3e3; border-radius: 4px; box-shadow: 0 2px 8px rgba(0,0,0,0.15);">
            <div style="font-weight: 500; margin-bottom: 2px; font-size: 13px;">${dateStr}</div>
            <div style="color: #666; font-size: 13px;">Значение: <span style="font-weight: 500; color: #000">${value}</span></div>
          </div>
        `
      }
    },
    series: series
  }

  try {
    chartInstance.value = new ApexCharts(chartContainer.value, options)
    await chartInstance.value.render()
  } catch (error) {
    console.error('Error rendering chart:', error)
  }
})

onBeforeUnmount(() => {
  if (chartInstance.value) {
    chartInstance.value.destroy()
  }
})
</script>

<style scoped lang="scss">
.heatmap {
  width: 100%;

  &__graph {
    width: 100%;
    min-height: 18.75rem;
  }

  &__legend {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-top: -1.5rem;
    justify-content: flex-end;

    &-color-block {
      width: 1.25rem;
      height: 1.25rem;
      border-radius: 2px;

      &_1-level {
        background-color: #EFF6FF;
      }

      &_2-level {
        background-color: #93C5FD;
      }

      &_3-level {
        background-color: #3B82F6;
      }

      &_4-level {
        background-color: #1E40AF;
      }
    }

    &-extremum-values {
      font-size: 0.75rem;
      color: #A1A1AA;
    }
  }
}
</style>