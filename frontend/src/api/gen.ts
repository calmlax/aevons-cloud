import { request } from '@/utils/request';

export interface GenTable {
  id: string;
  tableName: string;
  tableComment: string;
  className: string;
  moduleName: string;
  author: string;
  router: string;
  base_class: string;
  menuId: string;
  permission: string;
  remark: string;
}

export interface GenTableColumn {
  id: string;
  tableId: string;
  columnName: string;
  columnComment: string;
  columnType: string;
  dataType: string;
  fieldName: string;
  json: string;
  isPrimaryKey: boolean;
  isAutoIncrement: boolean;
  isRequired: boolean;
  isInsert: boolean;
  isEdit: boolean;
  queryType: string;
  dictType: string;
  sort: number;
  component: string;
  defaultValue: string;
  dataLength: number;
  dataPrecision: number;
}

export interface TablePageQuery {
  pageNum?: number;
  pageSize?: number;
  tableName?: string;
  tableComment?: string;
  moduleName?: string;
}

export interface PageResult<T> {
  rows: T[];
  total: number;
}

export interface DBTableInfo {
  tableName: string;
  tableComment: string;
}

// ── Table API ──────────────────────────────────────────
export const getTablePage = (params: TablePageQuery) =>
  request<PageResult<GenTable>>({ url: '/v1/gen/table/page', method: 'get', params });

export const getTableList = (params?: { tableName?: string; tableComment?: string; moduleName?: string }) =>
  request<GenTable[]>({ url: '/v1/gen/table/list', method: 'get', params });

export const getTableById = (id: string) =>
  request<GenTable>({ url: `/v1/gen/table/${id}`, method: 'get' });

export const updateTable = (id: string, data: Partial<GenTable>) =>
  request({ url: `/v1/gen/table/${id}`, method: 'put', data });

export const deleteTable = (ids: string[]) =>
  request({ url: `/v1/gen/table/${ids.join(',')}`, method: 'delete' });

// ── Import API ─────────────────────────────────────────
export const getDBTableList = () =>
  request<DBTableInfo[]>({ url: '/v1/gen/table/db', method: 'get' });

export const importTables = (tableNames: string[]) =>
  request({ url: '/v1/gen/table/import', method: 'post', data: tableNames });

// ── TableColumn API ────────────────────────────────────
export const getColumnList = (params: { tableId: string; columnName?: string; columnComment?: string }) =>
  request<GenTableColumn[]>({ url: '/v1/gen/table/column/list', method: 'get', params });

export const getColumnById = (id: string) =>
  request<GenTableColumn>({ url: `/v1/gen/table/column/${id}`, method: 'get' });

export const batchUpdateColumn = (data: Partial<GenTableColumn>[]) =>
  request({ url: `/v1/gen/table/column/batch-update`, method: 'put', data });

export const deleteColumn = (ids: string[]) =>
  request({ url: `/v1/gen/table/column/${ids.join(',')}`, method: 'delete' });

// ── Preview & Download API ─────────────────────────────
export interface PreviewFile {
  name: string;
  code: string;
  fileType: string;
}

export const previewCode = (tableId: string) =>
  request<PreviewFile[]>({ url: '/v1/gen/table/preview', method: 'get', params: { tableId } });

export const synchTable = (tableId: string) =>
  request({ url: '/v1/gen/table/synch', method: 'get', params: { tableId } });

export const downloadCode = (tableIds: string[]) =>
  request<Blob>({ url: '/v1/gen/table/download', method: 'get', params: { tableIds: tableIds.join(',') }, responseType: 'blob' });
