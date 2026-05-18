// 完全对应后端 DictDataDTO 结构体
export interface DictOption {
  id?: string | number;        // 字典编号（后端返回string/number）
  dictType?: string;           // 字典类型
  dictValue?: string | number; // 字典键值 → 作为组件匹配的value
  value?: string | number;     // 兼容前端静态选项
  sort?: string | number;      // 顺序
  tagType?: string;           // 标签风格（success/warning/danger/info/default）
  color?: string;              // 兼容前端静态选项
  tagClass?: string;          // 样式类名
  icon?: string;              // 图标
  langCode?: string;           // 语言标识
  label?: string;              // 标签翻译 → 作为组件显示的label
  text?: string;               // 兼容旧示例数据
  tip?: string;               // 提示翻译
}

export type DictMap = Record<string, DictOption[]>;

// 内置预设颜色（不变）
export const PRESET_COLORS: Record<string, { bg: string; text: string; border: string }> = {
  success: {
    bg: 'rgba(16, 185, 129, 0.12)',
    text: '#059669',
    border: 'rgba(16, 185, 129, 0.28)',
  },
  warning: {
    bg: 'rgba(245, 158, 11, 0.12)',
    text: '#d97706',
    border: 'rgba(245, 158, 11, 0.28)',
  },
  danger: {
    bg: 'rgba(239, 68, 68, 0.12)',
    text: '#dc2626',
    border: 'rgba(239, 68, 68, 0.28)',
  },
  info: {
    bg: 'rgba(37, 99, 235, 0.1)',
    text: '#2563eb',
    border: 'rgba(37, 99, 235, 0.24)',
  },
  default: {
    bg: 'rgba(148, 163, 184, 0.12)',
    text: 'var(--text-subtle)',
    border: 'rgba(148, 163, 184, 0.24)',
  },
  green: {
    bg: 'rgba(16, 185, 129, 0.12)',
    text: '#059669',
    border: 'rgba(16, 185, 129, 0.28)',
  },
  orange: {
    bg: 'rgba(245, 158, 11, 0.12)',
    text: '#d97706',
    border: 'rgba(245, 158, 11, 0.28)',
  },
  red: {
    bg: 'rgba(239, 68, 68, 0.12)',
    text: '#dc2626',
    border: 'rgba(239, 68, 68, 0.28)',
  },
  blue: {
    bg: 'rgba(37, 99, 235, 0.1)',
    text: '#2563eb',
    border: 'rgba(37, 99, 235, 0.24)',
  },
  purple: {
    bg: 'rgba(139, 92, 246, 0.1)',
    text: '#7c3aed',
    border: 'rgba(139, 92, 246, 0.24)',
  },
  cyan: {
    bg: 'rgba(6, 182, 212, 0.1)',
    text: '#0891b2',
    border: 'rgba(6, 182, 212, 0.24)',
  },
};
