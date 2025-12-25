<template>
  <div class="resource-card" :class="item.action === 'Заблокировано' ? 'resource-card_transparent' : ''">
    <div class="resource-card__header">
      <div class="resource-card__title-row">
        <span class="resource-card__label">Domain</span>
      </div>
      <div class="resource-card__title-block">
        <div class="resource-card__title-icon-block">
          <DetectionsIcon class="resource-card__title-icon" />
        </div>
        <span class="resource-card__title">{{ item.domain || 'Неизвестно' }}</span>
        <span :class="['resource-card__category', `resource-card__category--${getModificatorByCategory(item.category)}`]">
          {{ item.category || 'Неизвестно' }}
        </span>
        <span class="resource-card__date">
          {{
            new Date(item.categorized_at).toLocaleString('ru-RU', {
              day: '2-digit',
              month: '2-digit',
              year: 'numeric',
              hour: '2-digit',
              minute: '2-digit'
            })
          }}
        </span>
      </div>
    </div>

    <div class="resource-card__section">
      <div class="resource-card__label">Описание</div>
      <div class="resource-card__description">{{ item.description || '–' }}</div>
    </div>

    <div class="resource-card__section">
      <div class="resource-card__label">Host</div>
      <div class="resource-card__host">
        <div v-if="item.location !== 'private'" class="resource-card__flag-icon-block">
          <EmojiFlag :code="item.location.substring(0, 2).toLowerCase()" />
        </div>
        <span class="resource-card__description resource-card__description-location">{{ item.location }}</span>
        <a :href="`https://${item.ip}`" class="resource-card__ip">{{ item.ip }}</a>
      </div>
      <div class="resource-card__meta">
        <span 
          :class="['resource-card__badge']"
          :style="{
            backgroundColor: deviceColors.background,
            color: deviceColors.color
          }"
        >
          {{ item.hostname || 'Неизвестно' }}
        </span>
        <span class="resource-card__requests">
          Количество обращений: <strong>{{ item.request_count }}</strong>
        </span>
      </div>
    </div>

    <div class="resource-card__footer">
      <div :class="['resource-card__status', `resource-card__status--${getStatusNameByAction(item.action)}`]">
        <span class="resource-card__status-dot"></span>
        {{ item.action }}
      </div>
      <div v-if="item.action === 'Не решено'" class="resource-card__actions">
        <BaseButton
          type="button"
          variant="primary"
          extras="small"
          @click="$emit('confirm')"
        >
          Подтвердить
        </BaseButton>
        <BaseButton
          type="button"
          variant="secondary"
          extras="small"
          @click="$emit('reject')"
        >
          Отклонить
        </BaseButton>
      </div>
      <div v-else-if="item.action === 'Разрешено'" class="resource-card__accepted">
        Доступ разрешен
      </div>
      <div v-else-if="item.action === 'Заблокировано'" class="resource-card__blocked">
        Заблокировано
      </div>
      <div v-else>
        Неизвестно
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getModificatorByCategory, getStatusNameByAction } from '~/helpers';
import BaseButton from '~/components/ui/BaseButton.vue';
import EmojiFlag from "~/components/ui/EmojiFlag.vue"
import { useDeviceColors } from '~/composables/useDeviceColors';
import type { Detection } from '~/types/detections';
import DetectionsIcon from "~/assets/img/detections.svg"


interface Props {
  item: Detection
}

const props = defineProps<Props>()

defineEmits<{
  confirm: []
  reject: []
}>()

const { generateColor } = useDeviceColors();
const deviceColors = generateColor(props.item.hostname || 'Неизвестно');

</script>

<style scoped lang="scss">
.resource-card {
  max-width: 32rem;
  width: 100%;
  background: #ffffff;
  border-radius: 28px;
  padding: 1.25rem;
  box-shadow: 0px 1px 2px 0px #0000000D;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;

  &_transparent {
    opacity: 0.5;
  }

  &__header {
    display: flex;
    flex-direction: column;
  }

  &__title-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  &__title-icon-block {
    width: 1.25rem;
    height: 1.25rem;
    color: #FF8800;

    & img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }

  &__title {
    margin-left: 0.375rem;
    margin-right: 0.75rem;
    font-weight: 600;
    font-size: 1rem;
    line-height: 1.5rem;
    color: #3F3F46;
    max-width: 8.75rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &__label {
    color: #A1A1AA;
    font-size: 0.875rem;
    font-weight: 400;
    line-height: 1.25rem;
  }

  &__category {
    padding: 0.25rem 0.5rem;
    border-radius: 0.375rem;
    font-size: 0.75rem;
    line-height: 1rem;
    font-weight: 500;

    &--aggression { background: #FEE2E2; color: #DC2626; }
    &--racism { background: #FECACA; color: #B91C1C; }
    &--terrorism { background: #FEE2E2; color: #DC2626; }
    &--drugs { background: #FECACA; color: #B91C1C; }
    &--malware { background: #FEE2E2; color: #EF4444; }
    &--malicious-software { background: #FECACA; color: #DC2626; }

    &--botnets { background: #FFEDD5; color: #EA580C; }
    &--cryptomining { background: #FDBA74; color: #C2410C; }
    &--cryptojacking { background: #FFEDD5; color: #D97706; }
    &--phishing { background: #FEF3C7; color: #D97706; }

    &--pornography { background: #F3E8FF; color: #9333EA; }
    &--adult { background: #E9D5FF; color: #7E22CE; }
    &--gambling { background: #F3E8FF; color: #7C3AED; }

    &--proxy { background: #DBEAFE; color: #2563EB; }
    &--torrents { background: #BFDBFE; color: #1D4ED8; }

    &--leisure { background: #DCFCE7; color: #16A34A; }
    &--games { background: #BBF7D0; color: #15803D; }
    &--positive { background: #BBF7D0; color: #15803D; }
    &--online-games { background: #DCFCE7; color: #15803D; }
    &--gaming-platforms { background: #BBF7D0; color: #16A34A; }
    &--movies { background: #DCFCE7; color: #059669; }

    &--social { background: #E0F2FE; color: #0284C7; }
    &--chats { background: #BAE6FD; color: #0369A1; }
    &--webmail { background: #E0F2FE; color: #0EA5E9; }

    &--shopping { background: #E0E7FF; color: #4F46E5; }
    &--advertising { background: #C7D2FE; color: #4338CA; }

    &--files { background: #F3F4F6; color: #6B7280; }
    &--blocked { background: #E5E7EB; color: #4B5563; }
    &--additional { background: #F3F4F6; color: #6B7280; }

    &--depressive { background: #FCE7F3; color: #DB2777; }
    
    &--alcohol { background: #FEF3C7; color: #D97706; }

    &--unknown { background: #f3f4f6; color: #6b7280; }
  }

  &__date {
    margin-left: auto;
    font-weight: 400;
    color: #A1A1AA;
    font-size: 1rem;
    line-height: 1.5rem;
  }

  &__title-block {
    display: flex;
    align-items: center;
  }

  &__icon {
    color: #f97316;
    flex-shrink: 0;
  }

  &__section {
    display: flex;
    flex-direction: column;
  }

  &__description {
    font-weight: 400;
    font-size: 1rem;
    line-height: 1.5rem;
    color: #3F3F46;

    &-location {
      margin-right: 0.75rem;
    }
  }

  &__host {
    display: flex;
    align-items: center;
  }

  &__flag-icon-block {
    width: 1.25rem;
    height: 1.25rem;
    margin-right: 0.375rem;

    & img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }

  &__ip {
    color: #2563EB;
    text-decoration: none;
    font-weight: 400;
    font-size: 1rem;
    line-height: 1.5rem;

    &:hover {
      text-decoration: underline;
    }
  }

  &__meta {
    margin-top: 0.5rem;
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  &__badge {
    padding: 0.25rem 0.5rem;
    border-radius: 6px;
    font-size: 0.75rem;
    line-height: 1rem;
    font-weight: 500;
    text-transform: uppercase;
  }

  &__requests {
    color: #A1A1AA;
    font-size: 1rem;
    line-height: 1.5rem;

    strong {
      color: #3F3F46;
      font-weight: 600;
    }
  }

  &__footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding-top: 0.75rem;
    border-top: 1px solid #E4E4E7;
  }

  &__status {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    padding: 0.25rem 0.5rem;
    border-radius: 6px;
    font-size: 0.75rem;
    line-height: 1rem;
    font-weight: 500;

    &--blocked {
      background: #FFE2E2;
      color: #C10007;
    }

    &--accepted {
      color: #059669;
      background-color: #DCFCE7;
    }

    &--verification {
      background: #FEF9C2;
      color: #894B00;
    }

    &-- {
      background: #E7E5E4;
      color: #57534E;
    }
  }

  &__status-dot {
    width: 0.375rem;
    height: 0.375rem;
    border-radius: 50%;
    background: currentColor;
  }

  &__actions {
    display: flex;
    gap: 0.5rem;
    margin-left: auto;
  }

  &__button {
    padding: 0.5rem 1.25rem;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 500;
    border: none;
    cursor: pointer;
    transition: all 0.2s;

    &--primary {
      background: #2563eb;
      color: white;

      &:hover {
        background: #1d4ed8;
      }
    }

    &--secondary {
      background: #f3f4f6;
      color: #374151;

      &:hover {
        background: #e5e7eb;
      }
    }
  }

  &__blocked {
    padding: 0.5rem 0.75rem;
    color: #C10007;
    background-color: #FFE2E2;
    font-size: 0.875rem;
    line-height: 1.25rem;
    font-weight: 600;
    border-radius: 6px;
  }

  &__accepted {
    padding: 0.5rem 0.75rem;
    color: #059669;
    background-color: #DCFCE7;
    font-size: 0.875rem;
    line-height: 1.25rem;
    font-weight: 600;
    border-radius: 6px;
  }
}
</style>