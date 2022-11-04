export enum DatasourceType {
    HttpRequest = 'HttpRequest',
    HttpServer = 'HttpServer',
    RandomValue = 'RandomValue',
    Sql = 'Sql',
}

export enum DataPointLoggingType {
    AllData = 'AllData',
    WhenValueChanges = 'WhenValueChanges',
    TimeStampChanges = 'TimeStampChanges',
}