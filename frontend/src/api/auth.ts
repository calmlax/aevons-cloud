import settings from '@/config/settings.json';
import seedUsers from '@/mock/auth-users.json';
import { request } from '@/utils/request';
import type {
  AuthSession,
  AuthUser,
  AevonsAccount,
  LoginPayload,
  RegisterPayload,
  ResetPasswordPayload,
  StoredUserRecord,
  UpdatePasswordPayload,
  VerificationPurpose,
  VerificationRecord,
} from '@/types/auth';

export interface TokenPair {
  access_token: string;
  refresh_token?: string;
  token_type: string;
  expires_in: number;
  scope?: string;
}

export interface PublicKeyResp {
  key_id: string;
  public_key: string;
}

export function buildBasicAuth(clientId: string, clientSecret: string): string {
  const token = btoa(`${clientId}:${clientSecret}`);
  return `Basic ${token}`;
}

export const getPublicKey = (): Promise<PublicKeyResp> => {
  return request({ url: '/auth/v1/public-key', method: 'get', ...{ isToken: false } as any });
};

export const backendLoginRequest = async (
  clientId: string,
  clientSecret: string,
  payload: any
): Promise<TokenPair> => {
  return request({
    url: '/auth/v1/login',
    method: 'post',
    data: payload,
    headers: {
      Authorization: buildBasicAuth(clientId, clientSecret),
    },
    ...{ isToken: false } as any
  });
};

export const sendBackendEmailCode = async (email: string, purpose: string) => {
  return request({ url: '/auth/v1/email/code', method: 'post', data: { email, purpose }, ...{ isToken: false } as any });
};

export const getUserInfoRequest = async (): Promise<AuthUser> => {
  const data = await request<AuthUser>({ url: '/auth/v1/user', method: 'get' });
  return data;
};

export const getAuthorizeInfo = async (
  clientId: string,
  redirectUri: string,
): Promise<{
  client_id: string;
  client_name: string;
  logo_uri: string;
  scope: string;
  redirect_uri: string;
  state: string;
  autoapprove: boolean;
}> => {
  return request({ url: '/auth/v1/authorize', method: 'get', params: { client_id: clientId, redirect_uri: redirectUri, response_type: 'code' } });
};

export const approveAuthorize = async (params: {
  state: string;
  access_token?: string;
  scopes?: string[];
}): Promise<{ redirect_uri: string }> => {
  return request({ url: '/auth/v1/authorize', method: 'post', data: params });
};


const USERS_STORAGE_KEY = 'aevo-auth-users';
const SESSION_STORAGE_KEY = 'aevo-auth-session';
const CODES_STORAGE_KEY = 'aevo-auth-codes';
const AUTH_REQUEST_DELAY = Number(settings.mock?.authRequestDelay ?? 180);
const CODE_EXPIRES_IN = Number(settings.mock?.codeExpiresIn ?? 5 * 60 * 1000);
const SESSION_EXPIRES_IN = Number(settings.mock?.sessionExpiresIn ?? 8 * 60 * 60 * 1000);

const isBrowser = () => typeof window !== 'undefined';

const wait = (duration: number) => new Promise((resolve) => window.setTimeout(resolve, duration));

const createId = (prefix: string) => `${prefix}-${Math.random().toString(36).slice(2, 10)}-${Date.now().toString(36)}`;

const normalizeEmail = (email: string) => email.trim().toLowerCase();

const toPublicUser = (user: StoredUserRecord): AuthUser => {
  return user as unknown as AuthUser;
};

const readStorage = <T>(key: string, fallback: T): T => {
  if (!isBrowser()) {
    return fallback;
  }

  try {
    const raw = window.localStorage.getItem(key);
    return raw ? (JSON.parse(raw) as T) : fallback;
  } catch {
    return fallback;
  }
};

const writeStorage = <T>(key: string, value: T) => {
  if (!isBrowser()) {
    return;
  }

  window.localStorage.setItem(key, JSON.stringify(value));
};

const removeStorage = (key: string) => {
  if (!isBrowser()) {
    return;
  }

  window.localStorage.removeItem(key);
};

const cloneSeedUsers = (): StoredUserRecord[] => {
  return (seedUsers as StoredUserRecord[]).map((item) => ({
    ...item,
    permissions: [...item.permissions],
  }));
};

const syncUsersWithSeeds = (storedUsers: StoredUserRecord[]) => {
  const seedRecords = cloneSeedUsers();
  const seedIds = new Set(seedRecords.map((item) => item.id));
  const customUsers = storedUsers.filter((item) => !seedIds.has(item.id));

  return [...seedRecords, ...customUsers];
};

const readUsers = (): StoredUserRecord[] => {
  const users = readStorage<StoredUserRecord[]>(USERS_STORAGE_KEY, []);
  if (users.length > 0) {
    const syncedUsers = syncUsersWithSeeds(users);
    writeStorage(USERS_STORAGE_KEY, syncedUsers);
    return syncedUsers;
  }

  const nextUsers = cloneSeedUsers();
  writeStorage(USERS_STORAGE_KEY, nextUsers);
  return nextUsers;
};

const writeUsers = (users: StoredUserRecord[]) => {
  writeStorage(USERS_STORAGE_KEY, users);
};

const readCodes = (): VerificationRecord[] => {
  const now = Date.now();
  const records = readStorage<VerificationRecord[]>(CODES_STORAGE_KEY, []).filter((item) => item.expiresAt > now);
  writeStorage(CODES_STORAGE_KEY, records);
  return records;
};

const writeCodes = (codes: VerificationRecord[]) => {
  writeStorage(CODES_STORAGE_KEY, codes);
};

const readSession = (): AuthSession | null => {
  const session = readStorage<AuthSession | null>(SESSION_STORAGE_KEY, null);
  if (!session) {
    return null;
  }

  if (session.expiresAt <= Date.now()) {
    removeStorage(SESSION_STORAGE_KEY);
    return null;
  }

  return session;
};

const writeSession = (user: StoredUserRecord | null) => {
  if (!user) {
    removeStorage(SESSION_STORAGE_KEY);
    return null;
  }

  const session: AuthSession = {
    token: createId('token'),
    userId: user.id,
    expiresAt: Date.now() + SESSION_EXPIRES_IN,
  };

  writeStorage(SESSION_STORAGE_KEY, session);
  return session;
};

const findUserByEmail = (users: StoredUserRecord[], email: string) => {
  const normalizedEmail = normalizeEmail(email);
  return users.find((item) => normalizeEmail(item.email) === normalizedEmail);
};

const verifyCode = (email: string, purpose: VerificationPurpose, code: string) => {
  const normalizedEmail = normalizeEmail(email);
  const records = readCodes();
  const matchedRecord = records.find((item) => item.purpose === purpose && normalizeEmail(item.email) === normalizedEmail);

  if (!matchedRecord) {
    throw new Error('auth.codeExpired');
  }

  if (matchedRecord.code !== code.trim()) {
    throw new Error('auth.codeInvalid');
  }

  writeCodes(
    records.filter((item) => !(item.purpose === purpose && normalizeEmail(item.email) === normalizedEmail))
  );
};

export const initializeMockAuth = () => {
  readUsers();
  readCodes();
  readSession();
};

export const getCurrentUserSnapshot = (): AuthUser | null => {
  initializeMockAuth();
  const session = readSession();
  if (!session) {
    return null;
  }

  const matchedUser = readUsers().find((item) => item.id === session.userId);
  return matchedUser ? toPublicUser(matchedUser) : null;
};

export const getCurrentPermissionsSnapshot = () => {
  return getCurrentUserSnapshot()?.permissions ?? [];
};

export const hasPermissionSnapshot = (required?: string | string[]) => {
  if (!required) {
    return true;
  }

  const currentUser = getCurrentUserSnapshot();
  if (!currentUser) {
    return false;
  }

  if (currentUser.roles?.some(r => ['admin', 'super_admin', 'super-admin'].includes(r.role_key))) {
    return true;
  }

  if (currentUser.permissions?.includes('*') || currentUser.permissions?.includes('*:*:*')) {
    return true;
  }

  const requiredList = Array.isArray(required) ? required : [required];
  return requiredList.every((item) => currentUser.permissions.includes(item));
};

export const fetchAevonsAccounts = (): AevonsAccount[] => {
  return cloneSeedUsers().map((item) => ({
    id: item.id,
    name: item.name,
    email: item.email,
    password: item.password,
    role: item.role,
    permissions: [...item.permissions],
  }));
};

export const loginRequest = async ({ email, password }: LoginPayload) => {
  initializeMockAuth();
  await wait(AUTH_REQUEST_DELAY);

  const users = readUsers();
  const matchedUser = findUserByEmail(users, email);

  if (!matchedUser || matchedUser.password !== password) {
    throw new Error('auth.invalidCredentials');
  }

  if (!matchedUser.verified) {
    throw new Error('auth.emailNotVerified');
  }

  matchedUser.lastLoginAt = new Date().toISOString();
  writeUsers(users);
  const session = writeSession(matchedUser);

  return {
    token: session?.token ?? '',
    user: toPublicUser(matchedUser),
  };
};

export const logoutRequest = async (): Promise<void> => {
  await request({ url: '/auth/v1/logout', method: 'post' });
};

export const sendEmailCodeRequest = async (email: string, purpose: VerificationPurpose) => {
  initializeMockAuth();
  await wait(AUTH_REQUEST_DELAY);

  const users = readUsers();
  const matchedUser = findUserByEmail(users, email);

  if (purpose === 'register' && matchedUser) {
    throw new Error('auth.emailAlreadyRegistered');
  }

  if (purpose === 'reset-password' && !matchedUser) {
    throw new Error('auth.emailNotFound');
  }

  const code = String(Math.floor(100000 + Math.random() * 900000));
  const normalizedEmail = normalizeEmail(email);
  const nextCodes = readCodes().filter((item) => !(item.purpose === purpose && normalizeEmail(item.email) === normalizedEmail));

  nextCodes.push({
    email: normalizedEmail,
    code,
    purpose,
    expiresAt: Date.now() + CODE_EXPIRES_IN,
  });

  writeCodes(nextCodes);
  return { code, expiresIn: CODE_EXPIRES_IN };
};

export const registerRequest = async (payload: RegisterPayload) => {
  return request({ url: '/auth/v1/register', method: 'post', data: payload, ...{ isToken: false } as any });
};

export const resetPasswordRequest = async (payload: ResetPasswordPayload) => {
  return request({ url: '/auth/v1/reset-password', method: 'post', data: payload, ...{ isToken: false } as any });
};

export const updatePasswordRequest = async (payload: UpdatePasswordPayload) => {
  return request({ url: '/auth/v1/user/password', method: 'put', data: payload, ...{ isToken: true } as any });
};

export interface UpdateProfilePayload {
  nickname: string;
  email: string;
}

export const updateProfileRequest = async (payload: UpdateProfilePayload) => {
  return request({ url: '/auth/v1/user/profile', method: 'put', data: payload, ...{ isToken: true } as any });
};