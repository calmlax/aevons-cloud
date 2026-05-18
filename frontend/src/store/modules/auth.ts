import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import {
  fetchAevonsAccounts,
  getCurrentPermissionsSnapshot,
  getCurrentUserSnapshot,
  hasPermissionSnapshot,
  backendLoginRequest,
  getUserInfoRequest,
  getPublicKey,
  logoutRequest,
  registerRequest,
  resetPasswordRequest,
  sendEmailCodeRequest,
  sendBackendEmailCode,
  updatePasswordRequest,
  updateProfileRequest,
  type UpdateProfilePayload,
} from '@/api/auth';
import type {
  AuthUser,
  AevonsAccount,
  LoginPayload,
  RegisterPayload,
  ResetPasswordPayload,
  UpdatePasswordPayload,
  VerificationPurpose,
} from '@/types/auth';
import { useMenuStore } from './menu';

export const useAuthStore = defineStore('auth', () => {
  const readStoredUser = (): AuthUser | null => {
    const raw = window.localStorage.getItem('aevo-user');
    return raw ? JSON.parse(raw) : null;
  };

  const authUser = ref<AuthUser | null>(readStoredUser());
  const accessToken = ref<string | null>(window.localStorage.getItem('aevo-token'));
  const initialized = ref(false);
  const aevoAccounts = ref<AevonsAccount[]>([]);

  const initializeAuth = () => {
    if (initialized.value) {
      return authUser.value;
    }

    aevoAccounts.value = fetchAevonsAccounts();
    initialized.value = true;
    return authUser.value;
  };

  const login = async (payload: LoginPayload & { encryptedPassword?: string; key_id?: string; code?: string; grant_type?: string }) => {
    // initializeAuth(); // No longer necessary since we use the real backend

    const clientId = import.meta.env.VITE_OAUTH_CLIENT_ID || '';
    const clientSecret = import.meta.env.VITE_OAUTH_CLIENT_SECRET || '';

    const reqPayload = {
      grant_type: payload.grant_type || 'password',
      username: payload.email, // using email as username
      password: payload.encryptedPassword || payload.password,
      key_id: payload.key_id,
      email: payload.email,
      code: payload.code,
    };

    const tokenPair = await backendLoginRequest(clientId, clientSecret, reqPayload);

    accessToken.value = tokenPair.access_token;
    window.localStorage.setItem('aevo-token', tokenPair.access_token);

    // Fetch real profile from backend
    authUser.value = await getUserInfoRequest();

    window.localStorage.setItem('aevo-user', JSON.stringify(authUser.value));

    // Web application assumes sync state. `accessToken` triggers request updates, 
    // but the actual interceptor looks at authStore.accessToken

    return authUser.value;
  };

  const clearAuth = (redirect = true) => {
    authUser.value = null;
    accessToken.value = null;
    window.localStorage.removeItem('aevo-token');
    window.localStorage.removeItem('aevo-user');
    const menuStore = useMenuStore();
    menuStore.resetMenu();
    // if (redirect && window.location.pathname !== '/login') {
    //   // 延迟跳转，让 toast 有时间显示
    //   setTimeout(() => { window.location.href = '/login'; }, 300);
    // }
  };

  const logout = async () => {
    try {
      await logoutRequest();
    } catch {
      // 即使后端返回 401 或其他错误，也要清除本地状态
    } finally {
      clearAuth(false); // 不走延迟，调用方自行跳转
      //window.location.href = '/login';
    }
  };

  const sendEmailCode = async (email: string, purpose: VerificationPurpose) => {
    await sendBackendEmailCode(email, purpose);
    return '';
  };

  const sendRealEmailCode = async (email: string, purpose: string) => {
    return await sendBackendEmailCode(email, purpose);
  };

  const register = async (payload: RegisterPayload) => {
    const user = await registerRequest(payload);
    return user;
  };

  const resetPassword = async (payload: ResetPasswordPayload) => {
    return resetPasswordRequest(payload);
  };

  const updatePassword = async (payload: UpdatePasswordPayload) => {
    await updatePasswordRequest(payload);
    // Password changed, server revokes all tokens, clear local state
    clearAuth();
  };

  const updateProfile = async (payload: UpdateProfilePayload) => {
    await updateProfileRequest(payload);
    await refreshUserInfo();
  };

  const refreshUserInfo = async () => {
    authUser.value = await getUserInfoRequest();
    window.localStorage.setItem('aevo-user', JSON.stringify(authUser.value));
  };

  const hasPermission = (required?: string | string[]) => {
    if (!required) {
      return true;
    }

    const user = authUser.value;
    if (!user) {
      return false;
    }

    if (user.permissions?.includes('*') || user.permissions?.includes('*:*:*')) {
      return true;
    }

    if (user.roles?.some(r => ['admin', 'super_admin', 'super-admin'].includes(r.role_key))) {
      return true;
    }

    const requiredList = Array.isArray(required) ? required : [required];
    return requiredList.every((item) => user.permissions.includes(item));
  };

  const hasAnyPermission = (required?: string[]) => {
    if (!required?.length) {
      return true;
    }
    return required.some((item) => hasPermission(item));
  };

  const isAuthenticated = computed(() => Boolean(authUser.value));

  const setToken = (token: string) => {
    accessToken.value = token;
    window.localStorage.setItem('aevo-token', token);
  };

  const getCurrentUser = () => authUser.value;

  const getAuthHomePath = () => {
    const user = authUser.value;

    if (!user) {
      return '/login';
    }
    // return '/account/profile';
    return '/dashboard';
  };

  const permissions = computed(() => authUser.value?.permissions ?? []);

  return {
    currentUser: authUser,
    initialized,
    aevoAccounts,
    permissions,
    isAuthenticated,
    accessToken,
    setToken,
    initializeAuth,
    login,
    logout,
    register,
    resetPassword,
    sendEmailCode,
    sendRealEmailCode,
    updatePassword,
    updateProfile,
    refreshUserInfo,
    hasPermission,
    hasAnyPermission,
    getCurrentUser,
    getAuthHomePath,
    clearAuth,
  };
});
