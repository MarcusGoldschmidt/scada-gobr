export interface Env {
    API_BASE_URL: string
}

// @ts-ignore
// Rollup replace plugin will create a json of type Env
const env: Env = ROLLUP_REPLACE_ENVIROMENT

export default env

