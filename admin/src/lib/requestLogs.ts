
export type RequestLog = {
    id: string
    url: string
    method: string
    proto: string
    headers: HTTPHeader[]
    body?: string
    timestamp: string
    response?: ResponseLog
}

export type ResponseLog = {
    proto: string
    statusCode: number
    statusReason: string
    body?: string
    headers: HTTPHeader[]
}

export type HTTPHeader = {
    key: string
    value: string
}