<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const summaryCards = computed(() => [
  { label: t('dashboard.summary.payment'), value: '¥ 286,400', growth: '+18.6%', tone: 'sunrise' },
  { label: t('dashboard.summary.users'), value: '1,284', growth: '+9.2%', tone: 'teal' },
  { label: t('dashboard.summary.conversion'), value: '68.9%', growth: '+3.5%', tone: 'amber' },
  { label: t('dashboard.summary.tickets'), value: '23', growth: '-12.0%', tone: 'slate' },
]);

const sourceData = computed(() => [
  { channel: t('dashboard.sources.organic'), share: 42, amount: '¥ 92,300' },
  { channel: t('dashboard.sources.ads'), share: 28, amount: '¥ 66,100' },
  { channel: t('dashboard.sources.repurchase'), share: 18, amount: '¥ 41,200' },
  { channel: t('dashboard.sources.partners'), share: 12, amount: '¥ 24,800' },
]);

const todoList = computed(() => [
  {
    title: t('dashboard.todos.landingPage'),
    deadline: t('dashboard.todos.landingPageDeadline'),
    status: t('dashboard.todos.inProgress'),
  },
  {
    title: t('dashboard.todos.shipmentIssue'),
    deadline: t('dashboard.todos.shipmentIssueDeadline'),
    status: t('dashboard.todos.pending'),
  },
  {
    title: t('dashboard.todos.membershipCopy'),
    deadline: t('dashboard.todos.membershipCopyDeadline'),
    status: t('dashboard.todos.scheduled'),
  },
]);

const activityFeed = computed(() => [
  { time: '09:20', text: t('dashboard.activity.growthRelease'), color: 'rgb(var(--green-6))' },
  { time: '10:10', text: t('dashboard.activity.deliveryRecovery'), color: 'rgb(var(--orange-6))' },
  { time: '11:45', text: t('dashboard.activity.roiTarget'), color: 'rgb(var(--arcoblue-6))' },
]);
</script>

<template>
  <div class="page-stack">
    <section class="hero-panel">
      <div>
        <p class="hero-kicker">{{ t('dashboard.heroKicker') }}</p>
        <h2>{{ t('dashboard.heroTitle') }}</h2>
        <p class="hero-copy">{{ t('dashboard.heroCopy') }}</p>
      </div>

      <a-space size="small">
        <a-button v-permission="'dashboard:create-task'" type="primary">{{ t('dashboard.createTask') }}</a-button>
        <a-button v-permission="'dashboard:export'">{{ t('dashboard.exportDaily') }}</a-button>
      </a-space>
    </section>

    <a-row :gutter="12" class="metric-grid">
      <a-col v-for="card in summaryCards" :key="card.label" :xs="24" :sm="12" :xl="6">
        <a-card :class="['metric-card', `metric-card--${card.tone}`]" :bordered="false">
          <p>{{ card.label }}</p>
          <strong>{{ card.value }}</strong>
          <span>{{ card.growth }} {{ t('dashboard.vsYesterday') }}</span>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="12">
      <a-col :xs="24" :xl="15">
        <a-card class="panel-card" :bordered="false" :title="t('dashboard.trafficSources')">
          <div class="source-list">
            <div v-for="item in sourceData" :key="item.channel" class="source-row">
              <div>
                <strong>{{ item.channel }}</strong>
                <span>{{ item.amount }}</span>
              </div>
              <div class="source-progress">
                <a-progress :percent="item.share / 100" :show-text="false" />
                <em>{{ item.share }}%</em>
              </div>
            </div>
          </div>
        </a-card>
      </a-col>

      <a-col :xs="24" :xl="9">
        <a-card class="panel-card" :bordered="false" :title="t('dashboard.todoTitle')">
          <div class="todo-list">
            <div v-for="item in todoList" :key="item.title" class="todo-item">
              <div>
                <strong>{{ item.title }}</strong>
                <span>{{ item.deadline }}</span>
              </div>
              <a-tag>{{ item.status }}</a-tag>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <a-card class="panel-card" :bordered="false" :title="t('dashboard.activityTitle')">
      <a-timeline>
        <a-timeline-item v-for="item in activityFeed" :key="item.time" :dot-color="item.color">
          <div class="timeline-row">
            <strong>{{ item.time }}</strong>
            <span>{{ item.text }}</span>
          </div>
        </a-timeline-item>
      </a-timeline>
    </a-card>
  </div>
</template>