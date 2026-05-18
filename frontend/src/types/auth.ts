export type AuthRole = 'super-admin' | 'operator' | 'auditor';
export type VerificationPurpose = 'register' | 'reset-password';

export interface StoredUserRecord {
  id: string;
  name: string;
  avatar?: string;
  email: string;
  password: string;
  role: AuthRole;
  permissions: string[];
  verified: boolean;
  createdAt: string;
  lastLoginAt: string;
}

export interface AuthSession {
  token: string;
  userId: string;
  expiresAt: number;
}

export interface VerificationRecord {
  email: string;
  code: string;
  purpose: VerificationPurpose;
  expiresAt: number;
}

export interface LoginPayload {
  email: string;
  password: string;
}

export interface RegisterPayload {
  name: string;
  email: string;
  password: string;
  code: string;
}

export interface ResetPasswordPayload {
  email: string;
  password: string;
  code: string;
}

export interface UpdatePasswordPayload {
  currentPassword: string;
  nextPassword: string;
}

export interface AevonsAccount {
  id: string;
  name: string;
  email: string;
  password: string;
  role: AuthRole;
  permissions: string[];
}

export interface AuthRoleItem {
  id: string;
  role_key: string;
  role_name: string;
  data_scope: string;
  dept_ids: string[];
}

export interface AuthDeptItem {
  dept_id: string;
  post_id: string;
}

export interface AuthUser {
  user_id: string;
  username: string;
  nickname: string;
  email: string;
  avatar: string;
  status: number;
  roles: AuthRoleItem[];
  depts: AuthDeptItem[];
  permissions: string[];
}