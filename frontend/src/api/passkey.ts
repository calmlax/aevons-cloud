import { buildBasicAuth } from '@/api/auth';
import { request } from '@/utils/request';

export interface PasskeyBeginResponse {
  options: string;
  session_key: string;
}

export interface PasskeyCredential {
  id: string;
  user_id: string;
  username: string;
  device_name: string;
  aaguid: string;
  attachment: string;
  transports: string;
  device_type: string;
  backup_state: boolean;
  is_revoked: boolean;
  last_used_at: string | null;
  created_at: string;
}

function passkeyClientHeaders() {
  const clientId = import.meta.env.VITE_OAUTH_CLIENT_ID || '';
  const clientSecret = import.meta.env.VITE_OAUTH_CLIENT_SECRET || '';
  return {
    Authorization: buildBasicAuth(clientId, clientSecret),
  };
}

// ── 注册流程 ──────────────────────────────────────────────────────────────────

export const passkeyRegisterBegin = (): Promise<PasskeyBeginResponse> =>
  request({ url: '/auth/v1/passkey/register/begin', method: 'post' });

export const passkeyRegisterFinish = (sessionKey: string, credential: PublicKeyCredential): Promise<void> =>
  request({
    url: '/auth/v1/passkey/register/finish',
    method: 'post',
    data: { session_key: sessionKey, response: credentialToJSON(credential) },
  });

// ── 认证流程 ──────────────────────────────────────────────────────────────────

export const passkeyLoginBegin = (username?: string): Promise<PasskeyBeginResponse> =>
  request({
    url: '/auth/v1/passkey/login/begin',
    method: 'post',
    headers: passkeyClientHeaders(),
    data: username ? { username } : {},
    ...{ isToken: false } as any,
  });

export const passkeyLoginFinish = (sessionKey: string, credential: PublicKeyCredential): Promise<any> =>
  request({
    url: '/auth/v1/passkey/login/finish',
    method: 'post',
    headers: passkeyClientHeaders(),
    data: { session_key: sessionKey, response: credentialToJSON(credential) },
    ...{ isToken: false } as any,
  });

// ── 凭据管理 ──────────────────────────────────────────────────────────────────

export const passkeyListCredentials = (): Promise<PasskeyCredential[]> =>
  request({ url: '/auth/v1/passkey/credentials', method: 'get' });

export const passkeyRevokeCredential = (id: string): Promise<void> =>
  request({ url: `/auth/v1/passkey/credentials/${id}`, method: 'delete' });

// ── 工具：ArrayBuffer ↔ base64url ─────────────────────────────────────────────

function bufferToBase64url(buffer: ArrayBuffer): string {
  const bytes = new Uint8Array(buffer);
  let str = '';
  for (const b of bytes) str += String.fromCharCode(b);
  return btoa(str).replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '');
}

function base64urlToBuffer(b64: string): ArrayBuffer {
  const padded = b64.replace(/-/g, '+').replace(/_/g, '/').padEnd(b64.length + (4 - (b64.length % 4)) % 4, '=');
  const binary = atob(padded);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) bytes[i] = binary.charCodeAt(i);
  return bytes.buffer;
}

// ── 工具：PublicKeyCredential → JSON ─────────────────────────────────────────

function credentialToJSON(cred: PublicKeyCredential): Record<string, unknown> {
  const response = cred.response as AuthenticatorAttestationResponse | AuthenticatorAssertionResponse;
  const base: Record<string, unknown> = {
    id: cred.id,
    rawId: bufferToBase64url(cred.rawId),
    type: cred.type,
    clientExtensionResults: cred.getClientExtensionResults?.() ?? {},
  };

  if ('attestationObject' in response) {
    const r = response as AuthenticatorAttestationResponse;
    const transports = r.getTransports?.() ?? [];
    base.response = {
      attestationObject: bufferToBase64url(r.attestationObject),
      clientDataJSON: bufferToBase64url(r.clientDataJSON),
      transports,
    };
    // authenticatorAttachment 放顶层
    if ((cred as any).authenticatorAttachment) {
      base.authenticatorAttachment = (cred as any).authenticatorAttachment;
    }
  } else {
    const r = response as AuthenticatorAssertionResponse;
    base.response = {
      authenticatorData: bufferToBase64url(r.authenticatorData),
      clientDataJSON: bufferToBase64url(r.clientDataJSON),
      signature: bufferToBase64url(r.signature),
      userHandle: r.userHandle ? bufferToBase64url(r.userHandle) : null,
    };
  }
  return base;
}

// ── 工具：服务端 options JSON → 浏览器 API 格式 ──────────────────────────────
// go-webauthn 返回的结构是 { publicKey: { challenge, user, ... } }

export function parseCreationOptions(optionsJSON: string): PublicKeyCredentialCreationOptions {
  const outer = JSON.parse(optionsJSON);
  const opts = outer.publicKey ?? outer;
  return {
    ...opts,
    challenge: base64urlToBuffer(opts.challenge),
    user: { ...opts.user, id: base64urlToBuffer(opts.user.id) },
    excludeCredentials: (opts.excludeCredentials ?? []).map((c: any) => ({
      ...c,
      id: base64urlToBuffer(c.id),
    })),
  };
}

export function parseRequestOptions(optionsJSON: string): PublicKeyCredentialRequestOptions {
  const outer = JSON.parse(optionsJSON);
  const opts = outer.publicKey ?? outer;
  return {
    ...opts,
    challenge: base64urlToBuffer(opts.challenge),
    allowCredentials: (opts.allowCredentials ?? []).map((c: any) => ({
      ...c,
      id: base64urlToBuffer(c.id),
    })),
  };
}
