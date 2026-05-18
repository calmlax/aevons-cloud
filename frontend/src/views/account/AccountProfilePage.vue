<script setup lang="ts">
import { computed, ref, reactive, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { Message } from '@arco-design/web-vue';
import { IconUser, IconSafe, IconLocation, IconClockCircle, IconIdcard, IconPlus, IconDelete, IconDesktop, IconLock } from '@arco-design/web-vue/es/icon';
import { 
  getUserProfileRequest, 
  updateUserProfileRequest, 
  getProfileLoginLogsRequest,
  type UserProfileResponse,
  type LoginLogItem 
} from '@/api/user-profile';
import { useAuthStore } from '@/store/modules/auth';
import {
  passkeyRegisterBegin,
  passkeyRegisterFinish,
  passkeyListCredentials,
  passkeyRevokeCredential,
  parseCreationOptions,
  type PasskeyCredential,
} from '@/api/passkey';

const proxy = getCurrentInstance()!.proxy as any
const dicts = proxy.$useDict('sys_success_status')
const { t } = useI18n();
const authStore = useAuthStore();

const loading = ref(true);
const saving = ref(false);

const profile = ref<UserProfileResponse | null>(null);
const loginLogs = ref<LoginLogItem[]>([]);

const visible = ref(false);
const formRef = ref();
const formData = reactive({
  nickname: '',
  email: '',
  mobile: '',
  sex: 0,
});

// ── Passkey 管理 ──────────────────────────────────────────────────────────────
const passkeySupported = typeof window !== 'undefined' && typeof window.PublicKeyCredential !== 'undefined';
const credentials = ref<PasskeyCredential[]>([]);
const passkeyRegistering = ref(false);
const passkeyLoadingList = ref(false);

const fetchCredentials = async () => {
  if (!passkeySupported) return;
  passkeyLoadingList.value = true;
  try {
    credentials.value = await passkeyListCredentials();
  } catch {
    // 静默失败
  } finally {
    passkeyLoadingList.value = false;
  }
};

const handleAddPasskey = async () => {
  if (!passkeySupported) {
    Message.warning(t('auth.passkeyNotSupported'));
    return;
  }
  passkeyRegistering.value = true;
  try {
    // Step 1: 开始注册流程
    const { options: optionsJSON, session_key } = await passkeyRegisterBegin();
    const creationOptions = parseCreationOptions(optionsJSON);
    
    // Step 2: 调用浏览器 WebAuthn API
    const credential = await navigator.credentials.create({ publicKey: creationOptions }) as PublicKeyCredential | null;
    if (!credential) {
      Message.warning(t('auth.passkeyCancelled'));
      return;
    }
    
    // Step 3: 完成注册
    await passkeyRegisterFinish(session_key, credential);
    Message.success(t('auth.passkey_register_success') || 'Passkey 添加成功');
    await fetchCredentials();
  } catch (err: any) {
    // 用户取消操作，静默处理
    if (err?.name === 'NotAllowedError') {
      console.log('User cancelled passkey registration');
      return;
    }
    
    // WebAuthn API 错误
    if (err?.name === 'InvalidStateError') {
      Message.error('此设备已注册过 Passkey，请勿重复添加');
      return;
    }
    
    if (err?.name === 'NotSupportedError') {
      Message.error('当前设备不支持此类型的 Passkey');
      return;
    }
    
    // 网络或服务器错误（已在 axios 拦截器中处理并显示）
    console.error('Passkey registration error:', err);
    
    // 如果是开发环境，显示详细错误信息
    if (import.meta.env.DEV) {
      console.group('🔍 Passkey Registration Debug Info');
      console.error('Error name:', err?.name);
      console.error('Error message:', err?.message);
      console.error('Error response:', err?.response?.data);
      console.error('Full error:', err);
      console.groupEnd();
    }
  } finally {
    passkeyRegistering.value = false;
  }
};

const handleRevokePasskey = async (cred: PasskeyCredential) => {
  try {
    await passkeyRevokeCredential(cred.id);
    Message.success(t('auth.passkey_revoke_success') || '已移除');
    credentials.value = credentials.value.filter(c => c.id !== cred.id);
  } catch (err: any) {
    // 错误消息已经在 request 拦截器中被翻译并显示
    console.error('Passkey revoke error:', err);
  }
};

const formatCredDate = (dateStr: string) => {
  if (!dateStr) return '--';
  return new Date(dateStr).toLocaleDateString();
};
// ─────────────────────────────────────────────────────────────────────────────

const fetchProfile = async () => {
  try {
    loading.value = true;
    const [profileData, logsData] = await Promise.all([
      getUserProfileRequest(),
      getProfileLoginLogsRequest()
    ]);
    profile.value = profileData;
    loginLogs.value = logsData || [];
  } catch (error: any) {
    Message.error(error.message || t('profile.fetchError') || 'Failed to fetch user profile');
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchProfile();
  fetchCredentials();
});

const profileName = computed(() => profile.value?.user.nickname || profile.value?.user.username || 'User');
const profileInitials = computed(() => profileName.value.slice(0, 2).toUpperCase());

const sexMap = computed<Record<number, string>>(() => ({
  0: t('profile.sexMale'),
  1: t('profile.sexFemale'),
  2: t('profile.sexUnknown'),
}));

const openEditModal = () => {
  formData.nickname = profile.value?.user.nickname || '';
  formData.email = profile.value?.user.email || '';
  formData.mobile = profile.value?.user.mobile || '';
  formData.sex = profile.value?.user.sex || 0;
  visible.value = true;
};

const handleBeforeOk = async () => {
  const errors = await formRef.value?.validate();
  if (errors) return false;

  saving.value = true;
  try {
    await updateUserProfileRequest({ ...formData });
    Message.success(t('profile.updateSuccess') || 'Profile updated successfully');
    await fetchProfile();
    authStore.refreshUserInfo(); 
    return true;
  } catch (error: any) {
    Message.error(error.message || t('profile.updateError') || 'Failed to update profile');
    return false;
  } finally {
    saving.value = false;
  }
};

const logsColumns = computed(() => [
  { title: t('profile.logBrowserOs'), dataIndex: 'browser', slotName: 'systemInfo' },
  { title: t('profile.logIpAddress'), dataIndex: 'ip' },
  { title: t('profile.logStatus'), dataIndex: 'status', slotName: 'status' },
  { title: t('profile.logMessage'), dataIndex: 'msg' },
  { title: t('profile.logDateTime'), dataIndex: 'loginAt', slotName: 'loginAt' }
]);

const parseDate = (dateStr: string) => {
  if (!dateStr) return '--';
  return new Date(dateStr).toLocaleString();
};
</script>

<template>
  <div class="page-stack profile-page-layout">
    <a-spin :loading="loading" style="width: 100%">
      <a-row :gutter="16">
        <!-- Main Identity Sidebar -->
        <a-col :xs="24" :lg="8" :xl="7">
          <a-card class="panel-card identity-card" :bordered="false">
            <div class="identity-header">
              <a-avatar :size="100" class="hero-avatar">
                <img v-if="profile?.user.avatar" :src="profile.user.avatar" alt="Avatar" />
                <span v-else>{{ profileInitials }}</span>
              </a-avatar>
              <h2>{{ profileName }}</h2>
              <p class="identity-email">{{ profile?.user.email || '--' }}</p>
              
              <div class="identity-tags">
                <a-tag v-for="role in profile?.roles" :key="role.id" color="arcoblue" bordered size="large">
                  <template #icon><icon-safe /></template>
                  {{ role.role_name }}
                </a-tag>
              </div>
            </div>

            <a-divider />

            <div class="identity-details">
              <div class="detail-row">
                <icon-idcard />
                <span>{{ t('profile.username') }}: <strong>{{ profile?.user.username }}</strong></span>
              </div>
              <div class="detail-row">
                <icon-user />
                <span>{{ t('profile.gender') }}: <strong>{{ sexMap[profile?.user.sex ?? 2] }}</strong></span>
              </div>
              <div class="detail-row">
                <icon-clockCircle />
                <span>{{ t('profile.status') }}: <strong class="status-active">{{ profile?.user.status === 1 ? t('profile.statusDisabled') : t('profile.statusActive') }}</strong></span>
              </div>
            </div>

            <a-button type="primary" long class="edit-btn" @click="openEditModal">
              {{ t('profile.editAction') }}
            </a-button>
          </a-card>
        </a-col>

        <!-- Right Content Panels -->
        <a-col :xs="24" :lg="16" :xl="17">
          <!-- Organizational Assignments -->
          <a-card class="panel-card mb-4" :bordered="false" :title="t('profile.orgStructure') || 'Organizational structure'">
            <div v-if="profile?.dept_posts?.length" class="org-list">
              <a-row :gutter="16">
                <a-col v-for="dp in profile.dept_posts" :key="dp.dept_id + dp.post_id" :xs="24" :sm="12">
                  <div class="org-item">
                    <div class="org-icon"><icon-location /></div>
                    <div class="org-info">
                      <div class="org-dept">{{ dp.dept_name }}</div>
                      <div class="org-post">{{ dp.post_name }}</div>
                    </div>
                  </div>
                </a-col>
              </a-row>
            </div>
            <a-empty v-else :description="t('profile.noDepartments') || 'No departments assigned'" />
          </a-card>

          <!-- Passkey 凭据 -->
          <a-card class="panel-card mb-4" :bordered="false">
            <template #title>
              <icon-safe /> Passkey 凭据
            </template>
            <template #extra>
              <a-button
                type="primary"
                size="small"
                :loading="passkeyRegistering"
                :disabled="!passkeySupported"
                @click="handleAddPasskey"
              >
                <template #icon><icon-plus /></template>
                添加设备
              </a-button>
            </template>

            <a-spin :loading="passkeyLoadingList" style="width: 100%;">
              <div v-if="!passkeySupported" class="passkey-unsupported">
                当前浏览器不支持 Passkey
              </div>
              <div v-else-if="credentials.length === 0" class="passkey-empty">
                <p>暂无凭据</p>
                <p class="passkey-hint">添加后可使用生物识别快速登录</p>
              </div>
              <div v-else class="passkey-list">
                <div v-for="cred in credentials" :key="cred.id" class="passkey-item">
                  <div class="passkey-item-icon">
                    <icon-safe v-if="cred.attachment === 'platform'" />
                    <icon-lock v-else />
                  </div>
                  <div class="passkey-item-info">
                    <div class="passkey-item-name">
                      {{ cred.device_name || (cred.attachment === 'platform' ? '本设备' : '安全密钥') }}
                      <a-tag v-if="cred.backup_state" size="small" color="arcoblue">可同步</a-tag>
                      <a-tag v-if="cred.is_revoked" size="small" color="red">已吊销</a-tag>
                    </div>
                    <div class="passkey-item-meta">
                      {{ formatCredDate(cred.created_at) }}
                      <span v-if="cred.last_used_at"> · 最后使用 {{ formatCredDate(cred.last_used_at) }}</span>
                    </div>
                  </div>
                  <a-popconfirm
                    v-if="!cred.is_revoked"
                    content="确认移除此凭据？"
                    @ok="handleRevokePasskey(cred)"
                  >
                    <a-button type="text" status="danger" size="small">
                      <template #icon><icon-delete /></template>
                    </a-button>
                  </a-popconfirm>
                </div>
              </div>
            </a-spin>
          </a-card>

          <!-- Login Activity Logs -->
          <a-card class="panel-card" :bordered="false" :title="t('profile.recentLoginActivity') || 'Recent Login Activity'">
            <template #extra>
              <a-typography-text type="secondary">{{ t('profile.latest10Records') || 'Latest 10 records' }}</a-typography-text>
            </template>
            <a-table :data="loginLogs" :columns="logsColumns" :pagination="false" :bordered="false" stripe hoverable>
              <template #systemInfo="{ record }">
                <div class="sys-info">
                  <icon-desktop />
                  <span>{{ record.browser }} ({{ record.os }})</span>
                </div>
              </template>
              <template #status="{ record }">
                <DictTag :options="dicts.sys_success_status" :value="record.status" dot />
              </template>
              <template #loginAt="{ record }">
                {{ parseDate(record.loginAt) }}
              </template>
            </a-table>
          </a-card>

        </a-col>
      </a-row>

      <!-- Edit Modal -->
      <a-modal v-model:visible="visible" :title="t('profile.editAction') || 'Edit Profile'" :on-before-ok="handleBeforeOk" :ok-loading="saving" width="480px">
        <a-form ref="formRef" :model="formData" layout="vertical" class="profile-form">
          <a-row :gutter="16">
            <a-col :span="24">
              <a-form-item field="nickname" :label="t('profile.nickname') || 'Nickname'" :rules="[{ required: true, message: t('profile.nicknamePlaceholder') }]">
                <a-input v-model="formData.nickname" :placeholder="t('profile.nicknamePlaceholder')" />
              </a-form-item>
            </a-col>
            <a-col :span="24">
              <a-form-item field="email" :label="t('profile.email') || 'Email'" :rules="[{ type: 'email', required: true, message: t('profile.emailRequired') }]">
                <a-input v-model="formData.email" placeholder="user@example.com" />
              </a-form-item>
            </a-col>
            <a-col :xs="24" :sm="12">
              <a-form-item field="mobile" :label="t('profile.mobile') || 'Mobile Number'">
                <a-input v-model="formData.mobile" :placeholder="t('profile.mobilePlaceholder')" />
              </a-form-item>
            </a-col>
            <a-col :xs="24" :sm="12">
              <a-form-item field="sex" :label="t('profile.gender') || 'Gender'">
                <a-select v-model="formData.sex" :placeholder="t('profile.selectGender') || 'Select gender'">
                  <a-option :value="0">{{ t('profile.sexMale') }}</a-option>
                  <a-option :value="1">{{ t('profile.sexFemale') }}</a-option>
                  <a-option :value="2">{{ t('profile.sexUnknown') }}</a-option>
                </a-select>
              </a-form-item>
            </a-col>
          </a-row>
        </a-form>
      </a-modal>
    </a-spin>
  </div>
</template>

<style scoped>
.mb-4 { margin-bottom: 16px; }

.identity-card {
  text-align: center;
  padding: 16px 0;
}

.hero-avatar {
  background: linear-gradient(135deg, rgb(var(--arcoblue-5)), rgb(var(--arcoblue-6)));
  font-size: 32px;
  font-weight: 600;
  margin-bottom: 20px;
  box-shadow: 0 8px 24px rgba(var(--arcoblue-6), 0.2);
}

.identity-header h2 {
  margin: 0 0 4px 0;
  font-size: 22px;
  color: var(--color-text-1);
}

.identity-email {
  color: var(--color-text-3);
  margin-bottom: 16px;
}

.identity-tags {
  margin-bottom: 12px;
}

.identity-details {
  text-align: left;
  margin: 16px 0 24px 0;
}

.detail-row {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  color: var(--color-text-2);
  border-radius: 6px;
  margin-bottom: 4px;
}

.detail-row:hover {
  background-color: var(--color-fill-1);
}

.detail-row svg {
  margin-right: 12px;
  font-size: 16px;
  color: var(--color-text-3);
}

.detail-row span {
  flex: 1;
}

.status-active {
  color: rgb(var(--green-6));
}

.edit-btn {
  margin-top: 8px;
  border-radius: 8px;
}

.org-item {
  display: flex;
  align-items: center;
  padding: 16px;
  background-color: var(--color-fill-1);
  border-radius: 8px;
  margin-bottom: 12px;
  border: 1px solid var(--color-neutral-3);
  transition: all 0.2s;
}

.org-item:hover {
  background-color: var(--color-fill-2);
  border-color: rgba(var(--arcoblue-5), 0.3);
}

.org-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rgba(var(--arcoblue-5), 0.1);
  color: rgb(var(--arcoblue-6));
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  margin-right: 16px;
}

.org-dept {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-1);
  margin-bottom: 2px;
}

.org-post {
  font-size: 13px;
  color: var(--color-text-3);
}

.sys-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* ── Passkey ── */
.passkey-unsupported,
.passkey-empty {
  width: 100%;
  padding: 40px 0;
  text-align: center;
  color: var(--color-text-3);
}
.passkey-empty p {
  margin: 0;
}
.passkey-hint {
  font-size: 12px;
  margin-top: 8px !important;
}
.passkey-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.passkey-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  border-radius: 8px;
  border: 1px solid var(--color-border-1);
  transition: background 0.15s;
}
.passkey-item:hover {
  background: var(--color-fill-1);
}
.passkey-item-icon {
  font-size: 24px;
  flex-shrink: 0;
}
.passkey-item-info {
  flex: 1;
  min-width: 0;
}
.passkey-item-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-1);
  display: flex;
  align-items: center;
}
.passkey-item-meta {
  font-size: 12px;
  color: var(--color-text-3);
  margin-top: 2px;
}
</style>