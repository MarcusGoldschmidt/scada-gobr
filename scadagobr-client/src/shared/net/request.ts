// TODO: refactor to computed values
export const PathsV1 = {
    // Auth
    Login: 'api/v1/auth/login',
    RefreshToken: 'api/v1/auth/refresh-token',
    WhoAmI: 'api/v1/auth/who-am-i',
    // User
    UserGet: 'api/v1/user',
    UserCreate: 'api/v1/user',
    UserDelete: 'api/v1/user/',
    UserUpdate: 'api/v1/user/',
    UserGetById: 'api/v1/user/',
    // Views
    ViewGet: 'api/v1/view',
    ViewCreate: 'api/v1/view',
    ViewDelete: 'api/v1/view/',
    ViewUpdate: 'api/v1/view/',
    ViewGetById: 'api/v1/view/',
    // View Components
    ViewComponentDelete: (viewId, componentId) => `/api/v1/view/${viewId}/component/${componentId}`,
    // DataSources
    DataSourceGet: 'api/v1/datasource',
    DataSourceCreate: 'api/v1/datasource',
    DataSourceDelete: 'api/v1/datasource/',
    DataSourceUpdate: 'api/v1/datasource/',
    DataSourceGetById: 'api/v1/datasource/',
    // DataPoints
    DataPointGet: (dataSourceId) => `/api/v1/datasource/${dataSourceId}/datapoint`,
    DataPointCreate: (dataSourceId) => `/api/v1/datasource/${dataSourceId}/datapoint`,
    DataPointDelete: (dataSourceId, dataPointId) => `/api/v1/datasource/${dataSourceId}/datapoint/${dataPointId}`,
    DataPointUpdate: (dataSourceId, dataPointId) => `/api/v1/datasource/${dataSourceId}/datapoint/${dataPointId}`,

} as const;

export const Request = {
    V1: PathsV1,
} as const

export enum PathsWsV1 {
    DataPoint = "/api/v1/datapoints/ws/",
}

