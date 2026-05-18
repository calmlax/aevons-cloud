import { createHighlighter, type Highlighter } from 'shiki';

// Singleton highlighter — created once, reused across calls
let _highlighter: Highlighter | null = null;

const LANGS = ['go', 'typescript', 'vue', 'sql', 'yaml', 'json', 'xml', 'html', 'javascript', 'bash', 'text'];

export async function getHighlighter(): Promise<Highlighter> {
  if (_highlighter) return _highlighter;
  _highlighter = await createHighlighter({
    themes: ['github-dark', 'github-light'],
    langs: LANGS,
  });
  return _highlighter;
}

/** Map common fileType strings to shiki language ids */
export function resolveLanguage(fileType: string): string {
  const map: Record<string, string> = {
    go: 'go',
    ts: 'typescript',
    typescript: 'typescript',
    vue: 'vue',
    sql: 'sql',
    yaml: 'yaml',
    yml: 'yaml',
    json: 'json',
    xml: 'xml',
    html: 'html',
    js: 'javascript',
    javascript: 'javascript',
    sh: 'bash',
    bash: 'bash',
  };
  return map[fileType.toLowerCase()] ?? 'text';
}
