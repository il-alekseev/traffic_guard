export interface StatItem {
  name: string
  stat: string
  previousStat: string
  change: string
  changeType: 'increase' | 'decrease' | 'neutral'
}
