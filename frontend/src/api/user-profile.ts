import { request } from '@/utils/request';
import type { AuthRoleItem } from '@/types/auth'; // reuse basic types

export interface SysUser {
  user_id: string; // The backend JSON outputs "id,string" for SysUser base structure? Let's check model.SysUser: "id,string" -> "id", mapped in UserProfileResponse -> User
  id: string;
  username: string;
  nickname: string;
  email: string;
  mobile: string;
  sex: number;
  avatar: string;
  status: number;
}

export interface ProfileDeptPost {
  dept_id: string;
  dept_name: string;
  post_id: string;
  post_name: string;
}

export interface UserProfileResponse {
  user: SysUser;
  roles: AuthRoleItem[];
  dept_posts: ProfileDeptPost[];
  permissions: string[];
}

export interface UpdateProfilePayload {
  nickname: string;
  email: string;
  mobile?: string;
  sex?: number;
  avatar?: string;
}

export interface LoginLogItem {
  id: string;
  username: string;
  ipaddr: string;
  browser: string;
  os: string;
  status: number;
  msg: string;
  login_at: string;
}

/**
 * 获取当前登录用户的完整个人资料
 */
export const getUserProfileRequest = async (): Promise<UserProfileResponse> => {
  return request<UserProfileResponse>({ url: '/v1/auth/user/profile', method: 'get' });
};

/**
 * 更新个人资料
 */
export const updateUserProfileRequest = async (payload: UpdateProfilePayload) => {
  return request({ url: '/v1/auth/user/profile', method: 'put', data: payload });
};

/**
 * 获取当前登录用户的最新10条登录日志
 */
export const getProfileLoginLogsRequest = async (): Promise<LoginLogItem[]> => {
  return request<LoginLogItem[]>({ url: '/v1/auth/user/login-logs', method: 'get' });
};
